package full_search

import (
	"context"
	"strings"
	"time"
)

type Service struct {
	repo Repository
}

// @client

func (s Service) ClientGetList(ctx context.Context, filter Filter) ([]ClientGetList, error) {
	list, err := s.repo.ClientGetList(ctx, filter)
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

	return list, err
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

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}
