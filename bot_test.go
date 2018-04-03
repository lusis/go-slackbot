package slackbot

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	slacktest "github.com/lusis/slack-test"
	"github.com/nlopes/slack"
	"github.com/stretchr/testify/require"
)

func testNewBot(t *testing.T) *slacktest.Server {
	s, err := slacktest.NewTestServer()
	require.NoError(t, err)
	return s
}

func testChannelJoinFunc(ctx context.Context, bot *Bot, channel *slack.Channel) {
	p := slack.PostMessageParameters{
		AsUser: true,
	}
	bot.logger.Printf("joined channel")
	_, _, _ = bot.Client.PostMessage(channel.ID, "Thanks for the invite!", p)
}

func testPingFunc(ctx context.Context, bot *Bot, evt *slack.MessageEvent) {
	bot.Reply(evt, "pong", WithoutTyping)
}

func testWithTypingPingFunc(ctx context.Context, bot *Bot, evt *slack.MessageEvent) {
	bot.Reply(evt, "pong", WithTyping)
}

func testPingAttachmentFunc(ctx context.Context, bot *Bot, evt *slack.MessageEvent) {
	fallback := "pong"
	attachment := slack.Attachment{
		Color:    "#000066",
		Fallback: fallback,
		Title:    "WarBot Info",
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Message",
				Value: "pong",
				Short: true,
			},
		},
	}
	bot.ReplyWithAttachments(evt, []slack.Attachment{attachment}, WithoutTyping)
}

func TestOnChannelJoin(t *testing.T) {
	s := testNewBot(t)
	slack.SLACK_API = s.GetAPIURL()
	go s.Start()
	bot, err := NewWithOpts(WithClient(slack.New("ABCDEFG")))
	require.NoError(t, err)
	bot.OnChannelJoin(testChannelJoinFunc)
	go bot.Run()
	defer func() {
		s.Stop()
		bot.Stop()
	}()

	maxWait := 5 * time.Second
	s.SendBotChannelInvite()
	time.Sleep(maxWait)
	messages := s.GetSeenOutboundMessages()
	require.Len(t, messages, 2, "should see 2 messages")
	var m = &slack.Message{}
	// second message in slice should be our notification
	_ = json.Unmarshal([]byte(messages[1]), m)
	require.Equal(t, "Thanks for the invite!", m.Text)
}

func TestOnGroupJoin(t *testing.T) {
	s := testNewBot(t)
	slack.SLACK_API = s.GetAPIURL()
	go s.Start()
	bot, err := NewWithOpts(WithClient(slack.New("ABCDEFG")))
	require.NoError(t, err)
	bot.OnChannelJoin(testChannelJoinFunc)
	go bot.Run()
	defer func() {
		s.Stop()
		bot.Stop()
	}()

	maxWait := 5 * time.Second
	s.SendBotGroupInvite()
	time.Sleep(maxWait)
	messages := s.GetSeenOutboundMessages()
	require.Len(t, messages, 2, "should see 2 messages")
	var m = &slack.Message{}
	// second message in slice should be our notification
	_ = json.Unmarshal([]byte(messages[1]), m)
	require.Equal(t, "Thanks for the invite!", m.Text)
}

func TestReply(t *testing.T) {
	s := testNewBot(t)
	slack.SLACK_API = s.GetAPIURL()
	go s.Start()
	bot, err := NewWithOpts(WithClient(slack.New("ABCDEFG")))
	require.NoError(t, err)
	tome := bot.Messages(DirectMention).Subrouter()
	tome.Hear(`^ping$`).MessageHandler(testPingFunc)
	go bot.Run()
	defer func() {
		s.Stop()
		bot.Stop()
	}()

	maxWait := 5 * time.Second
	s.SendMessageToBot("#test", "ping")
	time.Sleep(maxWait)
	require.True(t, s.SawMessage("pong"))
}

func TestReplyWithTyping(t *testing.T) {
	s := testNewBot(t)
	slack.SLACK_API = s.GetAPIURL()
	go s.Start()
	bot, err := NewWithOpts(WithClient(slack.New("ABCDEFG")))
	require.NoError(t, err)
	tome := bot.Messages(DirectMention).Subrouter()
	tome.Hear(`^ping$`).MessageHandler(testWithTypingPingFunc)
	go bot.Run()
	defer func() {
		s.Stop()
		bot.Stop()
	}()

	maxWait := 5 * time.Second
	s.SendMessageToBot("#test", "ping")
	time.Sleep(maxWait)
	require.True(t, s.SawMessage("pong"))
}

func TestReplyWithAttachment(t *testing.T) {
	s := testNewBot(t)
	slack.SLACK_API = s.GetAPIURL()
	go s.Start()
	bot, err := NewWithOpts(WithClient(slack.New("ABCDEFG")))
	require.NoError(t, err)
	tome := bot.Messages(DirectMention).Subrouter()
	tome.Hear(`^ping$`).MessageHandler(testPingAttachmentFunc)
	go bot.Run()
	defer func() {
		s.Stop()
		bot.Stop()
	}()

	maxWait := 5 * time.Second
	s.SendMessageToBot("#test", "ping")
	time.Sleep(maxWait)
	msgs := s.GetSeenOutboundMessages()
	slackMsg := &slack.Message{}
	_ = json.Unmarshal([]byte(msgs[1]), slackMsg)
	require.Contains(t, slackMsg.Attachments[0].Fallback, "pong")
}
