package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"github.com/restaurant/foundation/web"
	"github.com/restaurant/internal/auth"
	"github.com/restaurant/internal/entity"
	"github.com/restaurant/internal/repository/postgres"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type CurrencyValue struct {
	ID        string   `json:"id"`
	Value     *float32 `json:"value"`
	PriceDate *string  `json:"price_date"`
	Currency  *string  `json:"currency"`
	Icon      *string  `json:"icon"`
}

// Config is the required properties to use the database.
type Config struct {
	User          string
	Password      string
	Host          string
	Name          string
	DisableTLS    bool
	ServerBaseUrl string
	DefaultLang   string
}

type Database struct {
	*bun.DB
	DBName        string
	DBPassword    string
	DBUser        string
	ServerBaseUrl string
	DefaultLang   string
}

func NewDB(cfg Config) *Database {
	dsn := fmt.Sprintf("postgres://%v:%v@localhost:5432/%v?sslmode=disable", cfg.User, cfg.Password, cfg.Name)

	sqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(sqlDB, pgdialect.New())

	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	return &Database{DB: db, DBName: cfg.Name, DBPassword: cfg.Password, DBUser: cfg.User, ServerBaseUrl: cfg.ServerBaseUrl, DefaultLang: cfg.DefaultLang}
}

func (d Database) DeleteRow(ctx context.Context, table string, id int64, role string) error {
	claims, err := d.CheckClaims(ctx, role)
	if err != nil {
		return err
	}

	//if _, err = uuid.Parse(id); err != nil {
	//	return web.NewRequestError(postgres.ErrInvalidID, http.StatusBadRequest)
	//}

	q := d.NewUpdate().
		Table(table).
		Where("id = ?", id).
		Set("deleted_at = ?", time.Now()).
		Set("deleted_by = ?", claims.UserId)

	//if claims.RestaurantID != nil {
	//	q.Where("restaurant_id = ?", *claims.RestaurantID)
	//}

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrapf(err, "deleting %s", table), http.StatusBadRequest)
	}

	return nil
}

func (d Database) Log(ctx context.Context, data entity.Log) error {
	_, err := d.NewInsert().Model(&data).Exec(ctx)

	return err
}

func (d Database) CheckClaims(ctx context.Context, role string) (auth.Claims, error) {
	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return auth.Claims{}, web.NewRequestError(errors.New("claims missing from context"), http.StatusBadRequest)
	}

	if strings.Compare(role, auth.RoleAdmin) == 0 && claims.RestaurantID == nil {
		return auth.Claims{}, web.NewRequestError(
			errors.New("role admin doesn't contain RestaurantID"),
			http.StatusUnauthorized,
		)
	} else if (strings.Compare(role, auth.RoleBranch) == 0 ||
		strings.Compare(role, auth.RoleCashier) == 0 ||
		strings.Compare(role, auth.RoleWaiter) == 0) && claims.BranchID == nil {
		return auth.Claims{}, web.NewRequestError(
			errors.New("role (branch, cashier, waiter) doesn't contain BranchID"),
			http.StatusUnauthorized,
		)
	}

	if ok = claims.Authorized(role); !ok {
		return auth.Claims{}, web.NewRequestError(postgres.ErrForbidden, http.StatusForbidden)
	}

	return claims, nil
}

func (d Database) GetLang(ctx context.Context) string {
	if value, ok := ctx.Value("lang").(string); ok {
		return value
	}

	return d.DefaultLang
}

func (d Database) LoadTimeZone() *time.Location {
	loc, _ := time.LoadLocation("Asia/Tashkent")

	return loc
}

