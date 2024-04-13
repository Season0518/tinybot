package tinybot

import "github.com/tidwall/gjson"

type APIRequest struct {
	Action string `json:"action,omitempty"`
	Params Params `json:"params,omitempty"`
	Echo   string `json:"echo,omitempty"`
}

type APIResponse struct {
	Status  string       `json:"status"`
	Data    gjson.Result `json:"data"`
	Message string       `json:"message"`
	Wording string       `json:"wording"`
	RetCode int64        `json:"retcode"`
	Echo    string       `json:"echo"`
}

type Api interface {
	SendGroupMsg(groupId int64, messageChain []MsgBody) error
	GetGroupInfo(groupId int64, noCache bool) error
	GetMemberList(groupId int64, noCache bool) error
	SetGroupAddRequest(flag, subType string, approve bool, reason string) error //缺少ctx
	GetStrangerInfo(userId int64, noCache bool) error
}
