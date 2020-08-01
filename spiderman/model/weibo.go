package model

import "net/http"

type WeiboConf struct {
	Comments map[string][]string `yaml:"comments"`
	Tick     int                 `yaml:"tick"`
	Token    WeiboToken          `yaml:"token"`
	Common   struct {
		Comments []string `yaml:"comments"`
	} `yaml:"common"`
}

type WeiboToken struct {
	SUB  string `yaml:"sub"`
	SUHB string `yaml:"suhb"`
}

type WeiboCommonResp struct {
	OK   int    `json:"ok"`
	Msg  string `json:"msg"`
	Code string `json:"errno"`
}

type WeiboAuthConf struct {
	WeiboCommonResp
	Cookies []*http.Cookie
	Data    struct {
		St string `json:"st"`
	} `json:"data"`
}

type WeiboUserInfoStatus struct {
	WeiboCommonResp
	MID string
}

type WeiboUserInfoResp struct {
	WeiboCommonResp
	Data struct {
		Statuses []WeiboUserInfoStatus `json:"statuses"`
	} `json:"data"`
}
