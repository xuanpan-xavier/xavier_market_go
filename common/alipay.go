package common

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

/*** 请填写以下配置信息 ***/
var appid string = "2021003146622660"                    //https://open.alipay.com 账户中心->密钥管理->开放平台密钥，填写添加了电脑网站支付的应用的APPID
var returnUrl string = "http://www.example.com/callback" //付款成功后的同步回调地址
var notifyUrl string = "http://www.example.com/notify"   //付款成功后的异步回调地址
var payAmount float64 = 0.01                             //付款金额，单位:元
var orderName string = "支付测试"                            //订单标题
var signType string = "RSA2"                             //签名算法类型，支持RSA2和RSA，推荐使用RSA2

//商户私钥，填写对应签名算法类型的私钥，如何生成密钥参考：https://docs.open.alipay.com/291/105971和https://docs.open.alipay.com/200/105310
var rsaPrivateKey string = "MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCuFR9HM8ePOiUzF2CkKj54K9ZiZKAmzB+pw67EHJ9BrUfqnPtZlN39pXLdT0ayzgp2zLEYrO0xxwHqKlpY6OCsUj3ostgF+gbqCoZGelP7yXUMVuZh4kPBRge9BpljIfvcjROuzwtL/j2IFK35/NYmcRn/RsfFy2pRx0yt3dcskBpMeYKs/pUL65Yw3jH366O0B/icndnM/F7aCKdDwgmaw40tNiddnJGMq0O40iS3EWX8NGAH/Ifpo0Or9vXzWVmWZNN/+Wm/cKp405XvuABCqWGqx7n33RnOf9iCCB7MUH7/IhvYKmFIo8Z3ZpVX8HgZLH5mNyl6KJUNLbk3LBZbAgMBAAECggEANCfJ5jhCOswt5v/DCYxKiRCBO5twTA9uERu6mbEFmdTDriv4Nl+enubZ4iplQR7wDaXIDV+J9RbtoO2tgePLpHVDbRtgE3JUbfxXPHhb/xzdf6PthDgZGqfTm7OFcyj0X6C42+/WL59taf7gHPy+YGgUZOIc6krGgNfYRmMnm4bU0N4vVIZTP9ZuF486NTEtCWFHwuiQ+CpLv8MqCI/zqacn6EDCjj1fcliy2mida7raSDTmwxW3dxLA3W3IcBALRJKOxnVgpXflNews2DYULDq59ynXcevvSrNsSHfOrpEAJXMIw/PiuprkHEaSVHN1bptnTiocjjgkroZUP0srkQKBgQD4/0KbQ+CuXnKAAGqvHfJQZHiDBVp1BZK8GyNK4XjEjOASp2r6dNUB+WZMSIyW6qltHvb5ntafksY8vqBqTih+DTjT8U52n3a0ulVnjvF7H77xXlDD85yUeVY9JgDHX6kVcTUDtrpC6EwsqRdyrikZiZ+V3uUhlu8MIRTWAjKs0wKBgQCy+n0cVd0/tixpVRVrWgO6CuKqZ7Wo59xRjsLQIBRpNq/J9R2umQOreumiFz4V26GCJymfNEXh6hJrEbSaWqWK92JFTAq5ndZvjASa+qMhTQRAUgwlcE44jpXZG8120J64QjMNlUWHITr/0AVv4vqfgHsPfdl1plQU1jTUsAVbWQKBgQCfFAXryGvtOu4QFTQTeSBseCP+IRdvmAXtI7Fij7ZmbFJj2I4iQS4tOrc+mUMoz/c6fgPv6qjX/h+kg21dPV3lWnEqBwrIAo0GXd1echaZQOjusAENj8BDo1mTa4on1lfiGRdzrfRF9r9j8QvH+rjLVf+c8nCVCybiKuhy6aiSfQKBgH3Hja8hJLjGpRVtRv1HqEDRr5PECKVoNKxYYM51DQMX/Jg4m6B9Rlf3rENOhEJGNI3MSfLpgIcQwdre+YZUV+x+iGBc3Viu9Ok1hHMJ2yTnBCrTvld9Qmt6wRJv+Igcer7RPyG3msIpd1YZ8gg8gpcfatxrPCzzNDnkhAJIpasZAoGAaVoaGr8dj08V4F47M6VwmIPEaMxiJHqba8iZqZQAInHbtC/FBv43D/gS/dvr7WCJSY2E9PsXcsW8hr9OBZehOsUw4G396NCqcy1dIO+TnqDbDFl9yaDV/LivKXyIJ4jqv9yMFUskztIJSmaHbEzk0CpnfSZzwUe7owagHWwDkaI="

