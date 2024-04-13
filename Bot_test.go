package tinybot_test

import (
	"testing"
	"time"
	"tinybot"
)

func TestBotFlow(t *testing.T) {
	mockClient := tinybot.NewMockDriver()

	event := tinybot.Event{PostType: "message", Message: tinybot.NewMsgBuilder().Text("TestEvent").Build()}
	mockClient.SendEvent(event)

	logger := func() func(ctx *tinybot.Context) {
		return func(ctx *tinybot.Context) {
			st := time.Now()
			ctx.Next()
			t.Logf("Processed Event, Type: %s, Elapsed time: %s", ctx.Event.PostType, time.Since(st))
		}
	}

	bot := tinybot.NewBot(mockClient)
	bot.Use(logger())

	plugin := bot.RegPlugin("testPlugin")

	cond := func(ctx *tinybot.Context) bool {
		// PostType 修改为Enum类型
		if ctx.Event.PostType == "message" {
			return true
		}
		return false
	}

	action := func(ctx *tinybot.Context) {
		t.Logf("Received Message: %s", ctx.Event.Message)
	}

	middleware := func() func(ctx *tinybot.Context) {
		return func(ctx *tinybot.Context) {
			t.Logf("Do something before action...")
		}
	}

	plugin.Use(middleware())
	plugin.Add(cond, action)

	bot.ListenAndServe()
}

func TestExceptionRecover(t *testing.T) {
	mockClient := tinybot.NewMockDriver()

	event := tinybot.Event{PostType: "message", Message: tinybot.NewMsgBuilder().Text("TestEvent").Build()}

	mockClient.SendEvent(event)

	bot := tinybot.NewBot(mockClient)
	bot.ListenAndServe()
}
