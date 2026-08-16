package bitgetv3

// ExecutorV2 - stateless typed subscription handle for a single v2 channel
type ExecutorV2[T any] struct {
	args          SubscriptionArgsV2
	subscriptions *SubscriptionsV2
}

func NewExecutorV2[T any](args SubscriptionArgsV2, subscriptions *SubscriptionsV2) *ExecutorV2[T] {
	o := new(ExecutorV2[T])
	o.args = args
	o.subscriptions = subscriptions
	return o
}

func (o *ExecutorV2[T]) Subscribe(onShot func(TopicV2[T])) {
	o.subscriptions.subscribe(o.args, func(raw RawTopicV2) error {
		topic, err := UnmarshalRawTopicV2[T](raw)
		if err == nil {
			onShot(topic)
		}
		return err
	})
}

func (o *ExecutorV2[T]) Unsubscribe() {
	o.subscriptions.unsubscribe(o.args)
}
