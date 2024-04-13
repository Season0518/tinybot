package tinybot

type Context struct {
	// bot data
	Event  *Event
	Bot    *Bot
	Plugin *Plugin

	// middleware
	handlers []func(ctx *Context)
	index    int
}

func (ctx *Context) Next() {
	ctx.index++
	for ; ctx.index < len(ctx.handlers); ctx.index++ {
		ctx.handlers[ctx.index](ctx)
	}
}
