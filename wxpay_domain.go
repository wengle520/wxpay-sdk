package wxpay

import "strconv"

type DomainInfo struct {
	domain        string
	primaryDomain bool
}

// SDK使用者调用
func NewDomainInfo(domain string, primaryDomain bool) DomainInfo {
	di := DomainInfo{domain, primaryDomain}
	return di
}

func (di *DomainInfo) String() string {
	return "DomainInfo{" + "domain='" +
		di.domain + "'" + ", primaryDomain=" +
		strconv.FormatBool(di.primaryDomain) + "}"
}

// SDK使用者负责实现这个接口
// 用于表示当前使用的支付api的网络性能
// constants.go 定义了多个域名，用于主备域名切换
// DOMAIN_API   = "api.mch.weixin.qq.com"
// DOMAIN_API2  = "api2.mch.weixin.qq.com"
// DOMAIN_APIHK = "apihk.mch.weixin.qq.com"
// DOMAIN_APIUS = "apius.mch.weixin.qq.com"
type IWXPayDomain interface {
	Report(domain string, elapsedTimeMillis int64, err error)
	GetDomain() DomainInfo
}

// TODO 实现一个IWXPayDomain接口
type DummyWXPayDomain struct {
}

func NewDummyWXPayDomain() IWXPayDomain {
	return new(DummyWXPayDomain)
}

func (dwp *DummyWXPayDomain) Report(domain string, elapsedTimeMillis int64, err error) {

}

func (dwp *DummyWXPayDomain) GetDomain() DomainInfo {
	return NewDomainInfo(DOMAIN_API, true)
}
