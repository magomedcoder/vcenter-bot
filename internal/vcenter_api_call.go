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

type List struct {
	Id   string
	Name string
}

func (v *VCenterApiCall) getListVM() ([]*List, error) {
	req, err := http.NewRequest("GET", v.Conf.Vcenter.Host+"/rest/vcenter/vm", nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	res, err := client(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode == 401 {
		return nil, err
	}
	var value *struct {
		Value []*struct {
			VM   string `json:"vm"`
			Name string `json:"name"`
		} `json:"value"`
	}
	if err = json.NewDecoder(res.Body).Decode(&value); err != nil {
		fmt.Println(err)
		return nil, err
	}
	items := make([]*List, 0)
	for _, item := range value.Value {
		items = append(items, &List{Id: item.VM, Name: item.Name})
	}

	return items, nil
}

type VM struct {
	Name       string `json:"name"`
	PowerState string `json:"power_state"`
}

func (v *VCenterApiCall) getVM(vm string) (*VM, error) {
	req, err := http.NewRequest("GET", v.Conf.Vcenter.Host+"/rest/vcenter/vm/"+vm, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	res, err := client(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode == 401 {
		return nil, err
	}
	var value *struct {
		Value *VM `json:"value"`
	}
	if err = json.NewDecoder(res.Body).Decode(&value); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &VM{Name: value.Value.Name, PowerState: value.Value.PowerState}, nil
}