//支付宝公钥，登录支付宝开放平台，账户中心->密钥管理->开放平台密钥，找到对应的应用，在接口内容加密方式处查看支付宝公钥
var alipayPublicKey string = "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAkX6GjcSRpso66FYRGBT7yXg34Xo6cJp6Ii4qb+RJxxxdzIjPwtZtu62XC28RtlSVxUqPjzgIOJdG5joB5SBYEIFxju4Mpfn8YCeOTlwwHmtKOhTbTm9YNkcTVS2OuBJtFuqDahwnIh1wnb6M8vnbvU0omF7YtryvW6oySw1EXo3Fs1im70OtkD3GRGep0PTwzWBzshzxB9I7jA0bAzV4e2SYdEw5CYeGRTqQ0o8k1nPSzmNvGBBiGCLAhk0Z1BxN86qD1onfc5ALFeI1mtO05lPhMbA9TjiTnSVCLWkCuerWWIz+LlftMr2Uhld5+wSWFnWLGee1P6/cMG2nJkfw1QIDAQAB"

/*** 配置结束 ***/

type AlipayService struct {
	appId           string
	returnUrl       string
	notifyUrl       string
	charset         string
	rsaPrivateKey   *rsa.PrivateKey
	alipayPublicKey *rsa.PublicKey
	totalFee        float64
	outTradeNo      string
	orderName       string
}

func (this *AlipayService) SetAppId(appId string) {
	this.appId = appId
}

func (this *AlipayService) SetReturnUrl(returnUrl string) {
	this.returnUrl = returnUrl
}

func (this *AlipayService) SetNotifyUrl(notifyUrl string) {
	this.notifyUrl = notifyUrl
}

func (this *AlipayService) SetCharset(charset string) {
	this.charset = charset
}

func (this *AlipayService) SetRsaPrivateKey(keyString string) {

	privateKey, err := ParsePrivateKey(FormatPrivateKey(keyString))
	if err != nil {
		panic(err)
	}
	this.rsaPrivateKey = privateKey
}

func (this *AlipayService) SetAlipayPublicKey(keyString string) {
	alipayPublicKey, err := ParsePublicKey(FormatPublicKey(keyString))
	if err != nil {
		panic(err)
	}
	this.alipayPublicKey = alipayPublicKey
}

func (this *AlipayService) SetTotalFee(totalFee float64) {
	this.totalFee = totalFee
}

func (this *AlipayService) SetOutTradeNo(outTradeNo string) {
	this.outTradeNo = outTradeNo
}

func (this *AlipayService) SetOrderName(orderName string) {
	this.orderName = orderName
}

// GenSign 产生签名
func (this *AlipayService) GenSign(m map[string]string) string {
	var data []string
	var encryptedBytes []byte
	for k, v := range m {
		if v != "" && k != "sign" {
			data = append(data, fmt.Sprintf(`%s=%s`, k, v))
		}
	}
	sort.Strings(data)
	signData := strings.Join(data, "&")
	s := sha256.New()
	_, err := s.Write([]byte(signData))
	if err != nil {
		panic(err)
	}
	hashByte := s.Sum(nil)
	hashs := crypto.SHA256
	rsaPrivateKey := this.rsaPrivateKey
	if encryptedBytes, err = rsa.SignPKCS1v15(rand.Reader, rsaPrivateKey, hashs, hashByte); err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(encryptedBytes)
}

