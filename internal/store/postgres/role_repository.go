package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"database/sql"

	"github.com/doug-martin/goqu/v9"
	newrelic "github.com/newrelic/go-agent"
	"github.com/odpf/shield/core/role"
	"github.com/odpf/shield/pkg/db"
)

type RoleRepository struct {
	dbc *db.Client
}

func NewRoleRepository(dbc *db.Client) *RoleRepository {
	return &RoleRepository{
		dbc: dbc,
	}
}

func (r RoleRepository) buildListQuery(dialect goqu.DialectWrapper) *goqu.SelectDataset {
	roleSelectStatement := dialect.Select(
		goqu.I("r.id"),
		goqu.I("r.name"),
		goqu.I("r.types"),
		goqu.I("r.namespace_id"),
		goqu.I("r.metadata"),
		goqu.I("namespaces.id").As(goqu.C("namespace.id")),
		goqu.I("namespaces.name").As(goqu.C("namespace.name")),
	).From(goqu.T(TABLE_ROLES).As("r"))
	return roleSelectStatement.Join(goqu.T(TABLE_NAMESPACES), goqu.On(
		goqu.I("namespaces.id").Eq(goqu.I("r.namespace_id"))))
}

func (r RoleRepository) Get(ctx context.Context, id string) (role.Role, error) {
	if strings.TrimSpace(id) == "" {
		return role.Role{}, role.ErrInvalidID
	}

	query, params, err := r.buildListQuery(dialect).
		Where(
			goqu.Ex{"r.id": id},
		).ToSQL()
	if err != nil {
		return role.Role{}, fmt.Errorf("%w: %s", queryErr, err)
	}

	var roleModel Role
	if err = r.dbc.WithTimeout(ctx, func(ctx context.Context) error {
		nrCtx := newrelic.FromContext(ctx)
		if nrCtx != nil {
			nr := newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: TABLE_ROLES,
				Operation:  "Get",
				StartTime:  nrCtx.StartSegmentNow(),
			}
			defer nr.End()
		}

		return r.dbc.GetContext(ctx, &roleModel, query, params...)
	}); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return role.Role{}, role.ErrNotExist
		}
		return role.Role{}, fmt.Errorf("%w: %s", dbErr, err)
	}

	transformedRole, err := roleModel.transformToRole()
	if err != nil {
		return role.Role{}, err
	}

	return transformedRole, nil
}

// TODO this is actually an upsert
func (r RoleRepository) Create(ctx context.Context, rl role.Role) (string, error) {
	if strings.TrimSpace(rl.ID) == "" {
		return "", role.ErrInvalidID
	}

	if strings.TrimSpace(rl.Name) == "" {
		return "", role.ErrInvalidDetail
	}

	marshaledMetadata, err := json.Marshal(rl.Metadata)
	if err != nil {
		return "", fmt.Errorf("%w: %s", parseErr, err)
	}

	//TODO we have to go with this manually populating data since goqu does not support insert array string
	query, _, err := dialect.Insert(TABLE_ROLES).Rows(
		goqu.Record{
			"id":           goqu.L("$1"),
			"name":         goqu.L("$2"),
			"types":        goqu.L("$3"),
			"namespace_id": goqu.L("$4"),
			"metadata":     goqu.L("$5"),
		}).OnConflict(
		goqu.DoUpdate("id", goqu.Record{
			"name": goqu.L("$2"),
		},
		)).Returning("id").ToSQL()
	if err != nil {
		return "", fmt.Errorf("%w: %s", queryErr, err)
	}

	types := strings.Join(rl.Types, ",")
	types = fmt.Sprintf("{%s}", types)

	var roleID string
	if err = r.dbc.WithTimeout(ctx, func(ctx context.Context) error {
		nrCtx := newrelic.FromContext(ctx)
		if nrCtx != nil {
			nr := newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: TABLE_ROLES,
				Operation:  "Create",
				StartTime:  nrCtx.StartSegmentNow(),
			}
			defer nr.End()
		}

		return r.dbc.QueryRowxContext(ctx, query, rl.ID, rl.Name, types, rl.NamespaceID, marshaledMetadata).Scan(&roleID)
	}); err != nil {
		err = checkPostgresError(err)
		switch {
		case errors.Is(err, errDuplicateKey):
			return "", role.ErrConflict
		case errors.Is(err, errForeignKeyViolation):
			return "", role.ErrInvalidDetail
		default:
			return "", err
		}
	}

	return roleID, nil
}

func (r RoleRepository) List(ctx context.Context) ([]role.Role, error) {
	query, params, err := r.buildListQuery(dialect).ToSQL()
	if err != nil {
		return []role.Role{}, fmt.Errorf("%w: %s", queryErr, err)
	}

	var fetchedRoles []Role
	if err = r.dbc.WithTimeout(ctx, func(ctx context.Context) error {
		nrCtx := newrelic.FromContext(ctx)
		if nrCtx != nil {
			nr := newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: TABLE_ROLES,
				Operation:  "List",
				StartTime:  nrCtx.StartSegmentNow(),
			}
			defer nr.End()
		}

		return r.dbc.SelectContext(ctx, &fetchedRoles, query, params...)
	}); err != nil {
		return []role.Role{}, fmt.Errorf("%w: %s", dbErr, err)
	}

	var transformedRoles []role.Role
	for _, o := range fetchedRoles {
		transformedOrg, err := o.transformToRole()
		if err != nil {
			return []role.Role{}, fmt.Errorf("%w: %s", parseErr, err)
		}
		transformedRoles = append(transformedRoles, transformedOrg)
	}

	return transformedRoles, nil
}

func (r RoleRepository) Update(ctx context.Context, rl role.Role) (string, error) {
	if strings.TrimSpace(rl.ID) == "" {
		return "", role.ErrInvalidID
	}

	if strings.TrimSpace(rl.Name) == "" {
		return "", role.ErrInvalidDetail
	}

	marshaledMetadata, err := json.Marshal(rl.Metadata)
	if err != nil {
		return "", fmt.Errorf("%w: %s", parseErr, err)
	}

	//TODO we have to go with this manually populating data since goqu does not support insert array string
	query, _, err := dialect.Update(TABLE_ROLES).Set(
		goqu.Record{
			"name":         goqu.L("$2"),
			"types":        goqu.L("$3"),
			"namespace_id": goqu.L("$4"),
			"metadata":     goqu.L("$5"),
			"updated_at":   goqu.L("now()"),
		}).Where(
		goqu.Ex{"id": goqu.L("$1")},
	).Returning("id").ToSQL()
	if err != nil {
		return "", fmt.Errorf("%w: %s", queryErr, err)
	}

	var roleID string
	if err = r.dbc.WithTimeout(ctx, func(ctx context.Context) error {
		nrCtx := newrelic.FromContext(ctx)
		if nrCtx != nil {
			nr := newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: TABLE_ROLES,
				Operation:  "Update",
				StartTime:  nrCtx.StartSegmentNow(),
			}
			defer nr.End()
		}

		return r.dbc.QueryRowxContext(ctx, query, rl.ID, rl.Name, rl.Types, rl.NamespaceID, marshaledMetadata).Scan(&roleID)
	}); err != nil {
		err = checkPostgresError(err)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return "", role.ErrNotExist
		case errors.Is(err, errForeignKeyViolation):
			return "", role.ErrInvalidDetail
		case errors.Is(err, errDuplicateKey):
			return "", role.ErrConflict
		default:
			return "", err
		}
	}

	return roleID, nil
}
