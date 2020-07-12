package xmatters

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"net/url"
)

type xmattersClient struct {
	url    string
	token  string
	client *http.Client
}

type GroupRosterResponse struct {
	Count int               `json:"count"`
	Total int               `json:"total"`
	Data  []GroupRosterBody `json:"data"`
	Links Links             `json:"links"`
}

type Links struct {
	Self string `json:"self"`
}

type Group struct {
	ID            string `json:"id"`
	TargetName    string `json:"targetName"`
	RecipientType string `json:"recipientType"`
	Links         Links  `json:"links"`
}

type Member struct {
	ID            string `json:"id"`
	TargetName    string `json:"targetName"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	RecipientType string `json:"recipientType"`
	Links         Links  `json:"links"`
}

type GroupRosterBody struct {
	Group  Group  `json:"group"`
	Member Member `json:"member,omitempty"`
}

func NewXmattersClient(url, token string) xmattersClient {
	return xmattersClient{
		url:    url,
		token:  token,
		client: &http.Client{},
	}
}

func (self xmattersClient) GetGroupRoster(group string) (*GroupRosterResponse, error) {
	uri := fmt.Sprintf("%s/groups/%s/members?limit=1000", self.url, url.PathEscape(group))
	req, err := http.NewRequest("GET", uri, nil)

	if err != nil {
		return nil, err
	}

	data, err := self.sendRequest(req)
	if err != nil {
		return nil, err
	}

	body := &GroupRosterResponse{}
	if err := json.Unmarshal(data, body); err != nil {
		log.Error().Err(err).Str("uri", uri).Str("body", string(data)).Msg("json parse failed")
		return nil, err
	}

	return body, nil
}

// TODO replace token with username + password (basic auth)
func (self xmattersClient) sendRequest(req *http.Request) ([]byte, error) {
	sublogger := log.With().Str("uri", req.URL.String()).Logger()
	req.Header.Add("Authorization", "Basic "+self.token)
	resp, err := self.client.Do(req)

	if err != nil {
		sublogger.Error().Err(err).Msg("xMatters request failed")
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		sublogger.Error().Err(err).Msg("Error reading xMatters response body")
		return nil, err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		err := errors.New("401 response, invalid credentials")
		sublogger.Error().Err(err).Str("body", string(data)).Send()
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		err := errors.New("NOK status code: " + string(resp.StatusCode))
		sublogger.Error().Err(err).Str("body", string(data)).Send()
		return nil, err
	}

	return data, err
}