func (this *AlipayService) buildRequestForm(m map[string]string) string {
	var buf []string
	buf = append(buf, fmt.Sprintf("正在跳转至支付页面...<form id='alipaysubmit' name='alipaysubmit' action='https://openapi.alipay.com/gateway.do?charset=%s' method='POST'>", this.charset))
	for k, v := range m {
		buf = append(buf, fmt.Sprintf("<input type='hidden' name='%s' value='%s'/>", k, v))
	}
	buf = append(buf, "<input type='submit' value='ok' style='display:none;''></form>")
	buf = append(buf, "<script>document.forms['alipaysubmit'].submit();</script>")
	return strings.Join(buf, "")
}

func (this AlipayService) DoPay() (string, error) {

	//请求参数
	var bizContent = make(map[string]string)
	bizContent["out_trade_no"] = this.outTradeNo
	bizContent["product_code"] = "FAST_INSTANT_TRADE_PAY"
	bizContent["total_amount"] = strconv.FormatFloat(float64(this.totalFee), 'f', 2, 64) //2表示保留2位小数
	bizContent["subject"] = this.orderName
	bizContentJson, err := json.Marshal(bizContent)
	if err != nil {
		return "", errors.New("json.Marshal: " + err.Error())
	}

	//公共参数
	var m = make(map[string]string)
	m["app_id"] = this.appId
	m["method"] = "alipay.trade.page.pay" //接口名称
	m["format"] = "JSON"
	m["return_url"] = this.returnUrl
	m["charset"] = this.charset
	m["sign_type"] = "RSA2"
	m["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
	m["version"] = "1.0"
	m["notify_url"] = this.notifyUrl
	m["biz_content"] = string(bizContentJson)

	//获取签名
	sign := this.GenSign(m)
	m["sign"] = sign
	return this.buildRequestForm(m), nil
}

func (this *AlipayService) VerifySign(data url.Values) (ok bool, err error) {
	return verifySign(data, this.alipayPublicKey)
}

func ParsePrivateKey(data []byte) (key *rsa.PrivateKey, err error) {
	var block *pem.Block
	block, _ = pem.Decode(data)
	if block == nil {
		return nil, errors.New("private key failed to load")
	}

	var priInterface interface{}
	priInterface, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	key, ok := priInterface.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("private key failed to load")
	}

	return key, err
}

func ParsePublicKey(data []byte) (key *rsa.PublicKey, err error) {
	var block *pem.Block
	block, _ = pem.Decode(data)
	if block == nil {
		return nil, errors.New("alipay public key failed to load")
	}

	var pubInterface interface{}
	pubInterface, err = x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	key, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("alipay public key failed to load")
	}

	return key, err
}

func FormatPublicKey(raw string) []byte {
	return formatKey(raw, "-----BEGIN PUBLIC KEY-----", "-----END PUBLIC KEY-----", 64)
}

func FormatPrivateKey(raw string) []byte {
	return formatKey(raw, "-----BEGIN RSA PRIVATE KEY-----", "-----END RSA PRIVATE KEY-----", 64)
}

func formatKey(raw, prefix, suffix string, lineCount int) []byte {
	if raw == "" {
		return nil
	}
	raw = strings.Replace(raw, prefix, "", 1)
	raw = strings.Replace(raw, suffix, "", 1)
	raw = strings.Replace(raw, " ", "", -1)
	raw = strings.Replace(raw, "\n", "", -1)
	raw = strings.Replace(raw, "\r", "", -1)
	raw = strings.Replace(raw, "\t", "", -1)

	var sl = len(raw)
	var c = sl / lineCount
	if sl%lineCount > 0 {
		c = c + 1
	}

	var buf bytes.Buffer
	buf.WriteString(prefix + "\n")
	for i := 0; i < c; i++ {
		var b = i * lineCount
		var e = b + lineCount
		if e > sl {
			buf.WriteString(raw[b:])
		} else {
			buf.WriteString(raw[b:e])
		}
		buf.WriteString("\n")
	}
	buf.WriteString(suffix)
	return buf.Bytes()
}

