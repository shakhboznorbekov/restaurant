package user_work_time

import (
	"github.com/uptrace/bun"
	"time"
)

type Filter struct {
	Limit    *int
	Offset   *int
	Page     *int
	Date     *string
	WaiterID *int64
}

type GetDetailByID struct {
	ID      int      `json:"id"`
	Date    *string  `json:"date"`
	Periods []Period `json:"periods"`
}

type GetDetailByWaiterIDAndDateResponse struct {
	ID              int              `json:"id"`
	Date            *string          `json:"date"`
	PeriodsResponse []PeriodResponse `json:"periods"`
	Periods         []Period         `json:"-"`
	TotalDuration   *int             `json:"total_duration"`
	Now             bool             `json:"now"`
}

type CreateResponse struct {
	bun.BaseModel `bun:"table:waiter_work_time"`

	ID       int        `json:"id" bun:"-"`
	WaiterID *int64     `json:"waiter_id" bun:"waiter_id"`
	Date     *time.Time `json:"date" bun:"date"`
	Periods  []Period   `json:"periods" bun:"periods"`
}

type Period struct {
	ComeTime *time.Time `json:"start" bun:"come_time"`
	GoneTime *time.Time `json:"finish" bun:"gone_time"`
	Duration *int       `json:"duration"`
}

type PeriodResponse struct {
	ComeTime *int `json:"start"  bun:"come_time"`
	GoneTime *int `json:"finish"  bun:"come_time"`
	Duration *int `json:"duration"`
}

type ListFilter struct {
	Limit    *int
	Offset   *int
	Page     *int
	WaiterID *int
}

type GetListResponse struct {
	ID              *int64           `json:"id"`
	Date            *string          `json:"date"`
	TotalDuration   *int             `json:"total_duration"`
	Now             bool             `json:"now"`
	Periods         []Period         `json:"-"`
	PeriodsResponse []PeriodResponse `json:"periods"`
}

type BranchFilter struct {
	Limit    *int
	Offset   *int
	Page     *int
	Date     *string
	WaiterID *int64
}

type BranchGetListResponse struct {
	ID              *int64           `json:"id"`
	WaiterID        *int64           `json:"waiter_id"`
	Waiter          *string          `json:"name"`
	Avatar          *string          `json:"avatar"`
	TotalDuration   *int             `json:"total_duration"`
	Now             bool             `json:"now"`
	StartTime       *int             `json:"start_time"`
	Duration        *int             `json:"duration"`
	Periods         []Period         `json:"-"`
	PeriodsResponse []PeriodResponse `json:"periods"`
}
