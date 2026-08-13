package bitgetv3

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

// Subscriptions - thread-safe subscription registry and topic router
type Subscriptions struct {
	c     SubscriptionClient
	mutex sync.Mutex
	funcs SubscriptionFuncs
}

func NewSubscriptions(c SubscriptionClient) *Subscriptions {
	o := new(Subscriptions)
	o.c = c
	o.funcs = make(SubscriptionFuncs)
	return o
}

func (o *Subscriptions) subscribe(args SubscriptionArgs, f SubscriptionFunc) {
	if o.c.Ready() {
		o.c.subscribe(args)
	}
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.funcs[args] = f
}

func (o *Subscriptions) unsubscribe(args SubscriptionArgs) {
	if o.c.Ready() {
		o.c.unsubscribe(args)
	}
	o.mutex.Lock()
	defer o.mutex.Unlock()
	delete(o.funcs, args)
}

// subscribeAll - resubscribe to all registered channels in a single batch request
func (o *Subscriptions) subscribeAll() {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	if len(o.funcs) == 0 {
		return
	}
	args := make([]SubscriptionArgs, 0, len(o.funcs))
	for a := range o.funcs {
		args = append(args, a)
	}
	o.c.subscribe(args...)
}

func (o *Subscriptions) processTopic(data []byte) (err error) {
	var topic RawTopic
	err = json.Unmarshal(data, &topic)
	if err == nil {
		f := o.getFunc(topic.Arg)
		if f == nil {
			err = fmt.Errorf("subscription of topic[%s %s %s %s] not found",
				topic.Arg.InstType, topic.Arg.Topic, topic.Arg.Symbol, topic.Arg.Interval)
		} else {
			err = f(topic)
		}
	}
	return
}

// getFunc - match the echoed arg field-wise, case-insensitive
func (o *Subscriptions) getFunc(passed SubscriptionArgs) (f SubscriptionFunc) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	for args, fn := range o.funcs {
		if strings.EqualFold(args.InstType, passed.InstType) &&
			strings.EqualFold(args.Topic, passed.Topic) &&
			strings.EqualFold(args.Symbol, passed.Symbol) &&
			strings.EqualFold(args.Interval, passed.Interval) {
			f = fn
			break
		}
	}
	return
}

type SubscriptionClient interface {
	Ready() bool
	subscribe(...SubscriptionArgs)
	unsubscribe(...SubscriptionArgs)
}

type SubscriptionFunc func(RawTopic) error

type SubscriptionFuncs map[SubscriptionArgs]SubscriptionFunc
