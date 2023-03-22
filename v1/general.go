package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	url2 "net/url"
	"reflect"
	"strconv"
)

type GetServersInfo struct {
	BaseRequest
	ServerPop string `url:"server_pop"`
}

func (r *GetServersInfo) ToURLValues() *url2.Values {
	values := url2.Values{}
	values = toURLValues(reflect.ValueOf(r))
	return &values
}

type ServerInfo struct {
	ServerName            string          `json:"server_name"`
	ServerShortname       string          `json:"server_shortname"`
	ServerHostname        string          `json:"server_hostname"`
	ServerIP              string          `json:"server_ip"`
	ServerCountry         string          `json:"server_country"`
	ServerPOP             VoIpMsStringInt `json:"server_pop"`
	ServerRecommended     bool
	ServerRecommendedText string `json:"server_recommended"`
}

type GetServersInfoResponse struct {
	BaseResponse
	Servers []ServerInfo `json:"servers"`
}

func ParseGetServersInfo(data *[]byte) (*GetServersInfoResponse, error) {
	var err error
	response := &GetServersInfoResponse{}
	if err = json.Unmarshal(*data, response); err != nil {
		return nil, err
	}
	return response, nil
}

func (vms *VoIpMsApi) GetServersInfo() (*GetServersInfoResponse, error) {
	var (
		err  error
		data *[]byte
	)
	data, err = vms.NewHttpRequest(http.MethodGet, "getServersInfo", &GetServersInfo{})
	if err != nil {
		return nil, err
	}

	return ParseGetServersInfo(data)
}

func (vms *VoIpMsApi) GetServersInfoForPopHostname(serverPopHostname string) (*ServerInfo, error) {
	var (
		err         error
		serversList *GetServersInfoResponse
	)

	if serversList, err = vms.GetServersInfo(); err == nil {
		for _, server := range serversList.Servers {
			if server.ServerHostname == serverPopHostname {
				return &server, nil
			}
		}
		return nil, fmt.Errorf("couldn't find server %s", serverPopHostname)
	} else {
		return nil, err
	}
}

func (vms *VoIpMsApi) GetServersInfoForPop(pop int) (*ServerInfo, error) {
	var (
		err         error
		data        *[]byte
		serversInfo *GetServersInfoResponse
	)

	data, err = vms.NewHttpRequest(http.MethodGet, "getServersInfo", &GetServersInfo{
		ServerPop: strconv.Itoa(pop),
	})

	if err != nil {
		return nil, err
	}

	if serversInfo, err = ParseGetServersInfo(data); err == nil {
		if len(serversInfo.Servers) == 1 {
			return &(serversInfo.Servers[0]), nil
		} else {
			return nil, fmt.Errorf("couldn't find exactly 1 server with that ID, found %v", len(serversInfo.Servers))
		}
	} else {
		return nil, err
	}
}
