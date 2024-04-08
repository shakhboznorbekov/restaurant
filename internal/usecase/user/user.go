package user

import (
	"context"
	"github.com/pkg/errors"
	"github.com/restaurant/foundation/web"
	"github.com/restaurant/internal/auth"
	"github.com/restaurant/internal/pkg/file"
	"github.com/restaurant/internal/service/cashier"
	"github.com/restaurant/internal/service/sms"
	"github.com/restaurant/internal/service/user"
	"github.com/restaurant/internal/service/waiter"
	wwt "github.com/restaurant/internal/service/waiter_work_time"
	"github.com/robfig/cron/v3"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

type UseCase struct {
	user              User
	sms               Sms
	waiter            Waiter
	servicePercentage ServicePercentage
	attendance        AttendanceService
	waiterWorkTime    WaiterWorkTimeService
	auth              *auth.Auth
	schedule          *cron.Cron
	cashier           Cashier
}

func NewUseCase(user User, sms Sms, waiter Waiter, servicePercentage ServicePercentage, attendance AttendanceService, waiterWorkTime WaiterWorkTimeService, auth *auth.Auth, schedule *cron.Cron, cashier Cashier) *UseCase {
	return &UseCase{user, sms, waiter, servicePercentage, attendance, waiterWorkTime, auth, schedule, cashier}
}

// #user

// @super-admin

func (uu UseCase) SuperAdminGetUserList(ctx context.Context, filter user.Filter) ([]user.SuperAdminGetList, int, error) {
	fields := make(map[string][]string)
	fields["users"] = []string{"id", "name", "phone", "role"}
	filter.Fields = fields

	list, count, err := uu.user.SuperAdminGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

func (uu UseCase) SuperAdminGetUserDetail(ctx context.Context, id int64) (user.SuperAdminGetDetail, error) {
	var detail user.SuperAdminGetDetail

	data, err := uu.user.SuperAdminGetDetail(ctx, id)
	if err != nil {
		return user.SuperAdminGetDetail{}, err
	}

	detail.ID = data.ID
	detail.Name = data.Name
	detail.Phone = data.Phone
	detail.Gender = data.Gender
	detail.Role = data.Role

	if data.BirthDate != nil {
		birthDate := data.BirthDate.Format("02.01.2006 15:04")
		detail.BirthDate = &birthDate
	}

	return detail, nil
}

func (uu UseCase) SuperAdminCreateUser(ctx context.Context, data user.SuperAdminCreateRequest) (user.SuperAdminCreateResponse, error) {
	exists, err := uu.user.IsPhoneExists(ctx, *data.Phone)
	if err != nil {
		return user.SuperAdminCreateResponse{}, err
	}

	if exists {
		return user.SuperAdminCreateResponse{}, errors.New("phone already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(*data.Password), bcrypt.DefaultCost)
	if err != nil {
		return user.SuperAdminCreateResponse{}, web.NewRequestError(errors.Wrap(err, "hashing password"), http.StatusInternalServerError)
	}
	hashedPassword := string(hash)
	data.Password = &hashedPassword

	return uu.user.SuperAdminCreate(ctx, data)
}

func (uu UseCase) SuperAdminUpdateUser(ctx context.Context, data user.SuperAdminUpdateRequest) error {
	return uu.user.SuperAdminUpdateAll(ctx, data)
}

func (uu UseCase) SuperAdminUpdateUserColumn(ctx context.Context, data user.SuperAdminUpdateRequest) error {
	return uu.user.SuperAdminUpdateColumns(ctx, data)
}

func (uu UseCase) SuperAdminDeleteUser(ctx context.Context, id int64) error {
	return uu.user.SuperAdminDelete(ctx, id)
}

// @client

func (uu UseCase) ClientGetUserMe(ctx context.Context, token string) (user.ClientDetail, error) {
	claims, err := uu.auth.GetTokenData(token)
	if err != nil {
		return user.ClientDetail{}, err
	}

	detail, err := uu.user.ClientGetMe(ctx, claims.UserId)
	if err != nil {
		return user.ClientDetail{}, err
	}

	var clientDetail user.ClientDetail
	clientDetail.ID = detail.ID
	clientDetail.Name = detail.Name
	clientDetail.Gender = detail.Gender
	clientDetail.Phone = detail.Phone
	clientDetail.Role = detail.Role
	if detail.BirthDate != nil {
		birthDate := detail.BirthDate.Format("02.01.2006")
		clientDetail.BirthDate = &birthDate
	}

	return clientDetail, nil
}

func (uu UseCase) ClientUpdateUserColumn(ctx context.Context, data user.ClientUpdateRequest) error {
	return uu.user.ClientUpdateColumn(ctx, data)
}

func (uu UseCase) ClientDeleteMe(ctx context.Context) error {
	return uu.user.ClientDeleteMe(ctx)
}

// #waiter

func (uu UseCase) WaiterCreateComeAttendance(ctx context.Context) error {
	attendanceResp, err := uu.attendance.WaiterCameCreate(ctx)
	if err != nil {
		return err
	}

	action := "ENTER"
	actionTime := time.Now()
	attendanceResp.Action = &action
	attendanceResp.ActionTime = &actionTime

	err = uu.waiterWorkTime.WaiterCreate(ctx, attendanceResp)
	if err != nil {
		return err
	}
	return err
}

func (uu UseCase) WaiterCreateGoneAttendance(ctx context.Context) error {
	attendanceResp, err := uu.attendance.WaiterGoneCreate(ctx)
	if err != nil {
		return err
	}
	action := "EXIT"
	actionTime := time.Now()
	attendanceResp.Action = &action
	attendanceResp.ActionTime = &actionTime

	err = uu.waiterWorkTime.WaiterCreate(ctx, attendanceResp)
	if err != nil {
		return err
	}
	return err
}

func (uu UseCase) WaiterGetListWorkTime(ctx context.Context, filter wwt.ListFilter) ([]wwt.GetListResponse, int, error) {
	return uu.waiterWorkTime.WaiterGetListWorkTime(ctx, filter)
}

// @admin

func (uu UseCase) AdminGetWaiterList(ctx context.Context, filter waiter.Filter) ([]waiter.AdminGetList, int, error) {
	list, count, err := uu.waiter.AdminGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

// @branch

func (uu UseCase) AdminGetCashierList(ctx context.Context, filter cashier.Filter) ([]cashier.AdminGetList, int, error) {
	list, count, err := uu.cashier.AdminGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

func (uu UseCase) AdminGetCashierDetail(ctx context.Context, id int64) (cashier.AdminGetDetail, error) {
	data, err := uu.cashier.AdminGetDetail(ctx, id)
	if err != nil {
		return cashier.AdminGetDetail{}, err
	}

	return data, nil
}

func (uu UseCase) AdminCreateCashier(ctx context.Context, data cashier.AdminCreateRequest) (cashier.AdminCreateResponse, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(*data.Password), bcrypt.DefaultCost)
	if err != nil {
		return cashier.AdminCreateResponse{}, web.NewRequestError(errors.Wrap(err, "hashing password"), http.StatusInternalServerError)
	}
	hashedPassword := string(hash)
	data.Password = &hashedPassword

	exists, err := uu.user.IsPhoneExists(ctx, *data.Phone)
	if err != nil {
		return cashier.AdminCreateResponse{}, err
	}
	if exists {
		return cashier.AdminCreateResponse{}, errors.New("phone already exists")
	}

	if data.Photo != nil {
		imageLink, _, err := file.UploadSingle(data.Photo, "cashier")
		if err != nil {
			return cashier.AdminCreateResponse{}, errors.Wrap(err, "cashier upload")
		}
		data.PhotoLink = &imageLink
	}

	return uu.cashier.AdminCreate(ctx, data)
}

func (uu UseCase) AdminUpdateCashier(ctx context.Context, data cashier.AdminUpdateRequest) error {

	if data.Photo != nil {
		imageLink, _, err := file.UploadSingle(data.Photo, "cashier")
		if err != nil {
			return errors.Wrap(err, "cashier upload")
		}
		data.PhotoLink = &imageLink
	}
	return uu.cashier.AdminUpdateAll(ctx, data)
}

func (uu UseCase) AdminUpdateCashierColumn(ctx context.Context, data cashier.AdminUpdateRequest) error {

	if data.Photo != nil {
		imageLink, _, err := file.UploadSingle(data.Photo, "cashier")
		if err != nil {
			return errors.Wrap(err, "cashier upload")
		}
		data.PhotoLink = &imageLink
	}
	return uu.cashier.AdminUpdateColumns(ctx, data)
}

func (uu UseCase) AdminDeleteCashier(ctx context.Context, id int64) error {
	return uu.cashier.AdminDelete(ctx, id)
}

func (uu UseCase) AdminUpdateCashierStatus(ctx context.Context, id int64, status string) error {
	return uu.cashier.AdminUpdateStatus(ctx, id, status)
}

func (uu UseCase) AdminUpdateCashierPassword(ctx context.Context, data cashier.AdminUpdatePassword) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(*data.Password), bcrypt.DefaultCost)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "hashing password"), http.StatusInternalServerError)
	}

	hashedPassword := string(hash)
	data.Password = &hashedPassword

	return uu.cashier.AdminUpdatePassword(ctx, data)
}

