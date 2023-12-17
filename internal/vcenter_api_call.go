package internal

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
)

type VCenterApiCall struct {
	Conf *Config
}

func NewVmwareApiCallHandler(conf *Config) *VCenterApiCall {
	return &VCenterApiCall{Conf: conf}
}

func client(req *http.Request) (*http.Response, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	return res, err
}

type Session struct {
	Value string `json:"value"`
}

func (v *VCenterApiCall) session() {
	req, err := http.NewRequest("POST", v.Conf.Vcenter.Host+"/rest/com/vmware/cis/session", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	req.SetBasicAuth(v.Conf.Vcenter.Username, v.Conf.Vcenter.Password)
	res, err := client(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer res.Body.Close()

	var session *Session
	if err = json.NewDecoder(res.Body).Decode(&session); err != nil {
		fmt.Println(err)
	}

	fmt.Println(session.Value)
}
