package slackbot

import (
	"context"

	"github.com/nlopes/slack"
)

// MessageType represents a message type
type MessageType string

const (
	// DirectMessage represents a message type
	DirectMessage MessageType = "direct_message"
	// DirectMention represents a direct message
	DirectMention MessageType = "direct_mention"
	// Mention is a mention
	Mention MessageType = "mention"
	// Ambient is ambient
	Ambient MessageType = "ambient"
)

// Handler is a handler
type Handler func(context.Context)

// MessageHandler is a message handler
type MessageHandler func(ctx context.Context, bot *Bot, msg *slack.MessageEvent)

// Preprocessor is a preprocessor
type Preprocessor func(context.Context) context.Context

// Matcher type for matching message routes
type Matcher interface {
	Match(context.Context) (bool, context.Context)
	SetBotID(botID string)
}