func (uu UseCase) AdminUpdateCashierPhone(ctx context.Context, data cashier.AdminUpdatePhone) error {
	exists, err := uu.user.IsPhoneExists(ctx, data.Phone)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("phone already exists")
	}

	send, err := uu.sms.CashierCheckSMSCode(ctx, sms.Check{
		Phone: data.Phone,
		Code:  data.SMSCode,
	})

	if err != nil {
		return err
	}
	if send {
		return errors.New("can not send sms for cashier")
	}

	return uu.cashier.AdminUpdatePhone(ctx, data)
}

func (uu UseCase) AdminSendSmsCashier(ctx context.Context, request waiter.SendSms) error {
	exists, err := uu.user.IsPhoneExists(ctx, request.Phone)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("phone already exists")
	}

	return uu.sms.SendSMS(ctx, sms.Send{
		Phone:   request.Phone,
		SmsType: 1,
	})
}

// @branch

func (uu UseCase) BranchGetWaiterList(ctx context.Context, filter waiter.Filter) ([]waiter.BranchGetList, int, error) {
	list, count, err := uu.waiter.BranchGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

func (uu UseCase) BranchGetWaiterDetail(ctx context.Context, id int64) (waiter.BranchGetDetail, error) {
	data, err := uu.waiter.BranchGetDetail(ctx, id)
	if err != nil {
		return waiter.BranchGetDetail{}, err
	}

	return data, nil
}

func (uu UseCase) BranchCreateWaiter(ctx context.Context, data waiter.BranchCreateRequest) (waiter.BranchCreateResponse, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(*data.Password), bcrypt.DefaultCost)
	if err != nil {
		return waiter.BranchCreateResponse{}, web.NewRequestError(errors.Wrap(err, "hashing password"), http.StatusInternalServerError)
	}
	hashedPassword := string(hash)
	data.Password = &hashedPassword

	exists, err := uu.user.IsPhoneExists(ctx, *data.Phone)
	if err != nil {
		return waiter.BranchCreateResponse{}, err
	}
	if exists {
		return waiter.BranchCreateResponse{}, errors.New("phone already exists")
	}

	if data.ServicePercentageID != nil {
		_, err = uu.servicePercentage.AdminGetDetail(ctx, *data.ServicePercentageID)
		if err != nil {
			return waiter.BranchCreateResponse{}, errors.New("service_percentage_id is incorrect for this branch!")
		}
	}

	if data.Photo != nil {
		imageLink, _, err := file.UploadSingle(data.Photo, "waiter")
		if err != nil {
			return waiter.BranchCreateResponse{}, errors.Wrap(err, "waiter upload")
		}
		data.PhotoLink = &imageLink
	}

	return uu.waiter.BranchCreate(ctx, data)
}

