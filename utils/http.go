package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func PSPostObj1(urlstr, ref string, header http.Header, obj interface{}, dataPayload func() any) (err error) {
	tres := new(tResObj1)
	resPayload := func() any {
		return tres
	}
	err = BasePostObj(urlstr, ref, header, obj, resPayload)
	if err != nil {
		return err
	}
	if tres.Code != 0 {
		return errors.New(tres.Message)
	}
	resData := dataPayload()
	err = json.Unmarshal(tres.Data, resData)
	if err != nil {
		return err
	}
	return
}

type tResObj1 struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func PSPostObj(urlstr, ref string, header http.Header, obj interface{}, dataPayload func() any) (err error) {
	tres := new(tResObj)
	resPayload := func() any {
		return tres
	}
	err = BasePostObj(urlstr, ref, header, obj, resPayload)
	if err != nil {
		return err
	}
	if tres.Code != "00000" {
		return errors.New(tres.Message)
	}
	resData := dataPayload()
	err = json.Unmarshal(tres.Data, resData)
	if err != nil {
		return err
	}
	return
}

type tResObj struct {
	Code    string          `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func BasePostObj(urlstr, ref string, header http.Header, obj interface{}, resPayload func() any) (err error) {
	if header == nil {
		header = http.Header{}
	}
	header.Set("Content-Type", "application/json; charset=utf-8")
	var bodyBs []byte
	bodyBs, err = json.Marshal(obj)
	if err != nil {
		return err
	}
	bs, err := BaseReq(urlstr, ref, "POST", header, bodyBs)
	if err != nil {
		return err
	}
	resObj := resPayload()
	err = json.Unmarshal(bs, resObj)
	return
}
func BaseReq(url, ref, method string, header http.Header, bodyBs []byte) (bs []byte, err error) {
	breader := bytes.NewReader(bodyBs)
	request, err := http.NewRequest(method, url, breader)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36")
	request.Header.Set("Referer", ref)
	for k, v := range header {
		request.Header.Set(k, v[0])
	}
	http.DefaultClient.Timeout = 10 * time.Second

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	if response == nil {
		err = errors.New("err resp nil")
	} else if response.StatusCode != 200 && response.StatusCode != 206 {
		bs, _ = ioutil.ReadAll(response.Body)
		err = fmt.Errorf("http status err %v content: %s", response.Status, string(bs))
	}
	if err != nil {
		return nil, err
	}
	bs, err = ioutil.ReadAll(response.Body)
	return
}
