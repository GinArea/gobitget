package bitgetv3

// Executor - stateless typed subscription handle for a single channel
type Executor[T any] struct {
	args          SubscriptionArgs
	subscriptions *Subscriptions
}

func NewExecutor[T any](args SubscriptionArgs, subscriptions *Subscriptions) *Executor[T] {
	o := new(Executor[T])
	o.args = args
	o.subscriptions = subscriptions
	return o
}

func (o *Executor[T]) Subscribe(onShot func(Topic[T])) {
	o.subscriptions.subscribe(o.args, func(raw RawTopic) error {
		topic, err := UnmarshalRawTopic[T](raw)
		if err == nil {
			onShot(topic)
		}
		return err
	})
}

func (o *Executor[T]) Unsubscribe() {
	o.subscriptions.unsubscribe(o.args)
}
