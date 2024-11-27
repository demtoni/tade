package managerapi

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Shadowsocks struct {
	Name   string `json:"name"`
	Method string `json:"method"`
	Plugin string `json:"plugin"`
}

func GetShadowsocks(addr, name string) (map[string]interface{}, error) {
	resp, err := http.Get(addr + name)
	// TODO: send that server is down if no status code
	if err != nil {
		return nil, errors.New("couldn't get configuration: server is down?")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("failed to read response body.")
	}

	data := make(map[string]interface{}, 0)

	log.Println(string(body))
	if err := json.Unmarshal(body, &data); err != nil {
		log.Println(err)
		return nil, errors.New("failed to unmarshal json.")
	}

	return data, nil
}

func DeployShadowsocks(addr, name, method, plugin string) error {
	resp, err := http.PostForm(addr, url.Values{
		"name": {name}, "method": {method}, "plugin": {plugin},
	})
	if err != nil {
		return err
	}

	// TODO: send that server is down if no status code
	if resp.StatusCode != http.StatusCreated {
		return errors.New("couldn't create new service: name is probably already taken.")
	}

	return nil
}

func DeleteShadowsocks(addr, name string) error {
	req, err := http.NewRequest(http.MethodDelete, addr+name, nil)
	if err != nil {
		return errors.New("couldn't build request")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.New("couldn't delete configuration: server is down?")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("couldn't delete service: it doesn't exist")
	}

	return nil
}
