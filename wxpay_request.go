package wxpay

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

/**
 * 设置请求的连接超时时间和读取数据超时时间
 *
 */
func TimeoutDialer(cTimeout time.Duration, rTimeout time.Duration) func(ctx context.Context, net string, addr string) (c net.Conn, err error) {
	return func(ctx context.Context, netw string, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, cTimeout)
		if err != nil {
			return nil, err
		}
		conn.SetReadDeadline(time.Now().Add(rTimeout))
		return conn, nil
	}
}

/**
 * 请求，只请求一次，不做重试
 * @param domain
 * @param urlSuffix
 * @param uuid
 * @param data
 * @param connectTimeoutMs
 * @param readTimeoutMs
 * @param useCert 是否使用证书，针对退款、撤销等操作
 * @return
 */
func requestOnce(domain string, urlSuffix string, uuid string, data string, connectTimeoutMs int, readTimeoutMs int, useCert bool) (string, error) {
	client := &http.Client{Transport:
	&http.Transport{DialContext: TimeoutDialer(time.Duration(connectTimeoutMs), time.Duration(readTimeoutMs))}}
	configIns := GetConfigInstance()

	if useCert {
		//证书
		pool := x509.NewCertPool()
		caCert, err := ioutil.ReadFile(configIns.cert)
		if err != nil {
			errMsg := fmt.Sprintf("ReadFile err: %s", err)
			return "", errors.New(errMsg)
		}
		pool.AppendCertsFromPEM(caCert)
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: pool},
			DialContext:     TimeoutDialer(time.Duration(connectTimeoutMs), time.Duration(readTimeoutMs)),
		}
		client = &http.Client{Transport: tr}
	}

	url := "https://" + domain + urlSuffix
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		errMsg := fmt.Sprintf("NewRequest err: %s", err)
		return "", errors.New(errMsg)
	}

	req.Header.Set("Content-Type", "text/xml")
	req.Header.Set("User-Agent", USER_AGENT+" "+configIns.mchID)

	resp, err := client.Do(req)
	if err != nil {
		errMsg := fmt.Sprintf("Request url: %s, err: %s", url, err)
		return "", errors.New(errMsg)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func request(urlSuffix string, uuid string, data string, connectTimeoutMs int, readTimeoutMs int, useCert bool, autoReport bool) (string, error) {
	configIns := GetConfigInstance()

	var elapsedTimeMillis = int64(0)
	var startTimestampMs = GetCurrentTimestampMs()
	var firstHasDnsErr = false
	var firstHasConnectTimeout = false
	var firstHasReadTimeout = false
	var domainInfo = configIns.wxPayDomain.GetDomain()
	if domainInfo == (DomainInfo{}) {
		errMsg := "GetConfigInstance().wxPayDomain.GetDomain() is nil"
		return "", errors.New(errMsg)
	}

	result, err := requestOnce(domainInfo.domain, urlSuffix, uuid, data, connectTimeoutMs, readTimeoutMs, useCert)
	if err != nil {
		configIns.wxPayDomain.Report(domainInfo.domain, elapsedTimeMillis, err)
		return "", err
	}
	elapsedTimeMillis = GetCurrentTimestampMs() - startTimestampMs
	configIns.wxPayDomain.Report(domainInfo.domain, elapsedTimeMillis, nil)
	GetReportInstance(configIns).Report(uuid, elapsedTimeMillis, domainInfo.domain,
		domainInfo.primaryDomain, connectTimeoutMs, readTimeoutMs,
		firstHasDnsErr, firstHasConnectTimeout, firstHasReadTimeout)
	return result, nil
}

func WXPayRequestWithoutCert(urlSuffix string, uuid string, data string, connectTimeoutMs int, readTimeoutMs int, autoReport bool) (string, error) {
	return request(urlSuffix, uuid, data, connectTimeoutMs, readTimeoutMs, false, autoReport)
}

func WXPayRequestWithCert(urlSuffix string, uuid string, data string, connectTimeoutMs int, readTimeoutMs int, autoReport bool) (string, error) {

}
