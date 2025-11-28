package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/lvjiaben/go-wheel/pkg/container"
	"github.com/lvjiaben/go-wheel/pkg/httpclient"
)

// 验证码过期时间（秒）
const SmsCodeExpire = 300 // 5分钟

// 发送间隔时间（秒）
const SmsCodeInterval = 60 // 60秒

// 验证码位数
const SmsCodeLength = 6 // 4位或6位

// 缓存键前缀
const SmsCachePrefix = "sms_code:"

// SmsService 短信服务
type SmsService struct {
	container   *container.Container
	configCache *ConfigCacheService
}

// NewSmsService 创建短信服务
func NewSmsService(c *container.Container) *SmsService {
	return &SmsService{
		container:   c,
		configCache: NewConfigCacheService(c),
	}
}

// Send 发送短信验证码
func (s *SmsService) Send(mobile string, code string, event string) error {
	if code == "" {
		code = s.generateCode()
	}
	if event == "" {
		event = "default"
	}

	ctx := context.Background()
	rdb := s.container.GetRDB()

	// 检查发送间隔
	intervalKey := s.getCacheKey(mobile, event)
	if exists, _ := rdb.Exists(ctx, intervalKey).Result(); exists > 0 {
		return errors.New("send_too_frequent")
	}

	// 获取短信类型
	smsType := s.configCache.Get("sms_type")
	if smsType == "" {
		return errors.New("sms_type_not_configured")
	}

	// 获取模板ID
	template := s.configCache.Get("sms_template_" + event)
	if template == "" {
		template = s.configCache.Get("sms_template")
	}
	if template == "" {
		return errors.New("sms_template_not_configured")
	}

	// 根据类型发送短信
	var err error
	switch smsType {
	case "aliyun":
		err = s.sendAliyun(mobile, template, code)
	case "tencent":
		err = s.sendTencent(mobile, template, code)
	case "yunpian":
		err = s.sendYunpian(mobile, template, code)
	case "smsbao":
		err = s.sendSmsbao(mobile, template, code)
	default:
		return errors.New("sms_type_unsupported")
	}

	if err != nil {
		return err
	}

	// 存储验证码到Redis
	codeKey := s.getCacheKey(mobile, event) + "_code"
	rdb.Set(ctx, intervalKey, "1", time.Duration(SmsCodeInterval)*time.Second)
	rdb.Set(ctx, codeKey, code, time.Duration(SmsCodeExpire)*time.Second)

	return nil
}

// Verify 验证短信验证码
func (s *SmsService) Verify(mobile string, code string, event string) bool {
	if event == "" {
		event = "default"
	}

	ctx := context.Background()
	codeKey := s.getCacheKey(mobile, event) + "_code"
	storedCode, err := s.container.GetRDB().Get(ctx, codeKey).Result()
	if err != nil {
		return false
	}

	return storedCode == code
}

// Delete 删除验证码
func (s *SmsService) Delete(mobile string, event string) {
	if event == "" {
		event = "default"
	}

	ctx := context.Background()
	rdb := s.container.GetRDB()
	intervalKey := s.getCacheKey(mobile, event)
	codeKey := intervalKey + "_code"

	rdb.Del(ctx, intervalKey, codeKey)
}

// getCacheKey 获取缓存键
func (s *SmsService) getCacheKey(mobile string, event string) string {
	return fmt.Sprintf("%s%s_%s", SmsCachePrefix, mobile, event)
}

// generateCode 生成随机验证码
func (s *SmsService) generateCode() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// 根据位数生成验证码
	switch SmsCodeLength {
	case 6:
		return fmt.Sprintf("%06d", r.Intn(1000000))
	default:
		return fmt.Sprintf("%04d", r.Intn(10000))
	}
}

// GetExpire 获取过期时间
func (s *SmsService) GetExpire() int {
	return SmsCodeExpire
}

// GetInterval 获取发送间隔
func (s *SmsService) GetInterval() int {
	return SmsCodeInterval
}

// sendAliyun 阿里云短信发送
func (s *SmsService) sendAliyun(mobile string, template string, code string) error {
	accessKeyId := s.configCache.Get("sms_id")
	accessKeySecret := s.configCache.Get("sms_key")
	signName := s.configCache.Get("sms_token")

	if accessKeyId == "" || accessKeySecret == "" || signName == "" {
		return errors.New("aliyun_sms_not_configured")
	}

	// 构建请求参数
	params := map[string]string{
		"AccessKeyId":      accessKeyId,
		"Action":           "SendSms",
		"Format":           "JSON",
		"PhoneNumbers":     mobile,
		"SignName":         signName,
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureNonce":   fmt.Sprintf("%d", time.Now().UnixNano()),
		"SignatureVersion": "1.0",
		"TemplateCode":     template,
		"TemplateParam":    fmt.Sprintf(`{"code":"%s"}`, code),
		"Timestamp":        time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		"Version":          "2017-05-25",
	}

	// 生成签名
	signature := s.aliyunSign(params, accessKeySecret)
	params["Signature"] = signature

	// 使用 httpclient 发送请求
	client := httpclient.NewClient()
	resp, err := client.Get("https://dysmsapi.aliyuncs.com").
		SetQueryParams(params).
		Send()
	if err != nil {
		return fmt.Errorf("aliyun_request_failed: %v", err)
	}
	defer resp.Close()

	var result map[string]any
	if err := resp.JSON(&result); err != nil {
		return fmt.Errorf("aliyun_response_parse_failed: %v", err)
	}

	if result["Code"] != "OK" {
		return fmt.Errorf("aliyun_send_failed: %v", result["Message"])
	}

	return nil
}

