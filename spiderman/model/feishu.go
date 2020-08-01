package model

type LarkCommonResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type LarkUploadFileResp struct {
	LarkCommonResp
	Data struct {
		ImageKey string `json:"image_key"`
	} `json:"data"`
}

type LarkMessage struct {
	ChatID  string         `json:"chat_id,omitempty"`
	OpenId  string         `json:"open_id,omitempty"`
	UserId  string         `json:"user_id,omitempty"`
	Email   string         `json:"email,omitempty"`
	MsgType string         `json:"msg_type"`
	Content LarkMsgContent `json:"content"`
}

type LarkMsgContent struct {
	Post LarkMsgContentPost `json:"post"`
}

type LarkMsgContentPost struct {
	ZhCN LarkMsgPostLanguage `json:"zh_cn"`
}

type LarkMsgPostLanguage struct {
	Title   string                         `json:"title"`
	Content [][]LarkMsgPostLanguageContent `json:"content"`
}

type LarkMsgPostLanguageContent struct {
	Tag      string `json:"tag"`
	ImageKey string `json:"image_key"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

type LarkTokenResp struct {
	LarkCommonResp
	Expire            int64  `json:"expire"`
	TenantAccessToken string `json:"tenant_access_token"`
}
