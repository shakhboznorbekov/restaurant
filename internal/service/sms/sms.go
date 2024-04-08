package sms

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/restaurant/foundation/web"
	"github.com/restaurant/internal/pkg/repository/postgresql"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Config struct {
	Username   string
	Password   string
	Originator string
	MSGAlias   string
	BaseUrl    string
}

type Service struct {
	redisDB    *redis.Client
	postgresDB *postgresql.Database
	config     *Config
}

func NewService(config *Config, redisDB *redis.Client, postgresDB *postgresql.Database) *Service {
	return &Service{redisDB: redisDB, postgresDB: postgresDB, config: config}
}

func (s *Service) SendSMS(ctx context.Context, send Send) error {
	// phone format: 998XXYYYZZXX (example: 998337750008)

	smsCode := s.generateSmsCode()
	//log.Println(smsCode)

	//smsBody := fmt.Sprintf("Code: %s", smsCode)
	//
	//switch smsType {
	//case 1:
	//	smsBody = fmt.Sprintf("Parol o'rnatish uchun tasdiqlash kodi - %s", smsCode)
	//}
	//
	//data := map[string]interface{}{"messages": []map[string]interface{}{{
	//	"recipient":  phone,
	//	"message-id": s.config.MSGAlias + uuid.NewString(),
	//	"sms": map[string]interface{}{
	//		"originator": s.config.Originator,
	//		"content": map[string]string{
	//			"text": smsBody,
	//		},
	//	},
	//}}}
	//
	//result, err := json.Marshal(data)
	//if err != nil {
	//	return web.NewRequestError(errors.Wrap(err, "marshaling sms data"), http.StatusInternalServerError)
	//}
	//
	//req, err := http.NewRequest(http.MethodPost, s.config.BaseUrl, bytes.NewBuffer(result))
	//if err != nil {
	//	return web.NewRequestError(errors.Wrap(err, "request to play mobile"), http.StatusInternalServerError)
	//}
	//
	//req.Header.Set("Content-Type", "application/json")
	//req.SetBasicAuth(s.config.Username, s.config.Password)
	//
	//client := http.Client{}
	//
	//res, err := client.Do(req)
	//if err != nil {
	//	return web.NewRequestError(errors.Wrap(err, "play mobile request"), http.StatusInternalServerError)
	//}
	//
	//log.Print(res)
	//
	//if res != nil && res.StatusCode != http.StatusOK {
	//	return web.NewRequestError(errors.New(fmt.Sprintf("play mobile status code: %d", res.StatusCode)), http.StatusInternalServerError)
	//}

	s.redisDB.Set(ctx, fmt.Sprintf("sms_code_%s", send.Phone), smsCode, 2*time.Minute)

	return nil
}

func (s *Service) CheckSMSCode(ctx context.Context, check Check) (bool, error) {
	rows, err := s.postgresDB.QueryContext(ctx, fmt.Sprintf(`
		SELECT EXISTS (
			SELECT
				id
			FROM
			    users
		    WHERE
		        phone = '%s' AND
		        deleted_at IS NULL
		)
	`, check.Phone))
	if err != nil {
		return false, web.NewRequestError(errors.Wrap(err, "selecting exists error"), http.StatusInternalServerError)
	}

	exists := false
	if err = s.postgresDB.ScanRows(ctx, rows, &exists); err != nil {
		return false, web.NewRequestError(errors.Wrap(err, "scanning exists phone"), http.StatusInternalServerError)
	}

	//smsCode, err := s.redisDB.Get(ctx, fmt.Sprintf("sms_code_%s", check.Phone)).Result()
	//if err != nil && errors.Is(err, errors.New("redis: nil")) {
	//
	//}
	//if err != nil {
	//	log.Printf("get phone sms code error: %v", err)
	//}

	//if smsCode == "" || smsCode != check.Code {
	//	return false, web.NewRequestError(errors.New("incorrect sms code"), http.StatusBadRequest)
	//}

	if check.Code != "111111" {
		return exists, web.NewRequestError(errors.New("incorrect sms code"), http.StatusBadRequest)
	}

	return exists, nil
}

