package tinybot

import "github.com/gofrs/uuid"

type pluginOptions struct {
	Help  map[string]string
	Admin []int64
}

type Plugin struct {
	Id     uuid.UUID
	Name   string
	Admin  []int64
	Bot    *Bot //所有插件共享Bot实例
	Allies *Plugin

	middlewares []func(ctx *Context) //事件处理的中间件
}

type Module struct {
	pluginId uuid.UUID
	Enable   bool
	Cond     func(ctx *Context) bool
	Action   func(ctx *Context)
}

// RegPlugin 注册一个插件
func (p *Plugin) RegPlugin(name string) *Plugin {
	bot := p.Bot
	plugin := &Plugin{
		Id:   uuid.Must(uuid.NewV4()),
		Bot:  bot,
		Name: name,
	}
	bot.plugins[plugin.Id] = plugin

	return plugin
}

func (p *Plugin) Add(cond CondFunc, handler HandlerFunc) {
	module := Module{
		pluginId: p.Id,
		Enable:   true,
		Cond:     cond,
		Action:   handler,
	}
	p.Bot.Modules = append(p.Bot.Modules, module)
}

// Use 添加全局中间件
func (p *Plugin) Use(middlewares ...func(ctx *Context)) {
	p.middlewares = append(p.middlewares, middlewares...)
}