func (uu UseCase) BranchUpdateWaiter(ctx context.Context, data waiter.BranchUpdateRequest) error {
	if data.ServicePercentageID != nil {
		_, err := uu.servicePercentage.AdminGetDetail(ctx, *data.ServicePercentageID)
		if err != nil {
			return errors.New("service_percentage_id is incorrect for this branch!")
		}
	}
	if data.Photo != nil {
		imageLink, _, err := file.UploadSingle(data.Photo, "waiter")
		if err != nil {
			return errors.Wrap(err, "waiter upload")
		}
		data.PhotoLink = &imageLink
	}
	return uu.waiter.BranchUpdateAll(ctx, data)
}

func (uu UseCase) BranchUpdateWaiterColumn(ctx context.Context, data waiter.BranchUpdateRequest) error {
	if data.ServicePercentageID != nil {
		_, err := uu.servicePercentage.AdminGetDetail(ctx, *data.ServicePercentageID)
		if err != nil {
			return errors.New("service_percentage_id is incorrect for this branch!")
		}
	}
	if data.Photo != nil {
		imageLink, _, err := file.UploadSingle(data.Photo, "waiter")
		if err != nil {
			return errors.Wrap(err, "waiter upload")
		}
		data.PhotoLink = &imageLink
	}
	return uu.waiter.BranchUpdateColumns(ctx, data)
}

