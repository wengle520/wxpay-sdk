package wxpay

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"sort"
	"strings"
	"time"
)

const SYMBOLS = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const XMLHEADER = `<?xml version="1.0" encoding="UTF-8"?>`

// 属性的键值对容器
type AttrMap map[string]string

/**
 * 生成签名. 注意，若含有sign_type字段，必须和signType参数保持一致。
 *
 * @param data 待签名数据
 * @param key API密钥
 * @param signType 签名方式
 * @return 签名
 */
func GenerateSignature(data Params, apiKey string, signType SignTypeEnum) (string, error) {
	keys := make([]string, 0, len(data))
	for k := range data {
		if k == FIELD_SIGN {
			continue
		}
		if len(strings.TrimSpace(data[k])) > 0 {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	var buf bytes.Buffer
	for _, k := range keys {
		buf.WriteString(k)
		buf.WriteString("=")
		buf.WriteString(strings.TrimSpace(data[k]))
		buf.WriteString("&")
	}

	// 加入apiKey作加密密钥
	buf.WriteString(`key=`)
	buf.WriteString(apiKey)

	var (
		md5Set    [16]byte
		sha256Set []byte
		str       string
	)
	switch signType {
	case MD5:
		md5Set = md5.Sum(buf.Bytes())
		str = hex.EncodeToString(md5Set[:])
	case HMACSHA256:
		h := hmac.New(sha256.New, []byte(apiKey))
		h.Write(buf.Bytes())
		sha256Set = h.Sum(nil)
		str = hex.EncodeToString(sha256Set[:])
	default:
		errMsg := fmt.Sprintf("invalid sign_type: %s", signType)
		return "", errors.New(errMsg)
	}
	return strings.ToUpper(str), nil
}

/**
 * 获取随机字符串 Nonce Str
 *
 * @return String 随机字符串
 */
func GenerateNonceStr() string {
	rand.Seed(time.Now().Unix())

	symbolsLen := len(SYMBOLS)
	arrLen := 32

	res := make([]byte, arrLen)
	for i := 0; i < arrLen; i++ {
		res[i] = SYMBOLS[rand.Intn(symbolsLen)]
	}
	return string(res)
}

// start()用来构建开始节点
func start(tag string, attrs AttrMap) xml.StartElement {
	var attrSet []xml.Attr
	for k, v := range attrs {
		attrSet = append(attrSet, xml.Attr{Name: xml.Name{Space: "", Local: k}, Value: v})
	}
	return xml.StartElement{Name: xml.Name{Space: "", Local: tag}, Attr: attrSet}
}

/**
 * 将Map转换为XML格式的字符串
 *
 * @param data Map类型数据
 * @return XML格式的字符串
 */
func MapToXml(data Params) string {
	// 创建编码器
	buffer := new(bytes.Buffer)
	enc := xml.NewEncoder(buffer)

	buffer.WriteString(XMLHEADER)
	buffer.WriteString("\n")
	// 设置缩进，这里为4个空格
	enc.Indent("", "    ")

	// 开始生成XML
	startExtension := start("xml", AttrMap{})
	enc.EncodeToken(startExtension)
	for k, v := range data {
		startKey := start(k, AttrMap{})
		enc.EncodeToken(startKey)
		enc.EncodeToken(xml.CharData(v))
		enc.EncodeToken(startKey.End())
	}
	enc.EncodeToken(startExtension.End())

	// 写入XML
	enc.Flush()
	return buffer.String()
}

/**
 * XML格式字符串转换为Map
 *
 * @param strXML XML字符串
 * @return XML数据转换后的Map
 * @throws Exception
 */
func XmlToMap(strXML string) (Params, error) {
	data := make(map[string]string)
	reader := bytes.NewReader([]byte(strXML))
	dec := xml.NewDecoder(reader)
	key, value := "", ""

	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			errMsg := fmt.Sprintf("Invalid XML, can not convert to map. Error message: %s. XML content: %s", err, strXML)
			return nil, errors.New(errMsg)
		}
		switch tk := tok.(type) {
		case xml.StartElement:
			key = tk.Name.Local
		case xml.CharData:
			value = string(tk)
		case xml.EndElement:
			data[key] = value
			key, value = "", ""
		}
	}
	if k, ok := data[""]; ok {
		fmt.Printf("delete map key {%s}\n", k)
		delete(data, "")
	}
	return data, nil
}

/**
 * 判断签名是否正确，必须包含sign字段，否则返回false。
 *
 * @param data Map类型数据
 * @param key API密钥
 * @param signType 签名方式
 * @return 签名是否正确
 * @throws Exception
 */
func IsSignatureValid(data Params, apiKey string, signType SignTypeEnum) (bool, error) {
	if _, ok := data[FIELD_SIGN]; !ok {
		return false, nil
	}
	sign := data[FIELD_SIGN]
	genSign, err := GenerateSignature(data, apiKey, signType)
	if err != nil {
		return false, err
	}
	return genSign == sign, nil
}