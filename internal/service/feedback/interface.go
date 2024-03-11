package feedback

import (
	"context"
	"github.com/restaurant/internal/entity"
)

type Repository interface {

	// @admin

	AdminGetList(ctx context.Context, filter Filter) ([]AdminGetList, int, error)
	AdminGetDetail(ctx context.Context, id int64) (entity.Feedback, error)
	AdminCreate(ctx context.Context, request AdminCreate) (entity.Feedback, error)
	AdminUpdateColumns(ctx context.Context, request AdminUpdate) error
	AdminDelete(ctx context.Context, id int64) error

	// @client

	ClientGetList(ctx context.Context, filter Filter) ([]ClientGetList, int, error)
}
