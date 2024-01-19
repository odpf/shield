package plan

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/raystack/frontier/pkg/metadata"

	"github.com/raystack/frontier/billing/product"
	"github.com/raystack/frontier/pkg/utils"
	"github.com/stripe/stripe-go/v75/client"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (Plan, error)
	GetByName(ctx context.Context, name string) (Plan, error)
	Create(ctx context.Context, plan Plan) (Plan, error)
	UpdateByName(ctx context.Context, plan Plan) (Plan, error)
	List(ctx context.Context, filter Filter) ([]Plan, error)
}

type ProductService interface {
	Create(ctx context.Context, p product.Product) (product.Product, error)
	GetByID(ctx context.Context, id string) (product.Product, error)
	Update(ctx context.Context, p product.Product) (product.Product, error)
	AddPlan(ctx context.Context, p product.Product, planID string) error

	CreatePrice(ctx context.Context, price product.Price) (product.Price, error)
	UpdatePrice(ctx context.Context, price product.Price) (product.Price, error)
	GetPriceByID(ctx context.Context, id string) (product.Price, error)
	GetPriceByProductID(ctx context.Context, id string) ([]product.Price, error)

	List(ctx context.Context, flt product.Filter) ([]product.Product, error)

	UpsertFeature(ctx context.Context, f product.Feature) (product.Feature, error)
	GetFeatureByID(ctx context.Context, id string) (product.Feature, error)
	GetFeatureByProductID(ctx context.Context, id string) ([]product.Feature, error)
}

type Service struct {
	planRepository Repository
	stripeClient   *client.API
	productService ProductService
}

func NewService(stripeClient *client.API, planRepository Repository, productService ProductService) *Service {
	return &Service{
		stripeClient:   stripeClient,
		planRepository: planRepository,
		productService: productService,
	}
}

func (s Service) Create(ctx context.Context, p Plan) (Plan, error) {
	p.Name = strings.ToLower(p.Name)
	p.Interval = strings.ToLower(p.Interval)
	return s.planRepository.Create(ctx, p)
}

func (s Service) GetByID(ctx context.Context, id string) (Plan, error) {
	var fetchedPlan Plan
	var err error
	if utils.IsValidUUID(id) {
		fetchedPlan, err = s.planRepository.GetByID(ctx, id)
	} else {
		fetchedPlan, err = s.planRepository.GetByName(ctx, id)
	}
	if err != nil {
		return Plan{}, err
	}

	// enrich with product
	products, err := s.productService.List(ctx, product.Filter{
		PlanID: fetchedPlan.ID,
	})
	if err != nil {
		return Plan{}, err
	}
	fetchedPlan.Products = products
	return fetchedPlan, nil
}

func (s Service) List(ctx context.Context, filter Filter) ([]Plan, error) {
	listedPlans, err := s.planRepository.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	// enrich with product
	for i, listedPlan := range listedPlans {
		// TODO(kushsharma): we can do this in one query
		products, err := s.productService.List(ctx, product.Filter{
			PlanID: listedPlan.ID,
		})
		if err != nil {
			return nil, err
		}
		listedPlans[i].Products = products
	}
	return listedPlans, nil
}

