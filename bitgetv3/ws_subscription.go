package bitgetv3

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

// subKey - normalised topic identity used as the registry key.
// The server echoes arg exactly as it was sent, so a raw struct key would work;
// instType and symbol are still folded as a cheap safety net against a future
// change of the echoed case (REST uses uppercase categories, the WS lowercase).
// The topic/channel name is deliberately NOT folded: its case is significant -
// the v2 candle channels candle1m (minute) and candle1M (month) differ only by it,
// and folding them would route both streams into whichever handler the map yields
// first. Insert, lookup and delete all go through this key, so they stay symmetric.
type subKey struct {
	instType string
	topic    string
	symbol   string
}

func newSubKey(args SubscriptionArgs) subKey {
	return subKey{
		instType: strings.ToLower(args.InstType),
		topic:    args.Topic,
		symbol:   strings.ToUpper(args.Symbol),
	}
}

// subscription - a registered handler together with the args as they go on the
// wire: the key is normalised, the request must keep the original case
type subscription struct {
	args SubscriptionArgs
	f    SubscriptionFunc
}

// Subscriptions - thread-safe subscription registry and topic router
type Subscriptions struct {
	c     SubscriptionClient
	mutex sync.Mutex
	funcs SubscriptionFuncs
	// unsubscribing - topics whose unsubscribe request was sent but not yet acked.
	// The server goes on pushing until it processes the request (231ms observed
	// live), so those in-flight pushes are expected and must not be reported as
	// an unknown topic
	unsubscribing map[subKey]struct{}
}

func NewSubscriptions(c SubscriptionClient) *Subscriptions {
	o := new(Subscriptions)
	o.c = c
	o.funcs = make(SubscriptionFuncs)
	o.unsubscribing = make(map[subKey]struct{})
	return o
}

func (o *Subscriptions) subscribe(args SubscriptionArgs, f SubscriptionFunc) {
	if o.c.Ready() {
		o.c.subscribe(args)
	}
	o.mutex.Lock()
	defer o.mutex.Unlock()
	k := newSubKey(args)
	o.funcs[k] = subscription{args: args, f: f}
	// resubscribed before the pending unsubscribe was acked: the topic is live again
	delete(o.unsubscribing, k)
}

func (o *Subscriptions) unsubscribe(args SubscriptionArgs) {
	ready := o.c.Ready()
	if ready {
		o.c.unsubscribe(args)
	}
	o.mutex.Lock()
	defer o.mutex.Unlock()
	k := newSubKey(args)
	delete(o.funcs, k)
	if ready {
		// nothing was sent when not ready, so there is no ack to wait for
		o.unsubscribing[k] = struct{}{}
	}
}

// unsubscribeConfirmed - the server acked the unsubscribe and stopped pushing:
// from here on a push for this topic is a genuine routing error again
func (o *Subscriptions) unsubscribeConfirmed(args SubscriptionArgs) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	delete(o.unsubscribing, newSubKey(args))
}

// resetUnsubscribing - the socket is gone: no ack will ever arrive for a pending
// unsubscribe, and the server-side subscription state died with the connection
func (o *Subscriptions) resetUnsubscribing() {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	clear(o.unsubscribing)
}

// subscribeAll - resubscribe to all registered channels in a single batch request
func (o *Subscriptions) subscribeAll() {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	if len(o.funcs) == 0 {
		return
	}
	args := make([]SubscriptionArgs, 0, len(o.funcs))
	for _, s := range o.funcs {
		args = append(args, s.args)
	}
	o.c.subscribe(args...)
}

func (o *Subscriptions) processTopic(data []byte) (err error) {
	var topic RawTopic
	err = json.Unmarshal(data, &topic)
	if err == nil {
		f, unsubscribing := o.route(topic.Arg)
		switch {
		case f != nil:
			err = f(topic)
		case unsubscribing:
			// in-flight push sent before the server processed our unsubscribe: expected, drop it
		default:
			err = fmt.Errorf("subscription of topic[%s %s %s] not found",
				topic.Arg.InstType, topic.Arg.Topic, topic.Arg.Symbol)
		}
	}
	return
}

// route - resolve the handler of an echoed arg; reports unsubscribing when the
// topic has no handler only because its unsubscribe is still awaiting the ack
func (o *Subscriptions) route(args SubscriptionArgs) (f SubscriptionFunc, unsubscribing bool) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	k := newSubKey(args)
	if s, ok := o.funcs[k]; ok {
		f = s.f
		return
	}
	_, unsubscribing = o.unsubscribing[k]
	return
}

type SubscriptionClient interface {
	Ready() bool
	subscribe(...SubscriptionArgs)
	unsubscribe(...SubscriptionArgs)
}

type SubscriptionFunc func(RawTopic) error

type SubscriptionFuncs map[subKey]subscription
