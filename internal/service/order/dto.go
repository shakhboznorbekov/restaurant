package order

import (
	"github.com/lib/pq"
	"github.com/restaurant/internal/service/waiter"
	"github.com/uptrace/bun"
	"time"
)

type Filter struct {
	Limit    *int
	Offset   *int
	Page     *int
	Name     *string
	Search   *string
	Archived *bool
	Whose    *string
}

// @admin

type AdminList struct {
	ID          int64   `json:"id" bun:"id"`
	Status      *string `json:"status" bun:"status"`
	TableID     *int64  `json:"table_id" bun:"table_id"`
	UserID      *int64  `json:"user_id" bun:"user_id"`
	Number      *int    `json:"table_number" bun:"number"`
	TableNumber *int    `json:"table_name" bun:"table_number"`
	TableStatus *string `json:"table_status" bun:"table_status"`
	UserName    *string `json:"user_name" bun:"user_name"`
	UserPhone   *string `json:"user_phone" bun:"user_phone"`
}

// @waiter

type WaiterGetListResponse struct {
	Id          int64        `json:"id"`
	Number      *int         `json:"number"`
	Status      *string      `json:"status"`
	TableNumber *int         `json:"table_number"`
	CreatedDate *string      `json:"created_date"`
	CreatedAt   *string      `json:"created_at"`
	WaiterId    *int64       `json:"waiter_id"`
	WaiterName  *string      `json:"waiter_name"`
	Accepted    *bool        `json:"accepted"`
	Menus       []WaiterMenu `json:"menus"`
}

type WaiterMenu struct {
	Id              int64   `json:"id" bun:"id"`
	OrderMenuId     int64   `json:"order_menu_id" bun:"order_menu_id"`
	OrderMenuStatus string  `json:"order_menu_status" bun:"order_menu_status"`
	Name            *string `json:"name" bun:"name"`
	Price           float64 `json:"price" bun:"price"`
	Count           int     `json:"count" bun:"count"`
}

type WaiterCreateRequest struct {
	Menus       []Menu `json:"menus" bun:"-" form:"menus"`
	TableID     int64  `json:"table_id" bun:"table_id" form:"table_id"`
	ClientCount int    `json:"client_count" bun:"client_count" form:"client_count"`
}

type WaiterCreateResponse struct {
	ID           int64   `json:"id" bun:"id" form:"-"`
	Price        float64 `json:"price" bun:"price" form:"-"`
	CreatedAt    string  `json:"created_at" bun:"created_at" form:"-"`
	ClientCount  int     `json:"client_count" bun:"client_count" form:"-"`
	BranchID     int64   `json:"-" bun:"-"`
	RestaurantID int64   `json:"-" bun:"-"`
	UserID       int64   `json:"-" bun:"-"`
}

type WaiterGetDetailResponse struct {
	Id          int64        `json:"id" bun:"id"`
	Number      *int         `json:"number" bun:"number"`
	ClientCount *int         `json:"client_count" bun:"client_count"`
	Status      *string      `json:"status" bun:"status"`
	Price       *float64     `json:"price" bun:"price"`
	WaiterId    *int64       `json:"waiter_id"`
	WaiterName  *string      `json:"waiter_name"`
	TableNumber *int         `json:"table_number" bun:"table_number"`
	CreatedDate *string      `json:"created_date"`
	CreatedAt   *string      `json:"created_at"`
	Menus       []WaiterMenu `json:"menus" bun:"menus"`
}

type WaiterUpdateRequest struct {
	Id    int64  `json:"id"`
	Menus []Menu `json:"menus"`
}

type WaiterAcceptOrderResponse struct {
	OrderNumber *int64  `json:"order_number" bun:"order_number"`
	TableID     *int64  `json:"table_id" bun:"table_id"`
	TableNumber *int64  `json:"table_number" bun:"table_number"`
	ClientID    *int64  `json:"client_id" bun:"client_id"`
	WaiterID    *int64  `json:"waiter_id" bun:"waiter_id"`
	WaiterName  *string `json:"waiter_name" bun:"waiter_name"`
	WaiterPhoto *string `json:"waiter_photo" bun:"waiter_photo"`
	AcceptedAt  *string `json:"accepted_at" bun:"accepted_at"`
}

// @client

//type MobileGetList struct {
//	ActiveOrder []struct {
//		ID         int64    `json:"id"`
//		Status     *string  `json:"status"`
//		CreatedAt  *string  `json:"created_at"`
//		BranchName *string  `json:"branch_name"`
//		Price      *float32 `json:"price"`
//		Address    *string  `json:"address"`
//	} `json:"active_order"`
//
//	InActiveOrder []struct {
//		ID         int64    `json:"id"`
//		Status     *string  `json:"status"`
//		CreatedAt  *string  `json:"created_at"`
//		BranchName *string  `json:"branch_name"`
//		Price      *float32 `json:"price"`
//		Address    *string  `json:"address"`
//	} `json:"in_active_order"`
//}

type MobileGetList struct {
	ID         int64    `json:"id"`
	Status     *string  `json:"status"`
	CreatedAt  *string  `json:"created_at"`
	BranchName *string  `json:"branch_name"`
	Price      *float32 `json:"price"`
	Address    *string  `json:"address"`
	BranchID   *int     `json:"branch_id"`
}

type MobileGetDetail struct {
	ID                int64      `json:"id"`
	Foods             []FoodList `json:"foods"`
	CreatedAt         *string    `json:"created_at"`
	BranchName        *string    `json:"branch_name"`
	Sum               *float32   `json:"sum"`
	Service           *float32   `json:"service"`
	OverAll           *float32   `json:"over_all"`
	ServicePercentage *int       `json:"service_percentage"`
	TableNumber       *int       `json:"table_number"`
	OrderNumber       *int       `json:"order_number"`
	Menus             []MenuList `json:"menus"`
}

