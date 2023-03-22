package v1

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	url2 "net/url"
	"reflect"
	"strconv"
	"time"
)

const (
	RestAPIURL           = "https://voip.ms/api/v1/rest.php"
	voipmsDateTimeFormat = "2006-01-02 15:04:05"
	voipmsDateFormat     = "2006-01-02"
)

type VoIpMsApi struct {
	ApiUsername string
	ApiPassword string
	ApiUrl      string
	ApiTimeout  time.Duration
}

type VoIpMsDateTime struct {
	time.Time
}

func (vmsDateTime *VoIpMsDateTime) UnmarshalJSON(b []byte) (err error) {
	s := string(b)
	s = s[1 : len(s)-1]
	vmsDateTime.Time, err = time.Parse(voipmsDateTimeFormat, s)
	return
}

type VoIpMsDate struct {
	time.Time
}

func (vmsTime *VoIpMsDate) UnmarshalJSON(b []byte) (err error) {
	s := string(b)
	s = s[1 : len(s)-1]
	vmsTime.Time, err = time.Parse(voipmsDateFormat, s)
	return
}

type VoIpMsStringBool bool

func (valueRef *VoIpMsStringBool) UnmarshalJSON(data []byte) error {
	var boolString string

	err := json.Unmarshal(data, &boolString)

	if err == nil {
		if boolString == "Yes" {
			*valueRef = true
		} else if boolString == "No" {
			*valueRef = false
		} else {
			return fmt.Errorf("value for bool was %v, expecting Yes or No", boolString)
		}
		return nil
	}

	var value bool
	err = json.Unmarshal(data, &value)
	if err != nil {
		return err
	}
	*valueRef = VoIpMsStringBool(value)
	return nil
}

type VoIpMsStringInt int64

func (valueRef *VoIpMsStringInt) UnmarshalJSON(data []byte) error {
	var (
		intString string
		i         int64
		err       error
	)

	err = json.Unmarshal(data, &intString)

	if err == nil {
		i, err = strconv.ParseInt(intString, 10, 64)
		if err == nil {
			*valueRef = VoIpMsStringInt(i)
			return nil
		}
	}

	err = json.Unmarshal(data, &i)
	if err != nil {
		return err
	}
	*valueRef = VoIpMsStringInt(i)
	return nil
}

func toURLValues(v reflect.Value) url2.Values {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	values := url2.Values{}

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		tag := field.Tag.Get("url")

		if tag != "" {
			values.Add(tag, fmt.Sprintf("%v", v.Field(i).Interface()))
		} else if field.Anonymous {
			embeddedValues := toURLValues(v.Field(i))
			for k, v := range embeddedValues {
				values[k] = v
			}
		}
	}

	return values
}

type RequestParams interface {
	ToURLValues() *url2.Values
	SetApiUser(username string)
	SetApiPassword(password string)
	SetApiMethod(method string)
}

type BaseRequest struct {
	ApiUser     string `url:"api_username"`
	ApiPassword string `url:"api_password"`
	Method      string `url:"method"`
}

func (r *BaseRequest) SetApiUser(username string) {
	r.ApiUser = username
}

func (r *BaseRequest) SetApiPassword(password string) {
	r.ApiPassword = password
}

func (r *BaseRequest) SetApiMethod(method string) {
	r.Method = method
}

func (r *BaseRequest) ToURLValues() *url2.Values {
	values := url2.Values{}
	values = toURLValues(reflect.ValueOf(r))
	return &values
}

type BaseResponse struct {
	Success string `json:"success"`
	Status  string `json:"status"`
	Message string `json:"message"`
	RawText string
}

func ParseBaseResponse(data *[]byte) (*BaseResponse, error) {
	response := &BaseResponse{}
	response.RawText = string(*data)
	if err := json.Unmarshal(*data, response); err != nil {
		return nil, err
	}
	return response, nil
}

func NewVoIpMsClient(username string, password string) *VoIpMsApi {
	return &VoIpMsApi{
		ApiUsername: username,
		ApiPassword: password,
		ApiUrl:      RestAPIURL,
	}
}

func (vms *VoIpMsApi) NewHttpRequest(httpMethod string, apiMethod string, requestData RequestParams) (*[]byte, error) {
	var (
		err          error
		url          *url2.URL
		request      *http.Request
		response     *http.Response
		responseBody []byte
		headers      http.Header
		httpClient   *http.Client
	)

	httpClient = &http.Client{
		Timeout: vms.ApiTimeout,
	}

	requestData.SetApiUser(vms.ApiUsername)
	requestData.SetApiPassword(vms.ApiPassword)
	requestData.SetApiMethod(apiMethod)

	queryParameters := requestData.ToURLValues().Encode()
	if url, err = url2.Parse(fmt.Sprintf("%s?%s", vms.ApiUrl, queryParameters)); err != nil {
		return nil, err
	}

	headers = http.Header{
		"Accept": []string{"text/json"},
	}

	request = &http.Request{
		Method:        httpMethod,
		URL:           url,
		Body:          nil,
		Header:        headers,
		GetBody:       nil,
		ContentLength: 0,
		Close:         true,
		Form:          nil,
		PostForm:      nil,
		Response:      nil,
	}

	if response, err = httpClient.Do(request); err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if responseBody, err = io.ReadAll(response.Body); err != nil {
		return nil, err
	}

	return &responseBody, nil
}
