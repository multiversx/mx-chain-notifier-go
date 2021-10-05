package rabbitmq

import (
	"context"
	"testing"
)

var ctx = context.Background()

func TestNewRabbitMqPublisher(t *testing.T) {
	t.Parallel()

	r := NewRabbitMqPublisher(ctx)
	go r.Run()

	c := make(chan int)
	<-c
}
