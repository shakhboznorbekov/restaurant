package full_search

import (
	"context"
	"fmt"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"net/http"
	"restu-backend/foundation/web"
	"restu-backend/internal/pkg/repository/postgresql"
	"restu-backend/internal/service/full_search"
	"restu-backend/internal/service/hashing"
	"strings"
	"time"
)

type Repository struct {
	*postgresql.Database
}

// client

func (r Repository) ClientGetList(ctx context.Context, filter full_search.Filter) ([]full_search.ClientGetList, error) {
	//_, err := r.CheckClaims(ctx, auth.RoleClient)
	//if err != nil {
	//	return nil, err
	//}

	whereQuery := fmt.Sprintf(`WHERE b.deleted_at IS NULL`)

	search := *filter.Search

	location := true
	food := true
	if filter.Lat == nil || filter.Lon == nil {
		var l float64
		filter.Lat = &l
		filter.Lon = &l
		location = false
	}

	queryOrder := ""
	if location {
		queryOrder += " ORDER BY distance"
	} else {
		queryOrder += " ORDER BY b.rate"
	}

	if len(search) > 2 {
		if filter.Menu == nil {
			whereQuery += fmt.Sprintf(` AND 
			(
				b.menu_names ilike '%s' OR
				b.name ilike '%s' OR
				b.name ilike '%s'
				
			)`,
				"% "+search+"%",
				"% "+search+"%",
				""+search+"%",
			)
		} else if strings.ToUpper(*filter.Menu) == "RESTAURANT" {
			whereQuery += fmt.Sprintf(` AND (b.name ilike '%s' OR b.name ilike '%s')`, "% "+search+"%", ""+search+"%")
			food = false
		} else if strings.ToUpper(*filter.Menu) == "FOOD" {
			whereQuery += fmt.Sprintf(` AND b.menu_names ilike '%s'`, "% "+search+"%")
		}
	} else {
		queryOrder = " ORDER BY b.search_count desc"
		food = false
		limit := 10
		if filter.Limit == nil {
			filter.Limit = &limit
		} else if *filter.Limit > 10 {
			filter.Limit = &limit
		}
	}

	todayInt := time.Now().Weekday()
	var today string
	switch todayInt {
	case 1:
		today = "Monday"
	case 2:
		today = "Tuesday"
	case 3:
		today = "Wednesday"
	case 4:
		today = "Thursday"
	case 5:
		today = "Friday"
	case 6:
		today = "Saturday"
	case 7:
		today = "Sunday"
	}

	var limitQuery, offsetQuery string
	if filter.Page != nil && filter.Limit != nil {
		offset := (*filter.Page - 1) * (*filter.Limit)
		filter.Offset = &offset
	}
	if filter.Limit != nil {
		limitQuery = fmt.Sprintf(" LIMIT '%d'", *filter.Limit)
	}
	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf(" OFFSET '%d'", *filter.Offset)
	}

	query := fmt.Sprintf(`
				SELECT
				    b.id,
				    b.status,
				    b.location,
				    b.photos as photos,
				    b.work_time->>'%s' as work_time_today,
				    b.name,
				    b.category_id,
				    rc.name as category_name,
				    CASE WHEN b.menu_names ilike '%s' THEN true ELSE false END menu_status,
				    CASE WHEN '%t' THEN ST_DistanceSphere(
				            ST_SetSRID(ST_MakePoint(CAST(b.location->>'lon' AS float), CAST(b.location->>'lat' AS float)), 4326),
				            ST_SetSRID(ST_MakePoint('%v', '%v'), 4326)
				        ) END as distance,
				    b.rate
				FROM
				    branches AS b
				LEFT OUTER JOIN restaurant_category AS rc ON rc.id = b.category_id
				%s
				%s %s %s
	`, today, "% "+search+"%", location, *filter.Lon, *filter.Lat, whereQuery, queryOrder, limitQuery, offsetQuery)

	list := make([]full_search.ClientGetList, 0)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, web.NewRequestError(errors.Wrap(err, "select branches"), http.StatusInternalServerError)
	}

	err = r.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, web.NewRequestError(errors.Wrap(err, "scanning branches"), http.StatusBadRequest)
	}

	for k, v := range list {
		if v.Photos != nil && len(*v.Photos) > 0 {
			var photoLink pq.StringArray
			for _, v2 := range *v.Photos {
				if v2 != "" {
					baseLink := hashing.GenerateHash(r.ServerBaseUrl, v2)
					photoLink = append(photoLink, baseLink)
				}
			}
			list[k].Photos = &photoLink
		}
		if v.MenuStatus && food {
			queryMenu := fmt.Sprintf(`
										SELECT 
										    m.id,
										    m.name,
										    m.photos,
										    m.new_price as price
										FROM 
										    menus AS m
										WHERE 
										    m.branch_id = %d AND m.deleted_at IS NULL AND m.deleted_at IS NULL `, v.ID)
			menus := make([]full_search.Menu, 0)

			rowsMenu, err := r.QueryContext(ctx, queryMenu)
			if err != nil {
				return nil, web.NewRequestError(errors.Wrap(err, "select branches"), http.StatusInternalServerError)
			}

			err = r.ScanRows(ctx, rowsMenu, &menus)
			if err != nil {
				return nil, web.NewRequestError(errors.Wrap(err, "scanning branches"), http.StatusBadRequest)
			}

			list[k].Menus = menus

			for k1, v1 := range menus {
				var photoFoodLink pq.StringArray
				if v1.Photos != nil {
					for _, v2 := range *v1.Photos {
						baseLink := hashing.GenerateHash(r.ServerBaseUrl, v2)
						photoFoodLink = append(photoFoodLink, baseLink)
					}
					list[k].Menus[k1].Photos = &photoFoodLink
				}
			}
		}
	}
	return list, nil
}

func NewRepository(DB *postgresql.Database) *Repository {
	return &Repository{DB}
}
