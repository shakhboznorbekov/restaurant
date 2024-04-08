package menu

import (
	"context"
	"github.com/restaurant/internal/entity"
	"strings"
	"time"
)

type Service struct {
	repo Repository
}

// @admin

func (s Service) AdminGetList(ctx context.Context, filter Filter) ([]AdminGetList, int, error) {
	return s.repo.AdminGetList(ctx, filter)
}

func (s Service) AdminGetDetail(ctx context.Context, id int64) (entity.Menu, error) {
	return s.repo.AdminGetDetail(ctx, id)
}

func (s Service) AdminCreate(ctx context.Context, request AdminCreateRequest) ([]AdminCreateResponse, error) {
	return s.repo.AdminCreate(ctx, request)
}

func (s Service) AdminUpdateAll(ctx context.Context, request AdminUpdateRequest) error {
	return s.repo.AdminUpdateAll(ctx, request)
}

func (s Service) AdminUpdateColumns(ctx context.Context, request AdminUpdateRequest) error {
	return s.repo.AdminUpdateColumns(ctx, request)
}

func (s Service) AdminDelete(ctx context.Context, id int64) error {
	return s.repo.AdminDelete(ctx, id)
}

func (s Service) AdminRemovePhoto(ctx context.Context, id int64, index *int) (*string, error) {
	return s.repo.AdminRemovePhoto(ctx, id, index)
}

// @branch

func (s Service) BranchGetList(ctx context.Context, filter Filter) ([]BranchGetList, int, error) {
	return s.repo.BranchGetList(ctx, filter)
}

func (s Service) BranchGetDetail(ctx context.Context, id int64) (entity.Menu, error) {
	return s.repo.BranchGetDetail(ctx, id)
}

func (s Service) BranchCreate(ctx context.Context, request BranchCreateRequest) (BranchCreateResponse, error) {
	return s.repo.BranchCreate(ctx, request)
}

func (s Service) BranchUpdateAll(ctx context.Context, request BranchUpdateRequest) error {
	return s.repo.BranchUpdateAll(ctx, request)
}

func (s Service) BranchUpdateColumns(ctx context.Context, request BranchUpdateRequest) error {
	return s.repo.BranchUpdateColumns(ctx, request)
}

func (s Service) BranchDelete(ctx context.Context, id int64) error {
	return s.repo.BranchDelete(ctx, id)
}

func (s Service) BranchUpdatePrinterID(ctx context.Context, request BranchUpdatePrinterIDRequest) error {
	return s.repo.BranchUpdatePrinterID(ctx, request)
}

func (s Service) BranchDeletePrinterID(ctx context.Context, menuID int64) error {
	return s.repo.BranchDeletePrinterID(ctx, menuID)
}

func (s Service) BranchRemovePhoto(ctx context.Context, id int64, index int) (*string, error) {
	return s.repo.BranchRemovePhoto(ctx, id, index)
}

// @client

func (s Service) ClientGetList(ctx context.Context, filter Filter) ([]ClientGetList, error) {
	return s.repo.ClientGetList(ctx, filter)
}

func (s Service) CashierGetList(ctx context.Context, filter Filter) ([]CashierGetList, int, error) {
	return s.repo.CashierGetList(ctx, filter)
}

func (s Service) CashierGetDetail(ctx context.Context, id int64) (entity.Menu, error) {
	return s.repo.CashierGetDetail(ctx, id)
}

func (s Service) CashierCreate(ctx context.Context, request CashierCreateRequest) (CashierCreateResponse, error) {
	return s.repo.CashierCreate(ctx, request)
}

func (s Service) CashierUpdateAll(ctx context.Context, request CashierUpdateRequest) error {
	return s.repo.CashierUpdateAll(ctx, request)
}

func (s Service) CashierUpdateColumn(ctx context.Context, request CashierUpdateRequest) error {
	return s.repo.CashierUpdateColumn(ctx, request)
}

func (s Service) CashierDelete(ctx context.Context, id int64) error {
	return s.repo.CashierDelete(ctx, id)
}

func (s Service) CashierUpdatePrinterID(ctx context.Context, request CashierUpdatePrinterIDRequest) error {
	return s.repo.CashierUpdatePrinterID(ctx, request)
}

func (s Service) CashierDeletePrinterID(ctx context.Context, menuID int64) error {
	return s.repo.CashierDeletePrinterID(ctx, menuID)
}

func (s Service) CashierRemovePhoto(ctx context.Context, id int64, index int) (*string, error) {
	return s.repo.CashierRemovePhoto(ctx, id, index)
}

// @client

func (s Service) ClientGetDetail(ctx context.Context, id int64) (ClientGetDetail, error) {
	return s.repo.ClientGetDetail(ctx, id)
}

func (s Service) ClientGetListByCategoryID(ctx context.Context, foodCategoryID int, filter Filter) ([]ClientGetListByCategoryID, error) {
	list, err := s.repo.ClientGetListByCategoryID(ctx, foodCategoryID, filter)
	if err != nil {
		return nil, err
	}

	for k, v := range list {
		if v.WorkTimeToday != nil {
			workTimeToday := strings.Split(*v.WorkTimeToday, "-")
			if len(workTimeToday) == 2 {
				if isTimeWithinPeriod(workTimeToday[0], workTimeToday[1]) {
					b := false
					list[k].IsClosed = &b
				} else {
					b := true
					list[k].IsClosed = &b
				}

				list[k].OpenTime = &workTimeToday[0]
				list[k].CloseTime = &workTimeToday[1]
			}
		}
	}
	return list, nil
}

// @cashier

func (s Service) CashierUpdateColumns(ctx context.Context, request CashierUpdateMenuStatus) error {
	return s.repo.CashierUpdateColumns(ctx, request)
}

// @waiter

func (s Service) WaiterGetMenuList(ctx context.Context, filter Filter) ([]WaiterGetMenuListResponse, error) {
	return s.repo.WaiterGetMenuList(ctx, filter)
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func isTimeWithinPeriod(startTimeStr, endTimeStr string) bool {
	currentTime := time.Now()

	// Parse the start and end times
	layout := "15:04"
	startTime, err := time.Parse(layout, startTimeStr)
	if err != nil {
		return false
	}

	endTime, err := time.Parse(layout, endTimeStr)
	if err != nil {
		return false
	}

	// Normalize start and end times to today's date
	startTime = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), startTime.Hour(), startTime.Minute(), 0, 0, currentTime.Location())
	endTime = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), endTime.Hour(), endTime.Minute(), 0, 0, currentTime.Location())

	// Check if current time is within the range
	return currentTime.After(startTime) && currentTime.Before(endTime)
}