func (s Service) UpsertPlans(ctx context.Context, planFile File) error {
	// keep a list of product to feature list to ensure features are only
	// attached to the product they belong to
	featureToProduct := make(map[string][]string)

	// ensure features
	for _, featureToCreate := range planFile.Features {
		featureOb, err := s.productService.UpsertFeature(ctx, product.Feature{
			ID:         featureToCreate.ID,
			Name:       featureToCreate.Name,
			ProductIDs: featureToCreate.ProductIDs,
			Metadata:   metadata.Build(featureToCreate.Metadata),
		})
		if err != nil {
			return err
		}
		featureToProduct[featureOb.ID] = []string{}
	}

	// create products
	for _, productToCreate := range planFile.Products {
		productOb, err := s.productService.GetByID(ctx, productToCreate.Name)
		if err != nil && errors.Is(err, product.ErrProductNotFound) {
			// create product
			if productOb, err = s.productService.Create(ctx, product.Product{
				Name:        productToCreate.Name,
				Title:       productToCreate.Title,
				Description: productToCreate.Description,
				Config:      productToCreate.Config,
				Behavior:    productToCreate.Behavior,
				Metadata:    metadata.Build(productToCreate.Metadata),
			}); err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else {
			// update product
			if _, err = s.productService.Update(ctx, product.Product{
				ID:          productOb.ID,
				ProviderID:  productOb.ProviderID,
				Name:        productToCreate.Name,
				Title:       productToCreate.Title,
				Description: productToCreate.Description,
				Config:      productToCreate.Config,
			}); err != nil {
				return err
			}
		}

		// ensure price exists
		for blobIdx, priceToCreate := range productToCreate.Prices {
			if priceToCreate.Name == "" {
				priceToCreate.Name = fmt.Sprintf("default_%d", blobIdx)
			}
			priceObs, err := s.productService.GetPriceByProductID(ctx, productOb.ID)
			if err != nil {
				return fmt.Errorf("failed to get price by product id: %w", err)
			}
			// find price by name
			var priceOb product.Price
			for _, p := range priceObs {
				if p.Name == priceToCreate.Name {
					priceOb = p
					break
				}
			}
			if priceOb.ID == "" {
				// create price
				if priceOb, err = s.productService.CreatePrice(ctx, product.Price{
					Name:             priceToCreate.Name,
					Amount:           priceToCreate.Amount,
					Currency:         priceToCreate.Currency,
					BillingScheme:    priceToCreate.BillingScheme,
					UsageType:        priceToCreate.UsageType,
					MeteredAggregate: priceToCreate.MeteredAggregate,
					Interval:         priceToCreate.Interval,
					ProductID:        productOb.ID,
					Metadata:         metadata.Build(priceToCreate.Metadata),
				}); err != nil {
					return err
				}
			} else {
				// update price
				if _, err = s.productService.UpdatePrice(ctx, product.Price{
					ID:         priceOb.ID,
					ProviderID: priceOb.ProviderID,
					ProductID:  priceOb.ProductID,
					Name:       priceOb.Name,
				}); err != nil {
					return err
				}
			}
		}

		// ensure feature exists
		for _, featureToCreate := range productToCreate.Features {
			featureOb, err := s.productService.UpsertFeature(ctx, product.Feature{
				ID:         featureToCreate.ID,
				Name:       featureToCreate.Name,
				ProductIDs: featureToCreate.ProductIDs,
			})
			if err != nil {
				return err
			}
			featureToProduct[featureOb.ID] = append(featureToProduct[featureOb.ID], productOb.ID)
		}
	}

	// ensure feature is added to product and removed from other products where
	// it's no longer needed
	for featureID, productIDs := range featureToProduct {
		featureOb, err := s.productService.GetFeatureByID(ctx, featureID)
		if err != nil {
			return err
		}
		featureOb.ProductIDs = productIDs
		if _, err = s.productService.UpsertFeature(ctx, featureOb); err != nil {
			return err
		}
	}

	// create plans
	for _, planToCreate := range planFile.Plans {
		// ensure plan exists
		planOb, err := s.GetByID(ctx, planToCreate.Name)
		if err != nil && errors.Is(err, ErrNotFound) {
			// create plan
			if planOb, err = s.planRepository.Create(ctx, Plan{
				Name:           planToCreate.Name,
				Title:          planToCreate.Title,
				Description:    planToCreate.Description,
				OnStartCredits: planToCreate.OnStartCredits,
				Interval:       planToCreate.Interval,
				Metadata:       metadata.Build(planToCreate.Metadata),
			}); err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else {
			// update plan
			if _, err = s.planRepository.UpdateByName(ctx, Plan{
				ID:             planOb.ID,
				Name:           planToCreate.Name,
				Title:          planToCreate.Title,
				OnStartCredits: planToCreate.OnStartCredits,
				Description:    planToCreate.Description,
			}); err != nil {
				return err
			}
		}

		// ensure only one product has user count behavior
		if len(utils.Filter(planToCreate.Products, func(f product.Product) bool {
			return f.Behavior == product.PerSeatBehavior
		})) > 1 {
			return fmt.Errorf("plan %s has more than one product with per_seat behavior", planOb.Name)
		}

		// ensure product exists, if not fail
		for _, productToCreate := range planToCreate.Products {
			productOb, err := s.productService.GetByID(ctx, productToCreate.Name)
			if err != nil {
				return err
			}

			// ensure plan can be added to product
			hasMatchingPrice := utils.ContainsFunc(productOb.Prices, func(p product.Price) bool {
				return p.Interval == planOb.Interval
			})
			if !hasMatchingPrice {
				return fmt.Errorf("product %s has no prices registered with this interval, plan %s has interval %s",
					productOb.Name, planOb.Name, planOb.Interval)
			}
			if err = s.productService.AddPlan(ctx, productOb, planOb.ID); err != nil {
				return err
			}
		}
	}

	return nil
}
