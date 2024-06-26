package entity

import (
	"github.com/lib/pq"
	"github.com/uptrace/bun"
	"time"
)

type Menu struct {
	bun.BaseModel `bun:"table:menus"`

	ID             int64                  `json:"id" bun:"id,pk,autoincrement"`
	FoodID         *pq.Int64Array         `json:"food_id" bun:"food_ids"`
	NewPrice       *float64               `json:"new_price" bun:"new_price"`
	OldPrice       *float64               `json:"old_price" bun:"old_price"`
	BranchID       *int64                 `json:"branch_id" bun:"branch_id"`
	Status         *string                `json:"status" bun:"status"`
	Description    map[string]interface{} `json:"description" bun:"description"`
	CategoryID     *int64                 `json:"category_id" bun:"category_id"`
	MenuCategoryID *int64                 `json:"menu_category_id" bun:"menu_category_id"`
	Photos         *pq.StringArray        `json:"photos" bun:"photos"`
	Name           *string                `json:"name" bun:"name"`

	CreatedAt *time.Time `json:"created_at" bun:"created_at"`
	CreatedBy *int64     `json:"created_by" bun:"created_by"`
	UpdatedAt *time.Time `json:"updated_at" bun:"updated_at"`
	UpdatedBy *int64     `json:"updated_by" bun:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" bun:"deleted_at"`
	DeletedBy *int64     `json:"deleted_by" bun:"deleted_by"`
}