func (uu UseCase) BranchDeleteWaiter(ctx context.Context, id int64) error {
	return uu.waiter.BranchDelete(ctx, id)
}

func (uu UseCase) BranchUpdateWaiterPassword(ctx context.Context, data waiter.BranchUpdatePassword) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(*data.Password), bcrypt.DefaultCost)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "hashing password"), http.StatusInternalServerError)
	}

	hashedPassword := string(hash)
	data.Password = &hashedPassword

	return uu.waiter.UpdatePassword(ctx, data)
}

func (uu UseCase) BranchUpdateWaiterPhone(ctx context.Context, data waiter.BranchUpdatePhone) error {
	exists, err := uu.user.IsPhoneExists(ctx, data.Phone)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("phone already exists")
	}

	send, err := uu.sms.WaiterCheckSMSCode(ctx, sms.Check{
		Phone: data.Phone,
		Code:  data.SMSCode,
	})

	if err != nil {
		return err
	}
	if send {
		return errors.New("can not send sms for waiter")
	}

	return uu.waiter.UpdatePhone(ctx, data)
}

func (uu UseCase) BranchSendSmsWaiter(ctx context.Context, request waiter.SendSms) error {
	exists, err := uu.user.IsPhoneExists(ctx, request.Phone)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("phone already exists")
	}

	return uu.sms.SendSMS(ctx, sms.Send{
		Phone:   request.Phone,
		SmsType: 1,
	})
}

func (uu UseCase) BranchUpdateWaiterStatus(ctx context.Context, id int64, status string) error {
	return uu.waiter.BranchUpdateStatus(ctx, id, status)
}

func (uu UseCase) BranchGetListWaiterWorkTime(ctx context.Context, filter wwt.BranchFilter) ([]wwt.BranchGetListResponse, int, error) {
	return uu.waiterWorkTime.BranchGetListWaiterWorkTime(ctx, filter)
}

func (uu UseCase) BranchGetDetailWaiterWorkTime(ctx context.Context, filter wwt.ListFilter) ([]wwt.GetListResponse, int, error) {
	return uu.waiterWorkTime.BranchGetDetailWaiterWorkTime(ctx, filter)
}

// @branch

