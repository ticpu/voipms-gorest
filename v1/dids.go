package v1

import (
	"fmt"
	"net/http"
	url2 "net/url"
	"reflect"
)

type SetDidPopRequest struct {
	BaseRequest
	Did string `url:"did"`
	Pop int    `url:"pop"`
}

func (r *SetDidPopRequest) ToURLValues() *url2.Values {
	values := url2.Values{}
	values = toURLValues(reflect.ValueOf(r))
	return &values
}

func (vms *VoIpMsApi) SetDidPop(did string, pop int) (*BaseResponse, error) {
	var (
		err  error
		data *[]byte
	)

	data, err = vms.NewHttpRequest(http.MethodPatch, "setDIDPOP", &SetDidPopRequest{
		Did: did,
		Pop: pop,
	})

	if err != nil {
		return nil, err
	}

	return ParseBaseResponse(data)
}

func (vms *VoIpMsApi) SetDidPopByHostname(did string, popHostname string) (*BaseResponse, error) {
	var (
		err    error
		server *ServerInfo
	)

	if server, err = vms.GetServersInfoForPopHostname(popHostname); err != nil {
		return nil, err
	}

	if server.ServerPOP < 0 {
		return nil, fmt.Errorf("couldn't find POP for %s", popHostname)
	}

	return vms.SetDidPop(did, server.ServerPOP)
}