func (d Database) ValidateStruct(s interface{}, requiredFields ...string) error {
	structVal := reflect.Value{}
	if reflect.Indirect(reflect.ValueOf(s)).Kind() == reflect.Struct {
		structVal = reflect.Indirect(reflect.ValueOf(s))
	} else {
		return errors.New("input param should be a struct")
	}

	errFields := make([]web.FieldError, 0)

	structType := reflect.Indirect(reflect.ValueOf(s)).Type()
	fieldNum := structVal.NumField()

	for i := 0; i < fieldNum; i++ {
		field := structVal.Field(i)
		fieldName := structType.Field(i).Name

		isSet := field.IsValid() && !field.IsZero()
		if !isSet {
			log.Print(isSet, fieldName, reflect.ValueOf(field))
			for _, f := range requiredFields {
				if f == fieldName {
					errFields = append(errFields, web.FieldError{
						Error: "field is required!",
						Field: fieldName,
					})
				}
			}
		}
	}

	if len(errFields) > 0 {
		return &web.Error{
			Err:    errors.New("required fields"),
			Fields: errFields,
			Status: http.StatusBadRequest,
		}
	}

	return nil
}

//func (d Database) GetCurrencyValue(ctx context.Context, priceDate, currencyId, role string) (CurrencyValue, error) {
//	claims, err := d.CheckClaims(ctx, role)
//	if err != nil {
//		return CurrencyValue{}, err
//	}
//
//	var currencyDetail entity.Currency
//
//	err = d.NewSelect().Model(&currencyDetail).Where("id = ?", currencyId).Scan(ctx)
//	if err == sql.ErrNoRows {
//		return CurrencyValue{}, web.NewRequestError(errors.Wrap(postgres.ErrNotFound, "currency not found in GetCurrencyValue function"), http.StatusNotFound)
//	}
//	if err != nil {
//		return CurrencyValue{}, web.NewRequestError(errors.Wrap(err, "selecting currency in GetCurrencyValue function"), http.StatusNotFound)
//	}
//
//	lang := d.GetLang(ctx)
//
//	query := fmt.Sprintf(`
//		SELECT
//			value,
//			TO_CHAR(price_date, 'DD.MM.YYYY')
//		FROM currency_prices
//		WHERE
//				currency_id = '%s' AND price_date <= to_date('%s', 'DD.MM.YYYY') AND deleted_at IS NULL AND restaurant_id = '%s'
//		GROUP BY id, value, price_date
//		ORDER BY
//			   to_date('%s', 'DD.MM.YYYY') - price_date
//		LIMIT 1
//	`, currencyId, priceDate, *claims.RestaurantID, priceDate)
//
//	rows, err := d.QueryContext(ctx, query)
//	if err != nil {
//		return CurrencyValue{}, web.NewRequestError(errors.Wrap(err, "selecting currency_price in GetCurrencyValue function"), http.StatusNotFound)
//	}
//
//	var response CurrencyValue
//
//	for rows.Next() {
//		err = rows.Scan(&response.Value, &response.PriceDate)
//		if err != nil {
//			return CurrencyValue{}, web.NewRequestError(errors.Wrap(err, "scanning currency price in GetCurrencyValue function"), http.StatusNotFound)
//		}
//	}
//
//	response.ID = currencyId
//	response.Icon = currencyDetail.Icon
//	if currencyDetail.Name != nil {
//		currencyName := currencyDetail.Name[lang]
//		response.Currency = &currencyName
//	}
//
//	if response.Value == nil {
//		response.Value = currencyDetail.Value
//	}
//
//	return response, nil
//}

func (d Database) DeleteImageIndex(ctx context.Context, id int64, table string, column string, images []string, index int, role string) error {
	claims, err := d.CheckClaims(ctx, role)
	if err != nil {
		return err
	}

	columnF := fmt.Sprintf(`%s = ?`, column)
	newImages, err := deleteImage(images, index)
	if err != nil {
		return err
	}

	q := d.NewUpdate().
		Table(table).
		Where("id = ?", id).
		Set(columnF, pgdialect.Array(newImages)).
		Set("updated_at = ?", time.Now()).
		Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrapf(err, "deleting %s", table), http.StatusBadRequest)
	}

	return nil
}

func deleteImage(slice []string, imageIndex int) ([]string, error) {

	if imageIndex < 0 || imageIndex >= len(slice) {
		return []string{}, errors.New("invalid image index")
	}

	result := append(slice[:imageIndex], slice[imageIndex+1:]...)
	if len(result) == 0 {
		return nil, nil
	}
	return result, nil
}
