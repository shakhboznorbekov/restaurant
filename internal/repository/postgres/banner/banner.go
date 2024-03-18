package banner

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/restaurant/foundation/web"
	"github.com/restaurant/internal/auth"
	"github.com/restaurant/internal/commands"
	"github.com/restaurant/internal/pkg/repository/postgresql"
	"github.com/restaurant/internal/repository/postgres"
	"github.com/restaurant/internal/service/banner"
	"github.com/restaurant/internal/service/hashing"
	"net/http"
	"strings"
	"time"
)

type Repository struct {
	*postgresql.Database
}

// admin

func (r Repository) BranchGetList(ctx context.Context, filter banner.Filter) ([]banner.BranchGetList, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return nil, 0, err
	}

	lang := r.GetLang(ctx)

	whereQuery := fmt.Sprintf(`WHERE deleted_at IS NULL AND branch_id = '%d'`, *claims.BranchID)
	if filter.Expired != nil {
		if *filter.Expired {
			whereQuery += fmt.Sprintf(` AND expired_at < now()`)
		} else {
			whereQuery += fmt.Sprintf(` AND expired_at > now()`)
		}
	}

	if filter.Status != nil {
		whereQuery += fmt.Sprintf(` AND status='%s'`, *filter.Status)
	}
	countWhereQuery := whereQuery

	var limitQuery, offsetQuery string
	if filter.Limit != nil {
		limitQuery = fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	query := fmt.Sprintf(`
		SELECT
		    id,
		    photo,
		    title->>'%s',
		    description->>'%s',
		    price,
		    old_price,
		    TO_CHAR(expired_at,'DD.MM.YYYY | HH24:MI'),
		    status
		FROM banners
		%s %s %s
`, lang, lang, whereQuery, limitQuery, offsetQuery)
	if err != nil {
		return nil, 0, errors.Wrap(err, "select query")
	}

	list := make([]banner.BranchGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select stories"), http.StatusInternalServerError)
	}

	for rows.Next() {
		var detail banner.BranchGetList
		err = rows.Scan(&detail.ID, &detail.Photo, &detail.Title, &detail.Description, &detail.Price, &detail.OldPrice, &detail.ExpiredAt, &detail.Status)
		if err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning stories"), http.StatusBadRequest)
		}

		if detail.Photo != nil {
			link := hashing.GenerateHash(r.ServerBaseUrl, *detail.Photo)
			detail.Photo = &link
		}

		if detail.OldPrice != nil && detail.Price != nil {
			discount := 100 - int((*detail.Price)/(*detail.OldPrice)*100)

			detail.Discount = &discount
		}

		expired := true
		if detail.ExpiredAt != nil {
			expiredAt, err := time.Parse("02.01.2006 | 15:04", *detail.ExpiredAt)
			if err != nil {
				return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning stories"), http.StatusBadRequest)
			}
			if !expiredAt.Before(time.Now().Add(5 * time.Hour)) {
				expired = false
			}
		}
		detail.Expired = &expired
		list = append(list, detail)
	}
	countQuery := fmt.Sprintf(`
		SELECT
			count(id)
		FROM banners
		%s
	`, countWhereQuery)

	countRows, err := r.QueryContext(ctx, countQuery)
	if err == sql.ErrNoRows {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting stories"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning user count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func (r Repository) BranchGetDetail(ctx context.Context, id int64) (*banner.BranchGetDetail, error) {
	_, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return nil, err
	}

	lang := r.GetLang(ctx)

	var (
		response banner.BranchGetDetail
		menus    []banner.Menu
	)

	var title, description []byte

	query := fmt.Sprintf(`SELECT 
    								b.id,
									b.photo,
									b.title,
									b.description,
									b.price,
									b.old_price,
									b.menu_ids,
									TO_CHAR(b.expired_at,'DD.MM.YYYY | HH24:MI'),
									b.status
								FROM banners AS b
									LEFT JOIN branches br ON b.branch_id = br.id
									LEFT JOIN restaurants AS r ON br.restaurant_id = r.id
								WHERE b.deleted_at ISNULL AND b.id='%d'`, id)
	if err = r.QueryRowContext(ctx, query).Scan(&response.ID, &response.Photo, &title, &description, &response.Price, &response.OldPrice, &response.MenuIds, &response.ExpiredAt, &response.Status); err != nil {
		return nil, err
	}

	if response.OldPrice != nil && response.Price != nil {
		discount := 100 - int((*response.Price)/(*response.OldPrice)*100)

		response.Discount = &discount
	}

	if response.Photo != nil {
		link := hashing.GenerateHash(r.ServerBaseUrl, *response.Photo)
		response.Photo = &link
	}

	if title != nil {
		err = json.Unmarshal(title, &response.Title)
		if err != nil {
			return nil, err
		}
	}

	if description != nil {
		err = json.Unmarshal(description, &response.Description)
		if err != nil {
			return nil, err
		}
	}

	if response.MenuIds != nil && len(*response.MenuIds) > 0 {
		for _, v := range *response.MenuIds {
			var menu banner.Menu
			queryMenu := fmt.Sprintf(`SELECT m.id, 
													f.name, 
													m.description->>'%s',
													f.photos[1],
													m.new_price
									 	 	 FROM menus m 
									 	 	     JOIN foods f 
									 	 	         ON m.food_id = f.id 
									 	 	 WHERE m.id='%d'`, lang, v)
			if err = r.QueryRowContext(ctx, queryMenu).Scan(&menu.Id, &menu.Title, &menu.Description, &menu.Photo, &menu.Price); err != nil && !errors.Is(err, sql.ErrNoRows) {
				return nil, err
			}

			if menu.Photo != nil {
				link := hashing.GenerateHash(r.ServerBaseUrl, *menu.Photo)
				menu.Photo = &link
			}

			menus = append(menus, menu)
		}

		response.Menus = menus
	}

	return &response, nil
}

