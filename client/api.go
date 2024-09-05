package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Gewinum/go-df-discord/server"
	"github.com/go-resty/resty/v2"
	"github.com/go-viper/mapstructure/v2"
	"net/http"
)

type Api struct {
	host        string
	accessToken string
}

func NewApi(host, accessToken string) (*Api, error) {
	inst := &Api{host, accessToken}
	if !inst.Test() {
		return nil, errors.New(fmt.Sprintf("can't access %s", host))
	}
	return inst, nil
}

func (a *Api) Test() bool {
	resp, err := getRequest().Get(a.host + "/test")
	if err != nil {
		return false
	}
	return resp.StatusCode() == http.StatusOK
}

func (a *Api) IssueCode(xuid string) (*server.CodeInformation, error) {
	var responsePayload server.Payload
	var response server.CodeInformation
	resp, err := getRequest().SetBody(xuid).Post(a.host + "/codes/issue")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resp.Body(), &responsePayload)
	if err != nil {
		return nil, err
	}
	if responsePayload.Error != nil {
		return nil, errors.New(responsePayload.Error.Message)
	}
	err = mapstructure.Decode(responsePayload, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (a *Api) CheckCode(code string) (*server.CodeInformation, error) {
	var responsePayload server.Payload
	var response server.CodeInformation
	resp, err := getRequest().SetBody(code).Post(a.host + "/codes/check")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resp.Body(), &responsePayload)
	if err != nil {
		return nil, err
	}
	if responsePayload.Error != nil {
		return nil, errors.New(responsePayload.Error.Message)
	}
	err = mapstructure.Decode(responsePayload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (a *Api) RevokeCode(code string) (*server.CodeInformation, error) {
	var responsePayload server.Payload
	var response server.CodeInformation
	resp, err := getRequest().SetBody(code).Post(a.host + "/codes/revoke")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resp.Body(), &responsePayload)
	if err != nil {
		return nil, err
	}
	if responsePayload.Error != nil {
		return nil, errors.New(responsePayload.Error.Message)
	}
	err = mapstructure.Decode(responsePayload, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (a *Api) GetUserByDiscord(discordId string) (*server.User, error) {
	var responsePayload server.Payload
	var response server.User
	resp, err := getRequest().SetPathParams(map[string]string{"{discord}": discordId}).Get(a.host + "/users/discord/{discord}")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resp.Body(), &responsePayload)
	if err != nil {
		return nil, err
	}
	if responsePayload.Error != nil {
		return nil, errors.New(responsePayload.Error.Message)
	}
	err = mapstructure.Decode(responsePayload, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (a *Api) GetUserByXUID(xuid string) (*server.User, error) {
	var responsePayload server.Payload
	var response server.User
	resp, err := getRequest().SetPathParams(map[string]string{"xuid": xuid}).Get(a.host + "/users/xuid/{xuid}")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resp.Body(), &responsePayload)
	if err != nil {
		return nil, err
	}
	if responsePayload.Error != nil {
		return nil, errors.New(responsePayload.Error.Message)
	}
	err = mapstructure.Decode(responsePayload, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func getRequest() *resty.Request {
	return resty.New().R()
}
