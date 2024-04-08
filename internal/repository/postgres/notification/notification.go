package notification

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"net/http"
	"restu-backend/foundation/web"
	"restu-backend/internal/auth"
	"restu-backend/internal/commands"
	"restu-backend/internal/pkg/repository/postgresql"
	"restu-backend/internal/repository/postgres"
	"restu-backend/internal/service/hashing"
	"restu-backend/internal/service/notification"
	"time"
)

type Repository struct {
	*postgresql.Database
}

func NewRepository(DB *postgresql.Database) *Repository {
	return &Repository{DB}
}

// @admin

func (r Repository) AdminGetList(ctx context.Context, filter notification.Filter) ([]notification.AdminGetList, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return nil, 0, err
	}
	lang := r.GetLang(ctx)
	whereQuery := fmt.Sprintf(`WHERE deleted_at IS NULL and restaurant_id='%d'`, *claims.RestaurantID)

	if filter.Status != nil {
		whereQuery += fmt.Sprintf(` AND status = '%s'`, *filter.Status)
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
		    title->>'%s',
		    description->>'%s',
		    photo,
		    status
		FROM notifications
		%s %s %s
`, lang, lang, whereQuery, limitQuery, offsetQuery)

	var list []notification.AdminGetList
	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select notifications"), http.StatusInternalServerError)
	}

	for rows.Next() {
		var detail notification.AdminGetList

		if err = rows.Scan(
			&detail.ID,
			&detail.Title,
			&detail.Description,
			&detail.Photo,
			&detail.Status,
		); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning notifications"), http.StatusBadRequest)
		}

		if detail.Photo != nil {
			link := hashing.GenerateHash(r.ServerBaseUrl, *detail.Photo)
			detail.Photo = &link
		}
		list = append(list, detail)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(id)
		FROM notifications
		%s
	`, countWhereQuery)

	countRows, err := r.QueryContext(ctx, countQuery)
	if err == sql.ErrNoRows {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting notifications"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning notifications count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func (r Repository) AdminGetDetail(ctx context.Context, id int64) (*notification.AdminGetDetail, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return nil, err
	}

	var (
		response    notification.AdminGetDetail
		title       *string
		description *string
	)

	query := fmt.Sprintf(`
		SELECT
		    id,
		    title,
		    description,
		    photo,
		    status
		FROM notifications
	WHERE deleted_at IS NULL AND id='%d' AND restaurant_id='%d'`, id, *claims.RestaurantID)
	if err = r.QueryRowContext(ctx, query).Scan(&response.ID, &title, &description, &response.Photo, &response.Status); err != nil {
		return nil, err
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

	return &response, nil
}

func (r Repository) AdminCreate(ctx context.Context, request notification.AdminCreateRequest) (*notification.AdminCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return nil, err
	}

	err = r.ValidateStruct(&request, "Title", "Description", "PhotoLink")
	if err != nil {
		return nil, err
	}

	response := notification.AdminCreateResponse{
		Title:        request.Title,
		Description:  request.Description,
		Photo:        request.PhotoLink,
		CreatedAt:    time.Now(),
		CreatedBy:    claims.UserId,
		RestaurantId: claims.RestaurantID,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return nil, web.NewRequestError(errors.Wrap(err, "creating notification"), http.StatusBadRequest)
	}

	if response.Photo != nil {
		link := hashing.GenerateHash(r.ServerBaseUrl, *response.Photo)
		response.Photo = &link
	}

	return &response, nil
}

func (r Repository) AdminUpdateAll(ctx context.Context, request notification.AdminUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID", "Title", "Description"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("notifications").Where("deleted_at IS NULL AND status = 'NEW' AND id = ? AND restaurant_id = ?",
		request.ID, *claims.RestaurantID)

	q.Set("title = ?", request.Title)
	q.Set("description = ?", request.Description)
	q.Set("photo = ?", request.PhotoLink)
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating notification"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) AdminUpdateColumn(ctx context.Context, request notification.AdminUpdateRequest) error {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	if err = r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	q := r.NewUpdate().Table("notifications").Where("deleted_at IS NULL AND status = 'NEW' AND id = ? AND restaurant_id = ?",
		request.ID, *claims.RestaurantID)

	if request.Title != nil {
		q.Set("title = ?", request.Title)
	}
	if request.Description != nil {
		q.Set("description = ?", request.Description)
	}
	if request.PhotoLink != nil {
		q.Set("photo = ?", *request.PhotoLink)
	}

	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating notification"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) AdminDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "notifications", id, auth.RoleAdmin)
}

