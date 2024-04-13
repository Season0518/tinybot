package tinybot

import (
	"errors"
	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"time"
)

type HandlerFunc func(ctx *Context)
type CondFunc func(ctx *Context) bool
type Params map[string]any

type Bot struct {
	*Plugin

	Driver  Driver
	Account int64
	Latency time.Duration
	Modules []Module
	plugins map[uuid.UUID]*Plugin
}

type Driver interface {
	Dial()
	CallApi(req APIRequest) (APIResponse, error)
	RetrieveEvent() (Event, error)
}

func NewBot(driver Driver) *Bot {
	bot := &Bot{Driver: driver, plugins: make(map[uuid.UUID]*Plugin)}

	// 默认插件将应用于所有事件
	defaultPlugin := &Plugin{Id: uuid.Must(uuid.NewV4()), Bot: bot}
	defaultModule := Module{
		pluginId: defaultPlugin.Id,
		Cond:     func(ctx *Context) bool { return true },
		Action:   func(ctx *Context) {},
	}

	bot.Plugin = defaultPlugin
	bot.Modules = append(bot.Modules, defaultModule)
	bot.plugins[defaultPlugin.Id] = defaultPlugin

	return bot
}

// ListenAndServe 启动机器人
func (b *Bot) ListenAndServe() {
	b.Driver.Dial()

	for {
		event, err := b.Driver.RetrieveEvent()
		if err != nil {
			var closeError *websocket.CloseError
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Warn().Msg("[Bot] 链接关闭")
				return
			} else if errors.As(err, &closeError) {
				log.Warn().Msgf("[Bot] 链接异常关闭，程序已自动退出。错误代码：%d", err.(*websocket.CloseError).Code)
				return
			} else {
				log.Warn().Msgf("[bot]: %s", err.Error())
			}
		}

		ctx := &Context{index: -1, Event: &event}
		for _, m := range b.Modules {
			if m.Cond(ctx) {
				//一次扩容，减少内部数组扩容的次数，提高效率
				handlers := make([]func(ctx *Context), len(b.plugins[m.pluginId].middlewares)+1)
				copy(handlers, b.plugins[m.pluginId].middlewares)
				handlers[len(b.plugins[m.pluginId].middlewares)] = m.Action
				ctx.handlers = append(ctx.handlers, handlers...)
			}
		}
		ctx.Next()
	}
}

func (b *Bot) CallApi(req APIRequest) (APIResponse, error) {
	return b.Driver.CallApi(req)
}
