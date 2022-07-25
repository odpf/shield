package authz

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/odpf/shield/api/handler"
	"github.com/odpf/shield/core/namespace"
	"github.com/odpf/shield/core/resource"
	"github.com/odpf/shield/core/user"
	"github.com/odpf/shield/hook"
	"github.com/odpf/shield/middleware"
	"github.com/odpf/shield/pkg/body_extractor"

	"github.com/mitchellh/mapstructure"
	"github.com/odpf/salt/log"
)

type Authz struct {
	log log.Logger

	// To go to next hook
	next hook.Service

	// To skip all the next hooks and just respond back
	escape hook.Service

	Deps handler.Deps
}

func New(log log.Logger, next, escape hook.Service, deps handler.Deps) Authz {
	return Authz{
		log:    log,
		next:   next,
		escape: escape,
		Deps:   deps,
	}
}

type Config struct {
	Action     string                    `yaml:"action" mapstructure:"action"`
	Attributes map[string]hook.Attribute `yaml:"attributes" mapstructure:"attributes"`
}

func (a Authz) Info() hook.Info {
	return hook.Info{
		Name:        "authz",
		Description: "hook to modify permissions for the resource",
	}
}

func (a Authz) ServeHook(res *http.Response, err error) (*http.Response, error) {
	if err != nil || res.StatusCode >= 400 {
		return a.escape.ServeHook(res, err)
	}

	ruleFromRequest, ok := hook.ExtractRule(res.Request)
	if !ok {
		return a.next.ServeHook(res, nil)
	}

	hookSpec, ok := hook.ExtractHook(res.Request, a.Info().Name)
	if !ok {
		return a.next.ServeHook(res, nil)
	}

	config := Config{}
	if err := mapstructure.Decode(hookSpec.Config, &config); err != nil {
		return a.next.ServeHook(res, nil)
	}

	if ruleFromRequest.Backend.Namespace == "" {
		return a.next.ServeHook(res, fmt.Errorf("namespace variable not defined in rules"))
	}

	attributes := map[string]interface{}{}
	attributes["namespace"] = ruleFromRequest.Backend.Namespace

	attributes["user"] = res.Request.Header.Get(a.Deps.V1beta1.IdentityProxyHeader)
	res.Request = res.Request.WithContext(user.SetEmailToContext(res.Request.Context(), res.Request.Header.Get(a.Deps.V1beta1.IdentityProxyHeader)))

	for id, attr := range config.Attributes {
		bdy, _ := middleware.ExtractRequestBody(res.Request)
		bodySource := &res.Body
		if attr.Source == string(hook.SourceRequest) {
			bodySource = &bdy
		}

		headerSource := &res.Header
		if attr.Source == string(hook.SourceRequest) {
			headerSource = &res.Request.Header
		}

		switch attr.Type {
		case hook.AttributeTypeGRPCPayload:
			if !strings.HasPrefix(res.Header.Get("Content-Type"), "application/grpc") {
				a.log.Error("middleware: not a grpc request", "attr", attr)
				return a.escape.ServeHook(res, fmt.Errorf("invalid header for http request: %s", res.Header.Get("Content-Type")))
			}

			payloadField, err := body_extractor.GRPCPayloadHandler{}.Extract(bodySource, attr.Index)
			if err != nil {
				a.log.Error("middleware: failed to parse grpc payload", "err", err)
				return a.escape.ServeHook(res, fmt.Errorf("unable to parse grpc payload"))
			}
			attributes[id] = payloadField

			a.log.Info("middleware: extracted", "field", payloadField, "attr", attr)
		case hook.AttributeTypeJSONPayload:
			if attr.Key == "" {
				a.log.Error("middleware: payload key field empty")
				return a.escape.ServeHook(res, fmt.Errorf("payload key field empty"))
			}

			payloadField, err := body_extractor.JSONPayloadHandler{}.Extract(bodySource, attr.Key)
			if err != nil {
				a.log.Error("middleware: failed to parse json payload", "err", err)
				return a.escape.ServeHook(res, fmt.Errorf("failed to parse json payload"))
			}
			attributes[id] = payloadField

			a.log.Info("middleware: extracted", "field", payloadField, "attr", attr)
		case hook.AttributeTypeHeader:
			if attr.Key == "" {
				a.log.Error("middleware: header key field empty", "err", err)
				return a.escape.ServeHook(res, fmt.Errorf("failed to parse json payload"))
			}
			headerAttr := headerSource.Get(attr.Key)
			if headerAttr == "" {
				a.log.Error(fmt.Sprintf("middleware: header %s is empty", attr.Key))
				return a.escape.ServeHook(res, fmt.Errorf("failed to parse json payload"))
			}

			attributes[id] = headerAttr
			a.log.Info("middleware: extracted", "field", headerAttr, "attr", attr)

		case hook.AttributeTypeQuery:
			if attr.Key == "" {
				a.log.Error("middleware: query key field empty")
				return a.escape.ServeHook(res, fmt.Errorf("failed to parse json payload"))
			}
			queryAttr := res.Request.URL.Query().Get(attr.Key)
			if queryAttr == "" {
				a.log.Error(fmt.Sprintf("middleware: query %s is empty", attr.Key))
				return a.escape.ServeHook(res, fmt.Errorf("failed to parse json payload"))
			}

			attributes[id] = queryAttr
			a.log.Info("middleware: extracted", "field", queryAttr, "attr", attr)

		case hook.AttributeTypeConstant:
			if attr.Value == "" {
				a.log.Error("middleware: constant value empty")
				return a.escape.ServeHook(res, fmt.Errorf("failed to parse json payload"))
			}

			attributes[id] = attr.Value
			a.log.Info("middleware: extracted", "constant_key", res, "attr", attributes[id])

		default:
			a.log.Error("middleware: unknown attribute type", "attr", attr)
			return a.escape.ServeHook(res, fmt.Errorf("unknown attribute type: %v", attr))
		}
	}

	paramMap, _ := middleware.ExtractPathParams(res.Request)
	for key, value := range paramMap {
		attributes[key] = value
	}

	resources, err := createResources(attributes)
	if err != nil {
		a.log.Error(err.Error())
		return a.escape.ServeHook(res, fmt.Errorf(err.Error()))
	}
	for _, resource := range resources {
		newResource, err := a.Deps.V1beta1.ResourceService.Create(res.Request.Context(), resource)
		if err != nil {
			a.log.Error(err.Error())
			return a.escape.ServeHook(res, fmt.Errorf(err.Error()))
		}
		a.log.Info(fmt.Sprintf("Resource %s created with ID %s", newResource.Urn, newResource.Idxa))
	}

	return a.next.ServeHook(res, nil)
}

