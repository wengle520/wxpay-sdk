微信支付 Go SDK
------

对[微信支付开发者文档](https://pay.weixin.qq.com/wiki/doc/api/index.html)中给出的API进行了封装。

提供了对应的方法：

|方法名 | 说明 |
|--------|--------|
|microPay| 刷卡支付 |
|unifiedOrder | 统一下单|
|orderQuery | 查询订单 |
|reverse | 撤销订单 |
|closeOrder|关闭订单|
|refund|申请退款|
|refundQuery|查询退款|
|downloadBill|下载对账单|
|report|交易保障|
|shortUrl|转换短链接|
|authCodeToOpenid|授权码查询openid|

* 注意:
* 证书文件不能放在web服务器虚拟目录，应放在有访问权限控制的目录中，防止被他人下载
* 建议将证书文件名改为复杂且不容易猜测的文件名
* 商户服务器要做好病毒和木马防护工作，不被非法侵入者窃取证书文件
* 请妥善保管商户支付密钥、公众帐号SECRET，避免密钥泄露
* 参数为`Map<String, String>`对象，返回类型也是`Map<String, String>`
* 方法内部会将参数会转换成含有`appid`、`mch_id`、`nonce_str`、`sign\_type`和`sign`的XML
* 可选HMAC-SHA256算法和MD5算法签名
* 通过HTTPS请求得到返回数据后会对其做必要的处理（例如验证签名，签名错误则抛出异常）
* 对于downloadBill，无论是否成功都返回Map，且都含有`return_code`和`return_msg`，若成功，其中`return_code`为`SUCCESS`，另外`data`对应对账单数据


## 示例
配置类WXPayConfig:
```go
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
```

统一下单：

```go
package wxpay

import (
	"testing"
)

func TestUnifiedOrder(t *testing.T) {
	iwxpd := new(DummyWXPayDomain)
	testConfigName := "test/conf/wxpay_config.yaml"
	InitConfig(testConfigName, iwxpd)
	config := GetConfigInstance()

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
```