func (uu UseCase) BranchGetCashierList(ctx context.Context, filter cashier.Filter) ([]cashier.BranchGetList, int, error) {
	list, count, err := uu.cashier.BranchGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

func (uu UseCase) BranchGetCashierDetail(ctx context.Context, id int64) (cashier.BranchGetDetail, error) {
	data, err := uu.cashier.BranchGetDetail(ctx, id)
	if err != nil {
		return cashier.BranchGetDetail{}, err
	}

	return data, nil
}

func (uu UseCase) BranchCreateCashier(ctx context.Context, data cashier.BranchCreateRequest) (cashier.BranchCreateResponse, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(*data.Password), bcrypt.DefaultCost)
	if err != nil {
		return cashier.BranchCreateResponse{}, web.NewRequestError(errors.Wrap(err, "hashing password"), http.StatusInternalServerError)
	}
	hashedPassword := string(hash)
	data.Password = &hashedPassword

	exists, err := uu.user.IsPhoneExists(ctx, *data.Phone)
	if err != nil {
		return cashier.BranchCreateResponse{}, err
	}
	if exists {
		return cashier.BranchCreateResponse{}, errors.New("phone already exists")
	}

	if data.Photo != nil {
		imageLink, _, err := file.UploadSingle(data.Photo, "cashier")
		if err != nil {
			return cashier.BranchCreateResponse{}, errors.Wrap(err, "cashier upload")
		}
		data.PhotoLink = &imageLink
	}

	return uu.cashier.BranchCreate(ctx, data)
}

func (uu UseCase) BranchUpdateCashier(ctx context.Context, data cashier.BranchUpdateRequest) error {

	if data.Photo != nil {
		imageLink, _, err := file.UploadSingle(data.Photo, "cashier")
		if err != nil {
			return errors.Wrap(err, "cashier upload")
		}
		data.PhotoLink = &imageLink
	}
	return uu.cashier.BranchUpdateAll(ctx, data)
}

func (uu UseCase) BranchUpdateCashierColumn(ctx context.Context, data cashier.BranchUpdateRequest) error {

	if data.Photo != nil {
		imageLink, _, err := file.UploadSingle(data.Photo, "cashier")
		if err != nil {
			return errors.Wrap(err, "cashier upload")
		}
		data.PhotoLink = &imageLink
	}
	return uu.cashier.BranchUpdateColumns(ctx, data)
}

func (uu UseCase) BranchDeleteCashier(ctx context.Context, id int64) error {
	return uu.cashier.BranchDelete(ctx, id)
}

func (uu UseCase) BranchUpdateCashierStatus(ctx context.Context, id int64, status string) error {
	return uu.cashier.BranchUpdateStatus(ctx, id, status)
}

func (uu UseCase) BranchUpdateCashierPassword(ctx context.Context, data cashier.BranchUpdatePassword) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(*data.Password), bcrypt.DefaultCost)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "hashing password"), http.StatusInternalServerError)
	}

	hashedPassword := string(hash)
	data.Password = &hashedPassword

	return uu.cashier.UpdatePassword(ctx, data)
}

func (uu UseCase) BranchUpdateCashierPhone(ctx context.Context, data cashier.BranchUpdatePhone) error {
	exists, err := uu.user.IsPhoneExists(ctx, data.Phone)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("phone already exists")
	}

	send, err := uu.sms.CashierCheckSMSCode(ctx, sms.Check{
		Phone: data.Phone,
		Code:  data.SMSCode,
	})

	if err != nil {
		return err
	}
	if send {
		return errors.New("can not send sms for cashier")
	}

	return uu.cashier.UpdatePhone(ctx, data)
}

func (uu UseCase) BranchSendSmsCashier(ctx context.Context, request waiter.SendSms) error {
	exists, err := uu.user.IsPhoneExists(ctx, request.Phone)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("phone already exists")
	}

	return uu.sms.SendSMS(ctx, sms.Send{
		Phone:   request.Phone,
		SmsType: 1,
	})
}

// @cashier

func (uu UseCase) CashierGetListWaiterWorkTime(ctx context.Context, filter wwt.BranchFilter) ([]wwt.BranchGetListResponse, int, error) {
	return uu.waiterWorkTime.CashierGetListWaiterWorkTime(ctx, filter)
}

func (uu UseCase) CashierGetDetailWaiterWorkTime(ctx context.Context, filter wwt.ListFilter) ([]wwt.GetListResponse, int, error) {
	return uu.waiterWorkTime.CashierGetDetailWaiterWorkTime(ctx, filter)
}

