package wxpay

import (
	"testing"
)

func TestGenerateNonceStr(t *testing.T) {
	length := 32
	str := GenerateNonceStr()
	if len(str) == length {
		t.Log(str)
	}else {
		t.Error(str)
	}
}

func TestGenerateSignature(t *testing.T) {
	nonceStr := GenerateNonceStr()
	data := Params{}
	var signType SignTypeEnum
	signType = HMACSHA256
	data.SetStringParam("appid", "12341231234")
	data.SetStringParam("mch_id", "wengle_5s")
	data.SetStringParam("nonce_str", nonceStr)
	data.SetStringParam("sign_type", HMACSHA256_STR)
	dataSign, _ := GenerateSignature(data, "wengle123", signType)
	data.SetStringParam("sign", dataSign)
	t.Log(data)
}

func TestMapToXml(t *testing.T) {
	nonceStr := GenerateNonceStr()
	data := Params{}
	var signType SignTypeEnum
	signType = HMACSHA256
	data.SetStringParam("appid", "12341231234")
	data.SetStringParam("mch_id", "wengle_5s")
	data.SetStringParam("nonce_str", nonceStr)
	data.SetStringParam("sign_type", HMACSHA256_STR)
	dataSign, _ := GenerateSignature(data, "wengle123", signType)
	data.SetStringParam("sign", dataSign)
	res := MapToXml(data)
	t.Log(res)
}


func TestXmlToMap(t *testing.T) {
	nonceStr := GenerateNonceStr()
	data := Params{}
	var signType SignTypeEnum
	signType = HMACSHA256
	data.SetStringParam("appid", "12341231234")
	data.SetStringParam("mch_id", "wengle_5s")
	data.SetStringParam("nonce_str", nonceStr)
	data.SetStringParam("sign_type", HMACSHA256_STR)
	dataSign, _ := GenerateSignature(data, "wengle123", signType)
	data.SetStringParam("sign", dataSign)
	res := MapToXml(data)
	t.Log(res)

	ret, err := XmlToMap(res)
	if err != nil {
		t.Log(err)
	}
	t.Log(len(ret))
	t.Log(ret)
}