// @super-admin

func (r Repository) SuperAdminUpdateStatus(ctx context.Context, id int64, status string) error {
	claims, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return err
	}

	_, err = r.ExecContext(ctx,
		fmt.Sprintf(`UPDATE notifications 
							SET status='%s',
								updated_by = '%d',
								updated_at = NOW()
							WHERE status = 'NEW'
							AND id = '%d'`, status, claims.UserId, id))
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "error new notifications"), http.StatusInternalServerError)
	}

	return nil
}

func (r Repository) SuperAdminGetList(ctx context.Context, filter notification.Filter) ([]notification.SuperAdminGetList, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return nil, 0, err
	}

	lang := r.GetLang(ctx)

	whereQuery := fmt.Sprintf(`WHERE deleted_at IS NULL `)

	if filter.Status != nil {
		whereQuery += fmt.Sprintf(` AND status = '%s'`, *filter.Status)
	}
	if filter.Whose != nil {
		if *filter.Whose == "MYSELF" {
			whereQuery += fmt.Sprintf(` AND created_by = '%d'`, claims.UserId)
		} else if *filter.Whose == "OTHERS" {
			whereQuery += fmt.Sprintf(` AND created_by != '%d'`, claims.UserId)
		}
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
		    title->>'%s' as title,
		    description->>'%s' as description,
		    photo,
		    status
		FROM notifications
		%s %s %s
`, lang, lang, whereQuery, limitQuery, offsetQuery)

	list := make([]notification.SuperAdminGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select notifications"), http.StatusInternalServerError)
	}

	if err = r.ScanRows(ctx, rows, &list); err != nil {
		return nil, 0, err
	}

	for i := range list {
		if list[i].Photo != nil {
			link := hashing.GenerateHash(r.ServerBaseUrl, *list[i].Photo)
			list[i].Photo = &link
		}
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(id)
		FROM notifications
		%s
	`, countWhereQuery)

	countRows, err := r.QueryContext(ctx, countQuery)
	if err == sql.ErrNoRows {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting notifications"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning notifications count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func (r Repository) SuperAdminGetDetail(ctx context.Context, id int64) (*notification.SuperAdminGetDetail, error) {
	_, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return nil, err
	}

	lang := r.GetLang(ctx)

	var (
		response notification.SuperAdminGetDetail
	)

	query := fmt.Sprintf(`
		SELECT
		    n.title->>'%s',
		    n.description->>'%s',
		    n.photo,
		    array (select device_token from devices d where d.deleted_at is null and not is_log_out) as tokens
		FROM notifications n
	WHERE n.deleted_at IS NULL AND n.id='%d' AND n.title->>'%s' is not null AND n.description->>'%s' is not null`, lang, lang, id, lang, lang)
	if err = r.QueryRowContext(ctx, query).Scan(&response.Title, &response.Description, &response.Photo, &response.DeviceTokens); err != nil {
		return nil, err
	}

	if response.Photo != nil {
		link := hashing.GenerateHash(r.ServerBaseUrl, *response.Photo)
		response.Photo = &link
	}

	return &response, nil
}

func (r Repository) SuperAdminSend(ctx context.Context, request notification.SuperAdminSendRequest) ([]notification.SuperAdminSendResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return nil, err
	}

	if err = r.ValidateStruct(&request, "Title", "Description", "Status"); err != nil {
		return nil, err
	}

	var (
		response           []notification.SuperAdminSendResponse
		title, description string
		query              string
	)
	if request.Title != nil {
		title, err = commands.MapToJson(request.Title)
		if err != nil {
			return nil, err
		}
	}
	if request.Description != nil {
		description, err = commands.MapToJson(request.Description)
		if err != nil {
			return nil, err
		}
	}

	if *request.Status == "SENT" {
		queryDevice := fmt.Sprintf(`SELECT 
    											device_token, 
    											device_lang 
											FROM devices 
											WHERE 
											    deleted_at isnull 
											  and 
											    not is_log_out`)
		if request.UserId != nil {
			queryDevice += fmt.Sprintf(` and user_id='%d'`, *request.UserId)
		}

		rows, err := r.QueryContext(ctx, queryDevice)
		if err != nil {
			return nil, err
		}

		if err = r.ScanRows(ctx, rows, &response); err != nil {
			return nil, err
		}
	}

	if request.PhotoLink != nil {
		query = fmt.Sprintf(`INSERT INTO notifications (title, description, photo, status, created_by) VALUES ('%s', '%s', '%s', '%s', '%d')`, title, description, *request.PhotoLink, *request.Status, claims.UserId)
	} else {
		query = fmt.Sprintf(`INSERT INTO notifications (title, description, status, created_by) VALUES ('%s', '%s', '%s', '%d')`, title, description, *request.Status, claims.UserId)
	}

	if _, err = r.ExecContext(ctx, query); err != nil {
		return nil, err
	}

	return response, nil
}

