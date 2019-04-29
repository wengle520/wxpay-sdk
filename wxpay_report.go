package wxpay

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"sync"
	"time"
)

type ReportInfo struct {
	// 基本信息
	version           string
	sdk               string
	uuid              string // 交易的标识
	timestamp         int64  // 上报时的时间戳，单位秒
	elapsedTimeMillis int64  // 耗时，单位 毫秒

	// 针对主域名
	firstDomain               string // 第1次请求的域名
	primaryDomain             bool   //是否主域名
	firstConnectTimeoutMillis int    // 第1次请求设置的连接超时时间，单位 毫秒
	firstReadTimeoutMillis    int    // 第1次请求设置的读写超时时间，单位 毫秒
	firstHasDnsError          int    // 第1次请求是否出现dns问题
	firstHasConnectTimeout    int    // 第1次请求是否出现连接超时
	firstHasReadTimeout       int    // 第1次请求是否出现连接超时
}

func NewReportInfo(uuid string, currentTimestamp int64, elapsedTimeMillis int64,
	firstDomain string, primaryDomain bool, firstConnectTimeoutMillis int, firstReadTimeoutMillis int,
	firstHasDnsError bool, firstHasConnectTimeout bool, firstHasReadTimeout bool) *ReportInfo {
	reportInfo := new(ReportInfo)
	reportInfo.version = "v1"
	reportInfo.sdk = WXPAYSDK_VERSION
	reportInfo.uuid = uuid
	reportInfo.timestamp = currentTimestamp
	reportInfo.elapsedTimeMillis = elapsedTimeMillis
	reportInfo.firstDomain = firstDomain
	reportInfo.primaryDomain = primaryDomain
	reportInfo.firstConnectTimeoutMillis = firstConnectTimeoutMillis
	reportInfo.firstReadTimeoutMillis = firstReadTimeoutMillis
	reportInfo.firstHasDnsError = 0
	if firstHasDnsError {
		reportInfo.firstHasDnsError = 1
	}
	reportInfo.firstHasConnectTimeout = 0
	if firstHasConnectTimeout {
		reportInfo.firstHasConnectTimeout = 1
	}
	reportInfo.firstHasReadTimeout = 0
	if firstHasReadTimeout {
		reportInfo.firstHasReadTimeout = 1
	}
}

func (ri *ReportInfo) String() string {
	return "ReportInfo{" +
		"version='" + ri.version + "'" +
		", sdk='" + ri.sdk + "'" +
		", uuid='" + ri.uuid + "'" +
		", timestamp=" + strconv.FormatInt(ri.timestamp, 10) +
		", elapsedTimeMillis=" + strconv.FormatInt(ri.elapsedTimeMillis, 10) +
		", firstDomain='" + ri.firstDomain + "'" +
		", primaryDomain=" + strconv.FormatBool(ri.primaryDomain) +
		", firstConnectTimeoutMillis=" + strconv.Itoa(ri.firstConnectTimeoutMillis) +
		", firstReadTimeoutMillis=" + strconv.Itoa(ri.firstReadTimeoutMillis) +
		", firstHasDnsError=" + strconv.Itoa(ri.firstHasDnsError) +
		", firstHasConnectTimeout=" + strconv.Itoa(ri.firstHasConnectTimeout) +
		", firstHasReadTimeout=" + strconv.Itoa(ri.firstHasReadTimeout) +
		"}"
}

func (ri *ReportInfo) toLineString(apiKey string) (string, error) {
	separator := ","
	buffer := new(bytes.Buffer)

	v := reflect.ValueOf(ri)
	count := v.NumField()
	for i := 0; i < count; i++ {
		f := v.Field(i)
		switch f.Kind() {
		case reflect.String:
			buffer.WriteString(f.String())
		case reflect.Int:
			buffer.WriteString(strconv.FormatInt(f.Int(), 10))
		case reflect.Int64:
			buffer.WriteString(strconv.FormatInt(f.Int(), 10))
		case reflect.Bool:
			buffer.WriteString(strconv.FormatBool(f.Bool()))
		default:
			buffer.WriteString(fmt.Sprintf("%v", f))
		}
		buffer.WriteString(separator)
	}

	sign := GenHMACSHA256(buffer, apiKey)
	buffer.WriteString(sign)

	return buffer.String(), nil
}