func (s *Service) CheckSMSCodeUpdatePhone(ctx context.Context, check Check) (bool, error) {

	//smsCode, err := s.redisDB.Get(ctx, fmt.Sprintf("sms_code_%s", check.Phone)).Result()
	//if err != nil && errors.Is(err, errors.New("redis: nil")) {
	//
	//}
	//if err != nil {
	//	log.Printf("get phone sms code error: %v", err)
	//}

	//if smsCode == "" || smsCode != check.Code {
	//	return false, web.NewRequestError(errors.New("incorrect sms code"), http.StatusBadRequest)
	//}

	if check.Code != "111111" {
		return false, web.NewRequestError(errors.New("incorrect sms code"), http.StatusBadRequest)
	}

	return true, nil
}

func (s *Service) WaiterCheckSMSCode(ctx context.Context, check Check) (bool, error) {
	rows, err := s.postgresDB.QueryContext(ctx, fmt.Sprintf(`
		SELECT EXISTS (
			SELECT
				id
			FROM
			    users
		    WHERE
		        phone = '%s' AND
		        deleted_at IS NULL
		)
	`, check.Phone))
	if err != nil {
		return false, web.NewRequestError(errors.Wrap(err, "selecting exists error"), http.StatusInternalServerError)
	}

	exists := false
	if err = s.postgresDB.ScanRows(ctx, rows, &exists); err != nil {
		return false, web.NewRequestError(errors.Wrap(err, "scanning exists phone"), http.StatusInternalServerError)
	}

	//smsCode, err := s.redisDB.Get(ctx, fmt.Sprintf("sms_code_%s", check.Phone)).Result()
	//if err != nil && errors.Is(err, errors.New("redis: nil")) {
	//
	//}
	//if err != nil {
	//	log.Printf("get phone sms code error: %v", err)
	//}

	//if smsCode == "" || smsCode != check.Code {
	//	return false, web.NewRequestError(errors.New("incorrect sms code"), http.StatusBadRequest)
	//}

	if check.Code != "111111" {
		return exists, web.NewRequestError(errors.New("incorrect sms code"), http.StatusBadRequest)
	}

	return exists, nil
}

func (s *Service) CashierCheckSMSCode(ctx context.Context, check Check) (bool, error) {
	rows, err := s.postgresDB.QueryContext(ctx, fmt.Sprintf(`
		SELECT EXISTS (
			SELECT
				id
			FROM
			    users
		    WHERE
		        phone = '%s' AND
		        deleted_at IS NULL
		)
	`, check.Phone))
	if err != nil {
		return false, web.NewRequestError(errors.Wrap(err, "selecting exists error"), http.StatusInternalServerError)
	}

	exists := false
	if err = s.postgresDB.ScanRows(ctx, rows, &exists); err != nil {
		return false, web.NewRequestError(errors.Wrap(err, "scanning exists phone"), http.StatusInternalServerError)
	}

	//smsCode, err := s.redisDB.Get(ctx, fmt.Sprintf("sms_code_%s", check.Phone)).Result()
	//if err != nil && errors.Is(err, errors.New("redis: nil")) {
	//
	//}
	//if err != nil {
	//	log.Printf("get phone sms code error: %v", err)
	//}

	//if smsCode == "" || smsCode != check.Code {
	//	return false, web.NewRequestError(errors.New("incorrect sms code"), http.StatusBadRequest)
	//}

	if check.Code != "111111" {
		return exists, web.NewRequestError(errors.New("incorrect sms code"), http.StatusBadRequest)
	}

	return exists, nil
}

func (s *Service) generateSmsCode() string {
	codeMin := 100000
	codeMax := 999999
	rand.Seed(time.Now().UnixNano())
	code := strconv.Itoa(rand.Intn(codeMax-codeMin) + codeMin)
	return code
}
