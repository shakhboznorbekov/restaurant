package basket

import "context"

type Service struct {
	repo Repository
}

// @client

func (s Service) SetBasket(ctx context.Context, data Create) error {
	return s.repo.SetBasket(ctx, data)
}

func (s Service) GetBasket(ctx context.Context, key string) (OrderStore, error) {
	return s.repo.GetBasket(ctx, key)
}

func (s Service) UpdateBasket(ctx context.Context, key string, value Update) error {
	return s.repo.UpdateBasket(ctx, key, value)
}

func (s Service) DeleteBasket(ctx context.Context, key string) error {
	return s.repo.DeleteBasket(ctx, key)
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}
