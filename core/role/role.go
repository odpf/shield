package role

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/odpf/shield/core/namespace"
)

var (
	ErrNotExist    = errors.New("role doesn't exist")
	ErrInvalidUUID = errors.New("invalid syntax of uuid")
)

type Store interface {
	CreateRole(ctx context.Context, role Role) (Role, error)
	GetRole(ctx context.Context, id string) (Role, error)
	ListRoles(ctx context.Context) ([]Role, error)
	UpdateRole(ctx context.Context, toUpdate Role) (Role, error)
}

type Role struct {
	ID          string
	Name        string
	Types       []string
	Namespace   namespace.Namespace
	NamespaceID string
	Metadata    map[string]any
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func GetOwnerRole(ns namespace.Namespace) Role {
	id := fmt.Sprintf("%s_%s", ns.ID, "owner")
	name := fmt.Sprintf("%s_%s", strings.Title(ns.ID), "Owner")
	return Role{
		ID:        id,
		Name:      name,
		Types:     []string{UserType},
		Namespace: ns,
	}
}
