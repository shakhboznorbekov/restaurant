package story

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"restu-backend/foundation/web"
	"restu-backend/internal/auth"
	"restu-backend/internal/entity"
	"restu-backend/internal/pkg/repository/postgresql"
	"restu-backend/internal/repository/postgres"
	"restu-backend/internal/service/hashing"
	"restu-backend/internal/service/story"
	"time"
)

type Repository struct {
	*postgresql.Database
}

// admin

func (r Repository) AdminGetList(ctx context.Context, filter story.Filter) ([]story.AdminGetList, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return nil, 0, err
	}

	whereQuery := fmt.Sprintf(`WHERE deleted_at IS NULL AND restaurant_id = '%d'`, *claims.RestaurantID)
	countWhereQuery := whereQuery

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
		    name,
		    file,
		    type,
		    TO_CHAR(expired_at,'DD.MM.YYYY | HH24:MI'),
		    duration,
		    status
		FROM stories
		%s %s %s
`, whereQuery, limitQuery, offsetQuery)
	if err != nil {
		return nil, 0, errors.Wrap(err, "select query")
	}

	list := make([]story.AdminGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "select stories"), http.StatusInternalServerError)
	}

	for rows.Next() {
		var detail story.AdminGetList
		err = rows.Scan(&detail.ID, &detail.Name, &detail.File, &detail.Type, &detail.ExpiredAt, &detail.Duration, &detail.Status)
		if err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning stories"), http.StatusBadRequest)
		}

		if detail.File != nil {
			link := hashing.GenerateHash(r.ServerBaseUrl, *detail.File)
			detail.File = &link
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
		FROM stories
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

func (r Repository) AdminGetDetail(ctx context.Context, id int64) (entity.Story, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return entity.Story{}, err
	}

	var detail entity.Story

	err = r.NewSelect().Model(&detail).Where("id = ? AND restaurant_id = ? AND deleted_at IS NULL", id, claims.RestaurantID).Scan(ctx)
	if err != nil {
		return entity.Story{}, err
	}

	return detail, nil
}

func (r Repository) AdminCreate(ctx context.Context, request story.AdminCreateRequest) (story.AdminCreateResponse, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return story.AdminCreateResponse{}, err
	}

	err = r.ValidateStruct(&request, "Name", "FileLink", "Duration", "Type")
	if err != nil {
		return story.AdminCreateResponse{}, err
	}

	response := story.AdminCreateResponse{
		Name:         request.Name,
		File:         request.FileLink,
		Type:         request.Type,
		Duration:     request.Duration,
		ExpiredAt:    time.Now(),
		CreatedAt:    time.Now(),
		CreatedBy:    claims.UserId,
		RestaurantID: *claims.RestaurantID,
	}

	_, err = r.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return story.AdminCreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating story"), http.StatusBadRequest)
	}

	if response.File != nil {
		link := hashing.GenerateHash(r.ServerBaseUrl, *response.File)
		response.File = &link
	}

	return response, nil
}

func (r Repository) AdminUpdateStatus(ctx context.Context, id int64) error {
	claims, err := r.CheckClaims(ctx, auth.RoleAdmin)
	if err != nil {
		return err
	}

	_, err = r.ExecContext(ctx, fmt.Sprintf(`INSERT 
														INTO stories (name, file, type, duration, restaurant_id, created_by, created_at, status) 
															select name, file, type, duration, restaurant_id, created_by, now(), 'DRAFT' from stories where id='%d' and restaurant_id='%d'`, id, *claims.RestaurantID))
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "scan active stories"), http.StatusInternalServerError)
	}

	return nil
}

func (r Repository) AdminDelete(ctx context.Context, id int64) error {
	return r.DeleteRow(ctx, "stories", id, auth.RoleAdmin)
}

// client

func (r Repository) ClientGetList(ctx context.Context, filter story.Filter) ([]story.ClientGetList, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleClient)
	if err != nil {
		claims.UserId = 0
	}

	whereQuery := fmt.Sprintf(`WHERE (select count(s.id) from stories s where s.restaurant_id=o.id and s.expired_at > now() and s.deleted_at isnull and s.status='APPROVED') != 0 AND o.deleted_at isnull`)

	var limitQuery, offsetQuery string
	if filter.Limit != nil {
		limitQuery = fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	query := fmt.Sprintf(`SELECT 
    								o.id,
    								o.name,
    								o.mini_logo,
    								(select count(s.id) from stories s where s.restaurant_id=o.id and s.expired_at > now() and s.deleted_at isnull and s.status='APPROVED') = (select count(v.id) from story_views v join stories s on s.id=v.story_id where v.created_by='%d' and s.restaurant_id=o.id and s.status='APPROVED' and s.expired_at > now()) as seen
								 FROM restaurants o %s ORDER BY seen`, claims.UserId, whereQuery)
	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	list := []story.ClientGetList{}
	for rows.Next() {
		var row story.ClientGetList
		if err = rows.Scan(&row.RestaurantId, &row.RestaurantName, &row.RestaurantLogo, &row.Seen); err != nil {
			return nil, 0, err
		}

		storyQuery := fmt.Sprintf(`SELECT 
    											s.id, 
    											s.file, 
    											s.type, 
    											s.duration, 
    											(select exists(select id from story_views where story_id=s.id and created_by='%d')) as seen 
										  FROM stories s 
										  WHERE 
										      s.deleted_at isnull 
										    and 
										      s.expired_at>now() and s.restaurant_id='%d'
										    and 
										      s.status='APPROVED'
										  ORDER BY s.created_at %s %s`, claims.UserId, row.RestaurantId, limitQuery, offsetQuery)
		rows, err := r.QueryContext(ctx, storyQuery)
		if err != nil {
			return nil, 0, err
		}

		var stories []story.Story
		for rows.Next() {
			var row story.Story
			if err = rows.Scan(&row.Id, &row.File, &row.Type, &row.Duration, &row.Seen); err != nil {
				return nil, 0, err
			}

			if row.File != nil {
				link := hashing.GenerateHash(r.ServerBaseUrl, *row.File)
				row.File = &link
			}

			stories = append(stories, row)
		}

		row.Stories = stories

		if row.RestaurantLogo != nil {
			link := hashing.GenerateHash(r.ServerBaseUrl, *row.RestaurantLogo)
			row.RestaurantLogo = &link
		}

		list = append(list, row)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(o.id)
		FROM restaurants o 
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

func (r Repository) ClientSetViewed(ctx context.Context, id int64) error {
	claims, err := r.CheckClaims(ctx, auth.RoleClient)
	if err != nil {
		return nil
	}

	query := fmt.Sprintf(`INSERT 
									INTO story_views (story_id, created_by) 
									VALUES ('%d', '%d') 
										ON CONFLICT 
										    ON CONSTRAINT u_story_views 
									DO UPDATE SET updated_at=now()`, id, claims.UserId)

	if _, err = r.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}

// @super-admin

func (r Repository) SuperAdminGetList(ctx context.Context, filter story.Filter) ([]story.SuperAdminGetListResponse, int, error) {
	claims, err := r.CheckClaims(ctx, auth.RoleSuperAdmin)
	if err != nil {
		return nil, 0, err
	}

	var whereQuery string
	if filter.Status != nil {
		whereQuery = fmt.Sprintf(`WHERE (select count(s.id) from stories s where s.restaurant_id=o.id and s.deleted_at isnull and s.status='%s') != 0`, *filter.Status)
	} else {
		whereQuery = fmt.Sprintf(`WHERE (select count(s.id) from stories s where s.restaurant_id=o.id and s.deleted_at isnull) != 0`)
	}

	var limitQuery, offsetQuery string
	if filter.Limit != nil {
		limitQuery = fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	query := fmt.Sprintf(`SELECT 
    								o.id,
    								o.name,
    								o.mini_logo,
    								(select count(s.id) from stories s where s.restaurant_id=o.id and s.expired_at > now() and s.deleted_at isnull) = (select count(v.id) from story_views v join stories s on s.id=v.story_id where v.created_by='%d' and s.restaurant_id=o.id) as seen
								 FROM restaurants o %s ORDER BY created_at DESC`, claims.UserId, whereQuery)
	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	list := []story.SuperAdminGetListResponse{}
	for rows.Next() {
		var row story.SuperAdminGetListResponse
		if err = rows.Scan(&row.RestaurantId, &row.RestaurantName, &row.RestaurantLogo, &row.Seen); err != nil {
			return nil, 0, err
		}

		storyQuery := fmt.Sprintf(`SELECT 
    											s.id, 
    											s.file, 
    											s.type, 
    											s.duration, 
    											(select exists(select id from story_views where story_id=s.id and created_by='%d')) as seen ,
												s.status
										  FROM stories s 
										  WHERE 
										      s.deleted_at isnull 
										    and 
										      s.restaurant_id='%d'
										  ORDER BY seen %s %s`, claims.UserId, row.RestaurantId, limitQuery, offsetQuery)
		rowsStory, err := r.QueryContext(ctx, storyQuery)
		if err != nil {
			return nil, 0, err
		}

		var stories []story.Story
		for rowsStory.Next() {
			var rowStory story.Story
			if err = rowsStory.Scan(&rowStory.Id, &rowStory.File, &rowStory.Type, &rowStory.Duration, &rowStory.Seen, &rowStory.Status); err != nil {
				return nil, 0, err
			}

			if rowStory.File != nil {
				link := hashing.GenerateHash(r.ServerBaseUrl, *rowStory.File)
				rowStory.File = &link
			}

			stories = append(stories, rowStory)
		}

		row.Stories = stories

		if row.RestaurantLogo != nil {
			link := hashing.GenerateHash(r.ServerBaseUrl, *row.RestaurantLogo)
			row.RestaurantLogo = &link
		}

		list = append(list, row)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(o.id)
		FROM restaurants o 
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

func (r Repository) SuperAdminUpdateStatus(ctx context.Context, id int64, status string) error {
	if status == "APPROVED" {
		countActive := 0
		err := r.QueryRowContext(ctx, fmt.Sprintf("SELECT count(id) FROM stories WHERE expired_at>now() AND deleted_at IS NULL AND status='APPROVED'")).Scan(&countActive)
		if err != nil {
			return web.NewRequestError(errors.Wrap(err, "scan active stories"), http.StatusInternalServerError)
		}
		if countActive > 20 {
			return web.NewRequestError(errors.New("the number of active stories has reached 20"), http.StatusBadRequest)
		}
	}

	query := fmt.Sprintf(`UPDATE stories SET expired_at=now() + interval '24 hours', status='%s' WHERE id='%d'`, status, id)

	if _, err := r.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}

func NewRepository(DB *postgresql.Database) *Repository {
	return &Repository{DB}
}
