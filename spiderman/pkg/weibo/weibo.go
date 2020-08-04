package weibo

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/fitzix/spider/model"
	"github.com/fitzix/spider/pkg/logger"
	"github.com/go-resty/resty/v2"
)

const (
	weiboUrl         = "https://m.weibo.cn/"
	weiboLoginUrl    = "https://passport.weibo.cn/sso/login"
	weiboStUrl       = "https://m.weibo.cn/api/config"
	weiboCommentUrl  = "https://m.weibo.cn/api/comments/create"
	weiboUserInfoUrl = "https://m.weibo.cn/profile/info"
)

var (
	httpClient    *resty.Client
	weiboAuthConf model.WeiboAuthConf
	stTimer       *time.Timer
	twmTimer      *time.Timer
	mutex         sync.Mutex
)

func Init(token model.WeiboToken) {
	httpClient = resty.New().SetTimeout(10*time.Second).
		SetHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.1.2 Safari/605.1.15").
		SetHeader("Referer", "https://m.weibo.cn/")
	weiboAuthConf = model.WeiboAuthConf{Cookies: []*http.Cookie{
		{
			Name:  "SUB",
			Value: token.SUB,
		},
		{
			Name:  "SUHB",
			Value: token.SUHB,
		},
		nil,
	}}
	getTwm()
	getSt()
}

func getTwm() {
	mutex.Lock()
	defer mutex.Unlock()
	resp, err := httpClient.R().Get(weiboUrl)
	if err != nil {
		logger.Errorf("获取twm失败---> %s", err)
		return
	}

	for _, v := range resp.Cookies() {
		if v.Name == "_T_WM" {
			weiboAuthConf.Cookies[2] = v
		}
	}

	twmTimer = time.AfterFunc(9*24*time.Hour, func() {
		getTwm()
	})
}

func getSt() {
	mutex.Lock()
	defer mutex.Unlock()
	resp, err := httpClient.R().SetCookies(weiboAuthConf.Cookies).Get(weiboStUrl)
	if err != nil {
		logger.Errorf("获取st失败---> %s", err)
		return
	}
	if err := json.Unmarshal(resp.Body(), &weiboAuthConf); err != nil {
		logger.Errorf("解析st失败---> %s", err)
		return
	}
	if weiboAuthConf.OK != 1 {
		logger.Errorf("获取st失败---> %s", weiboAuthConf.Msg)
		return
	}
	stTimer = time.AfterFunc(10*time.Hour, func() {
		getSt()
	})
}

func CreateComment(mid, content string) error {
	resp, err := httpClient.R().SetCookies(weiboAuthConf.Cookies).SetFormData(map[string]string{
		"mid":     mid,
		"content": content,
		"st":      weiboAuthConf.Data.St,
	}).Post(weiboCommentUrl)
	if err != nil {
		return err
	}

	var commentResp model.WeiboCommonResp
	if err := json.Unmarshal(resp.Body(), &commentResp); err != nil {
		return err
	}

	if commentResp.OK != 1 {
		if commentResp.Code == "100006" {
			refreshToken()
		}
		return errors.New(commentResp.Msg)
	}

	return nil
}

func GetLatestWeibo(uid string) (string, error) {
	resp, err := httpClient.R().SetQueryParam("uid", uid).Get(weiboUserInfoUrl)
	if err != nil {
		return "", err
	}
	var userInfo model.WeiboUserInfoResp
	if err := json.Unmarshal(resp.Body(), &userInfo); err != nil {
		return "", err
	}

	if userInfo.OK == 1 && userInfo.Data.Statuses != nil && len(userInfo.Data.Statuses) > 0 {
		return userInfo.Data.Statuses[0].MID, nil
	}
	return "", errors.New(userInfo.Msg)
}

func refreshToken() {
	logger.Info("刷新token")
	twmTimer.Stop()
	stTimer.Stop()
	getTwm()
	getSt()
}
