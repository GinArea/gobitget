package bitgetv3

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

// newSubKeyV2 - normalised registry key of a legacy v2 channel; see subKey.
// The channel name carries the timeframe and its case is significant
// (candle1m is a minute, candle1M a month), so it is matched verbatim
func newSubKeyV2(args SubscriptionArgsV2) subKey {
	return subKey{
		instType: strings.ToLower(args.InstType),
		topic:    args.Channel,
		symbol:   strings.ToUpper(args.InstId),
	}
}

// subscriptionV2 - a registered handler together with the args as they go on the wire
type subscriptionV2 struct {
	args SubscriptionArgsV2
	f    SubscriptionFuncV2
}

// SubscriptionsV2 - thread-safe subscription registry and topic router of the legacy v2 WS
type SubscriptionsV2 struct {
	c     SubscriptionClientV2
	mutex sync.Mutex
	funcs SubscriptionFuncsV2
	// unsubscribing - channels whose unsubscribe request was sent but not yet
	// acked; the server goes on pushing until it processes the request, and those
	// in-flight pushes are expected rather than an unknown topic
	unsubscribing map[subKey]struct{}
}

func NewSubscriptionsV2(c SubscriptionClientV2) *SubscriptionsV2 {
	o := new(SubscriptionsV2)
	o.c = c
	o.funcs = make(SubscriptionFuncsV2)
	o.unsubscribing = make(map[subKey]struct{})
	return o
}

func (o *SubscriptionsV2) subscribe(args SubscriptionArgsV2, f SubscriptionFuncV2) {
	if o.c.Ready() {
		o.c.subscribe(args)
	}
	o.mutex.Lock()
	defer o.mutex.Unlock()
	k := newSubKeyV2(args)
	o.funcs[k] = subscriptionV2{args: args, f: f}
	// resubscribed before the pending unsubscribe was acked: the channel is live again
	delete(o.unsubscribing, k)
}

func (o *SubscriptionsV2) unsubscribe(args SubscriptionArgsV2) {
	ready := o.c.Ready()
	if ready {
		o.c.unsubscribe(args)
	}
	o.mutex.Lock()
	defer o.mutex.Unlock()
	k := newSubKeyV2(args)
	delete(o.funcs, k)
	if ready {
		// nothing was sent when not ready, so there is no ack to wait for
		o.unsubscribing[k] = struct{}{}
	}
}

// unsubscribeConfirmed - the server acked the unsubscribe and stopped pushing:
// from here on a push for this channel is a genuine routing error again
func (o *SubscriptionsV2) unsubscribeConfirmed(args SubscriptionArgsV2) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	delete(o.unsubscribing, newSubKeyV2(args))
}

// resetUnsubscribing - the socket is gone: no ack will ever arrive for a pending
// unsubscribe, and the server-side subscription state died with the connection
func (o *SubscriptionsV2) resetUnsubscribing() {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	clear(o.unsubscribing)
}

// subscribeAll - resubscribe to all registered channels in a single batch request
func (o *SubscriptionsV2) subscribeAll() {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	if len(o.funcs) == 0 {
		return
	}
	args := make([]SubscriptionArgsV2, 0, len(o.funcs))
	for _, s := range o.funcs {
		args = append(args, s.args)
	}
	o.c.subscribe(args...)
}

func (o *SubscriptionsV2) processTopic(data []byte) (err error) {
	var topic RawTopicV2
	err = json.Unmarshal(data, &topic)
	if err == nil {
		f, unsubscribing := o.route(topic.Arg)
		switch {
		case f != nil:
			err = f(topic)
		case unsubscribing:
			// in-flight push sent before the server processed our unsubscribe: expected, drop it
		default:
			err = fmt.Errorf("subscription of v2 topic[%s %s %s] not found",
				topic.Arg.InstType, topic.Arg.Channel, topic.Arg.InstId)
		}
	}
	return
}

// route - resolve the handler of an echoed arg; reports unsubscribing when the
// channel has no handler only because its unsubscribe is still awaiting the ack
func (o *SubscriptionsV2) route(args SubscriptionArgsV2) (f SubscriptionFuncV2, unsubscribing bool) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	k := newSubKeyV2(args)
	if s, ok := o.funcs[k]; ok {
		f = s.f
		return
	}
	_, unsubscribing = o.unsubscribing[k]
	return
}

type SubscriptionClientV2 interface {
	Ready() bool
	subscribe(...SubscriptionArgsV2)
	unsubscribe(...SubscriptionArgsV2)
}

type SubscriptionFuncV2 func(RawTopicV2) error

type SubscriptionFuncsV2 map[subKey]subscriptionV2
