package slackbot

import (
	"context"

	"github.com/nlopes/slack"
)

const (
	BOT_CONTEXT     = "__BOT_CONTEXT__"
	MESSAGE_CONTEXT = "__MESSAGE_CONTEXT__"
	// NamedCaptureContextKey is the key for named captures
	NamedCaptureContextKey = "__NAMED_CAPTURES__"
)

// BotFromContext creates a Bot from provided Context
func BotFromContext(ctx context.Context) *Bot {
	if result, ok := ctx.Value(BOT_CONTEXT).(*Bot); ok {
		return result
	}
	return nil
}

// AddBotToContext sets the bot reference in context and returns the newly derived context
func AddBotToContext(ctx context.Context, bot *Bot) context.Context {
	return context.WithValue(ctx, BOT_CONTEXT, bot)
}

func MessageFromContext(ctx context.Context) *slack.MessageEvent {
	if result, ok := ctx.Value(MESSAGE_CONTEXT).(*slack.MessageEvent); ok {
		return result
	}
	return nil
}

// AddMessageToContext sets the Slack message event reference in context and returns the newly derived context
func AddMessageToContext(ctx context.Context, msg *slack.MessageEvent) context.Context {
	return context.WithValue(ctx, MESSAGE_CONTEXT, msg)
}

// NamedCapturesFromContext returns any NamedCaptures parsed from regexp
func NamedCapturesFromContext(ctx context.Context) NamedCaptures {
	return ctx.Value(NamedCaptureContextKey).(NamedCaptures)
}
