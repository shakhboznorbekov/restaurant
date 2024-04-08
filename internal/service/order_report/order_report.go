package order_report

import "golang.org/x/net/context"

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo}
}

func (s *Service) CashierOrderReport(ctx context.Context) error {
	return s.repo.CashierOrderReport(ctx)
}
