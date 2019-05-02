package wxpay

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

func TestInitConfig(t *testing.T) {
	// expected
	iwxpd := new(DummyWXPayDomain)
	configExpected := WXPayConfig{appID:"wx2421b1c4370ec43b", mchID:"10000100", autoReport:true, reportWorkerNum:6, wxPayDomain:iwxpd}

	// test
	testConfigName := "test/conf/wxpay_config.yaml"
	InitConfig(testConfigName, iwxpd)
	configTest := GetConfigInstance()
	if configTest.appID != configExpected.appID {
		t.Errorf("appID test: %s, expected: %s\n", configTest.appID, configExpected.appID)
	}
	if configTest.mchID != configExpected.mchID {
		t.Errorf("mchID test: %s, expected: %s\n", configTest.mchID, configExpected.mchID)
	}
	if configTest.autoReport != configExpected.autoReport {
		t.Errorf("autoReport test: %s, expected: %s\n", strconv.FormatBool(configTest.autoReport), strconv.FormatBool(configExpected.autoReport))
	}
	if configTest.reportWorkerNum != configExpected.reportWorkerNum {
		t.Errorf("reportWorkerNum test: %d, expected: %d\n", configTest.reportWorkerNum, configExpected.reportWorkerNum)
	}
	fmt.Println("type:", reflect.TypeOf(iwxpd))
	t.Log("success")
}