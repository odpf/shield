package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/raystack/frontier/core/domain"
	"github.com/raystack/frontier/pkg/db"
	"github.com/raystack/salt/log"
)

type DomainRepository struct {
	log log.Logger
	dbc *db.Client
	Now func() time.Time
}

func NewDomainRepository(logger log.Logger, dbc *db.Client) *DomainRepository {
	return &DomainRepository{
		dbc: dbc,
		log: logger,
		Now: func() time.Time {
			return time.Now().UTC()
		},
	}
}

func (s *DomainRepository) Create(ctx context.Context, toCreate domain.Domain) (domain.Domain, error) {
	query, params, err := dialect.Insert(TABLE_DOMAINS).Rows(
		goqu.Record{
			"org_id": toCreate.OrgID,
			"name":   toCreate.Name,
			"token":  toCreate.Token,
		}).Returning(&Domain{}).ToSQL()
	if err != nil {
		return domain.Domain{}, fmt.Errorf("%w: %s", parseErr, err)
	}

	var domainModel Domain
	if err = s.dbc.WithTimeout(ctx, TABLE_DOMAINS, "Create", func(ctx context.Context) error {
		return s.dbc.QueryRowxContext(ctx, query, params...).StructScan(&domainModel)
	}); err != nil {
		err = checkPostgresError(err)
		if errors.Is(err, ErrDuplicateKey) {
			return domain.Domain{}, domain.ErrDuplicateKey
		}
		return domain.Domain{}, fmt.Errorf("%w: %s", dbErr, err)
	}

	dmn := domainModel.transform()
	return dmn, nil
}

func (s *DomainRepository) List(ctx context.Context, flt domain.Filter) ([]domain.Domain, error) {
	stmt := dialect.Select().From(TABLE_DOMAINS)
	if flt.OrgID != "" && flt.State != "" {
		stmt = stmt.Where(goqu.Ex{
			"org_id": flt.OrgID,
			"state":  flt.State,
		})
	} else if flt.Name != "" && flt.State != "" {
		stmt = stmt.Where(goqu.Ex{
			"name":  flt.Name,
			"state": flt.State,
		})
	} else if flt.OrgID != "" {
		stmt = stmt.Where(goqu.Ex{
			"org_id": flt.OrgID,
		})
	}

	query, params, err := stmt.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", parseErr, err)
	}

	var domains []Domain
	if err = s.dbc.WithTimeout(ctx, TABLE_DOMAINS, "List", func(ctx context.Context) error {
		return s.dbc.SelectContext(ctx, &domains, query, params...)
	}); err != nil {
		err = checkPostgresError(err)
		return nil, fmt.Errorf("%w: %s", dbErr, err)
	}

	var result []domain.Domain
	for _, d := range domains {
		transformedDomain := d.transform()
		result = append(result, transformedDomain)
	}

	return result, nil
}

func (s *DomainRepository) Get(ctx context.Context, id string) (domain.Domain, error) {
	query, params, err := dialect.From(TABLE_DOMAINS).Where(goqu.Ex{
		"id": id,
	}).ToSQL()

	if err != nil {
		return domain.Domain{}, fmt.Errorf("%w: %s", queryErr, err)
	}

	var domainModel Domain
	if err = s.dbc.WithTimeout(ctx, TABLE_DOMAINS, "Get", func(ctx context.Context) error {
		return s.dbc.QueryRowxContext(ctx, query, params...).StructScan(&domainModel)
	}); err != nil {
		err = checkPostgresError(err)
		return domain.Domain{}, fmt.Errorf("%w: %s", dbErr, err)
	}

	domain := domainModel.transform()
	return domain, nil
}

func (s *DomainRepository) Delete(ctx context.Context, id string) error {
	query, params, err := dialect.Delete(TABLE_DOMAINS).Where(goqu.Ex{
		"id": id,
	}).Returning(&Domain{}).ToSQL()
	if err != nil {
		return fmt.Errorf("%w: %s", queryErr, err)
	}

	return s.dbc.WithTimeout(ctx, TABLE_DOMAINS, "Delete", func(ctx context.Context) error {
		result, err := s.dbc.ExecContext(ctx, query, params...)
		if err != nil {
			err = checkPostgresError(err)
			return fmt.Errorf("%w: %s", dbErr, err)
		}

		if count, _ := result.RowsAffected(); count > 0 {
			return nil
		}

		return domain.ErrNotExist
	})
}

func (s *DomainRepository) Update(ctx context.Context, toUpdate domain.Domain) (domain.Domain, error) {
	if strings.TrimSpace(toUpdate.ID) == "" {
		return domain.Domain{}, domain.ErrInvalidId
	}

	query, params, err := dialect.Update(TABLE_DOMAINS).Set(
		goqu.Record{
			"token":      toUpdate.Token,
			"state":      toUpdate.State,
			"updated_at": goqu.L("now()"),
		}).Where(goqu.Ex{
		"id": toUpdate.ID,
	}).Returning(&Domain{}).ToSQL()
	if err != nil {
		return domain.Domain{}, fmt.Errorf("%w: %s", queryErr, err)
	}

	var domainModel Domain
	if err = s.dbc.WithTimeout(ctx, TABLE_DOMAINS, "Update", func(ctx context.Context) error {
		return s.dbc.QueryRowxContext(ctx, query, params...).StructScan(&domainModel)
	}); err != nil {
		err = checkPostgresError(err)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return domain.Domain{}, domain.ErrNotExist
		default:
			return domain.Domain{}, fmt.Errorf("%w: %s", dbErr, err)
		}
	}

	domain := domainModel.transform()
	return domain, nil
}

func (s *DomainRepository) DeleteExpiredDomainRequests(ctx context.Context) error {
	query, params, err := dialect.Delete(TABLE_DOMAINS).Where(goqu.Ex{
		"created_at": goqu.Op{"lte": s.Now().Add(-domain.DefaultTokenExpiry)},
		"state":      domain.Pending,
	}).ToSQL()
	if err != nil {
		return fmt.Errorf("%w: %s", queryErr, err)
	}

	return s.dbc.WithTimeout(ctx, TABLE_DOMAINS, "DeleteExpiredDomain", func(ctx context.Context) error {
		result, err := s.dbc.ExecContext(ctx, query, params...)
		if err != nil {
			err = checkPostgresError(err)
			return fmt.Errorf("%w: %s", dbErr, err)
		}

		count, _ := result.RowsAffected()
		s.log.Debug("DeleteExpiredDomains", "expired_domain_count", count)

		return nil
	})
}
