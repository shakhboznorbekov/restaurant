package menu

import (
	"time"

	"github.com/lib/pq"
	"github.com/uptrace/bun"
)

type Filter struct {
	Limit      *int
	Offset     *int
	Page       *int
	Search     *string
	CategoryId *int
	Lat        *float64
	Lon        *float64
	BranchID   *int64
	PrinterID  *int
	Printer    *bool
}

// @admin

type AdminGetList struct {
	//ID         int64    `json:"id" bun:"id"`
	//FoodID     *int64   `json:"food_id" bun:"food_id"`
	//BranchID   *int64   `json:"branch_id" bun:"branch_id"`
	//Status     *string  `json:"status" bun:"status"`
	//FoodName   *string  `json:"food_name" bun:"food_name"`
	//BranchName *string  `json:"branch_name" bun:"branch_name"`
	//OldPrice   *float64 `json:"old_price" bun:"old_price"`
	//NewPrice   *float64 `json:"new_price" bun:"new_price"`

	CategoryID   *int64  `json:"category_id"`
	CategoryName *string `json:"category_name"`
	UserID       *int64  `json:"-"`
	Menus        []struct {
		ID     int64           `json:"id"`
		Name   string          `json:"name"`
		Photos *pq.StringArray `json:"photos"`
		Price  *float32        `json:"price"`
		Count  *int            `json:"count"`
		Status *string         `json:"status"`
	} `json:"menus"`
}

type AdminGetDetail struct {
	ID       int64    `json:"id"`
	FoodID   *int64   `json:"food_id"`
	BranchID *int64   `json:"branch_id"`
	Status   *string  `json:"status"`
	OldPrice *float64 `json:"old_price"`
	NewPrice *float64 `json:"new_price"`
}

type AdminCreateRequest struct {
	FoodID      *int64                 `json:"food_id" form:"food_id"`
	NewPrice    *float64               `json:"new_price" form:"new_price"`
	BranchID    []int64                `json:"branch_id" form:"branch_id"`
	Description map[string]interface{} `json:"description" form:"description"`
}

type AdminCreateResponse struct {
	bun.BaseModel `bun:"table:menus"`

	ID          int64                  `json:"id" bun:"id,pk,autoincrement"`
	FoodID      *int64                 `json:"food_id" bun:"food_id"`
	BranchID    *int64                 `json:"branch_id" bun:"branch_id"`
	Status      *string                `json:"status" bun:"status"`
	NewPrice    *float64               `json:"new_price" bun:"new_price"`
	Description map[string]interface{} `json:"description" bun:"description"`

	CreatedAt time.Time `json:"created_at" bun:"created_at"`
	CreatedBy int64     `json:"created_by" bun:"created_by"`
}

type AdminUpdateRequest struct {
	ID          int64                  `json:"id" form:"id"`
	FoodID      *int64                 `json:"food_id" form:"food_id"`
	BranchID    *int64                 `json:"branch_id" form:"branch_id"`
	NewPrice    *float64               `json:"new_price" form:"new_price"`
	Status      *string                `json:"status" form:"status"`
	Description map[string]interface{} `json:"description" form:"description"`
}

// @branch

type BranchGetList struct {
	//ID         int64    `json:"id" bun:"id"`
	//FoodID     *int64   `json:"food_id" bun:"food_id"`
	//BranchID   *int64   `json:"branch_id" bun:"branch_id"`
	//Status     *string  `json:"status" bun:"status"`
	//FoodName   *string  `json:"food_name" bun:"food_name"`
	//BranchName *string  `json:"branch_name" bun:"branch_name"`
	//OldPrice   *float64 `json:"old_price" bun:"old_price"`
	//NewPrice   *float64 `json:"new_price" bun:"new_price"`
	//Photo      *string  `json:"photo" bun:"photo"`

	CategoryID   *int64  `json:"category_id"`
	CategoryName *string `json:"category_name"`
	UserID       *int64  `json:"-"`
	Menus        []struct {
		ID      int64           `json:"id"`
		Name    string          `json:"name"`
		Photos  *pq.StringArray `json:"photos"`
		Price   *float32        `json:"price"`
		Count   *int            `json:"count"`
		Status  *string         `json:"status"`
		Printer *bool           `json:"printer"`
	} `json:"menus"`
}

type BranchGetDetail struct {
	ID       int64    `json:"id"`
	FoodID   *int64   `json:"food_id"`
	BranchID *int64   `json:"branch_id"`
	Status   *string  `json:"status"`
	OldPrice *float64 `json:"old_price"`
	NewPrice *float64 `json:"new_price"`
}

type BranchCreateRequest struct {
	FoodID      *int64                 `json:"food_id" form:"food_id"`
	NewPrice    *float64               `json:"new_price" form:"new_price"`
	Description map[string]interface{} `json:"description" form:"description"`
}

type BranchCreateResponse struct {
	bun.BaseModel `bun:"table:menus"`

	ID          int64                  `json:"id" bun:"id,pk,autoincrement"`
	FoodID      *int64                 `json:"food_id" bun:"food_id"`
	NewPrice    *float64               `json:"new_price" bun:"new_price"`
	BranchID    *int64                 `json:"branch_id" bun:"branch_id"`
	Status      *string                `json:"status" bun:"status"`
	Description map[string]interface{} `json:"description" bun:"description"`

	CreatedAt time.Time `json:"created_at" bun:"created_at"`
	CreatedBy int64     `json:"created_by" bun:"created_by"`
}