type MobileCreateRequest struct {
	Menus       []Menu `json:"menus" bun:"-" form:"menus"`
	TableID     int64  `json:"table_id" bun:"table_id" form:"table_id"`
	ClientCount int    `json:"client_count" bun:"client_count" form:"client_count"`
}

type MobileCreateResponse struct {
	ID           int64          `json:"id" bun:"id" form:"-"`
	Price        float64        `json:"price" bun:"price" form:"-"`
	CreatedAt    string         `json:"created_at" bun:"created_at" form:"-"`
	ClientCount  int            `json:"client_count" bun:"client_count" form:"-"`
	Menus        []ResponseMenu `json:"menus" bun:"menus" form:"menus"`
	BranchID     int64          `json:"-" bun:"-"`
	RestaurantID int64          `json:"-" bun:"-"`
	UserID       int64          `json:"-" bun:"-"`
}

type ResponseMenu struct {
	ID    int64   `json:"id" bun:"id" form:"id"`
	Count int     `json:"count" bun:"count" form:"count"`
	Price float64 `json:"price" bun:"price" form:"price"`
	Name  string  `json:"name" bun:"name" form:"name"`
}

type MobileUpdateRequest struct {
	Id     int64   `json:"id"`
	Menus  []Menu  `json:"menus"`
	Status *string `json:"status"`
}

type ClientGetDetail struct {
	ID          int64    `json:"id" bun:"id"`
	OrderNumber *int     `json:"order_number"  bun:"order_number"`
	Status      *string  `json:"status"  bun:"status"`
	TableID     *int64   `json:"table_id"  bun:"table_id"`
	TableNumber *int     `json:"table_number"  bun:"table_number"`
	BranchID    *int64   `json:"-"  bun:"branch_id"`
	Price       *float32 `json:"price" bun:"price"`
	Accept      bool     `json:"-" bun:"accept"`
}

type ClientReviewRequest struct {
	Star        int     `json:"star"`
	Description *string `json:"description"`
	OrderId     int64   `json:"order_id"`
}

type ClientReviewResponse struct {
	bun.BaseModel `bun:"table:waiter_reviews"`

	ID          int64   `json:"id" bun:"id,pk,autoincrement"`
	Star        int     `json:"star" bun:"star"`
	Score       int     `json:"score" bun:"score"`
	Description *string `json:"description" bun:"description"`
	OrderId     int64   `json:"order_id" bun:"order_id"`

	CreatedAt time.Time `json:"created_at" bun:"created_at"`
	CreatedBy int64     `json:"created_by" bun:"created_by"`
}

// @cashier

type CashierUpdateStatusRequest struct {
	Id     int64   `json:"id"`
	Status *string `json:"status"`
}

type CashierGetList struct {
	ID          int64    `json:"id"`
	Number      *int     `json:"number"`
	Status      *string  `json:"status"`
	TableID     *int64   `json:"table_id"`
	TableNumber *int     `json:"table_number"`
	CreatedAt   *string  `json:"created_at"`
	Price       *float32 `json:"price"`
}

type CashierGetDetail struct {
	ID          int64                   `json:"id"`
	Number      *int                    `json:"number"`
	Status      *string                 `json:"status"`
	TableID     *int64                  `json:"table_id"`
	TableNumber *int                    `json:"table_number"`
	CreatedAt   *string                 `json:"created_at"`
	Price       *float32                `json:"price"`
	Menus       []MenuList              `json:"menus"`
	Waiter      waiter.CashierGetDetail `json:"waiter"`
	WaiterID    *int64                  `json:"-"`
}

// others--------------------------------------

type Menu struct {
	ID    int64 `json:"id" bun:"id" form:"id"`
	Count int   `json:"count" bun:"count" form:"count"`
}

type FoodList struct {
	ID     int64           `json:"id"`
	Name   *string         `json:"name"`
	Price  *float32        `json:"price"`
	Count  int             `json:"count"`
	Status *string         `json:"status"`
	Photos *pq.StringArray `json:"photos"`
}

type MenuList struct {
	ID     int64          `json:"id"`
	Name   *string        `json:"name"`
	Price  *float32       `json:"price"`
	Count  int            `json:"count"`
	Status *string        `json:"status"`
	Photos pq.StringArray `json:"photos"`
}

type Order struct {
	ID      int64 `json:"id" bun:"id" form:"id"`
	Number  int   `json:"number" bun:"number" form:"number"`
	TableID int64 `json:"table_id" bun:"table_id" form:"table_id"`
	UserID  int64 `json:"user_id" bun:"user_id" form:"user_id"`
}

type GetWsMessageResponse struct {
	UserId       int64    `json:"user_id"`
	RestaurantId int64    `json:"restaurant_id"`
	BranchId     int64    `json:"branch_id"`
	OrderID      int64    `json:"order_id"`
	OrderNumber  int64    `json:"order_number"`
	TableID      int64    `json:"table_id"`
	TableNumber  int64    `json:"table_number"`
	Price        *float32 `json:"price"`
	CreatedAt    *string  `json:"created_at"`
	WaiterID     *int64   `json:"waiter_id"`
}

type GetWsOrderMenusResponse struct {
	OrderNumber *int64   `json:"order_number"`
	Ip          *string  `json:"ip"`
	Waiter      *string  `json:"waiter"`
	TableNumber *int64   `json:"table_number"`
	Foods       []WsFood `json:"foods"`
}

type WsFood struct {
	Name  *string `json:"name"`
	Count *int    `json:"count"`
}

type GetWsWaiterResponse struct {
	ID         int64   `json:"id"`
	OrderCount int     `json:"order_count"`
	Grade      float32 `json:"grade"`
}