func createResources(permissionAttributes map[string]interface{}) ([]resource.Resource, error) {
	var resources []resource.Resource
	projects, err := getAttributesValues(permissionAttributes["project"])
	if err != nil {
		return nil, err
	}

	orgs, err := getAttributesValues(permissionAttributes["organization"])
	if err != nil {
		return nil, err
	}

	teams, err := getAttributesValues(permissionAttributes["team"])
	if err != nil {
		return nil, err
	}

	resourceList, err := getAttributesValues(permissionAttributes["resource"])
	if err != nil {
		return nil, err
	}

	backendNamespace, err := getAttributesValues(permissionAttributes["namespace"])
	if err != nil {
		return nil, err
	}

	resourceType, err := getAttributesValues(permissionAttributes["resource_type"])
	if err != nil {
		return nil, err
	}

	if len(projects) < 1 || len(orgs) < 1 || len(resourceList) < 1 || (backendNamespace[0] == "") || (resourceType[0] == "") {
		return nil, fmt.Errorf("namespace, resource type, projects, organizations, resource, and team are required")
	}

	for _, org := range orgs {
		for _, project := range projects {
			for _, res := range resourceList {
				if len(teams) > 0 {
					for _, team := range teams {
						resources = append(resources, resource.Resource{
							Name:           res,
							OrganizationId: org,
							ProjectId:      project,
							GroupId:        team,
							NamespaceId:    namespace.CreateID(backendNamespace[0], resourceType[0]),
						})
					}
				} else {
					resources = append(resources, resource.Resource{
						Name:           res,
						OrganizationId: org,
						ProjectId:      project,
						NamespaceId:    namespace.CreateID(backendNamespace[0], resourceType[0]),
					})
				}
			}
		}
	}
	return resources, nil
}

func getAttributesValues(attributes interface{}) ([]string, error) {
	var values []string
	switch attributes.(type) {
	case []string:
		for _, i := range attributes.([]string) {
			values = append(values, i)
		}
	case string:
		values = append(values, attributes.(string))
	case []interface{}:
		for _, i := range attributes.([]interface{}) {
			values = append(values, i.(string))
		}
	case interface{}:
		values = append(values, attributes.(string))
	case nil:
		return values, nil
	default:
		return values, fmt.Errorf("unsuported attribute type %v", attributes)
	}
	return values, nil
}
