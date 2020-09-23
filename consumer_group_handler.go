package titan

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/ramadani/saturn"
)

type consumerGroupHandler struct {
	evLst eventListenersMap
	ready chan struct{}
}

func (h *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	close(h.ready)
	return nil
}

func (h *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if listeners, ok := h.evLst[msg.Topic]; ok {
			for _, listener := range listeners {
				go func(ctx context.Context, listener saturn.Listener, val []byte) {
					_ = listener.Handle(ctx, val)
				}(session.Context(), listener, msg.Value)
			}

			session.MarkMessage(msg, "")
		}
	}

	return nil
}
