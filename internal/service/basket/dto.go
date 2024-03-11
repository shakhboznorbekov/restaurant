package basket

import (
	"github.com/restaurant/internal/service/menu"
	"time"
)

//type CreateRequest struct {
//	Value      OrderStore `json:"value" form:"value"`
//	Expiration *int       `json:"expiration" form:"expiration"`
//	BranchID   *int64     `json:"branch_id" form:"branch_id"`
//}

type Create struct {
	Key        *string        `json:"key"`
	Value      OrderStore     `json:"value"`
	Expiration *time.Duration `json:"expiration"`
	BranchID   *int64         `json:"branch_id"`
}

type Update struct {
	Food     Food   `json:"food" form:"food"`
	BranchID *int64 `json:"branch_id"`
	TableID  *int64 `json:"table_id" form:"table_id"`
	UserID   *int64 `json:"-" form:"-"`
}

type Detail struct {
	Foods             []menu.ClientGetDetail `json:"foods"`
	TableID           *int64                 `json:"table_id" form:"table_id"`
	UserID            *int64                 `json:"user_id" form:"user_id"`
	Sum               *float32               `json:"sum"`
	Service           *float32               `json:"service"`
	OverAll           *float32               `json:"over_all"`
	ServicePercentage *int                   `json:"service_percentage"`
}

type OrderStore struct {
	Foods   []Food `json:"foods" form:"foods"`
	TableID *int64 `json:"table_id" form:"table_id"`
	UserID  *int64 `json:"user_id" form:"user_id"`
}

type Food struct {
	ID    int64 `json:"id" form:"id"`
	Count *int  `json:"count" form:"count"`
}
