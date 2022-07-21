package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/doug-martin/goqu/v9"

	"github.com/odpf/shield/internal/relation"
	"github.com/odpf/shield/model"
	"github.com/odpf/shield/pkg/utils"
)

type Relation struct {
	Id                 string         `db:"id"`
	SubjectNamespaceId string         `db:"subject_namespace_id"`
	SubjectNamespace   Namespace      `db:"subject_namespace"`
	SubjectId          string         `db:"subject_id"`
	ObjectNamespaceId  string         `db:"object_namespace_id"`
	ObjectNamespace    Namespace      `db:"object_namespace"`
	ObjectId           string         `db:"object_id"`
	RoleId             sql.NullString `db:"role_id"`
	Role               Role           `db:"role"`
	NamespaceId        sql.NullString `db:"namespace_id"`
	CreatedAt          time.Time      `db:"created_at"`
	UpdatedAt          time.Time      `db:"updated_at"`
	DeletedAt          sql.NullTime   `db:"deleted_at"`
}

type relationCols struct {
	Id                 string         `db:"id"`
	SubjectNamespaceId string         `db:"subject_namespace_id"`
	SubjectId          string         `db:"subject_id"`
	ObjectNamespaceId  string         `db:"object_namespace_id"`
	ObjectId           string         `db:"object_id"`
	RoleId             sql.NullString `db:"role_id"`
	NamespaceId        sql.NullString `db:"namespace_id"`
	CreatedAt          time.Time      `db:"created_at"`
	UpdatedAt          time.Time      `db:"updated_at"`
}

func buildCreateRelationQuery(dialect goqu.DialectWrapper) (string, error) {
	// TODO: Look for a better way to implement goqu.OnConflict with multiple columns

	createRelationQuery, _, err := dialect.Insert(TABLE_RELATION).Rows(
		goqu.Record{
			"subject_namespace_id": goqu.L("$1"),
			"subject_id":           goqu.L("$2"),
			"object_namespace_id":  goqu.L("$3"),
			"object_id":            goqu.L("$4"),
			"role_id":              goqu.L("$5"),
			"namespace_id":         goqu.L("$6"),
		}).OnConflict(goqu.DoUpdate("subject_namespace_id, subject_id, object_namespace_id,  object_id, COALESCE(role_id, ''), COALESCE(namespace_id, '')", goqu.Record{
		"subject_namespace_id": goqu.L("$1"),
	})).Returning(&relationCols{}).ToSQL()

	return createRelationQuery, err
}

func buildListRelationQuery(dialect goqu.DialectWrapper) (string, error) {
	listRelationQuery, _, err := dialect.Select(&relationCols{}).From(TABLE_RELATION).ToSQL()

	return listRelationQuery, err
}

func buildGetRelationsQuery(dialect goqu.DialectWrapper) (string, error) {
	getRelationsQuery, _, err := dialect.Select(&relationCols{}).From(TABLE_RELATION).Where(goqu.Ex{
		"id": goqu.L("$1"),
	}).ToSQL()

	return getRelationsQuery, err
}

func buildUpdateRelationQuery(dialect goqu.DialectWrapper) (string, error) {
	updateRelationQuery, _, err := goqu.Update(TABLE_RELATION).Set(
		goqu.Record{
			"subject_namespace_id": goqu.L("$2"),
			"subject_id":           goqu.L("$3"),
			"object_namespace_id":  goqu.L("$4"),
			"object_id":            goqu.L("$5"),
			"role_id":              goqu.L("$6"),
			"namespace_id":         goqu.L("$7"),
		}).Where(goqu.Ex{
		"id": goqu.L("$1"),
	}).Returning(&relationCols{}).ToSQL()

	return updateRelationQuery, err
}

func buildGetRelationByFieldsQuery(dialect goqu.DialectWrapper) (string, error) {
	getRelationByFieldsQuery, _, err := dialect.Select(&relationCols{}).From(TABLE_RELATION).Where(goqu.Ex{
		"subject_namespace_id": goqu.L("$1"),
		"subject_id":           goqu.L("$2"),
		"object_namespace_id":  goqu.L("$3"),
		"object_id":            goqu.L("$4"),
	}, goqu.And(
		goqu.Or(
			goqu.C("role_id").IsNull(),
			goqu.C("role_id").Eq(goqu.L("$5")),
		)),
		goqu.And(
			goqu.Or(
				goqu.C("namespace_id").IsNull(),
				goqu.C("namespace_id").Eq(goqu.L("$6")),
			),
		)).ToSQL()

	return getRelationByFieldsQuery, err
}

func buildDeleteRelationByIdQuery(dialect goqu.DialectWrapper) (string, error) {
	deleteRelationByIdQuery, _, err := dialect.Delete(TABLE_RELATION).Where(goqu.Ex{
		"id": goqu.L("$1"),
	}).ToSQL()

	return deleteRelationByIdQuery, err
}

