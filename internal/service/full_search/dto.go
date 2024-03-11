package full_search

import (
	"github.com/lib/pq"
	"github.com/restaurant/internal/pkg/utils"
)

type Filter struct {
	Limit          *int
	Offset         *int
	Page           *int
	FoodCategoryID *int64
	Search         *string
	Lon            *float64
	Menu           *string
	Lat            *float64
	BranchID       *int64
	Fields         map[string][]string
	Joins          map[string]utils.Joins
}

type ClientGetList struct {
	ID            int64              `json:"id" bun:"id"`
	Location      map[string]float32 `json:"location" bun:"location"`
	Photos        *pq.StringArray    `json:"photos" bun:"photos"`
	Status        *string            `json:"status" bun:"status"`
	OpenTime      *string            `json:"open_time" bun:"-"`
	CloseTime     *string            `json:"close_time" bun:"-"`
	Name          *string            `json:"name" bun:"name"`
	CategoryID    *int64             `json:"category_id" bun:"category_id"`
	CategoryName  *string            `json:"category_name" bun:"category_name"`
	Point         *int               `json:"point" bun:"point"`
	Rate          *float32           `json:"rate" bun:"rate"`
	Distance      *string            `json:"distance" bun:"distance"`
	WorkTimeToday *string            `json:"-" bun:"work_time_today"`
	IsLiked       *bool              `json:"is_liked" bun:"is_liked"`
	IsClosed      *bool              `json:"is_closed" bun:"-"`
	Menus         []Menu             `json:"menus" bun:"-"`
	MenuStatus    bool               `json:"-" bun:"menu_status"`
}

type Menu struct {
	ID     int64           `json:"id"`
	Name   *string         `json:"name"`
	Photos *pq.StringArray `json:"photos"`
	Price  *float32        `json:"price"`
}
