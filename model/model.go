package model

import "time"

type Project struct {
	Id           string
	Name         string
	Slug         string
	Organization Organization
	Metadata     map[string]string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Organization struct {
	Id        string
	Name      string
	Slug      string
	Metadata  map[string]string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Group struct {
	Id           string
	Name         string
	Slug         string
	Organization Organization
	Metadata     map[string]string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Role struct {
	Id        string
	Name      string
	Types     []string
	Namespace string
	Metadata  map[string]string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Action struct {
	Id        string
	Name      string
	Slug      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Namespace struct {
	Id        string
	Name      string
	Slug      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Policy struct {
	Id        string
	Role      Role
	Namespace Namespace
	Action    Action
	CreatedAt time.Time
	UpdatedAt time.Time
}
