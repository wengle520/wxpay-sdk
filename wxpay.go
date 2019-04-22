package wxpay

import (
	"errors"
	"fmt"
)

type WXPay struct {
	config     *WXPayConfig
	signType   SignTypeEnum
	useSandbox bool
	autoReport bool
	notifyUrl  string
}

/**
* 作用：统一下单
* 场景：公共号支付、扫码支付、APP支付
* @param reqData 向wxpay post的请求数据
* @return API返回数据
* @throws Exception
*/
func (wxp *WXPay) unifiedOrder(reqData Params) (Params, error) {
	var url string
	if wxp.useSandbox {
		url = SANDBOX_UNIFIEDORDER_URL_SUFFIX
	} else {
		url = UNIFIEDORDER_URL_SUFFIX
	}

	if wxp.notifyUrl != "" {
		reqData["notify_url"] = wxp.notifyUrl
	}

	connectTimeoutMs := wxp.config.httpConnectTimeoutMs
	readTimeoutMs := wxp.config.httpReadTimeoutMs

	if _, err := wxp.fillRequestData(reqData, wxp.config.apiKey, wxp.signType); err != nil {
		return nil, err
	}

	msgUUID := reqData.GetStringParam("nonce_str")
	reqBody := MapToXml(reqData)
	//var respXml string
	respXml, err := WXPayRequestWithoutCert(url, msgUUID, reqBody, connectTimeoutMs, readTimeoutMs, wxp.autoReport)
	if err != nil {
		return nil, err
	}
	return processResponseXml(respXml)
}

/**
* 向 Map 中添加 appid、mch_id、nonce_str、sign_type、sign
* 该函数适用于商户统一下单等接口，不适用于红包、代金券接口
*
* @param reqData
* @return
*/
func (wxp *WXPay) fillRequestData(reqData Params, apiKey string, signType SignTypeEnum) (Params, error) {
	reqData["appid"] = wxp.config.appID
	reqData["mch_id"] = wxp.config.mchID
	reqData["nonce_str"] = GenerateNonceStr()
	if wxp.signType == MD5 {
		reqData["sign_type"] = MD5_STR
	} else if wxp.signType == HMACSHA256 {
		reqData["sign_type"] = HMACSHA256_STR
	}
	signRet, err := GenerateSignature(reqData, wxp.config.apiKey, wxp.signType)
	reqData["sign"] = signRet
	return reqData, err
}

/**
* 处理 HTTPS API返回数据，转换成Map对象。return_code为SUCCESS时，验证签名。
* @param xmlStr API返回的XML格式数据
* @return Map类型数据
* @throws Exception
*/
func (wxp *WXPay) processResponseXml(xmlStr string) (Params, error) {
	RETURN_CODE := "return_code"
	returnCode := ""
	respData, err := XmlToMap(xmlStr)
	if err != nil {
		return nil, err
	}

	if _, ok := respData[RETURN_CODE]; !ok {
		errMsg := fmt.Sprintf("No `return_code` in XML: %s", xmlStr)
		return nil, errors.New(errMsg)
	}
	returnCode = respData[RETURN_CODE]
	if returnCode == FAIL {
		return respData, nil
	} else if returnCode == SUCCESS {
		if valid, err := IsSignatureValid(respData, wxp.config.apiKey, wxp.signType); !valid || err != nil {
			errMsg := fmt.Sprintf("Invalid sign value in XML: %s", xmlStr)
			return nil, errors.New(errMsg)
		}
	} else {
		errMsg := fmt.Sprintf("return_code value %s is invalid in XML: %s", returnCode, xmlStr)
		return nil, errors.New(errMsg)
	}
	return respData, nil
}