// aliyunSign 阿里云签名
func (s *SmsService) aliyunSign(params map[string]string, secret string) string {
	// 按key排序
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 构建待签名字符串
	var queryParts []string
	for _, k := range keys {
		queryParts = append(queryParts, url.QueryEscape(k)+"="+url.QueryEscape(params[k]))
	}
	queryString := strings.Join(queryParts, "&")
	stringToSign := "GET&" + url.QueryEscape("/") + "&" + url.QueryEscape(queryString)

	// HMAC-SHA1 签名
	mac := hmac.New(sha256.New, []byte(secret+"&"))
	mac.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// sendTencent 腾讯云短信发送
func (s *SmsService) sendTencent(mobile string, template string, code string) error {
	appId := s.configCache.Get("sms_id")
	appKey := s.configCache.Get("sms_key")
	sign := s.configCache.Get("sms_token")

	if appId == "" || appKey == "" {
		return errors.New("tencent_sms_not_configured")
	}

	// 构建请求
	timestamp := time.Now().Unix()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomNum := fmt.Sprintf("%d", r.Intn(999999))

	// 构建签名
	sigStr := fmt.Sprintf("appkey=%s&random=%s&time=%d&mobile=%s", appKey, randomNum, timestamp, mobile)
	h := sha256.New()
	h.Write([]byte(sigStr))
	sig := hex.EncodeToString(h.Sum(nil))

	// 构建请求体
	reqBody := map[string]any{
		"ext":    "",
		"extend": "",
		"params": []string{code, fmt.Sprintf("%d", SmsCodeExpire/60)},
		"sig":    sig,
		"sign":   sign,
		"tel": map[string]string{
			"mobile":     mobile,
			"nationcode": "86",
		},
		"time":   timestamp,
		"tpl_id": template,
	}

	reqUrl := fmt.Sprintf("https://yun.tim.qq.com/v5/tlssmssvr/sendsms?sdkappid=%s&random=%s", appId, randomNum)

	// 使用 httpclient 发送请求
	client := httpclient.NewClient()
	resp, err := client.Post(reqUrl).
		SetJSON(reqBody).
		Send()
	if err != nil {
		return fmt.Errorf("tencent_request_failed: %v", err)
	}
	defer resp.Close()

	var result map[string]any
	if err := resp.JSON(&result); err != nil {
		return fmt.Errorf("tencent_response_parse_failed: %v", err)
	}

	if result["result"].(float64) != 0 {
		return fmt.Errorf("tencent_send_failed: %v", result["errmsg"])
	}

	return nil
}

// sendYunpian 云片短信发送
func (s *SmsService) sendYunpian(mobile string, template string, code string) error {
	apiKey := s.configCache.Get("sms_key")

	if apiKey == "" {
		return errors.New("yunpian_sms_not_configured")
	}

	// 使用 httpclient 发送请求
	client := httpclient.NewClient()
	resp, err := client.Post("https://sms.yunpian.com/v2/sms/tpl_single_send.json").
		SetForm("apikey", apiKey).
		SetForm("mobile", mobile).
		SetForm("tpl_id", template).
		SetForm("tpl_value", url.QueryEscape(fmt.Sprintf("#code#=%s&#minute#=%d", code, SmsCodeExpire/60))).
		Send()
	if err != nil {
		return fmt.Errorf("yunpian_request_failed: %v", err)
	}
	defer resp.Close()

	var result map[string]any
	if err := resp.JSON(&result); err != nil {
		return fmt.Errorf("yunpian_response_parse_failed: %v", err)
	}

	if result["code"].(float64) != 0 {
		return fmt.Errorf("yunpian_send_failed: %v", result["msg"])
	}

	return nil
}

// sendSmsbao 短信宝短信发送
func (s *SmsService) sendSmsbao(mobile string, template string, code string) error {
	username := s.configCache.Get("sms_id")
	password := s.configCache.Get("sms_key") // MD5后的密码或ApiKey
	sign := s.configCache.Get("sms_token")   // 短信签名

	if username == "" || password == "" {
		return errors.New("smsbao_sms_not_configured")
	}

	// 构建短信内容：【签名】您的验证码是xxx
	content := fmt.Sprintf("【%s】您的验证码是%s，请尽快验证", sign, code)

	// 使用 httpclient 发送请求
	client := httpclient.NewClient()
	resp, err := client.Get("https://api.smsbao.com/sms").
		SetQuery("u", username).
		SetQuery("p", password).
		SetQuery("m", mobile).
		SetQuery("c", content).
		Send()
	if err != nil {
		return fmt.Errorf("smsbao_request_failed: %v", err)
	}
	defer resp.Close()

	// 短信宝返回 "0" 表示成功
	body, err := resp.String()
	if err != nil {
		return fmt.Errorf("smsbao_response_read_failed: %v", err)
	}

	if body != "0" {
		// 错误码对应
		errMsg := s.smsbaoErrorMsg(body)
		return fmt.Errorf("smsbao_send_failed: %s", errMsg)
	}

	return nil
}

// smsbaoErrorMsg 短信宝错误码转换
func (s *SmsService) smsbaoErrorMsg(code string) string {
	switch code {
	case "30":
		return "密码错误"
	case "40":
		return "账号不存在"
	case "41":
		return "余额不足"
	case "43":
		return "IP地址限制"
	case "50":
		return "内容含有敏感词"
	case "51":
		return "手机号码不正确"
	default:
		return "未知错误: " + code
	}
}
