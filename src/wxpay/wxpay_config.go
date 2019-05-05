package wxpay

import (
	"fmt"
	"github.com/kylelemons/go-gypsy/yaml"
)

type WXPayConfig struct {
	appID                string
	mchID                string
	apiKey               string
	cert                 string // 证书文件路径
	httpConnectTimeoutMs int
	httpReadTimeoutMs    int
	wxPayDomain          IWXPayDomain // SDK的使用者负责赋值
	autoReport           bool
	reportWorkerNum      int
	reportQueueMaxSize   int
	reportBatchSize      int
}

var configIns *WXPayConfig

func InitConfig(configFile string, wxpd IWXPayDomain) {
	configIns = &WXPayConfig{}

	config, err := yaml.ReadFile(configFile)
	if err != nil {
		fmt.Printf("read %s failed: %s\n", configFile, err)
	}
	configIns.appID, _ = config.Get("appID")
	configIns.mchID, _ = config.Get("mchID")
	configIns.apiKey, _ = config.Get("apiKey")
	configIns.cert, _ = config.Get("cert")
	data, _ := config.GetInt("httpConnectTimeoutMs")
	configIns.httpConnectTimeoutMs = int(data)
	data, _ = config.GetInt("httpReadTimeoutMs")
	configIns.httpReadTimeoutMs = int(data)
	configIns.autoReport, _ = config.GetBool("autoReport")
	data, _ = config.GetInt("reportWorkerNum")
	configIns.reportWorkerNum = int(data)
	data, _ = config.GetInt("reportQueueMaxSize")
	configIns.reportQueueMaxSize = int(data)
	data, _ = config.GetInt("reportBatchSize")
	configIns.reportBatchSize = int(data)
	configIns.wxPayDomain = wxpd
}

func GetConfigInstance() *WXPayConfig {
	return configIns
}
