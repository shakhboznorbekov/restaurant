package partner

import (
	"context"
	"github.com/restaurant/internal/service/partner"
)

type UseCase struct {
	partner Partner
}

func NewUseCase(partner Partner) *UseCase {
	return &UseCase{partner: partner}
}

// partner

// @admin

func (uu UseCase) AdminGetPartnerList(ctx context.Context, filter partner.Filter) ([]partner.AdminGetList, int, error) {
	return uu.partner.AdminGetList(ctx, filter)
}

func (uu UseCase) AdminGetPartnerDetail(ctx context.Context, id int64) (partner.AdminGetDetail, error) {
	return uu.partner.AdminGetDetail(ctx, id)
}

func (uu UseCase) AdminCreatePartner(ctx context.Context, data partner.AdminCreateRequest) (partner.AdminCreateResponse, error) {
	return uu.partner.AdminCreate(ctx, data)
}

func (uu UseCase) AdminUpdatePartner(ctx context.Context, data partner.AdminUpdateRequest) error {
	return uu.partner.AdminUpdateAll(ctx, data)
}

func (uu UseCase) AdminUpdatePartnerColumn(ctx context.Context, data partner.AdminUpdateRequest) error {
	return uu.partner.AdminUpdateColumns(ctx, data)
}

func (uu UseCase) AdminDeletePartner(ctx context.Context, id int64) error {
	return uu.partner.AdminDelete(ctx, id)
}

// @admin

func (uu UseCase) BranchGetPartnerList(ctx context.Context, filter partner.Filter) ([]partner.BranchGetList, int, error) {
	return uu.partner.BranchGetList(ctx, filter)
}
