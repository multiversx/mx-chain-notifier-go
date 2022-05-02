package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/notifier-go/common"
)

const minTTLInMinutes = 0

type ArgsRedlockWrapper struct {
	Client       RedLockClient
	TTLInMinutes uint32
}

type redlockWrapper struct {
	client RedLockClient
	ttl    time.Duration
}

// NewRedlockWrapper create a new redLock based on a cache instance
func NewRedlockWrapper(args ArgsRedlockWrapper) (*redlockWrapper, error) {
	if check.IfNil(args.Client) {
		return nil, ErrNilRedlockClient
	}
	if args.TTLInMinutes <= minTTLInMinutes {
		return nil, fmt.Errorf("%w for TTL in minutes, got: %d, minimum: %d",
			common.ErrInvalidValue, args.TTLInMinutes, minTTLInMinutes)
	}

	ttl := time.Minute * time.Duration(args.TTLInMinutes)

	return &redlockWrapper{
		client: args.Client,
		ttl:    ttl,
	}, nil
}

// IsEventProcessed returns wether the item is already locked
func (r *redlockWrapper) IsEventProcessed(ctx context.Context, blockHash string) (bool, error) {
	return r.client.SetEntry(ctx, blockHash, true, r.ttl)
}

// HasConnection return true if the redis client is connected
func (r *redlockWrapper) HasConnection(ctx context.Context) bool {
	return r.client.IsConnected(ctx)
}

// IsInterfaceNil returns true if there is no value under the interface
func (r *redlockWrapper) IsInterfaceNil() bool {
	return r == nil
}
