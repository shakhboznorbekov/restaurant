package auth

import "github.com/restaurant/internal/service/device"

// #sign-in

type SignInRequest struct {
	Phone    string        `json:"phone" form:"phone"`
	Password string        `json:"password" form:"password"`
	Device   device.Create `json:"device" form:"device"`
}

type SignInWaiter struct {
	Phone    string        `json:"phone" form:"phone"`
	Password string        `json:"password" form:"password"`
	Device   device.Create `json:"device" form:"device"`
}

type SignInClient struct {
	Phone   string        `json:"phone" form:"phone"`
	SMSCode string        `json:"sms_code" form:"sms_code"`
	Device  device.Create `json:"device" form:"device"`
}

// #sign-up

type SignUpRequest struct {
	Name      *string `json:"name" form:"name"`
	BirthDate *string `json:"birth_date" form:"birth_date"`
	Gender    *string `json:"gender" form:"gender"`
	Token     string  `json:"-" from:"-"`
}

// #others

type SendSms struct {
	Phone string `json:"phone" form:"phone"`
}

type ClaimsAuth struct {
	Roles        string
	ID           int64
	RestaurantID *int64
	BranchID     *int64
}

type UpdatePhone struct {
	Phone   string        `json:"phone" form:"phone"`
	SMSCode string        `json:"sms_code" form:"sms_code"`
	Device  device.Create `json:"device" form:"device"`
}
