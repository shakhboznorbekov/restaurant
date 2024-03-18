package user

import (
	"context"
	"github.com/restaurant/internal/entity"
)

type Service struct {
	repo Repository
}

func (s Service) IsWaiterPhoneExists(ctx context.Context, phone string) (bool, error) {
	return s.repo.IsWaiterPhoneExists(ctx, phone)
}

func (s Service) IsSABCPhoneExists(ctx context.Context, phone string) (bool, error) {
	return s.repo.IsSABCPhoneExists(ctx, phone)
}

func (s Service) IsPhoneExists(ctx context.Context, phone string) (bool, error) {
	return s.repo.IsPhoneExists(ctx, phone)
}

// @branch

func (s Service) BranchCreate(ctx context.Context, request BranchCreateRequest) (BranchCreateResponse, error) {
	return s.repo.BranchCreate(ctx, request)
}

// @admin

func (s Service) AdminGetList(ctx context.Context, filter Filter) ([]AdminGetList, int, error) {
	return s.repo.AdminGetList(ctx, filter)
}

func (s Service) AdminCreate(ctx context.Context, request AdminCreateRequest) (AdminCreateResponse, error) {
	return s.repo.AdminCreate(ctx, request)
}

func (s Service) AdminUpdateColumns(ctx context.Context, request AdminUpdateRequest) error {
	return s.repo.AdminUpdateColumns(ctx, request)
}

func (s Service) AdminGetDetailByRestaurantID(ctx context.Context, restaurantID int64) (entity.User, error) {
	return s.repo.AdminGetDetailByRestaurantID(ctx, restaurantID)
}

func (s Service) AdminUpdateColumnsByRestaurantID(ctx context.Context, request AdminUpdateByRestaurantIDRequest) error {
	return s.repo.AdminUpdateColumnsByRestaurantID(ctx, request)
}

// @client

func (s Service) GetByPhone(ctx context.Context, phone string) (entity.User, error) {
	return s.repo.GetByPhone(ctx, phone)
}

func (s Service) ClientCreate(ctx context.Context, request ClientCreateRequest) (ClientCreateResponse, error) {
	return s.repo.ClientCreate(ctx, request)
}

func (s Service) ClientUpdateAll(ctx context.Context, request ClientUpdateRequest) error {
	return s.repo.ClientUpdateAll(ctx, request)
}

func (s Service) ClientUpdateColumn(ctx context.Context, request ClientUpdateRequest) error {
	return s.repo.ClientUpdateColumn(ctx, request)
}

func (s Service) ClientGetMe(ctx context.Context, id int64) (entity.User, error) {
	return s.repo.ClientGetMe(ctx, id)
}

func (s Service) ClientDeleteMe(ctx context.Context) error {
	return s.repo.ClientDeleteMe(ctx)
}

func (s Service) ClientUpdateMePhone(ctx context.Context, newPhone string) error {
	return s.repo.ClientUpdateMePhone(ctx, newPhone)
}

// @super-admin

func (s Service) SuperAdminGetList(ctx context.Context, filter Filter) ([]SuperAdminGetList, int, error) {
	return s.repo.SuperAdminGetList(ctx, filter)
}

func (s Service) SuperAdminGetDetail(ctx context.Context, id int64) (entity.User, error) {
	return s.repo.SuperAdminGetDetail(ctx, id)
}

func (s Service) SuperAdminCreate(ctx context.Context, request SuperAdminCreateRequest) (SuperAdminCreateResponse, error) {
	return s.repo.SuperAdminCreate(ctx, request)
}

func (s Service) SuperAdminUpdateAll(ctx context.Context, request SuperAdminUpdateRequest) error {
	return s.repo.SuperAdminUpdateAll(ctx, request)
}

func (s Service) SuperAdminUpdateColumns(ctx context.Context, request SuperAdminUpdateRequest) error {
	return s.repo.SuperAdminUpdateColumns(ctx, request)
}

func (s Service) SuperAdminDelete(ctx context.Context, id int64) error {
	return s.repo.SuperAdminDelete(ctx, id)
}

// @waiter

func (s Service) WaiterUpdateMePhone(ctx context.Context, newPhone string) error {
	return s.repo.WaiterUpdateMePhone(ctx, newPhone)
}

func (s Service) WaiterUpdatePassword(ctx context.Context, password string, waiterId int64) error {
	return s.repo.WaiterUpdatePassword(ctx, password, waiterId)
}

// @cashier

func (s Service) CashierGetMe(ctx context.Context) (*CashierGetMeResponse, error) {
	return s.repo.CashierGetMe(ctx)
}

// general

func (s Service) GetMe(ctx context.Context, userID int64) (*GetMeResponse, error) {
	return s.repo.GetMe(ctx, userID)
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}