func (s Store) CreateRelation(ctx context.Context, relationToCreate model.Relation) (model.Relation, error) {
	var newRelation Relation

	subjectNamespaceId := utils.DefaultStringIfEmpty(relationToCreate.SubjectNamespace.Id, relationToCreate.SubjectNamespaceId)
	objectNamespaceId := utils.DefaultStringIfEmpty(relationToCreate.ObjectNamespace.Id, relationToCreate.ObjectNamespaceId)
	roleId := utils.DefaultStringIfEmpty(relationToCreate.Role.Id, relationToCreate.RoleId)
	var nsId string

	if relationToCreate.RelationType == model.RelationTypes.Namespace {
		nsId = roleId
		roleId = ""
	}

	createRelationQuery, err := buildCreateRelationQuery(dialect)
	if err != nil {
		return model.Relation{}, fmt.Errorf("%w: %s", queryErr, err)
	}

	err = s.DB.WithTimeout(ctx, func(ctx context.Context) error {
		return s.DB.GetContext(
			ctx,
			&newRelation,
			createRelationQuery,
			subjectNamespaceId,
			relationToCreate.SubjectId,
			objectNamespaceId,
			relationToCreate.ObjectId,
			sql.NullString{String: roleId, Valid: roleId != ""},
			sql.NullString{String: nsId, Valid: nsId != ""},
		)
	})

	if err != nil {
		return model.Relation{}, err
	}

	transformedRelation, err := transformToRelation(newRelation)

	if err != nil {
		return model.Relation{}, err
	}

	return transformedRelation, nil
}

func (s Store) ListRelations(ctx context.Context) ([]model.Relation, error) {
	var fetchedRelations []Relation
	listRelationQuery, err := buildListRelationQuery(dialect)
	if err != nil {
		return []model.Relation{}, fmt.Errorf("%w: %s", queryErr, err)
	}

	err = s.DB.WithTimeout(ctx, func(ctx context.Context) error {
		return s.DB.SelectContext(ctx, &fetchedRelations, listRelationQuery)
	})

	if errors.Is(err, sql.ErrNoRows) {
		return []model.Relation{}, relation.RelationDoesntExist
	}

	if err != nil {
		return []model.Relation{}, fmt.Errorf("%w: %s", dbErr, err)
	}

	var transformedRelations []model.Relation

	for _, r := range fetchedRelations {
		transformedRelation, err := transformToRelation(r)
		if err != nil {
			return []model.Relation{}, fmt.Errorf("%w: %s", parseErr, err)
		}

		transformedRelations = append(transformedRelations, transformedRelation)
	}

	return transformedRelations, nil
}

func (s Store) GetRelation(ctx context.Context, id string) (model.Relation, error) {
	var fetchedRelation Relation
	getRelationsQuery, err := buildGetRelationsQuery(dialect)
	if err != nil {
		return model.Relation{}, fmt.Errorf("%w: %s", queryErr, err)
	}

	err = s.DB.WithTimeout(ctx, func(ctx context.Context) error {
		return s.DB.GetContext(ctx, &fetchedRelation, getRelationsQuery, id)
	})

	if errors.Is(err, sql.ErrNoRows) {
		return model.Relation{}, relation.RelationDoesntExist
	} else if err != nil && fmt.Sprintf("%s", err.Error()[0:38]) == "pq: invalid input syntax for type uuid" {
		// TODO: this uuid syntax is a error defined in db, not in library
		// need to look into better ways to implement this
		return model.Relation{}, relation.InvalidUUID
	} else if err != nil {
		return model.Relation{}, fmt.Errorf("%w: %s", dbErr, err)
	}

	if err != nil {
		return model.Relation{}, err
	}

	transformedRelation, err := transformToRelation(fetchedRelation)
	if err != nil {
		return model.Relation{}, err
	}

	return transformedRelation, nil
}

func (s Store) DeleteRelationById(ctx context.Context, id string) error {
	deleteRelationByIdQuery, err := buildDeleteRelationByIdQuery(dialect)
	if err != nil {
		return fmt.Errorf("%w: %s", queryErr, err)
	}

	err = s.DB.WithTimeout(ctx, func(ctx context.Context) error {
		result, err := s.DB.ExecContext(ctx, deleteRelationByIdQuery, id)
		if err == nil {
			count, err := result.RowsAffected()
			if err == nil {
				if count == 1 {
					return nil
				}
			}
		}
		return err
	})
	return err
}

