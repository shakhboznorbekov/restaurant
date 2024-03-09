package user

import (
	"context"
	"github.com/restaurant/internal/entity"
)

type Repository interface {
	// @super-admin
	SuperAdminCreate(ctx context.Context, request SuperAdminCreateRequest) (SuperAdminCreateResponse, error)
	SuperAdminGetList(ctx context.Context, filter Filter) ([]SuperAdminGetList, int, error)
	SuperAdminGetDetail(ctx context.Context, id int64) (entity.User, error)
	SuperAdminUpdateAll(ctx context.Context, request SuperAdminUpdateRequest) error
	SuperAdminUpdateColumns(ctx context.Context, request SuperAdminUpdateRequest) error
	SuperAdminDelete(ctx context.Context, id int64) error

	// others

	IsPhoneExists(ctx context.Context, phone string) (bool, error)
	IsWaiterPhoneExists(ctx context.Context, phone string) (bool, error)
}
