package wxpay

import (
	"testing"
)

func TestUnifiedOrder(t *testing.T) {
	iwxpd := new(DummyWXPayDomain)
	testConfigName := "test/conf/wxpay_config.yaml"
	InitConfig(testConfigName, iwxpd)
	config := GetConfigInstance()

	// recommend to using NewWXPay
	wxp := new(WXPay)
	wxp.config = config
	wxp.signType = MD5
	wxp.useSandbox = true
	wxp.notifyUrl = "http://www.example.com/wxpay/notify"

	reqData := make(Params)
	reqData["out_trade_no"] = "1415659990"
	reqData["total_fee"] = "1"
	reqData["trade_type"] = "APP"

	respData, err := wxp.UnifiedOrder(reqData)
	t.Log(reqData)
	if err != nil {
		t.Error(err)
	}
	t.Log(respData)
}
