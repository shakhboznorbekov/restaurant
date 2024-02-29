package user

import (
	"github.com/postgresql-restaurant/internal/pkg/utils"
	"github.com/uptrace/bun"
	"time"
)

type Filter struct {
	Limit        *int
	Offset       *int
	Name         *string
	Role         *string
	RestaurantID *int64
	BranchID     *int64
	Fields       map[string][]string
	Joins        map[string]utils.Joins
}

// @super-admin
type SuperAdminGetList struct {
	ID    int64   `json:"id"`
	Name  *string `json:"name"`
	Phone *string `json:"phone"`
	Role  *string `json:"role"`
}

type SuperAdminGetDetail struct {
	ID        int64   `json:"id"`
	Name      *string `json:"name"`
	Phone     *string `json:"phone"`
	BirthDate *string `json:"birth_date"`
	Gender    *string `json:"gender"`
	Role      *string `json:"role"`
}

type SuperAdminCreateRequest struct {
	Name      *string `json:"name" form:"name"`
	Phone     *string `json:"phone" form:"phone"`
	BirthDate *string `json:"birth_date" form:"birth_date"`
	Gender    *string `json:"gender" form:"gender"`
	Role      *string `json:"-" form:"-"`
}

type SuperAdminCreateResponse struct {
	bun.BaseModel `bun:"table:users"`

	ID        int64      `json:"id" bun:"id,pk,autoincrement"`
	Name      *string    `json:"name" bun:"name"`
	Phone     *string    `json:"phone" bun:"phone"`
	BirthDate *time.Time `json:"birth_date" bun:"birth_date"`
	Gender    *string    `json:"gender" bun:"gender"`
	Role      *string    `json:"-" bun:"role"`
	CreatedAt time.Time  `json:"-" bun:"created_at"`
	CreatedBy int64      `json:"-" bun:"created_by"`
}

type SuperAdminUpdateRequest struct {
	ID        int64   `json:"id" form:"id"`
	Name      *string `json:"name" form:"name"`
	Phone     *string `json:"phone" form:"phone"`
	BirthDate *string `json:"birth_date" form:"birth_date"`
	Gender    *string `json:"gender" form:"gender"`
	Role      *string `json:"-" form:"-"`
}

// @client

type ClientCreateRequest struct {
	Phone *string `json:"phone" form:"phone"`
	Role  *string `json:"-" form:"-"`
}

type ClientCreateResponse struct {
	bun.BaseModel `bun:"table:users"`

	ID        int64     `json:"id" bun:"id,pk,autoincrement"`
	Phone     *string   `json:"phone" bun:"phone"`
	Role      *string   `json:"role" bun:"role"`
	CreatedAt time.Time `json:"-" bun:"created_at"`
}

type ClientUpdateRequest struct {
	ID        int64   `json:"id" form:"id"`
	Name      *string `json:"name" form:"name"`
	BirthDate *string `json:"birth_date" form:"birth_date"`
	Gender    *string `json:"gender" form:"gender"`
}

type ClientDetail struct {
	ID        int64   `json:"id"`
	Name      *string `json:"name"`
	Phone     *string `json:"phone"`
	BirthDate *string `json:"birth_date"`
	Gender    *string `json:"gender"`
	Role      *string `json:"role"`
}

// @admin

type AdminCreateRequest struct {
	Name         *string `json:"name" form:"name"`
	Phone        *string `json:"phone" form:"phone"`
	BirthDate    *string `json:"birth_date" form:"birth_date"`
	Gender       *string `json:"gender" form:"gender"`
	Role         *string `json:"-" form:"-"`
	RestaurantID *int64  `json:"-" form:"-"`
	CreatedBy    int64   `json:"-" form:"-"`
}

type AdminCreateResponse struct {
	bun.BaseModel `bun:"table:users"`

	ID           int64      `json:"id" bun:"id,pk,autoincrement"`
	Name         *string    `json:"name" bun:"name"`
	Phone        *string    `json:"phone" bun:"phone"`
	BirthDate    *time.Time `json:"birth_date" bun:"birth_date"`
	Gender       *string    `json:"gender" bun:"gender"`
	Role         *string    `json:"role" bun:"role"`
	RestaurantID *int64     `json:"-" bun:"restaurant_id"`
	CreatedAt    time.Time  `json:"-" bun:"created_at"`
	CreatedBy    int64      `json:"-" bun:"created_by"`
}

type AdminUpdateRequest struct {
	ID           int64   `json:"id" form:"id"`
	Name         *string `json:"name" form:"name"`
	Phone        *string `json:"phone" form:"phone"`
	BirthDate    *string `json:"birth_date" form:"birth_date"`
	Gender       *string `json:"gender" form:"gender"`
	Role         *string `json:"-" form:"-"`
	RestaurantID *int64  `json:"-" form:"-"`
	UpdatedBy    int64   `json:"-" form:"-"`
}

type AdminGetList struct {
	ID        int64   `json:"id" bun:"id"`
	Name      *string `json:"name" bun:"name"`
	Phone     *string `json:"phone" bun:"phone"`
	Role      *string `json:"role" bun:"role"`
	BirthDate *string `json:"birth_date" bun:"birth_date"`
	Gender    *string `json:"gender" bun:"gender"`
}

type AdminUpdateByRestaurantIDRequest struct {
	RestaurantID int64   `json:"-" form:"-"`
	Name         *string `json:"name" form:"name"`
	Phone        *string `json:"phone" form:"phone"`
	BirthDate    *string `json:"birth_date" form:"birth_date"`
	Gender       *string `json:"gender" form:"gender"`
}

// @branch

type BranchCreateRequest struct {
	Name      *string `json:"name" form:"name"`
	Phone     *string `json:"phone" form:"phone"`
	BirthDate *string `json:"birth_date" form:"birth_date"`
	Gender    *string `json:"gender" form:"gender"`
	Role      *string `json:"-" form:"-"`
	BranchID  *int64  `json:"-" form:"-"`
	CreatedBy int64   `json:"-" form:"-"`
}

type BranchCreateResponse struct {
	bun.BaseModel `bun:"table:users"`

	ID        int64      `json:"id" bun:"id,pk,autoincrement"`
	Name      *string    `json:"name" bun:"name"`
	Phone     *string    `json:"phone" bun:"phone"`
	BirthDate *time.Time `json:"birth_date" bun:"birth_date"`
	Gender    *string    `json:"gender" bun:"gender"`
	Role      *string    `json:"role" bun:"role"`
	BranchID  *int64     `json:"-" bun:"branch_id"`
	CreatedAt time.Time  `json:"-" bun:"created_at"`
	CreatedBy int64      `json:"-" bun:"created_by"`
}

type UpdatePasswordRequest struct {
	Password string `json:"password"`
	SMSCode  string `json:"sms_code"`
}

// --Feedback

type FeedBack struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	FullName    *string `json:"full_name"`
	PhoneNumber *string `json:"phone_number"`
	Email       *string `json:"email"`
}
