package wxpay

import "io"

type WXPayConfig struct {
	appID                string
	mchID                string
	apiKey               string
	certStream           io.Reader
	httpConnectTimeoutMs int
	httpReadTimeoutMs    int
	//wxPayDomain          IWXPayDomain
	autoReport           bool
	reportWorkerNum      int
	reportQueueMaxSize   int
	reportBatchSize      int
}

///**
// * 获取 App ID
// *
// * @return App ID
// */
//func (wxpc *WXPayConfig) getAppID() string {
//	return wxpc.appID
//}
//
///**
// * 获取 Mch ID
// *
// * @return Mch ID
// */
//func (wxpc *WXPayConfig) getMchID() string {
//	return wxpc.mchID
//}
//
///**
// * 获取 API 密钥
// *
// * @return API密钥
// */
//func (wxpc *WXPayConfig) getKey() string {
//	return wxpc.apiKey
//}
//
///**
// * HTTP(S) 连接超时时间，单位毫秒
// *
// * @return
// */
//func (wxpc *WXPayConfig) getHttpConnectTimeoutMs() int {
//	return wxpc.httpConnectTimeoutMs
//}
//
///**
// * HTTP(S) 读数据超时时间，单位毫秒
// *
// * @return
// */
//
//func (wxpc *WXPayConfig) getHttpReadTimeoutMs() int {
//	return wxpc.httpReadTimeoutMs
//}