// @client

func (r Repository) ClientGetList(ctx context.Context, filter notification.Filter) ([]notification.ClientGetListResponse, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleClient)
	if err != nil {
		return []notification.ClientGetListResponse{}, 0, nil
	}

	lang := r.GetLang(ctx)

	whereQuery := fmt.Sprintf(`WHERE n.deleted_at IS NULL and n.status='SENT' and n.title->>'%s' is not null and n.description->>'%s' is not null`, lang, lang)

	countWhereQuery := whereQuery

	var limitQuery, offsetQuery string
	if filter.Limit != nil {
		limitQuery = fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Page != nil {
		offset := (*filter.Page - 1) * (*filter.Limit)
		offsetQuery = fmt.Sprintf(" OFFSET '%d'", offset)
	}

	query := fmt.Sprintf(`
		SELECT
		    n.id,
		    n.title->>'%s' as title,
		    n.description->>'%s' as description,
		    n.photo,
		    TO_CHAR(n.created_at, 'DD.MM.YYYY HH24:MI') as created_at,
		    exists(select v.id from notification_views v where v.notification_id = n.id and v.created_by='%d') as seen
		FROM notifications n
		%s ORDER BY created_at DESC %s %s
		`, lang, lang, claims.UserId, whereQuery, limitQuery, offsetQuery)

	list := make([]notification.ClientGetListResponse, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select notifications"), http.StatusInternalServerError)
	}

	if err = r.ScanRows(ctx, rows, &list); err != nil {
		return nil, 0, err
	}

	for i := range list {
		if list[i].Photo != nil {
			link := hashing.GenerateHash(r.ServerBaseUrl, *list[i].Photo)
			list[i].Photo = &link
		}
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(n.id)
		FROM notifications n 
		%s
	`, countWhereQuery)

	countRows, err := r.QueryContext(ctx, countQuery)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
		}
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting notifications"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning notifications count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func (r Repository) ClientGetUnseenCount(ctx context.Context) (int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleClient)
	if err != nil {
		return 0, nil
	}

	lang := r.GetLang(ctx)

	var count int

	query := fmt.Sprintf(`SELECT
    								(SELECT count(n.id) FROM notifications n WHERE n.status = 'SENT' and n.deleted_at isnull and n.title->>'%s' is not null and n.description->>'%s' is not null) 
    								    - 
    								(SELECT count(v.id) FROM notification_views v join notifications n on n.id=v.notification_id WHERE v.created_by='%d' and n.deleted_at isnull and n.status='SENT' and n.title->>'%s' is not null and n.description->>'%s' is not null) 
    								    as count;`, lang, lang, claims.UserId, lang, lang)
	if err = r.QueryRowContext(ctx, query).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (r Repository) ClientSetAsViewed(ctx context.Context, id int64) error {
	claims, err := r.CheckClaims(ctx, auth.RoleClient)
	if err != nil {
		return nil
	}

	query := fmt.Sprintf(`INSERT INTO 
    									notification_views (created_by, notification_id) 
								 VALUES ('%d', '%d') 
								 	ON CONFLICT 
								 	    ON CONSTRAINT u_notification_view 
								 	    DO UPDATE SET updated_at=now()`, claims.UserId, id)
	if _, err = r.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}

func (r Repository) ClientSetAllAsViewed(ctx context.Context) error {
	claims, err := r.CheckClaims(ctx, auth.RoleClient)
	if err != nil {
		return nil
	}

	var ids *pq.Int64Array
	idQuery := fmt.Sprintf(`SELECT ARRAY(SELECT id FROM notifications WHERE deleted_at isnull and status='SENT') as ids`)

	if err = r.QueryRowContext(ctx, idQuery).Scan(&ids); err != nil {
		return err
	}

	if ids != nil && len(*ids) > 0 {
		for _, v := range *ids {
			query := fmt.Sprintf(`INSERT INTO 
												notification_views (created_by, notification_id) 
										 VALUES ('%d', '%d') 
											ON CONFLICT 
												ON CONSTRAINT u_notification_view 
												DO UPDATE SET updated_at=now()`, claims.UserId, v)
			if _, err = r.ExecContext(ctx, query); err != nil {
				return err
			}
		}
	}

	return nil
}
