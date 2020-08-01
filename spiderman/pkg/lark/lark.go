package lark

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/allegro/bigcache/v2"
	"github.com/fitzix/spider/model"
	"github.com/go-resty/resty/v2"
)

var (
	httpClient *resty.Client
	cache      *bigcache.BigCache
)

const (
	// URL
	larkUrlTenantToken = "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal"
	larkUrlUpload      = "https://open.feishu.cn/open-apis/image/v4/put"
	larkUrlMsgSend     = "https://open.feishu.cn/open-apis/message/v4/send"
	// CACHE
	larkCacheKeyTenantToken = "lark-token-tenant"
)

func Init() {
	httpClient = resty.New().SetTimeout(10 * time.Second)
	cache, _ = bigcache.NewBigCache(bigcache.DefaultConfig(2 * time.Hour))
}

func getTenantToken() (string, error) {
	if token, err := cache.Get(larkCacheKeyTenantToken); err == nil {
		return string(token), nil
	}
	resp, err := httpClient.R().SetBody(map[string]string{
		"app_id":     "cli_9f8c308fbdae900c",
		"app_secret": "gCRmXFr39z1a3eQMzclfWbvv0CwiHkTY",
	}).Post(larkUrlTenantToken)
	if err != nil {
		return "", err
	}
	var tokenResp model.LarkTokenResp
	if err := json.Unmarshal(resp.Body(), &tokenResp); err != nil {
		return "", err
	}
	if tokenResp.Code != 0 {
		return "", errors.New(tokenResp.Msg)
	}
	if tokenResp.Expire < 7200 {
		cache, _ = bigcache.NewBigCache(bigcache.DefaultConfig(time.Duration(tokenResp.Expire) * time.Second))
	}

	token := fmt.Sprintf("Bearer %s", tokenResp.TenantAccessToken)
	_ = cache.Set(larkCacheKeyTenantToken, []byte(token))
	log.Printf("token: %s", token)
	return token, nil
}

func UploadImage(buf []byte) (*model.LarkUploadFileResp, error) {
	authToken, err := getTenantToken()
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.R().SetAuthToken(authToken).
		SetFileReader("image", "image", bytes.NewReader(buf)).
		SetFormData(map[string]string{"image_type": "message"}).
		Post(larkUrlUpload)
	if err != nil {
		return nil, err
	}
	var uploadResp model.LarkUploadFileResp
	if err = json.Unmarshal(resp.Body(), &uploadResp); err != nil {
		return nil, err
	}

	if uploadResp.Code != 0 {
		return nil, errors.New(uploadResp.Msg)
	}
	return &uploadResp, err
}

func SendMsg(msg model.LarkMessage) error {
	authToken, err := getTenantToken()
	if err != nil {
		return err
	}
	resp, err := httpClient.R().SetAuthToken(authToken).SetBody(msg).Post(larkUrlMsgSend)
	if err != nil {
		return err
	}
	var r model.LarkCommonResp
	if err := json.Unmarshal(resp.Body(), &r); err != nil {
		return err
	}
	if r.Code != 0 {
		return errors.New(r.Msg)
	}
	return nil
}
