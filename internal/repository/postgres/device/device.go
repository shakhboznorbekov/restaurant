package device

import (
	"context"
	"github.com/pkg/errors"
	"github.com/restaurant/internal/entity"
	"github.com/restaurant/internal/pkg/repository/postgresql"
	"github.com/restaurant/internal/service/device"
	"github.com/uptrace/bun"
	"time"
)

type Repository struct {
	*postgresql.Database
}

func (r Repository) ChangeDeviceLang(ctx context.Context, data device.ChangeDeviceLang) error {
	var detail entity.Device

	err := r.NewSelect().Model(&detail).Where("device_id = ? AND deleted_at isnull", data.DeviceID).Scan(ctx)
	if err != nil {
		return err
	}

	if data.Lang != nil {
		detail.DeviceLang = data.Lang
	}

	_, err = r.NewUpdate().Model(&detail).Where("id = ? AND deleted_at isnull", detail.ID).Exec(ctx)
	if err != nil {
		return err
	}

	return err
}

func (r Repository) Create(ctx context.Context, data device.Create) (entity.Device, error) {
	timeNow := time.Now()
	var (
		err error
	)

	//-------------------start-------------------------------------------

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return entity.Device{}, errors.Wrap(err, "device create")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	//-------------------check-if-exists---------------------------------

	exists, err := tx.NewSelect().Table("devices").Where("device_id = ?", data.DeviceID).Exists(ctx)
	if err != nil {
		return entity.Device{}, errors.Wrap(err, "select device")
	}

	if exists {
		var deviceLang string
		if data.DeviceLang != nil {
			deviceLang = *data.DeviceLang
		} else {
			deviceLang = r.DefaultLang
		}

		var detail entity.Device
		_, err = tx.NewUpdate().Table("devices").
			Set("device_token = ?", data.DeviceToken).
			Set("user_id = ?", data.UserID).
			Set("device_lang = ?", deviceLang).
			Set("deleted_at = ?", nil).
			Set("is_log_out = false").
			Set("updated_at = ?", time.Now()).
			Where("device_id = ?", data.DeviceID).
			Returning("*").
			Exec(ctx, &detail)
		if err != nil {
			return entity.Device{}, err
		}

		return detail, nil
	} else {
		var deviceLang string
		if data.DeviceLang != nil {
			deviceLang = *data.DeviceLang
		} else {
			deviceLang = r.DefaultLang
		}

		detail := entity.Device{
			Name:        data.Name,
			UserID:      data.UserID,
			DeviceID:    data.DeviceID,
			DeviceLang:  &deviceLang,
			DeviceToken: data.DeviceToken,
			CreatedAt:   &timeNow,
		}

		_, err = tx.NewInsert().Model(&detail).Exec(ctx)
		if err != nil {
			return entity.Device{}, err
		}

		return detail, nil
	}

	//--------------------end-of-process--------------------------------

}

func (r Repository) Update(ctx context.Context, data device.Update) error {
	var detail entity.Device

	err := r.NewSelect().Model(&detail).Where("id = ? AND deleted_at isnull", data.ID).Scan(ctx)
	if err != nil {
		return err
	}

	if data.Name != nil {
		detail.Name = data.Name
	}
	if data.UserID != nil {
		detail.UserID = data.UserID
	}
	if data.DeviceID != nil {
		detail.DeviceID = data.DeviceID
	}
	if data.IsLogOut != nil {
		detail.IsLogOut = data.IsLogOut
	}
	if data.DeviceLang != nil {
		detail.DeviceLang = data.DeviceLang
	}
	if data.DeviceToken != nil {
		detail.DeviceToken = data.DeviceToken
	}

	_, err = r.NewUpdate().Model(&detail).Where("id = ? AND deleted_at isnull", detail.ID).Exec(ctx)
	if err != nil {
		return err
	}

	return err
}

func (r Repository) List(ctx context.Context, filter device.Filter) ([]entity.Device, int, error) {
	var list []entity.Device

	q := r.NewSelect().Model(&list)

	q.Where(" deleted_at isnull ")
	if filter.IsLogOut != nil && *filter.IsLogOut {
		q.WhereGroup("and ", func(query *bun.SelectQuery) *bun.SelectQuery {
			query.Where("is_log_out = ?", *filter.IsLogOut)
			return query
		})
	} else if filter.IsLogOut != nil && !*filter.IsLogOut {
		q.WhereGroup("and ", func(query *bun.SelectQuery) *bun.SelectQuery {
			query.Where("is_log_out = ?", *filter.IsLogOut)
			return query
		})
	}
	if filter.DeviceID != nil {
		q.WhereGroup("and ", func(query *bun.SelectQuery) *bun.SelectQuery {
			query.Where("device_id = ?", *filter.DeviceID)
			return query
		})
	}
	if filter.UserID != nil {
		q.WhereGroup("and ", func(query *bun.SelectQuery) *bun.SelectQuery {
			query.Where("user_id = ?", *filter.UserID)
			return query
		})
	}
	if filter.Limit != nil {
		q.Limit(*filter.Limit)
	}
	if filter.Offset != nil {
		q.Offset(*filter.Offset)
	}

	count, err := q.ScanAndCount(ctx)

	return list, count, err
}

func (r Repository) Detail(ctx context.Context, id int64) (entity.Device, error) {
	var detail entity.Device

	err := r.NewSelect().Model(&detail).Where("id = ? AND deleted_at isnull", id).Scan(ctx)

	return detail, err
}

func (r Repository) Delete(ctx context.Context, id string) error {
	_, err := r.NewUpdate().Table("devices").
		Set("deleted_at = ?", time.Now()).
		Set("is_log_out = true").
		Where("device_id = ?", id).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func NewRepository(DB *postgresql.Database) *Repository {
	return &Repository{DB}
}
