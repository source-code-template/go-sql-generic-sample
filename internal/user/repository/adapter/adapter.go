package adapter

import (
	"context"
	"database/sql"
	"fmt"
	s "github.com/core-go/sql"
	"reflect"

	"go-service/internal/user/model"
)

func NewUserAdapter(db *sql.DB, buildQuery func(*model.UserFilter) (string, []interface{})) (*UserAdapter, error) {
	userType := reflect.TypeOf(model.User{})
	parameters, err := s.CreateParameters(userType, db)
	if err != nil {
		return nil, err
	}
	return &UserAdapter{DB: db, Parameters: parameters, BuildQuery: buildQuery}, nil
}

type UserAdapter struct {
	DB         *sql.DB
	BuildQuery func(*model.UserFilter) (string, []interface{})
	*s.Parameters
}

func (r *UserAdapter) All(ctx context.Context) ([]model.User, error) {
	query := `select * from users`
	var users []model.User
	err := s.Query(ctx, r.DB, r.Map, &users, query)
	return users, err
}

func (r *UserAdapter) Load(ctx context.Context, id string) (*model.User, error) {
	var users []model.User
	query := fmt.Sprintf("select %s from users where id = %s limit 1", r.Fields, r.BuildParam(1))
	err := s.Query(ctx, r.DB, r.Map, &users, query, id)
	if err != nil {
		return nil, err
	}
	if len(users) > 0 {
		return &users[0], nil
	}
	return nil, nil
}

func (r *UserAdapter) Create(ctx context.Context, user *model.User) (int64, error) {
	query, args := s.BuildToInsert("users", user, r.BuildParam)
	tx := s.GetTx(ctx, r.DB)
	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func (r *UserAdapter) Update(ctx context.Context, user *model.User) (int64, error) {
	query, args := s.BuildToUpdate("users", user, r.BuildParam)
	tx := s.GetTx(ctx, r.DB)
	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func (r *UserAdapter) Patch(ctx context.Context, user map[string]interface{}) (int64, error) {
	colMap := s.JSONToColumns(user, r.JsonColumnMap)
	query, args := s.BuildToPatch("users", colMap, r.Keys, r.BuildParam)
	tx := s.GetTx(ctx, r.DB)
	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func (r *UserAdapter) Delete(ctx context.Context, id string) (int64, error) {
	query := fmt.Sprintf("delete from users where id = %s", r.BuildParam(1))
	tx := s.GetTx(ctx, r.DB)
	res, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func (r *UserAdapter) Search(ctx context.Context, filter *model.UserFilter, limit int64, offset int64) ([]model.User, int64, error) {
	var users []model.User
	if limit <= 0 {
		return users, 0, nil
	}
	query, params := r.BuildQuery(filter)
	pagingQuery := s.BuildPagingQuery(query, limit, offset)
	countQuery := s.BuildCountQuery(query)

	row := r.DB.QueryRowContext(ctx, countQuery, params...)
	if row.Err() != nil {
		return users, 0, row.Err()
	}
	var total int64
	err := row.Scan(&total)
	if err != nil || total == 0 {
		return users, total, err
	}

	err = s.Query(ctx, r.DB, r.Map, &users, pagingQuery, params...)
	return users, total, err
}
