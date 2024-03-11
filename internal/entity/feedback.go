package entity

import "github.com/uptrace/bun"

type Feedback struct {
	bun.BaseModel `bun:"table:feedback"`

	ID   int64             `json:"id" bun:"id,pk,autoincrement"`
	Name map[string]string `json:"name" bun:"name"`
}

//for ----->  biror restaran, branch, yoki taomga taqriz berish