func verifySign(data url.Values, key *rsa.PublicKey) (ok bool, err error) {
	sign := data.Get("sign")

	var keys = make([]string, 0, 0)
	for key := range data {
		if key == "sign" || key == "sign_type" {
			continue
		}
		keys = append(keys, key)
	}

	sort.Strings(keys)

	var pList = make([]string, 0, 0)
	for _, key := range keys {
		pList = append(pList, key+"="+data.Get(key))
	}
	var s = strings.Join(pList, "&")

	return verifyData([]byte(s), sign, key)
}

func verifyData(data []byte, sign string, key *rsa.PublicKey) (ok bool, err error) {
	signBytes, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return false, err
	}

	var h = crypto.SHA256.New()
	h.Write(data)
	var hashed = h.Sum(nil)

	if err = rsa.VerifyPKCS1v15(key, crypto.SHA256, hashed, signBytes); err != nil {
		return false, err
	}
	return true, nil
}

//生成订单号
func Uniqid() string {
	now := time.Now()
	return fmt.Sprintf("%s%08x%05x", "", now.Unix(), now.UnixNano()%0x100000)
}

//发起支付
func pcPay(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	if params.Get("total_fee") != "" {
		payAmount, _ = strconv.ParseFloat(params.Get("total_fee"), 64)
	}
	outTradeNo := Uniqid()
	if params.Get("out_trade_no") != "" {
		outTradeNo = params.Get("out_trade_no")
	}
	if params.Get("order_name") != "" {
		orderName = params.Get("order_name")
	}
	aliPay := AlipayService{}
	aliPay.SetAppId(appid)
	aliPay.SetReturnUrl(returnUrl)
	aliPay.SetNotifyUrl(notifyUrl)
	aliPay.SetCharset("utf-8")
	aliPay.SetRsaPrivateKey(rsaPrivateKey)
	aliPay.SetTotalFee(payAmount)
	aliPay.SetOutTradeNo(outTradeNo)
	aliPay.SetOrderName(orderName)
	sHtml, err := aliPay.DoPay()
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, sHtml)
}

//异步回调通知处理
func notify(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	//接收post参数
	r.ParseForm()
	formdata := r.PostForm

	aliPay := AlipayService{}
	aliPay.SetCharset("utf-8")
	aliPay.SetAlipayPublicKey(alipayPublicKey)
	_, err := aliPay.VerifySign(formdata)
	if err != nil {
		fmt.Fprintf(w, "error:"+err.Error())
		return
	}
	//处理你的逻辑，例如获取订单号formdata["out_trade_no"][0]，订单金额formdata["total_amount"][0]等
	//例如这里将回调数据写入到nofity.txt文本中

	ioutil.WriteFile("nofity.txt", []byte(formdata.Encode()), 0644)
	//程序执行完后必须打印输出“success”（不包含引号）。如果商户反馈给支付宝的字符不是success这7个字符，支付宝服务器会不断重发通知，直到超过24小时22分钟。一般情况下，25小时以内完成8次通知（通知的间隔频率一般是：4m,10m,10m,1h,2h,6h,15h）；
	fmt.Fprintf(w, "success")
}

//同步回调
func callback(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	//接收get参数
	params := r.URL.Query()

	aliPay := AlipayService{}
	aliPay.SetCharset("utf-8")
	aliPay.SetAlipayPublicKey(alipayPublicKey)
	_, err := aliPay.VerifySign(params)
	if err != nil {
		fmt.Fprintf(w, "error:"+err.Error())
		return
	}
	fmt.Fprintf(w, "支付成功！订单号：%s", params.Get("out_trade_no"))
}

//func main() {
//	//开启web服务器
//	port := "8080"
//	http.HandleFunc("/", pcPay)
//	http.HandleFunc("/notify", notify)
//	http.HandleFunc("/callback", callback)
//	log.Println("服务器启动成功，监听端口：", port)
//	err := http.ListenAndServe(":"+port, nil)
//	if err != nil {
//		log.Fatal("ListenAndServe:", err)
//	}
//}
