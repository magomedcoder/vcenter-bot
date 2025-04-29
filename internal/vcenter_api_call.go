package internal

import (
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type VCenterApiCall struct {
	Conf *Config
	Db   *sql.DB
}

func NewVmwareApiCallHandler(conf *Config, db *sql.DB) *VCenterApiCall {
	return &VCenterApiCall{Conf: conf, Db: db}
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
		log.Panicf("client-%s", err)
	}

	return res, err
}

type Session struct {
	Value string `json:"value"`
}

func (v *VCenterApiCall) session(userId int64) bool {
	req, err := http.NewRequest("POST", v.Conf.VCenter.Host+"/rest/com/vmware/cis/session", nil)
	if err != nil {
		log.Printf("session-NewRequest-%s", err)
		return false
	}

	var username string
	var password string
	err = v.Db.QueryRow("SELECT username, password FROM users WHERE user_id = ?", userId).Scan(&username, &password)
	if err != nil {
		log.Printf("session-QueryRow-users-%s", err)
	}

	req.SetBasicAuth(username, password)
	res, err := client(req)
	if err != nil {
		log.Printf("session-client-%s", err)
		return false
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Printf("session-client-%s", err)
		}
	}(res.Body)

	var session *Session
	if err = json.NewDecoder(res.Body).Decode(&session); err != nil {
		log.Printf("session-Decode-%s", err)
	}

	if res.StatusCode == 200 {
		_, _err := v.Db.Exec("UPDATE users SET session_id = ? WHERE user_id = ?", session.Value, userId)
		if _err != nil {
			log.Printf("update-users-session_id-%s", err)
			return false
		}
		return true
	}

	return false
}

type List struct {
	Id   string
	Name string
}

func (v *VCenterApiCall) getListVM(userId int64) ([]*List, error) {
	req, err := http.NewRequest("GET", v.Conf.VCenter.Host+"/rest/vcenter/vm", nil)
	if err != nil {
		log.Printf("getListVM-%s", err)
		return nil, err
	}

	var sessionId string
	err = v.Db.QueryRow("SELECT session_id FROM users WHERE user_id = ?", userId).Scan(&sessionId)
	if err != nil {
		log.Printf("getListVM-users-%s", err)
	}

	req.Header.Add("vmware-api-session-id", sessionId)

	res, err := client(req)
	if err != nil {
		log.Printf("getListVM-client-%s", err)
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Printf("getListVM-Close-%s", err)
		}
	}(res.Body)

	if res.StatusCode == 401 {
		log.Printf("StatusCode == 401-%s", err)
		if v.session(userId) {
			return v.getListVM(userId)
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
		log.Printf("getListVM-NewDecoder-%s", err)
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
	Cpu        int    `json:"cpu"`
	Ram        int    `json:"memory"`
	PowerState string `json:"power_state"`
}

func (v *VCenterApiCall) getVM(userId int64, vm string) (*VM, error) {
	req, err := http.NewRequest("GET", v.Conf.VCenter.Host+"/rest/vcenter/vm/"+vm, nil)
	if err != nil {
		log.Printf("getVM-NewRequest-%s", err)
		return nil, err
	}

	var sessionId string
	err = v.Db.QueryRow("SELECT session_id FROM users WHERE user_id = ?", userId).Scan(&sessionId)
	if err != nil {
		log.Printf("getVM-users-user_id-%s", err)
		return nil, err
	}

	req.Header.Add("vmware-api-session-id", sessionId)

	res, err := client(req)
	if err != nil {
		log.Printf("getVM-client-%s", err)
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Printf("getVM-Close-%s", err)
		}
	}(res.Body)

	if res.StatusCode == 401 {
		log.Printf("getVM-StatusCode == 401-%s", err)
		if v.session(userId) {
			return v.getVM(userId, vm)
		}
		return nil, err
	}

	var value *struct {
		Value *struct {
			Name string `json:"name"`
			Cpu  struct {
				Count int `json:"count"`
			} `json:"cpu"`
			Memory struct {
				SizeMiB int `json:"size_MiB"`
			} `json:"memory"`
			PowerState string `json:"power_state"`
		} `json:"value"`
	}

	if err = json.NewDecoder(res.Body).Decode(&value); err != nil {
		log.Printf("getVM-NewDecoder-%s", err)
		return nil, err
	}

	return &VM{
		Name:       value.Value.Name,
		Cpu:        value.Value.Cpu.Count,
		Ram:        value.Value.Memory.SizeMiB,
		PowerState: value.Value.PowerState,
	}, nil
}

func (v *VCenterApiCall) StartVM(userId int64, vm string) bool {
	req, err := http.NewRequest("POST", v.Conf.VCenter.Host+"/rest/vcenter/vm/"+vm+"/power/start", nil)
	if err != nil {
		log.Printf("StartVM-NewRequest-%s", err)
		return false
	}

	var sessionId string
	err = v.Db.QueryRow("SELECT session_id FROM users WHERE user_id = ?", userId).Scan(&sessionId)
	if err != nil {
		log.Printf("StartVM-users-user_id-%s", err)
		return false
	}

	req.Header.Add("vmware-api-session-id", sessionId)

	res, err := client(req)
	if err != nil {
		log.Printf("StartVM-client-%s", err)
		return false
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Printf("StartVM-Close-%s", err)
		}
	}(res.Body)

	if res.StatusCode == 401 {
		log.Printf("StopVM-StatusCode == 401-%s", err)
		if v.session(userId) {
			return v.StartVM(userId, vm)
		}
		return false
	}

	return true
}

func (v *VCenterApiCall) StopVM(userId int64, vm string) bool {
	req, err := http.NewRequest("POST", v.Conf.VCenter.Host+"/rest/vcenter/vm/"+vm+"/power/stop", nil)
	if err != nil {
		log.Printf("StopVM-NewRequest-%s", err)
		return false
	}

	var sessionId string
	err = v.Db.QueryRow("SELECT session_id FROM users WHERE user_id = ?", userId).Scan(&sessionId)
	if err != nil {
		log.Printf("StopVM-QueryRow-user_id-%s", err)
		return false
	}

	req.Header.Add("vmware-api-session-id", sessionId)

	res, err := client(req)
	if err != nil {
		log.Printf("StopVM-client-%s", err)
		return false
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Printf("StopVM-Close-%s", err)
		}
	}(res.Body)

	if res.StatusCode == 401 {
		log.Printf("StopVM-StatusCode == 401-%s", err)
		if v.session(userId) {
			return v.StopVM(userId, vm)
		}
		return false
	}

	return true
}

func (v *VCenterApiCall) RebootVM(userId int64, vm string) bool {
	req, err := http.NewRequest("POST", v.Conf.VCenter.Host+"/rest/vcenter/vm/"+vm+"/power/reset", nil)
	if err != nil {
		log.Printf("RebootVM-NewRequest-%s", err)
		return false
	}

	var sessionId string
	err = v.Db.QueryRow("SELECT session_id FROM users WHERE user_id = ?", userId).Scan(&sessionId)
	if err != nil {
		log.Printf("RebootVM-QueryRow-users-user_id-%s", err)
		return false
	}

	req.Header.Add("vmware-api-session-id", sessionId)

	res, err := client(req)
	if err != nil {
		log.Printf("RebootVM-client-%s", err)
		return false
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Printf("RebootVM-Close-%s", err)
		}
	}(res.Body)

	if res.StatusCode == 401 {
		log.Printf("RebootVM-StatusCode == 401-%s", err)
		if v.session(userId) {
			return v.StopVM(userId, vm)
		}
		return false
	}

	return true
}