func (r Repository) BranchCreate(ctx context.Context, request banner.BranchCreateRequest) (*banner.BranchCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return nil, err
	}

	err = r.ValidateStruct(&request, "Title", "Description", "Price", "PhotoLink", "MenuIds", "ExpiredAt")
	if err != nil {
		return nil, err
	}

	expiredAt, err := time.Parse("02.01.2006 | 15:04", *request.ExpiredAt)
	if err != nil {
		return nil, web.NewRequestError(errors.Wrap(err, "parse expired_at"), http.StatusBadRequest)
	}

	response := banner.BranchCreateResponse{
		Title:       request.Title,
		Description: request.Description,
		Photo:       request.PhotoLink,
		Price:       request.Price,
		ExpiredAt:   expiredAt,
		CreatedAt:   time.Now(),
		CreatedBy:   claims.UserId,
		BranchID:    *claims.BranchID,
	}

	if request.MenuIds != nil && len(*request.MenuIds) > 0 {
		var oldPrice *float64

		// making in query array
		in := "("
		for i, v := range *request.MenuIds {
			if i == 0 {
				in += fmt.Sprintf("%d", v)
			} else {
				in += fmt.Sprintf(",%d", v)
			}
		}
		in += ")"

		priceQuery := fmt.Sprintf(`SELECT sum(m.new_price) FROM menus m WHERE m.deleted_at is null and m.id in %s`, in)
		if err = r.QueryRowContext(ctx, priceQuery).Scan(&oldPrice); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		if oldPrice != nil {
			response.OldPrice = oldPrice
		}
		response.MenuIds = request.MenuIds
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return nil, web.NewRequestError(errors.Wrap(err, "creating banner"), http.StatusBadRequest)
	}

	if response.Photo != nil {
		link := hashing.GenerateHash(r.ServerBaseUrl, *response.Photo)
		response.Photo = &link
	}

	return &response, nil
}

