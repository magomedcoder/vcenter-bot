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

func (v *VCenterApiCall) session() bool {
	req, err := http.NewRequest("POST", v.Conf.Vcenter.Host+"/rest/com/vmware/cis/session", nil)
	if err != nil {
		fmt.Println(err)
		return false
	}

	req.SetBasicAuth(v.Conf.Vcenter.Username, v.Conf.Vcenter.Password)
	res, err := client(req)
	if err != nil {
		fmt.Println(err)
		return false
	}

	defer res.Body.Close()

	var session *Session
	if err = json.NewDecoder(res.Body).Decode(&session); err != nil {
		fmt.Println(err)
	}
	if res.StatusCode == 200 {
		WriteTokenToFile(session.Value)
		return true
	}

	return false
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

	readToken, err := readTokenFromFile()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	req.Header.Add("vmware-api-session-id", readToken)
	res, err := client(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode == 401 {
		fmt.Println(err)
		if v.session() {
			return v.getListVM()
		}
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

	readToken, err := readTokenFromFile()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	req.Header.Add("vmware-api-session-id", readToken)
	res, err := client(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode == 401 {
		fmt.Println(err)
		if v.session() {
			return v.getVM(vm)
		}
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

func (v *VCenterApiCall) StartVM(vm string) bool {
	req, err := http.NewRequest("POST", v.Conf.Vcenter.Host+"/rest/vcenter/vm/"+vm+"/power/start", nil)
	if err != nil {
		fmt.Println(err)
		return false
	}

	readToken, err := readTokenFromFile()
	if err != nil {
		fmt.Println(err)
		return false
	}

	req.Header.Add("vmware-api-session-id", readToken)
	res, err := client(req)
	if err != nil {
		fmt.Println(err)
		return false
	}

	defer res.Body.Close()

	if res.StatusCode == 401 {
		fmt.Println(err)
		if v.session() {
			return v.StartVM(vm)
		}
		return false
	}

	return true
}

func (v *VCenterApiCall) StopVM(vm string) bool {
	req, err := http.NewRequest("POST", v.Conf.Vcenter.Host+"/rest/vcenter/vm/"+vm+"/power/stop", nil)
	if err != nil {
		fmt.Println(err)
		return false
	}

	readToken, err := readTokenFromFile()
	if err != nil {
		fmt.Println(err)
		return false
	}

	req.Header.Add("vmware-api-session-id", readToken)
	res, err := client(req)
	if err != nil {
		fmt.Println(err)
		return false
	}

	defer res.Body.Close()

	if res.StatusCode == 401 {
		fmt.Println(err)
		if v.session() {
			return v.StopVM(vm)
		}
		return false
	}

	return true
}
