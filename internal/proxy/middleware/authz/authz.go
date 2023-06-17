package authz

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/raystack/shield/core/authenticate"

	"github.com/raystack/shield/core/relation"

	"github.com/mitchellh/mapstructure"
	"github.com/raystack/salt/log"

	"github.com/raystack/shield/core/resource"
	"github.com/raystack/shield/internal/proxy/middleware"
	"github.com/raystack/shield/pkg/body_extractor"
)

type ResourceService interface {
	CheckAuthz(ctx context.Context, rel relation.Object, permissionName string) (bool, error)
}

type AuthnService interface {
	GetPrincipal(ctx context.Context) (authenticate.Principal, error)
}

type Authz struct {
	log             log.Logger
	userIDHeaderKey string
	next            http.Handler
	resourceService ResourceService
	userService     AuthnService
}

type Config struct {
	Actions     []string                        `yaml:"actions" mapstructure:"actions"`
	Permissions []Permission                    `yaml:"permissions" mapstructure:"permissions"`
	Attributes  map[string]middleware.Attribute `yaml:"attributes" mapstructure:"attributes"`
}

type Permission struct {
	Name      string `yaml:"name" mapstructure:"name"`
	Namespace string `yaml:"namespace" mapstructure:"namespace"`
	Attribute string `yaml:"attribute" mapstructure:"attribute"`
}

func New(
	log log.Logger,
	next http.Handler,
	userIDHeaderKey string,
	resourceService ResourceService,
	principalService AuthnService) *Authz {
	return &Authz{
		log:             log,
		userIDHeaderKey: userIDHeaderKey,
		next:            next,
		resourceService: resourceService,
		userService:     principalService,
	}
}

func (c Authz) Info() *middleware.MiddlewareInfo {
	return &middleware.MiddlewareInfo{
		Name:        "authz",
		Description: "rule based authorization using spicedb",
	}
}

func (c *Authz) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	rule, ok := middleware.ExtractRule(req)
	if !ok {
		c.next.ServeHTTP(rw, req)
		return
	}
	wareSpec, ok := rule.Middlewares.Get(c.Info().Name)
	if !ok {
		c.next.ServeHTTP(rw, req)
		return
	}

	usr, err := c.userService.GetPrincipal(req.Context())
	if err != nil {
		c.log.Error("middleware: failed to get user details", "err", err.Error())
		c.notAllowed(rw, nil)
		return
	}
	req.Header.Set(c.userIDHeaderKey, usr.ID)

	if rule.Backend.Namespace == "" {
		c.log.Error("namespace is not defined for this rule")
		c.notAllowed(rw, nil)
		return
	}

	// TODO: should cache it
	config := Config{}
	if err := mapstructure.Decode(wareSpec.Config, &config); err != nil {
		c.log.Error("middleware: failed to decode authz config", "config", wareSpec.Config)
		c.notAllowed(rw, nil)
		return
	}

	permissionAttributes := map[string]interface{}{}

	permissionAttributes["namespace"] = rule.Backend.Namespace

	permissionAttributes["user"] = req.Header.Get(c.userIDHeaderKey)

	for res, attr := range config.Attributes {
		_ = res

		switch attr.Type {
		case middleware.AttributeTypeGRPCPayload:
			// check if grpc request
			if !strings.HasPrefix(req.Header.Get("Content-Type"), "application/grpc") {
				c.log.Error("middleware: not a grpc request", "attr", attr)
				c.notAllowed(rw, nil)
				return
			}

			// TODO: we can optimise this by parsing all field at once
			payloadField, err := body_extractor.GRPCPayloadHandler{}.Extract(&req.Body, attr.Index)
			if err != nil {
				c.log.Error("middleware: failed to parse grpc payload", "err", err)
				return
			}

			permissionAttributes[res] = payloadField
			c.log.Info("middleware: extracted", "field", payloadField, "attr", attr)

		case middleware.AttributeTypeJSONPayload:
			if attr.Key == "" {
				c.log.Error("middleware: payload key field empty")
				c.notAllowed(rw, nil)
				return
			}
			payloadField, err := body_extractor.JSONPayloadHandler{}.Extract(&req.Body, attr.Key)
			if err != nil {
				c.log.Error("middleware: failed to parse grpc payload", "err", err)
				c.notAllowed(rw, nil)
				return
			}

			permissionAttributes[res] = payloadField
			c.log.Info("middleware: extracted", "field", payloadField, "attr", attr)

		case middleware.AttributeTypeHeader:
			if attr.Key == "" {
				c.log.Error("middleware: header key field empty")
				c.notAllowed(rw, nil)
				return
			}
			headerAttr := req.Header.Get(attr.Key)
			if headerAttr == "" {
				c.log.Error(fmt.Sprintf("middleware: header %s is empty", attr.Key))
				c.notAllowed(rw, nil)
				return
			}

			permissionAttributes[res] = headerAttr
			c.log.Info("middleware: extracted", "field", headerAttr, "attr", attr)

		case middleware.AttributeTypeQuery:
			if attr.Key == "" {
				c.log.Error("middleware: query key field empty")
				c.notAllowed(rw, nil)
				return
			}
			queryAttr := req.URL.Query().Get(attr.Key)
			if queryAttr == "" {
				c.log.Error(fmt.Sprintf("middleware: query %s is empty", attr.Key))
				c.notAllowed(rw, nil)
				return
			}

			permissionAttributes[res] = queryAttr
			c.log.Info("middleware: extracted", "field", queryAttr, "attr", attr)

		case middleware.AttributeTypeConstant:
			if attr.Value == "" {
				c.log.Error("middleware: constant value empty")
				c.notAllowed(rw, nil)
				return
			}

			permissionAttributes[res] = attr.Value
			c.log.Info("middleware: extracted", "constant_key", res, "attr", permissionAttributes[res])

		default:
			c.log.Error("middleware: unknown attribute type", "attr", attr)
			c.notAllowed(rw, nil)
			return
		}
	}

	paramMap, mapExists := middleware.ExtractPathParams(req)
	if !mapExists {
		c.log.Error("middleware: path param map doesn't exist")
		c.notAllowed(rw, nil)
		return
	}

	for key, value := range paramMap {
		permissionAttributes[key] = value
	}

	isAuthorized := false
	for _, perm := range config.Permissions {
		isAuthorized, err = c.resourceService.CheckAuthz(req.Context(), relation.Object{
			ID:        permissionAttributes[perm.Attribute].(string),
			Namespace: perm.Namespace,
		}, perm.Name)
		if err != nil {
			c.log.Error("error while performing authz permission check", "err", err)
			c.notAllowed(rw, err)
			return
		}
		if isAuthorized {
			break
		}
	}

	c.log.Info("authz check successful", "user", permissionAttributes["user"], "resource", permissionAttributes["resource"], "result", isAuthorized)
	if !isAuthorized {
		c.log.Info("user not allowed to make request", "user", permissionAttributes["user"], "resource", permissionAttributes["resource"], "result", isAuthorized)
		c.notAllowed(rw, nil)
		return
	}

	c.next.ServeHTTP(rw, req)
}

func (w Authz) notAllowed(rw http.ResponseWriter, err error) {
	if err != nil {
		switch {
		case errors.Is(err, resource.ErrNotExist):
			rw.WriteHeader(http.StatusNotFound)
			return
		}
	}
	rw.WriteHeader(http.StatusUnauthorized)
}