func (r Repository) BranchUpdateAll(ctx context.Context, request banner.BranchUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID", "Title", "Description", "Price", "PhotoLink"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("banners").Where("deleted_at IS NULL AND id = ? AND branch_id = ?",
		request.ID, claims.BranchID)

	if request.MenuIds != nil && len(*request.MenuIds) > 0 {
		var oldPrice *float64

		strSlice := make([]string, len(*request.MenuIds))
		for i, v := range *request.MenuIds {
			strSlice[i] = fmt.Sprintf("%d", v)
		}

		query := fmt.Sprintf(`SELECT sum(new_price) FROM menus WHERE deleted_at is null AND status='active' AND id IN (%s)`, strings.Join(strSlice, ","))
		if err = r.QueryRowContext(ctx, query).Scan(&oldPrice); err != nil {
			return err
		}
		q.Set("old_price = ?", oldPrice)
		q.Set("menu_ids = ?", request.MenuIds)
	} else {
		q.Set("old_price = 0")
		q.Set("menu_ids='{}'")
	}
	q.Set("title = ?", request.Title)
	q.Set("description = ?", request.Description)
	q.Set("photo = ?", request.PhotoLink)
	q.Set("price = ?", request.Price)
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)
	q.Set("status='DRAFT'")

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating banner"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) BranchUpdateColumn(ctx context.Context, request banner.BranchUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("banners").Where("deleted_at IS NULL AND id = ? AND branch_id = ?",
		request.ID, claims.BranchID)

	if request.Photo != nil {
		q.Set("photo = ?", request.PhotoLink)
	}
	if request.Title != nil {
		q.Set("title = ?", request.Title)
	}
	if request.Description != nil {
		q.Set("description = ?", request.Description)
	}
	if request.Price != nil {
		q.Set("price = ?", request.Price)
	}
	if request.MenuIds != nil && len(*request.MenuIds) > 0 {
		var oldPrice *float64

		strSlice := make([]string, len(*request.MenuIds))
		for i, v := range *request.MenuIds {
			strSlice[i] = fmt.Sprintf("%d", v)
		}

		query := fmt.Sprintf(`SELECT sum(new_price) FROM menus WHERE deleted_at is null AND status='active' AND id IN (%s)`, strings.Join(strSlice, ","))
		if err = r.QueryRowContext(ctx, query).Scan(&oldPrice); err != nil {
			return err
		}
		q.Set("old_price = ?", oldPrice)
		q.Set("menu_ids = ?", request.MenuIds)
	}
	q.Set("status='DRAFT'")
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating banner"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) BranchUpdateStatus(ctx context.Context, id int64, expireAt string) error {
	claims, err := r.CheckClaims(ctx, auth.RoleBranch)
	if err != nil {
		return err
	}

	expiredAt, err := time.Parse("02.01.2006 | 15:04", expireAt)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "parse expired_at"), http.StatusBadRequest)
	}

	if !expiredAt.After(time.Now()) {
		err = errors.New("expire_at is comes before time.now()")
		return err
	}

	_, err = r.ExecContext(ctx, fmt.Sprintf("UPDATE banners SET expired_at='%v', status='DRAFT' WHERE id = '%d' AND branch_id = '%d'", expiredAt.Format("2006-01-02 15:04:05"), id, *claims.BranchID))
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "scan active banners"), http.StatusInternalServerError)
	}

	return nil
}

func (r Repository) BranchDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "banners", id, auth.RoleBranch)
}

// @client

