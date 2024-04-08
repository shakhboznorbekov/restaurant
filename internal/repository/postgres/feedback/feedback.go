package feedback

import (
	"context"
	"github.com/pkg/errors"
	"net/http"
	"restu-backend/foundation/web"
	"restu-backend/internal/auth"
	"restu-backend/internal/entity"
	"restu-backend/internal/pkg/repository/postgresql"
	"restu-backend/internal/service/feedback"
)

type Repository struct {
	*postgresql.Database
}

// @admin

func (r Repository) AdminGetList(ctx context.Context, filter feedback.Filter) ([]feedback.AdminGetList, int, error) {
	_, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return nil, 0, err
	}

	lang := r.DefaultLang
	var data []entity.Feedback

	q := r.NewSelect().Model(&data)

	if filter.Limit != nil {
		q.Limit(*filter.Limit)
	}

	if filter.Offset != nil {
		q.Offset(*filter.Offset)
	}

	count, err := q.ScanAndCount(ctx)

	var list []feedback.AdminGetList
	for _, i := range data {
		var detail feedback.AdminGetList
		detail.ID = i.ID
		n := i.Name[lang]
		detail.Name = &n

		list = append(list, detail)
	}

	return list, count, err
}

func (r Repository) AdminGetDetail(ctx context.Context, id int64) (entity.Feedback, error) {
	_, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return entity.Feedback{}, err
	}

	var detail entity.Feedback

	err = r.NewSelect().Model(&detail).Where("id = ?", id).Scan(ctx)

	return detail, err
}

func (r Repository) AdminCreate(ctx context.Context, request feedback.AdminCreate) (entity.Feedback, error) {
	_, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return entity.Feedback{}, err
	}

	detail := entity.Feedback{
		Name: request.Name,
	}

	_, err = r.NewInsert().Model(&detail).Exec(ctx)

	return detail, err
}

func (r Repository) AdminUpdateColumns(ctx context.Context, request feedback.AdminUpdate) error {
	_, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("feedback").Where("id = ?", request.ID)

	if request.Name != nil {
		q.Set("name = ?", request.Name)
	}

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating feedback"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) AdminDelete(ctx context.Context, id int64) error {
	_, err := r.NewDelete().Table("feedback").Where("id = ?", id).Exec(ctx)

	return err
}

// client

func (r Repository) ClientGetList(ctx context.Context, filter feedback.Filter) ([]feedback.ClientGetList, int, error) {
	_, err := r.CheckClaims(ctx, auth.RoleClient)
	if err != nil {
		return nil, 0, err
	}

	lang := r.DefaultLang

	var data []entity.Feedback

	q := r.NewSelect().Model(&data)

	if filter.Limit != nil {
		q.Limit(*filter.Limit)
	}

	if filter.Offset != nil {
		q.Offset(*filter.Offset)
	}

	count, err := q.ScanAndCount(ctx)

	var list []feedback.ClientGetList
	for _, i := range data {
		var detail feedback.ClientGetList
		detail.ID = i.ID
		n := i.Name[lang]
		detail.Name = &n

		list = append(list, detail)
	}

	return list, count, err
}

func NewRepository(DB *postgresql.Database) *Repository {
	return &Repository{DB}
}