type BranchUpdateRequest struct {
	ID          int64                  `json:"id" form:"id"`
	FoodID      *int64                 `json:"food_id" form:"food_id"`
	NewPrice    *float64               `json:"new_price" form:"new_price"`
	Status      *string                `json:"status" form:"status"`
	Description map[string]interface{} `json:"description" form:"description"`
}

type BranchUpdatePrinterIDRequest struct {
	PrinterID *int64  `json:"printer_id"`
	MenuIds   []int64 `json:"menu_ids"`
}

// @client

type ClientGetList struct {
	CategoryID   *int64  `json:"category_id"`
	CategoryName *string `json:"category_name"`
	UserID       *int64  `json:"-"`
	Menus        []struct {
		ID     int64           `json:"id"`
		Name   string          `json:"name"`
		Photos *pq.StringArray `json:"photos"`
		Price  *float32        `json:"price"`
		Count  *int            `json:"count"`
	} `json:"menus"`
}

type ClientGetDetail struct {
	ID     int64           `json:"id"`
	Name   string          `json:"name"`
	Photos *pq.StringArray `json:"photos"`
	Price  *float32        `json:"price"`
	Count  *int            `json:"count"`
}

type ClientGetListByCategoryID struct {
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
}

type Menu struct {
	ID     int64           `json:"id"`
	Name   *string         `json:"name"`
	Photos *pq.StringArray `json:"photos"`
	Price  *float32        `json:"price"`
}

// @cashier

type CashierUpdateMenuStatus struct {
	ID     int64   `json:"id" form:"id"`
	Status *string `json:"status" form:"status"`
}

type CashierGetList struct {
	//ID         int64    `json:"id" bun:"id"`
	//FoodID     *int64   `json:"food_id" bun:"food_id"`
	//BranchID   *int64   `json:"branch_id" bun:"branch_id"`
	//Status     *string  `json:"status" bun:"status"`
	//FoodName   *string  `json:"food_name" bun:"food_name"`
	//BranchName *string  `json:"branch_name" bun:"branch_name"`
	//OldPrice   *float64 `json:"old_price" bun:"old_price"`
	//NewPrice   *float64 `json:"new_price" bun:"new_price"`
	//Photo      *string  `json:"photo" bun:"photo"`

	CategoryID   *int64  `json:"category_id"`
	CategoryName *string `json:"category_name"`
	UserID       *int64  `json:"-"`
	Menus        []struct {
		ID      int64           `json:"id"`
		Name    string          `json:"name"`
		Photos  *pq.StringArray `json:"photos"`
		Price   *float32        `json:"price"`
		Count   *int            `json:"count"`
		Status  *string         `json:"status"`
		Printer *bool           `json:"printer"`
	} `json:"menus"`
}

type CashierGetDetail struct {
	ID          int64                  `json:"id"`
	FoodID      *int64                 `json:"food_id"`
	BranchID    *int64                 `json:"branch_id"`
	Status      *string                `json:"status"`
	OldPrice    *float64               `json:"old_price"`
	NewPrice    *float64               `json:"new_price"`
	Description map[string]interface{} `json:"description" form:"description"`
}

type CashierCreateRequest struct {
	FoodID      *int64                 `json:"food_id" form:"food_id"`
	NewPrice    *float64               `json:"new_price" form:"new_price"`
	Description map[string]interface{} `json:"description" form:"description"`
}

type CashierCreateResponse struct {
	bun.BaseModel `bun:"table:menus"`

	ID          int64                  `json:"id" bun:"id,pk,autoincrement"`
	FoodID      *int64                 `json:"food_id" bun:"food_id"`
	NewPrice    *float64               `json:"new_price" bun:"new_price"`
	BranchID    *int64                 `json:"branch_id" bun:"branch_id"`
	Status      *string                `json:"status" bun:"status"`
	Description map[string]interface{} `json:"description" bun:"description"`

	CreatedAt time.Time `json:"created_at" bun:"created_at"`
	CreatedBy int64     `json:"created_by" bun:"created_by"`
}

type CashierUpdateRequest struct {
	ID          int64                  `json:"id" form:"id"`
	FoodID      *int64                 `json:"food_id" form:"food_id"`
	NewPrice    *float64               `json:"new_price" form:"new_price"`
	Status      *string                `json:"status" form:"status"`
	Description map[string]interface{} `json:"description" form:"description"`
}

type CashierUpdatePrinterIDRequest struct {
	PrinterID *int64  `json:"printer_id"`
	MenuIds   []int64 `json:"menu_ids"`
}

// @waiter

type WaiterGetMenuListResponse struct {
	Id    int64        `json:"id" bun:"id"`
	Name  *string      `json:"name" bun:"name"`
	Menus []WaiterMenu `json:"menus" bun:"-"`
}

type WaiterMenu struct {
	Id     int64    `json:"id" bun:"id"`
	Name   *string  `json:"name" bun:"name"`
	Price  *float64 `json:"price" bun:"price"`
	Photo  *string  `json:"photo" bun:"photo"`
	Status *string  `json:"status" bun:"status"`
}
