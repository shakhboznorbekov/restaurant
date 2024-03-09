package device

type Filter struct {
	Limit    *int
	Offset   *int
	UserID   *int64
	DeviceID *string
	IsLogOut *bool
}

type List struct {
	ID          int64   `json:"id"`
	Name        *string `json:"name"`
	UserID      *int64  `json:"user_id"`
	DeviceID    *string `json:"device_id"`
	IsLogOut    *bool   `json:"is_log_out"`
	DeviceToken *string `json:"device_token"`
}

type Detail struct {
	ID          int64   `json:"id"`
	Name        *string `json:"name"`
	UserID      *int64  `json:"user_id"`
	DeviceID    *string `json:"device_id"`
	IsLogOut    *bool   `json:"is_log_out"`
	DeviceToken *string `json:"device_token"`
}

type Create struct {
	Name        *string `json:"name" form:"name"`
	UserID      *int64  `json:"-" form:"-"`
	DeviceID    *string `json:"device_id" form:"device_id"`
	DeviceLang  *string `json:"device_lang" form:"device_lang"`
	DeviceToken *string `json:"device_token" form:"device_token"`
}

type Update struct {
	ID          int64   `json:"id" form:"id"`
	Name        *string `json:"name" form:"name"`
	UserID      *int64  `json:"user_id" form:"user_id"`
	DeviceID    *string `json:"device_id" form:"device_id"`
	IsLogOut    *bool   `json:"is_log_out" form:"is_log_out"`
	DeviceLang  *string `json:"device_lang" form:"device_lang"`
	DeviceToken *string `json:"device_token" form:"device_token"`
}

type ChangeDeviceLang struct {
	DeviceID *string `json:"device_id" form:"device_id"`
	Lang     *string `json:"lang" form:"lang"`
}
