package order_report

import "golang.org/x/net/context"

type Repository interface {
	CashierOrderReport(ctx context.Context) error
}
