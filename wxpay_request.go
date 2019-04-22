package wxpay

func requestOnce(domain string, urlSuffix string, uuid string, data string, connectTimeoutMs int, readTimeoutMs int, useCert bool)  {

}

func request(urlSuffix string, uuid string, data string, connectTimeoutMs int, readTimeoutMs int, useCert bool, autoReport bool) (string, error) {

}

func WXPayRequestWithoutCert(urlSuffix string, uuid string, data string, connectTimeoutMs int, readTimeoutMs int, autoReport bool) (string, error) {

}

func WXPayRequestWithCert(urlSuffix string, uuid string, data string, connectTimeoutMs int, readTimeoutMs int, autoReport bool) (string, error) {

}