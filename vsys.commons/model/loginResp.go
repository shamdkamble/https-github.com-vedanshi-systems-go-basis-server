package model

type ApiResp struct {
	Code         int         `json:"code"`
	Token        string      `json:"token"`
	Expire       string      `json:"expire"`
	EndUserRoles interface{} `json:"endUserRoles"`
	EndUsers     interface{} `json:"endUser"`
}
