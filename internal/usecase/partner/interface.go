package partner

import (
	"context"
	"github.com/restaurant/internal/service/partner"
)

type Partner interface {
	// @admin

	AdminGetList(ctx context.Context, filter partner.Filter) ([]partner.AdminGetList, int, error)
	AdminGetDetail(ctx context.Context, id int64) (partner.AdminGetDetail, error)
	AdminCreate(ctx context.Context, request partner.AdminCreateRequest) (partner.AdminCreateResponse, error)
	AdminUpdateAll(ctx context.Context, request partner.AdminUpdateRequest) error
	AdminUpdateColumns(ctx context.Context, request partner.AdminUpdateRequest) error
	AdminDelete(ctx context.Context, id int64) error

	// @admin

	BranchGetList(ctx context.Context, filter partner.Filter) ([]partner.BranchGetList, int, error)
}
