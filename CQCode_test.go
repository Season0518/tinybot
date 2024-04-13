package tinybot

import (
	"core/models"
	"fmt"
	"testing"
)

func TestSerializeCQCode(t *testing.T) {
	var messageChain []models.MessageBody

	messageChain = append(messageChain, models.MessageBody{
		Type: "at",
		Data: map[string]string{
			"qq": "all",
		},
	})

	platform := "测试平台"
	text := "测试内容✨🙏😩*&^%$$#你好Helloこんにちは안녕하세요🙏✨💗❤️"
	refer := "https://www.douyin.com/user/MS4wLjABAAAAu-A0s9aIathifzLqcPvwBvMaOIA5XGicTEU8wc1dilk"

	messageChain = append(messageChain, models.MessageBody{
		Type: "text",
		Data: map[string]string{
			"text":    fmt.Sprintf("公主在%s发布了内容！\\n\\n%s\\n\\n快点击%s围观吧！", platform, text, refer),
			"subType": "0",
		},
	})

	messageChain = append(messageChain, models.MessageBody{
		Type: "image",
		Data: map[string]string{
			"file":    "https://ok",
			"subType": "0",
		},
	})

	cqMessage := SerializeCQCode(messageChain)
	fmt.Println(cqMessage)

	jsonStr := fmt.Sprintf(`{"group_id": %d, "message": "%s", "auto_escape": %v}`, 865444787, cqMessage, false)

	fmt.Println(jsonStr)

}
