package branch

import (
	"context"
	"github.com/restaurant/internal/entity"
	"time"
)

type Repository interface {

	// @admin

	AdminGetList(ctx context.Context, filter Filter) ([]AdminGetList, int, error)
	AdminGetDetail(ctx context.Context, id int64) (entity.Branch, error)
	AdminCreate(ctx context.Context, request AdminCreateRequest) (AdminCreateResponse, error)
	AdminUpdateAll(ctx context.Context, request AdminUpdateRequest) error
	AdminUpdateColumns(ctx context.Context, request AdminUpdateRequest) error
	AdminDelete(ctx context.Context, id int64) error
	AdminDeleteImage(ctx context.Context, request AdminDeleteImageRequest) error

	// @client

	ClientGetList(ctx context.Context, filter Filter) ([]ClientGetList, int, error)
	ClientGetMapList(ctx context.Context, filter Filter) ([]ClientGetMapList, int, error)
	ClientGetDetail(ctx context.Context, id int64) (ClientGetDetail, error)
	ClientNearlyBranchGetList(ctx context.Context, filter Filter) ([]ClientGetList, int, error)
	ClientUpdateColumns(ctx context.Context, request ClientUpdateRequest) error
	ClientAddSearchCount(ctx context.Context, branchID int64) error
	ClientGetListOrderSearchCount(ctx context.Context, filter Filter) ([]ClientGetList, int, error)
	ClientGetListByCategoryID(ctx context.Context, filter Filter, CategoryID int64) ([]ClientGetList, int, error)

	// @branch

	BranchGetDetail(ctx context.Context, id int64) (BranchGetDetail, error)

	// @cashier

	CashierGetDetail(ctx context.Context, id int64) (CashierGetDetail, error)

	// @token

	BranchGetToken(ctx context.Context) (BranchGetToken, error)
	WsGetByToken(ctx context.Context, token string) (WsGetByTokenResponse, error)
	WsUpdateTokenExpiredAt(ctx context.Context, id int64) (string, error)
}

type RedisRepository interface {
	SetBranch(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	GetBranch(ctx context.Context, key string) (string, error)
}
