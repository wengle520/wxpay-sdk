package wxpay

import (
	"testing"
)

func TestGenerateNonceStr(t *testing.T) {
	length := 32
	str := GenerateNonceStr()
	if len(str) == length {
		t.Log(str)
	} else {
		t.Error(str)
	}
}

func TestGenerateSignature(t *testing.T) {
	nonceStr := GenerateNonceStr()
	data := Params{}
	var signType SignTypeEnum
	signType = HMACSHA256
	data["appid"] = "12341231234"
	data["mch_id"] = "xs max"
	data["nonce_str"] = nonceStr
	data["sign_type"] = HMACSHA256_STR
	dataSign, _ := GenerateSignature(data, "bob123", signType)
	data["sign"] = dataSign
	t.Log(data)
}

func TestMapToXml(t *testing.T) {
	nonceStr := GenerateNonceStr()
	data := Params{}
	var signType SignTypeEnum
	signType = HMACSHA256
	data["appid"] = "12341231234"
	data["mch_id"] = "xs max"
	data["nonce_str"] = nonceStr
	data["sign_type"] = HMACSHA256_STR
	dataSign, _ := GenerateSignature(data, "bob123", signType)
	data["sign"] = dataSign
	res := MapToXml(data)
	t.Log(res)
}

func TestXmlToMap(t *testing.T) {
	nonceStr := GenerateNonceStr()
	data := Params{}
	var signType SignTypeEnum
	signType = HMACSHA256
	data["appid"] = "12341231234"
	data["mch_id"] = "xs max"
	data["nonce_str"] = nonceStr
	data["sign_type"] = HMACSHA256_STR
	dataSign, _ := GenerateSignature(data, "bob123", signType)
	data["sign"] = dataSign
	res := MapToXml(data)
	t.Log(res)

	ret, err := XmlToMap(res)
	if err != nil {
		t.Log(err)
	}
	t.Log(len(ret))
	t.Log(ret)
}

func TestIsSignatureValid(t *testing.T) {
	nonceStr := GenerateNonceStr()
	data := Params{}
	var signType SignTypeEnum
	signType = HMACSHA256
	data["appid"] = "12341231234"
	data["mch_id"] = "xs max"
	data["nonce_str"] = nonceStr
	data["sign_type"] = HMACSHA256_STR
	dataSign, _ := GenerateSignature(data, "bob123", signType)
	data["sign"] = dataSign

	valid, err := IsSignatureValid(data, "bob123", signType)
	if !valid {
		t.Error(err)
	}
	t.Log(valid)
}
