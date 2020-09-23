package titan

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/ramadani/saturn"
)

type eventListenersMap map[string][]saturn.Listener

type consumerGroup struct {
	consumer sarama.ConsumerGroup
	evLst    eventListenersMap
	handler  *consumerGroupHandler
	ready    chan struct{}
}

func (e *consumerGroup) On(_ context.Context, header string, listeners []saturn.Listener) (err error) {
	if lts, ok := e.evLst[header]; ok {
		e.evLst[header] = append(lts, listeners...)
	} else {
		e.evLst[header] = listeners
	}
	return
}

func (e *consumerGroup) Listeners(_ context.Context, header string) (res []saturn.Listener, err error) {
	res = make([]saturn.Listener, 0)

	if lts, ok := e.evLst[header]; ok {
		res = lts
	}
	return
}

func (e *consumerGroup) Listen(ctx context.Context) (err error) {
	topics := make([]string, 0)
	for key := range e.evLst {
		topics = append(topics, key)
	}

	for {
		if err = e.consumer.Consume(ctx, topics, e.handler); err != nil {
			return
		}

		if er := ctx.Err(); er != nil {
			return
		}
	}
}

func (e *consumerGroup) Ready() <-chan struct{} {
	return e.handler.ready
}

func NewConsumerGroupEventListener(consumer sarama.ConsumerGroup, evLst eventListenersMap) EventListener {
	handler := &consumerGroupHandler{
		evLst: evLst,
		ready: make(chan struct{}),
	}

	return &consumerGroup{
		consumer: consumer,
		evLst:    evLst,
		handler:  handler,
	}
}

func DefaultConsumerGroupEventListener(consumer sarama.ConsumerGroup) EventListener {
	return NewConsumerGroupEventListener(consumer, make(eventListenersMap))
}
