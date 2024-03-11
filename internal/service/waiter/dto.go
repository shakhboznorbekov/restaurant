package waiter

import (
	"github.com/uptrace/bun"
	"mime/multipart"
	"time"
)

type Filter struct {
	Limit    *int
	Offset   *int
	BranchID *int64
}

// @admin

type AdminGetList struct {
	ID         int64   `json:"id"`
	Name       *string `json:"name"`
	Phone      *string `json:"phone"`
	Role       *string `json:"role"`
	BranchName *string `json:"branch_name"`
}

// @branch

type BranchGetList struct {
	ID         int64   `json:"id"`
	Name       *string `json:"name"`
	Phone      *string `json:"phone"`
	Role       *string `json:"role"`
	BranchName *string `json:"branch_name"`
	Status     *string `json:"status"`
	Photo      *string `json:"photo"`
	Address    *string `json:"address"`
}

type BranchGetDetail struct {
	ID        int64    `json:"id"`
	Name      *string  `json:"name"`
	Phone     *string  `json:"phone"`
	BirthDate *string  `json:"birth_date"`
	Gender    *string  `json:"gender"`
	Role      *string  `json:"role"`
	Photo     *string  `json:"photo"`
	Rating    *float64 `json:"rating"`
	Address   *string  `json:"address"`
}

type BranchCreateRequest struct {
	Name                *string               `json:"name" form:"name"`
	Password            *string               `json:"password" form:"password"`
	Phone               *string               `json:"phone" form:"phone"`
	BirthDate           *string               `json:"birth_date" form:"birth_date"`
	Gender              *string               `json:"gender" form:"gender"`
	ServicePercentageID *int64                `json:"service_percent_id" form:"service_percent_id"`
	Photo               *multipart.FileHeader `json:"photo" form:"photo"`
	PhotoLink           *string               `json:"-" form:"-"`
	Address             *string               `json:"address" form:"address"`
}

type BranchCreateResponse struct {
	bun.BaseModel `bun:"table:users"`

	ID             int64      `json:"id" bun:"id,pk,autoincrement"`
	Name           *string    `json:"name" bun:"name"`
	Password       *string    `json:"-" bun:"password"`
	Phone          *string    `json:"phone" bun:"phone"`
	BirthDate      *time.Time `json:"birth_date" bun:"birth_date"`
	Gender         *string    `json:"gender" bun:"gender"`
	Role           *string    `json:"role" bun:"role"`
	BranchID       *int64     `json:"branch_id" bun:"branch_id"`
	ServicePercent *int64     `json:"service_percent" bun:"service_percent"`
	Photo          *string    `json:"photo" bun:"photo"`
	Address        *string    `json:"address" bun:"address"`

	CreatedAt time.Time `json:"created_at" bun:"created_at"`
	CreatedBy int64     `json:"created_by" bun:"created_by"`
}

type BranchUpdateRequest struct {
	ID                  int64                 `json:"id" form:"id"`
	Name                *string               `json:"name" form:"name"`
	BirthDate           *string               `json:"birth_date" form:"birth_date"`
	Gender              *string               `json:"gender" form:"gender"`
	ServicePercentageID *int64                `json:"rating" form:"rating"`
	Photo               *multipart.FileHeader `json:"photo" form:"photo"`
	PhotoLink           *string               `json:"-" form:"-"`
	Address             *string               `json:"address" form:"address"`
}

type BranchUpdatePhone struct {
	ID      int64  `json:"id" form:"id"`
	Phone   string `json:"phone" form:"phone"`
	SMSCode string `json:"sms_code" form:"sms_code"`
}

type BranchUpdatePassword struct {
	ID       int64   `json:"id" form:"id"`
	Password *string `json:"password" form:"password"`
}

type SendSms struct {
	Phone string `json:"phone" form:"phone"`
}

// @waiter

type GetMeResponse struct {
	Id               int64    `json:"id" bun:"id"`
	Name             *string  `json:"name" bun:"name"`
	Photo            *string  `json:"photo" bun:"photo"`
	Profit           *float64 `json:"profit" bun:"profit"`
	Rating           *float32 `json:"rating" bun:"rating"`
	OrderCount       *int     `json:"order_count" bun:"order_count"`
	BirthDate        *string  `json:"birth_date" bun:"birth_date"`
	Phone            *string  `json:"phone" bun:"phone"`
	Address          *string  `json:"address" bun:"address"`
	AttendanceStatus *bool    `json:"attendance_status" bun:"attendance_status"`
}

type GetPersonalInfoResponse struct {
	Id        int64   `json:"id" bun:"id"`
	Name      *string `json:"name" bun:"name"`
	BirthDate *string `json:"birth_date" bun:"birth_date"`
	Phone     *string `json:"phone" bun:"phone"`
	Address   *string `json:"address" bun:"address"`
}

type WaiterPhotoUpdateRequest struct {
	Photo     *multipart.FileHeader `json:"photo" form:"photo"`
	PhotoLink *string               `json:"-" form:"-"`
}

// @cashier

type CashierGetList struct {
	ID         int64   `json:"id"`
	Name       *string `json:"name"`
	Phone      *string `json:"phone"`
	Role       *string `json:"role"`
	BranchName *string `json:"branch_name"`
}

type CashierGetDetail struct {
	ID     int64   `json:"id"`
	Name   *string `json:"name"`
	Avatar *string `json:"avatar"`
}
