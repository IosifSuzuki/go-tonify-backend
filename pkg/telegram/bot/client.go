package bot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go-tonify-backend/pkg/curl"
	"go-tonify-backend/pkg/telegram/bot/model"
	"io"
	"log"
	"net/http"
	"strings"
)

type Client interface {
	Execute(model any, method string) error
	ParseResponse(reader io.Reader) (*model.Update, error)
}

const telegramURL = "https://api.telegram.org"

type client struct {
	token string
}

func NewClient(token string) Client {
	return &client{
		token: token,
	}
}

func (c *client) baseURL() string {
	return fmt.Sprintf("%s/bot%s", telegramURL, c.token)
}

func (c *client) url(suffixPath string) string {
	return strings.Join([]string{c.baseURL(), suffixPath}, "/")
}

func (c *client) ParseResponse(reader io.Reader) (*model.Update, error) {
	var update model.Update
	if err := json.NewDecoder(reader).Decode(&update); err != nil {
		return nil, err
	}
	return &update, nil
}

func (c *client) Execute(reqModel any, method string) error {
	bodyBytes, err := json.Marshal(reqModel)
	if err != nil {
		return err
	}
	body := bytes.NewReader(bodyBytes)
	req, err := http.NewRequest(http.MethodPost, c.url(method), body)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return err
	}
	log.Println(curl.GetCurlCommand(req))
	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	var result model.Result
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return err
	}
	if !result.OK {
		return errors.New(result.Description)
	}
	return nil
}