func (r Repository) ClientGetList(ctx context.Context, filter banner.Filter) ([]banner.ClientGetList, int, error) {
	//_, err := r.CheckClaims(ctx, auth.RoleClient)
	//if err != nil {
	//	return nil, 0, err
	//}

	lang := r.GetLang(ctx)

	whereQuery := fmt.Sprintf(`WHERE b.deleted_at IS NULL AND b.expired_at > now() AND b.status='APPROVED'`)

	var limitQuery, offsetQuery string
	if filter.Limit != nil {
		limitQuery = fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
		if filter.Page != nil {
			offset := (*filter.Page - 1) * (*filter.Limit)
			offsetQuery = fmt.Sprintf(" OFFSET '%d'", offset)
		}
	}

	var (
		location = true
		lat, lon float64
	)
	if filter.Lat == nil || filter.Lon == nil {
		location = false
	} else {
		lat = *filter.Lat
		lon = *filter.Lon
	}

	query := fmt.Sprintf(`
		SELECT
		    b.id,
		    b.photo,
		    b.title->>'%s',
		    b.description->>'%s',
		    b.price,
		    b.old_price,
		    br.restaurant_id,
		    r.name,
		    r.logo,
		    CASE WHEN '%t' THEN ST_DistanceSphere(
				            ST_SetSRID(ST_MakePoint(CAST(br.location->>'lon' AS float), CAST(br.location->>'lat' AS float)), 4326),
				            ST_SetSRID(ST_MakePoint('%v', '%v'), 4326)
				        ) END as distance
		FROM banners AS b
		    LEFT JOIN branches br ON b.branch_id = br.id
			LEFT JOIN restaurants AS r ON br.restaurant_id = r.id
		%s %s %s
`, lang, lang, location, lat, lon, whereQuery, limitQuery, offsetQuery)

	list := make([]banner.ClientGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select stories"), http.StatusInternalServerError)
	}

	for rows.Next() {
		var (
			detail banner.ClientGetList
		)
		err = rows.Scan(&detail.Id, &detail.Photo, &detail.Title, &detail.Description, &detail.Price, &detail.OldPrice, &detail.RestaurantID, &detail.Restaurant, &detail.RestaurantPhoto, &detail.Distance)
		if err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning stories"), http.StatusBadRequest)
		}

		if detail.OldPrice != nil && detail.Price != nil {
			discount := 100 - int((*detail.Price)/(*detail.OldPrice)*100)

			detail.Discount = &discount
		}

		if detail.Photo != nil {
			link := hashing.GenerateHash(r.ServerBaseUrl, *detail.Photo)
			detail.Photo = &link
		}

		if detail.RestaurantPhoto != nil {
			link := hashing.GenerateHash(r.ServerBaseUrl, *detail.RestaurantPhoto)
			detail.RestaurantPhoto = &link
		}
		list = append(list, detail)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(*)
		FROM banners AS b
			LEFT JOIN branches br ON b.branch_id = br.id
			LEFT JOIN restaurants AS r ON br.restaurant_id = r.id
		%s
	`, whereQuery)

	countRows, err := r.QueryContext(ctx, countQuery)
	if err == sql.ErrNoRows {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting stories"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning user count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func (r Repository) ClientGetDetail(ctx context.Context, id int64) (*banner.ClientGetDetail, error) {
	//_, err := r.CheckClaims(ctx, auth.RoleClient)
	//if err != nil {
	//	return nil, err
	//}

	lang := r.GetLang(ctx)

	var (
		response banner.ClientGetDetail
		menus    []banner.Menu
	)

	query := fmt.Sprintf(`SELECT
									b.photo,
									b.title->>'%s',
									b.description->>'%s',
									b.price,
									b.old_price,
									b.menu_ids,
									r.id,
									r.logo,
									r.name
								FROM banners AS b
									LEFT JOIN branches br ON b.branch_id = br.id
									LEFT JOIN restaurants AS r ON br.restaurant_id = r.id
								WHERE b.deleted_at ISNULL AND b.expired_at > now() AND b.id='%d' AND b.status='APPROVED'`, lang, lang, id)
	if err := r.QueryRowContext(ctx, query).Scan(&response.Photo, &response.Title, &response.Description, &response.Price, &response.OldPrice, &response.MenuIds, &response.RestaurantID, &response.RestaurantPhoto, &response.Restaurant); err != nil {
		return nil, err
	}

	if response.OldPrice != nil && response.Price != nil {
		discount := 100 - int((*response.Price)/(*response.OldPrice)*100)

		response.Discount = &discount
	}

	if response.Photo != nil {
		link := hashing.GenerateHash(r.ServerBaseUrl, *response.Photo)
		response.Photo = &link
	}

	if response.RestaurantPhoto != nil {
		link := hashing.GenerateHash(r.ServerBaseUrl, *response.RestaurantPhoto)
		response.RestaurantPhoto = &link
	}

	if response.MenuIds != nil && len(*response.MenuIds) > 0 {
		for _, v := range *response.MenuIds {
			var menu banner.Menu
			queryMenu := fmt.Sprintf(`SELECT m.id, 
													f.name, 
													m.description->>'%s',
													f.photos[1]
									 	 	 FROM menus m 
									 	 	     JOIN foods f 
									 	 	         ON m.food_id = f.id 
									 	 	 WHERE m.id='%d'`, lang, v)
			if err := r.QueryRowContext(ctx, queryMenu).Scan(&menu.Id, &menu.Title, &menu.Description, &menu.Photo); err != nil && !errors.Is(err, sql.ErrNoRows) {
				return nil, err
			}

			if menu.Photo != nil {
				link := hashing.GenerateHash(r.ServerBaseUrl, *menu.Photo)
				menu.Photo = &link
			}

			menus = append(menus, menu)
		}

		response.Menus = menus
	}

	return &response, nil
}

// @super-admin

func (r Repository) SuperAdminGetList(ctx context.Context, filter banner.Filter) ([]banner.SuperAdminGetListResponse, int, error) {
	_, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return nil, 0, err
	}

	lang := r.GetLang(ctx)

	whereQuery := fmt.Sprintf(`WHERE b.deleted_at IS NULL AND b.expired_at > now()`)

	if filter.Status != nil {
		whereQuery += fmt.Sprintf(` AND b.status='%s'`, *filter.Status)
	}

	var limitQuery, offsetQuery string
	if filter.Limit != nil {
		limitQuery = fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
		if filter.Page != nil {
			offset := (*filter.Page - 1) * (*filter.Limit)
			offsetQuery = fmt.Sprintf(" OFFSET '%d'", offset)
		}
	}

	query := fmt.Sprintf(`
		SELECT
		    b.id,
		    b.photo,
		    b.title->>'%s',
		    b.description->>'%s',
		    b.price,
		    b.old_price,
		    br.restaurant_id,
		    r.name,
		    r.logo,
		    b.status
		FROM banners AS b
		    LEFT JOIN branches br ON b.branch_id = br.id
			LEFT JOIN restaurants AS r ON br.restaurant_id = r.id
		%s %s %s
`, lang, lang, whereQuery, limitQuery, offsetQuery)
	if err != nil {
		return nil, 0, errors.Wrap(err, "select query")
	}

	list := make([]banner.SuperAdminGetListResponse, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select stories"), http.StatusInternalServerError)
	}

	for rows.Next() {
		var (
			detail banner.SuperAdminGetListResponse
		)
		err = rows.Scan(&detail.Id, &detail.Photo, &detail.Title, &detail.Description, &detail.Price, &detail.OldPrice, &detail.RestaurantID, &detail.Restaurant, &detail.RestaurantPhoto, &detail.Status)
		if err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning stories"), http.StatusBadRequest)
		}

		if detail.OldPrice != nil && detail.Price != nil {
			discount := 100 - int((*detail.Price)/(*detail.OldPrice)*100)

			detail.Discount = &discount
		}

		if detail.Photo != nil {
			link := hashing.GenerateHash(r.ServerBaseUrl, *detail.Photo)
			detail.Photo = &link
		}

		if detail.RestaurantPhoto != nil {
			link := hashing.GenerateHash(r.ServerBaseUrl, *detail.RestaurantPhoto)
			detail.RestaurantPhoto = &link
		}
		list = append(list, detail)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(*)
		FROM banners AS b
			LEFT JOIN branches br ON b.branch_id = br.id
			LEFT JOIN restaurants AS r ON br.restaurant_id = r.id
		%s
	`, whereQuery)

	countRows, err := r.QueryContext(ctx, countQuery)
	if err == sql.ErrNoRows {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting stories"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning user count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func (r Repository) SuperAdminGetDetail(ctx context.Context, id int64) (*banner.SuperAdminGetDetailResponse, error) {
	_, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return nil, err
	}

	lang := r.GetLang(ctx)

	var (
		response    banner.SuperAdminGetDetailResponse
		menus       []banner.Menu
		title       *string
		description *string
	)

	query := fmt.Sprintf(`SELECT
									b.photo,
									b.title,
									b.description,
									b.price,
									b.old_price,
									b.menu_ids,
									r.id,
									r.logo,
									r.name
								FROM banners AS b
									LEFT JOIN branches br ON b.branch_id = br.id
									LEFT JOIN restaurants AS r ON br.restaurant_id = r.id
								WHERE b.deleted_at ISNULL AND b.expired_at > now() AND b.id='%d'`, id)
	if err = r.QueryRowContext(ctx, query).Scan(&response.Photo, &title, &description, &response.Price, &response.OldPrice, &response.MenuIds, &response.RestaurantID, &response.RestaurantPhoto, &response.Restaurant); err != nil {
		return nil, err
	}

	if response.OldPrice != nil && response.Price != nil {
		discount := 100 - int((*response.Price)/(*response.OldPrice)*100)

		response.Discount = &discount
	}

	if title != nil {
		m, err := commands.JsonToMap(*title)
		if err != nil {
			return nil, err
		}

		response.Title = m
	}

	if description != nil {
		m, err := commands.JsonToMap(*description)
		if err != nil {
			return nil, err
		}

		response.Description = m
	}

	if response.Photo != nil {
		link := hashing.GenerateHash(r.ServerBaseUrl, *response.Photo)
		response.Photo = &link
	}

	if response.RestaurantPhoto != nil {
		link := hashing.GenerateHash(r.ServerBaseUrl, *response.RestaurantPhoto)
		response.RestaurantPhoto = &link
	}

	if response.MenuIds != nil && len(*response.MenuIds) > 0 {
		for _, v := range *response.MenuIds {
			var menu banner.Menu
			queryMenu := fmt.Sprintf(`SELECT m.id,
													f.name,
													m.description->>'%s',
													f.photos[1],
													m.new_price
									 	 	 FROM menus m
									 	 	     JOIN foods f
									 	 	         ON m.food_id = f.id
									 	 	 WHERE m.id='%d'`, lang, v)
			if err = r.QueryRowContext(ctx, queryMenu).Scan(&menu.Id, &menu.Title, &menu.Description, &menu.Photo, &menu.Price); err != nil && !errors.Is(err, sql.ErrNoRows) {
				return nil, err
			}

			if menu.Photo != nil {
				link := hashing.GenerateHash(r.ServerBaseUrl, *menu.Photo)
				menu.Photo = &link
			}

			menus = append(menus, menu)
		}

		response.Menus = menus
	}

	return &response, nil
}

func (r Repository) SuperAdminUpdateStatus(ctx context.Context, id int64, status string) error {
	_, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return err
	}

	if status == "APPROVED" {
		countActive := 0
		err = r.QueryRowContext(ctx, fmt.Sprintf("SELECT count(id) FROM banners WHERE expired_at>now() AND status='APPROVED' AND deleted_at IS NULL")).Scan(&countActive)
		if err != nil {
			return web.NewRequestError(errors.Wrap(err, "scan active banners"), http.StatusInternalServerError)
		}
		if countActive > 20 {
			return web.NewRequestError(errors.New("the number of active banners has reached 20"), http.StatusBadRequest)
		}
	}

	_, err = r.ExecContext(ctx, fmt.Sprintf("UPDATE banners SET status='%s' WHERE id = '%d'", status, id))
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "update status banners"), http.StatusInternalServerError)
	}

	return nil
}

func NewRepository(DB *postgresql.Database) *Repository {
	return &Repository{DB}
}