func (s Store) GetRelationByFields(ctx context.Context, rel model.Relation) (model.Relation, error) {
	var fetchedRelation Relation

	subjectNamespaceId := utils.DefaultStringIfEmpty(rel.SubjectNamespace.Id, rel.SubjectNamespaceId)
	objectNamespaceId := utils.DefaultStringIfEmpty(rel.ObjectNamespace.Id, rel.ObjectNamespaceId)
	roleId := utils.DefaultStringIfEmpty(rel.Role.Id, rel.RoleId)
	var nsId string

	if rel.RelationType == model.RelationTypes.Namespace {
		nsId = roleId
		roleId = ""
	}

	getRelationByFieldsQuery, err := buildGetRelationByFieldsQuery(dialect)
	if err != nil {
		return model.Relation{}, fmt.Errorf("%w: %s", queryErr, err)
	}

	err = s.DB.WithTimeout(ctx, func(ctx context.Context) error {
		return s.DB.GetContext(ctx,
			&fetchedRelation,
			getRelationByFieldsQuery,
			subjectNamespaceId,
			rel.SubjectId,
			objectNamespaceId,
			rel.ObjectId,
			sql.NullString{String: roleId, Valid: roleId != ""},
			sql.NullString{String: nsId, Valid: nsId != ""},
		)
	})

	if errors.Is(err, sql.ErrNoRows) {
		return model.Relation{}, relation.RelationDoesntExist
	} else if err != nil && fmt.Sprintf("%s", err.Error()[0:38]) == "pq: invalid input syntax for type uuid" {
		// TODO: this uuid syntax is a error defined in db, not in library
		// need to look into better ways to implement this
		return model.Relation{}, relation.InvalidUUID
	} else if err != nil {
		return model.Relation{}, fmt.Errorf("%w: %s", dbErr, err)
	}

	if err != nil {
		return model.Relation{}, err
	}

	transformedRelation, err := transformToRelation(fetchedRelation)
	if err != nil {
		return model.Relation{}, err
	}

	return transformedRelation, nil
}

func (s Store) UpdateRelation(ctx context.Context, id string, toUpdate model.Relation) (model.Relation, error) {
	var updatedRelation Relation

	subjectNamespaceId := utils.DefaultStringIfEmpty(toUpdate.SubjectNamespace.Id, toUpdate.SubjectNamespaceId)
	objectNamespaceId := utils.DefaultStringIfEmpty(toUpdate.ObjectNamespace.Id, toUpdate.ObjectNamespaceId)
	roleId := utils.DefaultStringIfEmpty(toUpdate.Role.Id, toUpdate.RoleId)
	var nsId string

	if toUpdate.RelationType == model.RelationTypes.Namespace {
		nsId = roleId
		roleId = ""
	}

	updateRelationQuery, err := buildUpdateRelationQuery(dialect)
	if err != nil {
		return model.Relation{}, fmt.Errorf("%w: %s", queryErr, err)
	}

	err = s.DB.WithTimeout(ctx, func(ctx context.Context) error {
		return s.DB.GetContext(
			ctx,
			&updatedRelation,
			updateRelationQuery,
			id,
			subjectNamespaceId,
			toUpdate.SubjectId,
			objectNamespaceId,
			toUpdate.ObjectId,
			sql.NullString{String: roleId, Valid: roleId != ""},
			sql.NullString{String: nsId, Valid: nsId != ""},
		)
	})

	if errors.Is(err, sql.ErrNoRows) {
		return model.Relation{}, relation.RelationDoesntExist
	} else if err != nil && fmt.Sprintf("%s", err.Error()[0:38]) == "pq: invalid input syntax for type uuid" {
		// TODO: this uuid syntax is a error defined in db, not in library
		// need to look into better ways to implement this
		return model.Relation{}, fmt.Errorf("%w: %s", relation.InvalidUUID, err)
	} else if err != nil {
		return model.Relation{}, fmt.Errorf("%w: %s", dbErr, err)
	}

	toUpdate, err = transformToRelation(updatedRelation)
	if err != nil {
		return model.Relation{}, fmt.Errorf("%s: %w", parseErr, err)
	}

	return toUpdate, nil
}

func transformToRelation(from Relation) (model.Relation, error) {
	relationType := model.RelationTypes.Role
	roleId := from.RoleId.String

	if from.NamespaceId.Valid {
		roleId = from.NamespaceId.String
		relationType = model.RelationTypes.Namespace
	}

	return model.Relation{
		Id:                 from.Id,
		SubjectNamespaceId: from.SubjectNamespaceId,
		SubjectId:          from.SubjectId,
		ObjectNamespaceId:  from.ObjectNamespaceId,
		ObjectId:           from.ObjectId,
		RoleId:             roleId,
		RelationType:       relationType,
		CreatedAt:          from.CreatedAt,
		UpdatedAt:          from.UpdatedAt,
	}, nil
}
