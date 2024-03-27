package branch

import (
	"context"
	"github.com/restaurant/internal/entity"
	"strconv"
	"strings"
	"time"
)

type Service struct {
	repo      Repository
	redisRepo RedisRepository
}

// @admin

func (s Service) AdminGetList(ctx context.Context, filter Filter) ([]AdminGetList, int, error) {
	return s.repo.AdminGetList(ctx, filter)
}

func (s Service) AdminGetDetail(ctx context.Context, id int64) (entity.Branch, error) {
	return s.repo.AdminGetDetail(ctx, id)
}

func (s Service) AdminCreate(ctx context.Context, request AdminCreateRequest) (AdminCreateResponse, error) {
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

func (s Service) AdminDeleteImage(ctx context.Context, request AdminDeleteImageRequest) error {
	return s.repo.AdminDeleteImage(ctx, request)
}

// @client

func (s Service) ClientGetList(ctx context.Context, filter Filter) ([]ClientGetList, int, error) {
	list, count, err := s.repo.ClientGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
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

	return list, count, err
}

func (s Service) ClientGetMapList(ctx context.Context, filter Filter) ([]ClientGetMapList, int, error) {
	return s.repo.ClientGetMapList(ctx, filter)
}

func (s Service) ClientGetDetail(ctx context.Context, id int64) (ClientGetDetail, error) {
	return s.repo.ClientGetDetail(ctx, id)
}

func (s Service) ClientNearlyBranchGetList(ctx context.Context, filter Filter) ([]ClientGetList, int, error) {
	list, count, err := s.repo.ClientNearlyBranchGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
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

	return list, count, nil
}

func (s Service) ClientUpdateColumns(ctx context.Context, request ClientUpdateRequest) error {
	return s.repo.ClientUpdateColumns(ctx, request)
}

func (s Service) ClientAddSearchCount(ctx context.Context, branchID int64) error {
	return s.repo.ClientAddSearchCount(ctx, branchID)
}

func (s Service) ClientGetListByCategoryID(ctx context.Context, filter Filter, CategoryID int64) ([]ClientGetList, int, error) {
	return s.repo.ClientGetListByCategoryID(ctx, filter, CategoryID)
}

func (s Service) ClientGetListOrderSearchCount(ctx context.Context, filter Filter) ([]ClientGetList, int, error) {
	return s.repo.ClientGetListOrderSearchCount(ctx, filter)
}

// @branch

func (s Service) BranchGetDetail(ctx context.Context, id int64) (BranchGetDetail, error) {
	return s.repo.BranchGetDetail(ctx, id)
}

// @cashier

func (s Service) CashierGetDetail(ctx context.Context, id int64) (CashierGetDetail, error) {
	return s.repo.CashierGetDetail(ctx, id)
}

// @token

func (s Service) BranchGetToken(ctx context.Context) (BranchGetToken, error) {
	return s.repo.BranchGetToken(ctx)
}

func (s Service) WsGetByToken(ctx context.Context, token string) (WsGetByTokenResponse, error) {
	return s.repo.WsGetByToken(ctx, token)
}

func (s Service) WsUpdateTokenExpiredAt(ctx context.Context, id int64) (string, error) {
	return s.repo.WsUpdateTokenExpiredAt(ctx, id)
}

func NewService(repo Repository, redisRepo RedisRepository) *Service {
	return &Service{repo, redisRepo}
}

// #redis

func (s Service) SetBranch(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return s.redisRepo.SetBranch(ctx, key, value, expiration)
}

func (s Service) GetBranch(ctx context.Context, key string) (string, error) {
	return s.redisRepo.GetBranch(ctx, key)
}

// additional

func StringToTime(number int) string {
	h := strconv.Itoa(number / 60)
	m := strconv.Itoa(number % 60)

	if i := number % 60; i < 10 {
		m = "0" + strconv.Itoa(number%60)
	}
	if i := number / 60; i < 10 {
		h = "0" + strconv.Itoa(number/60)
	}
	hm := h + ":" + m
	return hm
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
