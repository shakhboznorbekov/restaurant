package fcm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
)

type ServiceFCM struct {
	client  *http.Client
	baseUrl string
	token   string
}

func (s *ServiceFCM) SendCloudMessage(model CloudMessage) (string, error) {
	model.ContentAvailable = true
	data, err := json.Marshal(model)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodPost, s.baseUrl, bytes.NewBuffer(data))
	if err != nil {
		return "", errors.Wrap(err, "new request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", s.token)

	res, err := s.client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "doing request")
	}

	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("error reading response body: %v\n", err)
	}

	log.Println(model.To, " -> ", responseBody)

	var response FCMResponse
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return "", err
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return "", errors.New("cannot send message")
	}

	if response.Results[0].Error != nil {
		return *response.Results[0].Error, nil
	} else {
		return response.Results[0].MessageID, nil
	}

	return "", nil
}

func NewFCMService(cfg *ConfigFCM) *ServiceFCM {
	token := fmt.Sprintf("key=%s", cfg.WebApi)
	return &ServiceFCM{
		client:  &http.Client{},
		baseUrl: cfg.Link,
		token:   token,
	}
}
