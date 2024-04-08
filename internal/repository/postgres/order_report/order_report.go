package order_report

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"restu-backend/foundation/web"
	"restu-backend/internal/auth"
	"restu-backend/internal/pkg/repository/postgresql"
)

type Repository struct {
	*postgresql.Database
}

func New(db *postgresql.Database) *Repository {
	return &Repository{db}
}

func (r *Repository) CashierOrderReport(ctx context.Context) error {
	claims, err := r.CheckClaims(ctx, auth.RoleCashier)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "extracting claims"), http.StatusInternalServerError)
	}
	if claims.BranchID == nil {
		err = errors.New("branch_id must not be nil in cashier")
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	where := fmt.Sprintf(`WHERE 
										o.status = 'PAID' 
									and 
										o.deleted_at isnull 
									and 
										om.status='PAID' 
									and 
										om.deleted_at isnull 
									and 
										t.branch_id = '%d' 
									and 
										case when frgh.from isnull then true else frgh.from < o.created_at end
									and 
										case when frgh.to isnull then true else frgh.to > o.created_at end`, *claims.BranchID)

	// TODO: get 'from_time' if branches.last_order_reported_at isnull from first order.created_at
	// from filter to where query
	fromQuery := fmt.Sprintf(`select last_order_reported_at from branches where id='%d'`, *claims.BranchID)
	var from *string
	if err = r.QueryRowContext(ctx, fromQuery).Scan(&from); err != nil {
		return err
	}
	if from != nil {
		where += fmt.Sprintf(` and o.created_at>='%s'`, *from)
	}

	selectQuery := fmt.Sprintf(`select 
    										p.id, 
    										sum(fr.amount * om.count)
										from orders o 
										    join order_menu om 
										        on o.id = om.order_id 
										    join menus m 
										        on m.id=om.menu_id 
										    join foods f 
										        on f.id = any (m.food_ids) 
										    join food_recipe_group_histories frgh 
										        on frgh.food_id = f.id 
										    join food_recipe_groups frg 
										        on frg.id = frgh.group_id
										    join food_recipe fr 
										        on fr.id = any (frg.recipe_ids)
										    join products p 
										        on p.id = fr.product_id 
										    join tables t 
										        on t.id = o.table_id %s 
										group by p.id`, where)

	rows, err := r.QueryContext(ctx, selectQuery)
	if err != nil {
		return err
	}

	// TODO: try to merge these two queries into one
	for rows.Next() {
		var (
			query     string
			productID int64
			amount    float64
		)
		if err = rows.Scan(&productID, &amount); err != nil {
			return err
		}

		if from != nil {
			query = fmt.Sprintf(`INSERT INTO order_report (from_time, product_id, amount, branch_id) VALUES('%s', '%d', '%f', '%d')`, *from, productID, amount, *claims.BranchID)
		} else {
			query = fmt.Sprintf(`INSERT INTO order_report (product_id, amount, branch_id) VALUES('%d', '%f', '%d')`, productID, amount, *claims.BranchID)
		}

		if _, err = r.ExecContext(ctx, query); err != nil {
			return err
		}
	}

	// TODO: do not skip returned.error; I did not create transaction for the performance sake
	_, _ = r.NewUpdate().Table("branches").Set("last_order_reported_at=now()").Where("id=?", *claims.BranchID).Exec(ctx)

	return nil
}