type WXPayReport struct {
	reportUrl               string
	defaultConnectTimeoutMs int
	defaultReadTimeoutMs    int
	reportMsgQueue          chan string
	config                  *WXPayConfig
}

var reportIns *WXPayReport
var reportMutex sync.Mutex

func runTask(ins *WXPayReport) {
	for {
		buffer := new(bytes.Buffer)
		firstMsg := <-ins.reportMsgQueue
		fmt.Printf("get first report msg: %s\n", firstMsg)
		buffer.WriteString(firstMsg)

		remainNumMsg := ins.config.reportBatchSize - 1
		for i := 0; i < remainNumMsg; i++ {
			fmt.Println("try get remain report msg")
			msg, ok := <-ins.reportMsgQueue
			if !ok { // chan has been closed
				break
			}
			fmt.Printf("get remain report msg: %s\n", firstMsg)
			buffer.WriteString("\n")
			buffer.WriteString(msg)
		}

		// 上报
		resp, err := ins.httpRequest(buffer.String(), ins.defaultConnectTimeoutMs, ins.defaultReadTimeoutMs)
		if err != nil {
			fmt.Printf("report fail. reason: %s\n", err)
		}
		fmt.Printf("response msg: %s\n", resp)
	}
}

func initReport(config *WXPayConfig) *WXPayReport {
	reportObj := &WXPayReport{}
	reportObj.config = config
	reportObj.reportUrl = "http://report.mch.weixin.qq.com/wxpay/report/default"
	reportObj.defaultConnectTimeoutMs = 6 * 1000
	reportObj.defaultReadTimeoutMs = 8 * 1000
	reportObj.reportMsgQueue = make(chan string, config.reportQueueMaxSize)

	if config.autoReport {
		fmt.Printf("report worker num: %d\n", config.reportWorkerNum)
		for i := 0; i < config.reportWorkerNum; i++ {
			go runTask(reportObj)
		}
	}
	return reportObj
}

func GetReportInstance(config *WXPayConfig) *WXPayReport {
	if reportIns == nil {
		reportMutex.Lock()
		defer reportMutex.Unlock()
		if reportIns == nil {
			reportIns = initReport(config)
		}
	}
	return reportIns
}

func (wxr *WXPayReport) Report(uuid string, elapsedTimeMillis int64,
	firstDomain string, primaryDomain bool, firstConnectTimeoutMillis int, firstReadTimeoutMillis int,
	firstHasDnsError bool, firstHasConnectTimeout bool, firstHasReadTimeout bool) {
	configIns := GetConfigInstance()
	currentTimestamp := GetCurrentTimestampMs()
	reportInfo := NewReportInfo(uuid, currentTimestamp, elapsedTimeMillis,
		firstDomain, primaryDomain, firstConnectTimeoutMillis, firstReadTimeoutMillis,
		firstHasDnsError, firstHasConnectTimeout, firstHasReadTimeout)
	data, err := reportInfo.toLineString(configIns.apiKey)
	if err != nil {
		errMsg := fmt.Sprintf("convert to cvs format failed: %s", err)
		fmt.Println(errMsg)
	}
	fmt.Printf("report %s\n", data)
	if data != "" {
		ok := false
		select {
		case wxr.reportMsgQueue <- data:
			ok = true
		default:
			ok = false
		}
		if !ok {
			fmt.Printf("current reportMsgQueue has full, discard report: %s\n", data)
		}
	}
}

func (wxr *WXPayReport) httpRequest(data string, connectTimeoutMs int, readTimeoutMs int) (string, error) {
	client := &http.Client{Transport:
	&http.Transport{DialContext: TimeoutDialer(time.Duration(connectTimeoutMs), time.Duration(readTimeoutMs))}}

	req, err := http.NewRequest(http.MethodPost, wxr.reportUrl, bytes.NewBuffer([]byte(data)))
	if err != nil {
		errMsg := fmt.Sprintf("NewRequest err: %s", err)
		return "", errors.New(errMsg)
	}

	req.Header.Set("Content-Type", "text/xml")
	req.Header.Set("User-Agent", USER_AGENT)

	resp, err := client.Do(req)
	if err != nil {
		errMsg := fmt.Sprintf("Request url: %s, err: %s", wxr.reportUrl, err)
		return "", errors.New(errMsg)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}
