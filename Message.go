package tinybot

import "encoding/json"

type MsgBody struct {
	Type string         `json:"type"`
	Data map[string]any `json:"data"`
}

// MsgBuilder https://github.com/botuniverse/onebot-11/tree/master/message/array.md#%E6%95%B0%E7%BB%84%E6%A0%BC%E5%BC%8F
// MsgBuilder: 构建消息串
type MsgBuilder []MsgBody

func NewMsgBuilder() *MsgBuilder {
	return &MsgBuilder{}
}

func (b *MsgBuilder) Text(t string) *MsgBuilder {
	*b = append(*b, MsgBody{Type: "text", Data: map[string]any{"text": t}})
	return b
}

func (b *MsgBuilder) Image(url string, cache bool) *MsgBuilder {
	*b = append(*b, MsgBody{Type: "image", Data: map[string]any{"file": url, "cache": cache}})
	return b
}

func (b *MsgBuilder) At(qq string) *MsgBuilder {
	*b = append(*b, MsgBody{Type: "at", Data: map[string]any{"qq": qq}})
	return b
}

func (b *MsgBuilder) ToByte() []byte {
	data, _ := json.Marshal(b)
	return data
}

func (b *MsgBuilder) Build() MsgBuilder {
	return *b
}
