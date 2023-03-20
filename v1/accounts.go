package v1

import (
	"encoding/json"
	"net/http"
	url2 "net/url"
	"reflect"
)

type GetRegistrationStatus struct {
	BaseRequest
	Account string `url:"account"`
}

type RegistrationStatus struct {
	Account           string     `json:"account"`
	ServerName        string     `json:"server_name"`
	ServerShortname   string     `json:"server_shortname"`
	ServerHostname    string     `json:"server_hostname"`
	ServerIP          string     `json:"server_ip"`
	ServerCountry     string     `json:"server_country"`
	ServerPOP         string     `json:"server_pop"`
	RegisterIP        string     `json:"register_ip"`
	RegisterPort      string     `json:"register_port"`
	RegisterNext      VoIpMsTime `json:"register_next"`
	RegisterProtocol  string     `json:"register_protocol"`
	RegisterTransport string     `json:"register_transport"`
	RegisterUseragent string     `json:"register_useragent"`
	Rerouted          int        `json:"rerouted"`
	FromServerPOP     int        `json:"from_server_pop"`
}

type GetRegistrationStatusResponse struct {
	BaseResponse
	Rerouted      int                  `json:"rerouted"`
	FromServerPOP int                  `json:"from_server_pop"`
	Registered    string               `json:"registered"`
	Registrations []RegistrationStatus `json:"registrations"`
}

func (r *GetRegistrationStatus) ToURLValues() *url2.Values {
	values := url2.Values{}
	values = toURLValues(reflect.ValueOf(r))
	return &values
}

func ParseGetRegistrationStatus(data *[]byte) (*GetRegistrationStatusResponse, error) {
	response := &GetRegistrationStatusResponse{}
	if err := json.Unmarshal(*data, response); err != nil {
		return nil, err
	}
	return response, nil
}

func (vms *VoIpMsApi) GetRegistrationStatus(account string) (*GetRegistrationStatusResponse, error) {
	var (
		err  error
		data *[]byte
	)
	data, err = vms.NewHttpRequest(http.MethodGet, "getRegistrationStatus", &GetRegistrationStatus{
		Account: account,
	})
	if err != nil {
		return nil, err
	}

	return ParseGetRegistrationStatus(data)
}