func (uu UseCase) CashierGetMe(ctx context.Context) (*user.CashierGetMeResponse, error) {
	return uu.user.CashierGetMe(ctx)
}

func (uu UseCase) CashierGetWaiterList(ctx context.Context, filter waiter.Filter) ([]waiter.CashierGetList, int, error) {
	list, count, err := uu.waiter.CashierGetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

// @waiter

func (uu UseCase) WaiterGetMe(ctx context.Context) (*waiter.GetMeResponse, error) {
	return uu.waiter.WaiterGetMe(ctx)
}

func (uu UseCase) WaiterGetPersonalInfo(ctx context.Context) (*waiter.GetPersonalInfoResponse, error) {
	return uu.waiter.WaiterGetPersonalInfo(ctx)
}

func (uu UseCase) WaiterUpdatePhoto(ctx context.Context, request waiter.WaiterPhotoUpdateRequest) error {
	if request.Photo != nil {
		imageLink, _, err := file.UploadSingle(request.Photo, "waiter")
		if err != nil {
			return errors.Wrap(err, "waiter upload")
		}
		request.PhotoLink = &imageLink
	}
	return uu.waiter.WaiterUpdatePhoto(ctx, request)
}

func (uu UseCase) WaiterGetActivityStatistics(ctx context.Context) (*waiter.GetActivityStatistics, error) {
	return uu.waiter.WaiterGetActivityStatistics(ctx)
}

func (uu UseCase) WaiterGetWeeklyActivityStatistics(ctx context.Context, filter waiter.EarnedMoneyFilter) (*waiter.GetEarnedMoneyStatistics, error) {
	return uu.waiter.WaiterGetWeeklyActivityStatistics(ctx, filter)
}

func (uu UseCase) WaiterGetWeeklyAcceptedOrdersStatistics(ctx context.Context, filter waiter.EarnedMoneyFilter) (*waiter.GetAcceptedOrdersStatistics, error) {
	return uu.waiter.WaiterGetWeeklyAcceptedOrdersStatistics(ctx, filter)
}

func (uu UseCase) WaiterGetWeeklyRatingStatistics(ctx context.Context, filter waiter.Filter) ([]waiter.GetWeeklyRating, error) {
	return uu.waiter.WaiterGetWeeklyRatingStatistics(ctx, filter)
}

// CalculateWaitersKPI : is run by robfig/cron in each day at 7:00 a.m
func (uu UseCase) CalculateWaitersKPI(ctx context.Context) error {
	_, err := uu.schedule.AddFunc("0 7 * * *", func() {
		if err := uu.waiter.CalculateWaitersKPI(ctx); err != nil {
			log.Printf("error>> daily cron schedule: waiter kpi error [%s]", err)
			return
		}
	})
	if err != nil {
		return err
	}

	uu.schedule.Start()

	return nil
}

// general get me

func (uu UseCase) GetMe(ctx context.Context, token string) (*user.GetMeResponse, error) {

	claims, err := uu.auth.GetTokenData(token)
	if err != nil {
		return nil, err
	}

	return uu.user.GetMe(ctx, claims.UserId)
}

//

// @cashier #waiter

func (uu UseCase) CashierGetWaiterLists(ctx context.Context, filter waiter.Filter) ([]waiter.CashierGetLists, int, error) {
	list, count, err := uu.waiter.CashierGetLists(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, err
}

func (uu UseCase) CashierGetWaiterDetails(ctx context.Context, id int64) (waiter.CashierGetDetails, error) {
	data, err := uu.waiter.CashierGetDetails(ctx, id)
	if err != nil {
		return waiter.CashierGetDetails{}, err
	}

	return data, nil
}

func (uu UseCase) CashierCreateWaiter(ctx context.Context, data waiter.CashierCreateRequest) (waiter.CashierCreateResponse, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(*data.Password), bcrypt.DefaultCost)
	if err != nil {
		return waiter.CashierCreateResponse{}, web.NewRequestError(errors.Wrap(err, "hashing password"), http.StatusInternalServerError)
	}
	hashedPassword := string(hash)
	data.Password = &hashedPassword

	exists, err := uu.user.IsPhoneExists(ctx, *data.Phone)
	if err != nil {
		return waiter.CashierCreateResponse{}, err
	}
	if exists {
		return waiter.CashierCreateResponse{}, errors.New("phone already exists")
	}

	if data.ServicePercentageID != nil {
		_, err = uu.servicePercentage.AdminGetDetail(ctx, *data.ServicePercentageID)
		if err != nil {
			return waiter.CashierCreateResponse{}, errors.New("service_percentage_id is incorrect for this branch!")
		}
	}

	if data.Photo != nil {
		imageLink, _, err := file.UploadSingle(data.Photo, "waiter")
		if err != nil {
			return waiter.CashierCreateResponse{}, errors.Wrap(err, "waiter upload")
		}
		data.PhotoLink = &imageLink
	}

	return uu.waiter.CashierCreate(ctx, data)
}

func (uu UseCase) CashierUpdateWaiter(ctx context.Context, data waiter.CashierUpdateRequest) error {
	if data.ServicePercentageID != nil {
		_, err := uu.servicePercentage.AdminGetDetail(ctx, *data.ServicePercentageID)
		if err != nil {
			return errors.New("service_percentage_id is incorrect for this branch!")
		}
	}
	if data.Photo != nil {
		imageLink, _, err := file.UploadSingle(data.Photo, "waiter")
		if err != nil {
			return errors.Wrap(err, "waiter upload")
		}
		data.PhotoLink = &imageLink
	}
	return uu.waiter.CashierUpdateAll(ctx, data)
}

func (uu UseCase) CashierUpdateWaiterColumn(ctx context.Context, data waiter.CashierUpdateRequest) error {
	if data.ServicePercentageID != nil {
		_, err := uu.servicePercentage.AdminGetDetail(ctx, *data.ServicePercentageID)
		if err != nil {
			return errors.New("service_percentage_id is incorrect for this branch!")
		}
	}
	if data.Photo != nil {
		imageLink, _, err := file.UploadSingle(data.Photo, "waiter")
		if err != nil {
			return errors.Wrap(err, "waiter upload")
		}
		data.PhotoLink = &imageLink
	}
	return uu.waiter.CashierUpdateColumns(ctx, data)
}

func (uu UseCase) CashierDeleteWaiter(ctx context.Context, id int64) error {
	return uu.waiter.CashierDelete(ctx, id)
}

func (uu UseCase) CashierUpdateWaiterPassword(ctx context.Context, data waiter.CashierUpdatePassword) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(*data.Password), bcrypt.DefaultCost)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "hashing password"), http.StatusInternalServerError)
	}

	hashedPassword := string(hash)
	data.Password = &hashedPassword

	return uu.waiter.CashierUpdatePassword(ctx, data)
}

func (uu UseCase) CashierUpdateWaiterPhone(ctx context.Context, data waiter.CashierUpdatePhone) error {
	exists, err := uu.user.IsPhoneExists(ctx, data.Phone)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("phone already exists")
	}

	send, err := uu.sms.WaiterCheckSMSCode(ctx, sms.Check{
		Phone: data.Phone,
		Code:  data.SMSCode,
	})

	if err != nil {
		return err
	}
	if send {
		return errors.New("can not send sms for waiter")
	}

	return uu.waiter.CashierUpdatePhone(ctx, data)
}

func (uu UseCase) CashierSendSmsWaiter(ctx context.Context, request waiter.SendSms) error {
	exists, err := uu.user.IsPhoneExists(ctx, request.Phone)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("phone already exists")
	}

	return uu.sms.SendSMS(ctx, sms.Send{
		Phone:   request.Phone,
		SmsType: 1,
	})
}

func (uu UseCase) CashierUpdateWaiterStatus(ctx context.Context, id int64, status string) error {
	return uu.waiter.CashierUpdateStatus(ctx, id, status)
}
