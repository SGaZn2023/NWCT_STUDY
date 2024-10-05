package httprouter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var _ HTTPRouter = &ApisixRouter{}

type ApisixConfig struct {
	Api string
	Key string
}

type ApisixRouter struct {
	conf *ApisixConfig
}

func NewApisixRouter(conf *ApisixConfig) *ApisixRouter {
	return &ApisixRouter{
		conf: conf,
	}
}

func (apisix *ApisixRouter) UpdateRoute(param map[string]interface{}) error {
	cli := &http.Client{
		Timeout: time.Second * 5,
	}

	url := fmt.Sprintf("%s/apisix/admin/routes", apisix.conf.Api)
	body, err := json.Marshal(param)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("X-API-KEY", apisix.conf.Key)

	resp, err := cli.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("invalid http code %d expected created", resp.StatusCode)
	}
	return nil
}

/*
func (apisix *ApisixRouter) GetRoutes() error {
	// TODO implement me
	panic("implement me")
}

func (apisix *ApisixRouter) DeleteRoute(id string) error {
	// TODO implement me
	panic("implement me")
}
*/
