package wxpay

import "runtime"

type SignTypeEnum int

const (
	MD5 = iota
	HMACSHA256
)

const (
	DOMAIN_API   = "api.mch.weixin.qq.com"
	DOMAIN_API2  = "api2.mch.weixin.qq.com"
	DOMAIN_APIHK = "apihk.mch.weixin.qq.com"
	DOMAIN_APIUS = "apius.mch.weixin.qq.com"
)

const (
	FAIL    = "FAIL"
	SUCCESS = "SUCCESS"
)

const (
	MD5_STR         = "MD5"
	HMACSHA256_STR  = "HMAC-SHA256"
	FIELD_SIGN      = "sign"
	FIELD_SIGN_TYPE = "sign_type"
)

const (
	WXPAYSDK_VERSION = "WXPaySDK/3.0.9"
)

const (
	USER_AGENT = WXPAYSDK_VERSION + " (" + runtime.GOARCH + runtime.GOOS + ") go/1.10.3 " + "httpclient"
)

const (
	MICROPAY_URL_SUFFIX         = "/pay/micropay"
	UNIFIEDORDER_URL_SUFFIX     = "/pay/unifiedorder"
	ORDERQUERY_URL_SUFFIX       = "/pay/orderquery"
	REVERSE_URL_SUFFIX          = "/secapi/pay/reverse"
	CLOSEORDER_URL_SUFFIX       = "/pay/closeorder"
	REFUND_URL_SUFFIX           = "/secapi/pay/refund"
	REFUNDQUERY_URL_SUFFIX      = "/pay/refundquery"
	DOWNLOADBILL_URL_SUFFIX     = "/pay/downloadbill"
	REPORT_URL_SUFFIX           = "/payitil/report"
	SHORTURL_URL_SUFFIX         = "/tools/shorturl"
	AUTHCODETOOPENID_URL_SUFFIX = "/tools/authcodetoopenid"
)

const (
	SANDBOX_MICROPAY_URL_SUFFIX         = "/sandboxnew/pay/micropay"
	SANDBOX_UNIFIEDORDER_URL_SUFFIX     = "/sandboxnew/pay/unifiedorder"
	SANDBOX_ORDERQUERY_URL_SUFFIX       = "/sandboxnew/pay/orderquery"
	SANDBOX_REVERSE_URL_SUFFIX          = "/sandboxnew/secapi/pay/reverse"
	SANDBOX_CLOSEORDER_URL_SUFFIX       = "/sandboxnew/pay/closeorder"
	SANDBOX_REFUND_URL_SUFFIX           = "/sandboxnew/secapi/pay/refund"
	SANDBOX_REFUNDQUERY_URL_SUFFIX      = "/sandboxnew/pay/refundquery"
	SANDBOX_DOWNLOADBILL_URL_SUFFIX     = "/sandboxnew/pay/downloadbill"
	SANDBOX_REPORT_URL_SUFFIX           = "/sandboxnew/payitil/report"
	SANDBOX_SHORTURL_URL_SUFFIX         = "/sandboxnew/tools/shorturl"
	SANDBOX_AUTHCODETOOPENID_URL_SUFFIX = "/sandboxnew/tools/authcodetoopenid"
)
