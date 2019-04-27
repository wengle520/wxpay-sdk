package wxpay

import "strconv"

type DomainInfo struct {
	domain        string
	primaryDomain bool
}

func NewDomainInfo(domain string, primaryDomain bool) DomainInfo {
	di := DomainInfo{domain, primaryDomain}
	return di
}

func (di *DomainInfo) String() string {
	return "DomainInfo{" + "domain='" +
		di.domain + "'" + ", primaryDomain=" +
		strconv.FormatBool(di.primaryDomain) + "}"
}

type IWXPayDomain interface {
	Report(domain string, elapsedTimeMillis int64, err error)
	GetDomain() DomainInfo
}
