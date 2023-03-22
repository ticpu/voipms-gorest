package v1

import (
	"encoding/json"
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

type GetDidInfoRequest struct {
	BaseRequest
	Client string `url:"client,omitempty"`
	Did    string `url:"did,omitempty"`
}

type DIDInfo struct {
	DID                   string          `json:"did"`
	Description           string          `json:"description"`
	Routing               string          `json:"routing"`
	FailoverBusy          string          `json:"failover_busy"`
	FailoverUnreachable   string          `json:"failover_unreachable"`
	FailoverNoAnswer      string          `json:"failover_noanswer"`
	Voicemail             string          `json:"voicemail"`
	Pop                   string          `json:"pop"`
	Dialtime              string          `json:"dialtime"`
	CNAM                  string          `json:"cnam"`
	E911                  string          `json:"e911"`
	CallerIDPrefix        string          `json:"callerid_prefix"`
	RecordCalls           string          `json:"record_calls"`
	Note                  string          `json:"note"`
	BillingType           string          `json:"billing_type"`
	NextBilling           VoIpMsDate      `json:"next_billing"`
	OrderDate             VoIpMsDateTime  `json:"order_date"`
	ResellerAccount       string          `json:"reseller_account"`
	ResellerNextBilling   string          `json:"reseller_next_billing"`
	ResellerMonthly       string          `json:"reseller_monthly"`
	ResellerMinute        string          `json:"reseller_minute"`
	ResellerSetup         string          `json:"reseller_setup"`
	SMSAvailable          VoIpMsStringInt `json:"sms_available"`
	SMSEnabled            string          `json:"sms_enabled"`
	Transcribe            string          `json:"transcribe"`
	TranscriptionLocale   string          `json:"transcription_locale"`
	TranscriptionEmail    string          `json:"transcription_email"`
	MMSAvailable          VoIpMsStringInt `json:"mms_available"`
	SMSEmail              string          `json:"sms_email"`
	SMSEmailEnabled       string          `json:"sms_email_enabled"`
	SMSForward            string          `json:"sms_forward"`
	SMSForwardEnabled     string          `json:"sms_forward_enabled"`
	SMSURLCallback        string          `json:"sms_url_callback"`
	SMSURLCallbackEnabled string          `json:"sms_url_callback_enabled"`
	SMSURLCallbackRetry   string          `json:"sms_url_callback_retry"`
	SMPPEnabled           string          `json:"smpp_enabled"`
	SMPPURL               string          `json:"smpp_url"`
	SMPPUser              string          `json:"smpp_user"`
	SMPPPass              string          `json:"smpp_pass"`
}

type GetDidInfoResponse struct {
	BaseResponse
	DIDs []DIDInfo `json:"dids"`
}

func ParseGetDidsInfo(data *[]byte) (*GetDidInfoResponse, error) {
	response := &GetDidInfoResponse{}
	if err := json.Unmarshal(*data, response); err != nil {
		return nil, err
	}
	return response, nil
}

func (vms *VoIpMsApi) GetAllDidInfo() (*GetDidInfoResponse, error) {
	var (
		err  error
		data *[]byte
	)

	data, err = vms.NewHttpRequest(http.MethodGet, "getDIDsInfo", &GetDidInfoRequest{})

	if err != nil {
		return nil, err
	}

	return ParseGetDidsInfo(data)
}

func (vms *VoIpMsApi) GetAllClientDidInfo(client string) (*GetDidInfoResponse, error) {
	var (
		err  error
		data *[]byte
	)

	data, err = vms.NewHttpRequest(http.MethodGet, "getDIDsInfo", &GetDidInfoRequest{
		Client: client,
	})

	if err != nil {
		return nil, err
	}

	return ParseGetDidsInfo(data)
}

func (vms *VoIpMsApi) GetDidInfo(client string, did string) (*DIDInfo, error) {
	var (
		err     error
		data    *[]byte
		didInfo *GetDidInfoResponse
	)

	data, err = vms.NewHttpRequest(http.MethodGet, "getDIDsInfo", &GetDidInfoRequest{
		Client: client,
		Did:    did,
	})

	if err != nil {
		return nil, err
	}

	if didInfo, err = ParseGetDidsInfo(data); err != nil {
		return nil, err
	}

	for i := range didInfo.DIDs {
		if didInfo.DIDs[i].DID == did {
			return &didInfo.DIDs[i], nil
		}
	}

	return nil, fmt.Errorf("couldn't find did %s", did)
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
