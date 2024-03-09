package entity

import (
	"github.com/uptrace/bun"
	"time"
)

type PartnerContract struct {
	bun.BaseModel `bun:"table:partner_contract"`

	ID             int64   `json:"id" bun:"id,pk,autoincrement"`
	Name           *string `json:"name" bun:"name"`
	ContactDate    *string `json:"contact_date" bun:"contact_date"`
	ContractNumber *int    `json:"contract_number" bun:"contact_number"`
	PaymentType    *string `json:"payment_type" bun:"payment_type"`

	PartnerID *int64     `json:"partner_id" bun:"partner_id"`
	CreatedAt *time.Time `json:"created_at" bun:"created_at"`
	CreatedBy *int64     `json:"created_by" bun:"created_by"`
	UpdatedAt *time.Time `json:"updated_at" bun:"updated_at"`
	UpdatedBy *int64     `json:"updated_by" bun:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" bun:"deleted_at"`
	DeletedBy *int64     `json:"deleted_by" bun:"deleted_by"`
}
