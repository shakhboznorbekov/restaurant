package utils

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"github.com/restaurant/foundation/web"
	"github.com/restaurant/internal/pkg/repository/postgresql"
	"github.com/restaurant/internal/repository/postgres"
	"net/http"
)

type Joins struct {
	JoinColumn *string
	MainColumn *string
}

type Count struct {
	TableName  *string
	WhereQuery *string
}

// fields that we gon' select.
// a := make(map[string][]string)
// a["region"] = []string{"id", "username"}
// a["users"] = []string{"number", "name"}
// a["client"] = []string{"phone", "address"}

// table for any joins.
// e,g: LEFT OUTER JOIN measure_units ON measure_units.id = main_table.measure_unit_id
// b := make(map[string]Joins)
// joinColumn := "id"
// mainColumn := "measure_unit_id"
// joinColumn1 := "id"
// mainColumn1 := "region_id"
// b["measure_units"] = Joins{JoinColumn: &joinColumn, MainColumn: &mainColumn}
// b["region"] = Joins{JoinColumn: &joinColumn1, MainColumn: &mainColumn1}

func SelectQuery(fields map[string][]string, joins map[string]Joins, mainTable *string, whereQuery *string) (string, error) {
	if fields == nil {
		return "", errors.New("fields is empty")
	}
	if mainTable == nil {
		return "", errors.New("table name is empty")
	}
	if whereQuery == nil {
		return "", errors.New("where query is empty")
	}

	var fieldQuery string
	lengthMap := len(fields)
	comma := ","
	level := 0
	for k, v := range fields {
		level += 1
		lengthRow := len(v)
		for i, r := range v {
			if lengthMap == level && lengthRow == i+1 {
				comma = ""
			}
			fieldQuery += fmt.Sprintf("%s.%s%s ", k, r, comma)
		}
	}

	var joinsQuery string
	if joins != nil {
		for k, v := range joins {
			if v.MainColumn == nil || v.JoinColumn == nil {
				return "", errors.New("mainColumn or joinColumn is empty")
			}

			joinsQuery += fmt.Sprintf("LEFT OUTER JOIN %s ON %s.%s = %s.%s ",
				k, k, *v.JoinColumn, *mainTable, *v.MainColumn)
		}
	}

	query := fmt.Sprintf(`SELECT %s FROM %s %s %s`, fieldQuery, *mainTable, joinsQuery, *whereQuery)
	return query, nil
}

func CountQuery(ctx context.Context, r *postgresql.Database, count Count) (int, error) {
	if count.TableName == nil {
		return 0, errors.New("count query: table name missing")
	}
	if count.WhereQuery == nil {
		return 0, errors.New("count query: where queries missing")
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(id)
		FROM
		    %s
		%s
	`, *count.TableName, *count.WhereQuery)

	countRows, err := r.QueryContext(ctx, countQuery)
	if err == sql.ErrNoRows {
		return 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	}
	if err != nil {
		return 0, web.NewRequestError(errors.Wrap(err, "selecting users"), http.StatusBadRequest)
	}

	countInt := 0
	for countRows.Next() {
		if err = countRows.Scan(&countInt); err != nil {
			return 0, web.NewRequestError(errors.Wrap(err, "scanning user count"), http.StatusBadRequest)
		}
	}

	return countInt, nil
}
