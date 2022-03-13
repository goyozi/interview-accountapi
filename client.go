package interviewaccountapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type AccountsClient struct {
	url    string
	client *http.Client
}

func NewAccountsClient(baseUrl string) *AccountsClient {
	return &AccountsClient{baseUrl + "/v1/organisation/accounts", &http.Client{}}
}

func (c AccountsClient) Create(data *AccountData) (*AccountData, error) {
	payload, err := json.Marshal(&AccountPayload{data})
	if err != nil {
		return nil, err
	}

	req, err := c.newRequest("POST", "", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	res, err := c.requestAccount(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c AccountsClient) Fetch(id string) (*AccountData, error) {
	req, err := c.newRequest("GET", "/"+id, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.requestAccount(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c AccountsClient) Delete(account *AccountData) error {
	req, err := c.newRequest("DELETE", "/"+account.ID+"?version="+strconv.FormatInt(*account.Version, 10), nil)
	if err != nil {
		return err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode == 404 {
		return errors.New("record " + account.ID + " does not exist")
	}

	return checkStatus(res)
}

func (c AccountsClient) newRequest(method string, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, c.url+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/vnd.api+json")
	req.Header.Add("Date", time.Now().GoString())

	return req, nil
}

func (c AccountsClient) requestAccount(req *http.Request) (*AccountData, error) {
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	err = checkStatus(res)
	if err != nil {
		return nil, err
	}

	return readAccountData(res)
}

func checkStatus(res *http.Response) error {
	if res.StatusCode < 400 {
		return nil
	}

	var errorBody ErrorData
	err := readBody(res, &errorBody)
	if err != nil {
		return err
	}

	return errors.New(errorBody.ErrorMessage)
}

func readAccountData(res *http.Response) (*AccountData, error) {
	var account AccountPayload
	err := readBody(res, &account)
	if err != nil {
		return nil, err
	}

	return account.Data, nil
}

func readBody(res *http.Response, v interface{}) error {
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, v)
}
