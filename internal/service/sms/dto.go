package sms

type Send struct {
	Phone   string
	SmsType int
}

type Check struct {
	Phone string
	Code  string
}
