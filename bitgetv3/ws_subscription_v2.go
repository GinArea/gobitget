package bitgetv3

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

// SubscriptionsV2 - thread-safe subscription registry and topic router of the legacy v2 WS
type SubscriptionsV2 struct {
	c     SubscriptionClientV2
	mutex sync.Mutex
	funcs SubscriptionFuncsV2
}

func NewSubscriptionsV2(c SubscriptionClientV2) *SubscriptionsV2 {
	o := new(SubscriptionsV2)
	o.c = c
	o.funcs = make(SubscriptionFuncsV2)
	return o
}

func (o *SubscriptionsV2) subscribe(args SubscriptionArgsV2, f SubscriptionFuncV2) {
	if o.c.Ready() {
		o.c.subscribe(args)
	}
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.funcs[args] = f
}

func (o *SubscriptionsV2) unsubscribe(args SubscriptionArgsV2) {
	if o.c.Ready() {
		o.c.unsubscribe(args)
	}
	o.mutex.Lock()
	defer o.mutex.Unlock()
	delete(o.funcs, args)
}

// subscribeAll - resubscribe to all registered channels in a single batch request
func (o *SubscriptionsV2) subscribeAll() {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	if len(o.funcs) == 0 {
		return
	}
	args := make([]SubscriptionArgsV2, 0, len(o.funcs))
	for a := range o.funcs {
		args = append(args, a)
	}
	o.c.subscribe(args...)
}

func (o *SubscriptionsV2) processTopic(data []byte) (err error) {
	var topic RawTopicV2
	err = json.Unmarshal(data, &topic)
	if err == nil {
		f := o.getFunc(topic.Arg)
		if f == nil {
			err = fmt.Errorf("subscription of v2 topic[%s %s %s] not found",
				topic.Arg.InstType, topic.Arg.Channel, topic.Arg.InstId)
		} else {
			err = f(topic)
		}
	}
	return
}

// getFunc - match the echoed arg field-wise, case-insensitive
func (o *SubscriptionsV2) getFunc(passed SubscriptionArgsV2) (f SubscriptionFuncV2) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	for args, fn := range o.funcs {
		if strings.EqualFold(args.InstType, passed.InstType) &&
			strings.EqualFold(args.Channel, passed.Channel) &&
			strings.EqualFold(args.InstId, passed.InstId) {
			f = fn
			break
		}
	}
	return
}

type SubscriptionClientV2 interface {
	Ready() bool
	subscribe(...SubscriptionArgsV2)
	unsubscribe(...SubscriptionArgsV2)
}

type SubscriptionFuncV2 func(RawTopicV2) error

type SubscriptionFuncsV2 map[SubscriptionArgsV2]SubscriptionFuncV2
