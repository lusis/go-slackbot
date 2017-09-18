package slackbot

import (
	"regexp"

	"context"
)

// Route represents a route
type Route struct {
	handler      Handler
	err          error
	matchers     []Matcher
	subrouter    Router
	preprocessor Preprocessor
	botUserID    string
}

func (r *Route) setBotID(botID string) {
	r.botUserID = botID
	for _, matcher := range r.matchers {
		matcher.SetBotID(botID)
	}
}

// RouteMatch stores information about a matched route.
type RouteMatch struct {
	Route   *Route
	Handler Handler
}

// Match matches
func (r *Route) Match(ctx context.Context, match *RouteMatch) (bool, context.Context) {
	if r.preprocessor != nil {
		ctx = r.preprocessor(ctx)
	}
	for _, m := range r.matchers {
		var matched bool
		matched, ctx = m.Match(ctx)
		if !matched {
			return false, ctx
		}
	}

	// if this route contains a subrouter, invoke the subrouter match
	if r.subrouter != nil {
		return r.subrouter.Match(ctx, match)
	}

	match.Route = r
	match.Handler = r.handler
	return true, ctx
}

// Hear adds a matcher for the message text
func (r *Route) Hear(regex string) *Route {
	r.err = r.addRegexpMatcher(regex)
	return r
}

func (r *Route) Messages(types ...MessageType) *Route {
	r.addTypesMatcher(types...)
	return r
}

// Handler sets a handler for the route.
func (r *Route) Handler(handler Handler) *Route {
	if r.err == nil {
		r.handler = handler
	}
	return r
}

// MessageHandler is a message handler
func (r *Route) MessageHandler(fn MessageHandler) *Route {
	return r.Handler(func(ctx context.Context) {
		bot := BotFromContext(ctx)
		msg := MessageFromContext(ctx)
		fn(ctx, bot, msg)
	})
}

// Preprocess preproccesses
func (r *Route) Preprocess(fn Preprocessor) *Route {
	if r.err == nil {
		r.preprocessor = fn
	}
	return r
}

// SubRouter creates a subrouter
func (r *Route) Subrouter() Router {
	if r.err == nil {
		r.subrouter = &SimpleRouter{}
	}
	return r.subrouter
}

// AddMatcher adds a matcher to the route.
func (r *Route) AddMatcher(m Matcher) *Route {
	if r.err == nil {
		r.matchers = append(r.matchers, m)
	}
	return r
}

// ============================================================================
// Regex Type Matcher
// ============================================================================
// RegexpMatcher is a regexp matcher
type RegexpMatcher struct {
	regex     string
	botUserID string
}

// Match matches
func (rm *RegexpMatcher) Match(ctx context.Context) (bool, context.Context) {
	msg := MessageFromContext(ctx)
	// A message be receded by a direct mention. For simplicity sake, strip out any potention direct mentions first
	text := StripDirectMention(msg.Text)
	// now consider stripped text against regular expression
	matched := regexp.MustCompile(rm.regex).MatchString(text)
	return matched, ctx
}

// SetBotID sets the bot id
func (rm *RegexpMatcher) SetBotID(botID string) {
	rm.botUserID = botID
}

// addRegexpMatcher adds a host or path matcher and builder to a route.
func (r *Route) addRegexpMatcher(regex string) error {
	if r.err != nil {
		return r.err
	}

	r.AddMatcher(&RegexpMatcher{regex: regex})
	return nil
}

// ============================================================================
// Message Type Matcher
// ============================================================================
// TypesMatcher is a type matcher
type TypesMatcher struct {
	types     []MessageType
	botUserID string
}

// Match matches
func (tm *TypesMatcher) Match(ctx context.Context) (bool, context.Context) {
	msg := MessageFromContext(ctx)
	for _, t := range tm.types {
		switch t {
		case DirectMessage:
			if IsDirectMessage(msg) {
				return true, ctx
			}
		case DirectMention:
			if IsDirectMention(msg, tm.botUserID) {
				return true, ctx
			}
		}
	}
	return false, ctx
}

// SetBotID sets the botid
func (tm *TypesMatcher) SetBotID(botID string) {
	tm.botUserID = botID
}

// addTypesMatcher adds a host or path matcher and builder to a route.
func (r *Route) addTypesMatcher(types ...MessageType) error {
	if r.err != nil {
		return r.err
	}

	r.AddMatcher(&TypesMatcher{types: types, botUserID: ""})
	return nil
}
