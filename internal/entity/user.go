package entity

import (
	"github.com/uptrace/bun"
	"time"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	ID               int64      `json:"id" bun:"id,pk,autoincrement"`
	Name             *string    `json:"name" bun:"name"`
	Phone            *string    `json:"phone" bun:"phone"`
	BirthDate        *time.Time `json:"birth_date" bun:"birth_date"`
	Gender           *string    `json:"gender" bun:"gender"`
	Role             *string    `json:"role" bun:"role"`
	Password         *string    `json:"password" bun:"password"`
	Status           *string    `json:"status" bun:"status"`
	Photo            *string    `json:"photo" bun:"photo"`
	Address          *string    `json:"address" bun:"address"`
	ServicePercentID *int64     `json:"service_percent" bun:"service_percent"`
	BranchID         *int64     `json:"branch_id" bun:"branch_id"`
	RestaurantID     *int64     `json:"restaurant_id" bun:"restaurant_id"`
	CreatedAt        *time.Time `json:"created_at" bun:"created_at"`
	CreatedBy        *int64     `json:"created_by" bun:"created_by"`
	UpdatedAt        *time.Time `json:"updated_at" bun:"updated_at"`
	UpdatedBy        *int64     `json:"updated_by" bun:"updated_by"`
	DeletedAt        *time.Time `json:"deleted_at" bun:"deleted_at"`
	DeletedBy        *int64     `json:"deleted_by" bun:"deleted_by"`
}
