package product

import (
	"github.com/restaurant/internal/pkg/utils"
	"github.com/uptrace/bun"
	"time"
)

type Filter struct {
	Limit  *int
	Offset *int
	Name   *string
	Fields map[string][]string
	Joins  map[string]utils.Joins
}

// @admin

type AdminGetList struct {
	ID            int64   `json:"id"`
	Name          *string `json:"name"`
	MeasureUnitID *int64  `json:"measure_unit_id"`
	MeasureUnit   *string `json:"measure_unit"`
	Barcode       *string `json:"barcode"`
}

type AdminGetDetail struct {
	ID            int64   `json:"id"`
	Name          *string `json:"name"`
	MeasureUnitID *int64  `json:"measure_unit_id"`
	Barcode       *string `json:"barcode"`
}

type AdminCreateRequest struct {
	Name          *string `json:"name" form:"name"`
	MeasureUnitID *int64  `json:"measure_unit_id" form:"measure_unit_id"`
	Barcode       *string `json:"barcode" form:"barcode"`
}

type AdminCreateResponse struct {
	bun.BaseModel `bun:"table:products"`

	ID            int64     `json:"id" bun:"id,pk,autoincrement"`
	Name          *string   `json:"name" bun:"name"`
	MeasureUnitID *int64    `json:"measure_unit_id" bun:"measure_unit_id"`
	Barcode       *string   `json:"barcode" bun:"barcode"`
	CreatedAt     time.Time `json:"-" bun:"created_at"`
	CreatedBy     int64     `json:"-" bun:"created_by"`
	RestaurantID  int64     `json:"-" bun:"restaurant_id"`
}

type AdminUpdateRequest struct {
	ID            int64   `json:"id" form:"id"`
	Name          *string `json:"name" form:"name"`
	Barcode       *string `json:"barcode" form:"barcode"`
	MeasureUnitID *int64  `json:"measure_unit_id" form:"measure_unit_id"`
}

// others

type SpendingFilter struct {
	FromDate *time.Time `json:"from_date"`
	ToDate   *time.Time `json:"to_date"`
	BranchId *int       `json:"branch_id"`
	Limit    *int       `json:"limit"`
	Offset   *int       `json:"offset"`
}

type AdminGetSpendingByBranchResponse struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`

	Amount      float64 `json:"amount"`
	MeasureUnit string  `json:"measure_unit"`
}

type CashierGetSpendingResponse struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`

	Amount      float64 `json:"amount"`
	MeasureUnit string  `json:"measure_unit"`
}

// @branch

type BranchGetList struct {
	ID            int64   `json:"id"`
	Name          *string `json:"name"`
	MeasureUnitID *int64  `json:"measure_unit_id"`
	MeasureUnit   *string `json:"measure_unit"`
	Barcode       *string `json:"barcode"`
}

type BranchGetDetail struct {
	ID            int64   `json:"id"`
	Name          *string `json:"name"`
	MeasureUnitID *int64  `json:"measure_unit_id"`
	Barcode       *string `json:"barcode"`
}

type BranchCreateRequest struct {
	Name          *string `json:"name" form:"name"`
	MeasureUnitID *int64  `json:"measure_unit_id" form:"measure_unit_id"`
	Barcode       *string `json:"barcode" form:"barcode"`
}

type BranchCreateResponse struct {
	bun.BaseModel `bun:"table:products"`

	ID            int64     `json:"id" bun:"id,pk,autoincrement"`
	Name          *string   `json:"name" bun:"name"`
	MeasureUnitID *int64    `json:"measure_unit_id" bun:"measure_unit_id"`
	Barcode       *string   `json:"barcode" bun:"barcode"`
	CreatedAt     time.Time `json:"-" bun:"created_at"`
	CreatedBy     int64     `json:"-" bun:"created_by"`
	RestaurantID  int64     `json:"-" bun:"restaurant_id"`
}

type BranchUpdateRequest struct {
	ID            int64   `json:"id" form:"id"`
	Name          *string `json:"name" form:"name"`
	Barcode       *string `json:"barcode" form:"barcode"`
	MeasureUnitID *int64  `json:"measure_unit_id" form:"measure_unit_id"`
}

// @cashier

type CashierGetList struct {
	ID            int64   `json:"id"`
	Name          *string `json:"name"`
	MeasureUnitID *int64  `json:"measure_unit_id"`
	MeasureUnit   *string `json:"measure_unit"`
}

type CashierGetDetail struct {
	ID            int64   `json:"id"`
	Name          *string `json:"name"`
	MeasureUnitID *int64  `json:"measure_unit_id"`
}

type CashierCreateRequest struct {
	Name          *string `json:"name" form:"name"`
	MeasureUnitID *int64  `json:"measure_unit_id" form:"measure_unit_id"`
}

type CashierCreateResponse struct {
	bun.BaseModel `bun:"table:products"`

	ID            int64     `json:"id" bun:"id,pk,autoincrement"`
	Name          *string   `json:"name" bun:"name"`
	MeasureUnitID *int64    `json:"measure_unit_id" bun:"measure_unit_id"`
	CreatedAt     time.Time `json:"-" bun:"created_at"`
	CreatedBy     int64     `json:"-" bun:"created_by"`
	RestaurantID  int64     `json:"-" bun:"restaurant_id"`
}

type CashierUpdateRequest struct {
	ID            int64   `json:"id" form:"id"`
	Name          *string `json:"name" bun:"name"`
	MeasureUnitID *int64  `json:"measure_unit_id" bun:"measure_unit_id"`
}
