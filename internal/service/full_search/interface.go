package full_search

import (
	"context"
)

type Repository interface {
	ClientGetList(ctx context.Context, filter Filter) ([]ClientGetList, error)
}
