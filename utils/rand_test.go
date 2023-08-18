package utils

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"
)

func TestGenCaptcha(t *testing.T) {
	InitCapcha()
	id, answer, img := GenCaptcha(true)
	log.Println(id, answer, img, VerifyCaptcha(id, answer))
	//Mail("abc", answer,"wxf4150@163.com")
}

func TestRandomDigits(t *testing.T) {
	got := RandomDigits(6)
	log.Println(got)
}
func TestHttp(t *testing.T) {
	//bodyBs := []byte(`{
	// "merNo":"104001001",
	// "merOrderNo":"asdfghjkl",
	// "version":"V3.0.0",
	// "sign":"9D6FDF4880B00B002B1F2AB61AE9A721",
	// "tradeNo": "DZ13576867867",
	// "queryType":"sales"
	//}`)
	////header := http.Header{}
	////header.Set("Content-Type", "application/json; charset=utf-8")
	////bs, err := BaseReq("https://payment.flyzeroc.xyz/order/query", "", "POST", header, bodyBs)
	//tdata := tResData{}
	//err := PSPostObj("https://payment.flyzeroc.xyz/order/query", "", nil, json.RawMessage(bodyBs), func() any {
	//	return tdata
	//})
	//if err != nil {
	//	log.Fatalln(err)
	//}
	bodyBs := []byte(`{
    "merNo": "104001001",
    "openId": "1691365203335311361",
    "icNo": "john123",
    "currency": "string"
}`)
	header := http.Header{}
	header.Set("Content-Type", "application/json; charset=utf-8")
	bs, err := BaseReq("http://116.204.78.80:9567/api/mg/bcard/add", "", "POST", header, json.RawMessage(bodyBs))

	//tdata := tResData{}
	//err := PSPostObj("https://payment.flyzeroc.xyz/order/query", "", nil, json.RawMessage(bodyBs), func() any {
	//	return tdata
	//})
	if err != nil {
		log.Fatalln(err)
	}
	bstr := string(bs)
	log.Println(bstr)

	//
	//bs := []byte(`{
	// "code": "00000",
	// "message": "SUCCESS",
	// "data": {
	//   "tradeNo": "DZ2202101809494673",
	//   "merOrderNo": "1644487789447",
	//   "merNo": 104001002,
	//   "state": "3"
	// }
	//}`)
	//
	//strbody := string(bs)
	//log.Println(strbody)
	//tres := new(tResObj)
	//err := json.Unmarshal(bs, tres)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//var tdata any
	//tdata = new(tResData)
	//err = json.Unmarshal(tres.Data, tdata)
	//if err != nil {
	//	log.Println(err)
	//}
}

type tResData struct {
	TradeNo    string `json:"tradeNo"`
	MerOrderNo string `json:"merOrderNo"`
	MerNo      int    `json:"merNo"`
	State      string `json:"state"`
}

func TestMd5(t *testing.T) {
	log.Println(fmt.Sprintf("%x", md5.Sum([]byte("123456"))))
}
