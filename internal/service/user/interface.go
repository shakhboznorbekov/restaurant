package user

import (
	"context"
	"github.com/postgresql-restaurant/internal/entity"
)

type Repository interface {
	SuperAdminCreate(ctx context.Context, request SuperAdminCreateRequest) (SuperAdminCreateResponse, error)
	SuperAdminGetList(ctx context.Context, filter Filter) ([]SuperAdminGetList, int, error)
	SuperAdminGetDetail(ctx context.Context, id int64) (entity.User, error)
	SuperAdminUpdateAll(ctx context.Context, request SuperAdminUpdateRequest) error
	SuperAdminUpdateColumns(ctx context.Context, request SuperAdminUpdateRequest) error
	SuperAdminDelete(ctx context.Context, id int64) error
}
