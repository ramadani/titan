package titan

import "github.com/ramadani/saturn"

type EventListener interface {
	saturn.EventListener
	Ready() <-chan struct{}
}
