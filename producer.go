package titan

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/ramadani/saturn"
)

type producer struct {
	prd sarama.SyncProducer
}

func (e *producer) Emit(_ context.Context, dispatchable saturn.Dispatchable) (err error) {
	topic := dispatchable.Header()
	body, err := dispatchable.Body()
	if err != nil {
		return
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(body),
	}

	_, _, err = e.prd.SendMessage(msg)
	return
}

func NewEmitter(prd sarama.SyncProducer) saturn.Emitter {
	return &producer{prd: prd}
}
