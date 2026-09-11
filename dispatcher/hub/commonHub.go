package hub

import (
	"sync"

	"github.com/google/uuid"
	"github.com/multiversx/mx-chain-core-go/core/check"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/data"
	"github.com/multiversx/mx-chain-notifier-go/dispatcher"
	"github.com/multiversx/mx-chain-notifier-go/filters"
)

var log = logger.GetOrCreate("hub")

// ArgsCommonHub defines the arguments needed for common hub creation
type ArgsCommonHub struct {
	Filter             filters.EventFilter
	SubscriptionMapper dispatcher.SubscriptionMapperHandler
}

type commonHub struct {
	filter             filters.EventFilter
	subscriptionMapper dispatcher.SubscriptionMapperHandler
	mutDispatchers     sync.RWMutex
	dispatchers        map[uuid.UUID]dispatcher.EventDispatcher
}

// NewCommonHub creates a new commonHub instance
func NewCommonHub(args ArgsCommonHub) (*commonHub, error) {
	err := checkArgs(args)
	if err != nil {
		return nil, err
	}

	return &commonHub{
		mutDispatchers:     sync.RWMutex{},
		filter:             args.Filter,
		subscriptionMapper: args.SubscriptionMapper,
		dispatchers:        make(map[uuid.UUID]dispatcher.EventDispatcher),
	}, nil
}

func checkArgs(args ArgsCommonHub) error {
	if check.IfNil(args.Filter) {
		return ErrNilEventFilter
	}
	if check.IfNil(args.SubscriptionMapper) {
		return ErrNilSubscriptionMapper
	}

	return nil
}

// Subscribe is used by a dispatcher to send a dispatcher.SubscribeEvent
func (ch *commonHub) Subscribe(event data.SubscribeEvent) {
	ch.subscriptionMapper.MatchSubscribeEvent(event)
}

// RegisterEvent will send event to a receive-only channel used to register dispatchers
func (ch *commonHub) RegisterEvent(event dispatcher.EventDispatcher) {
	ch.registerDispatcher(event)
}

// UnregisterEvent will send event to a receive-only channel used by a dispatcher to signal it has disconnected
func (ch *commonHub) UnregisterEvent(event dispatcher.EventDispatcher) {
	ch.unregisterDispatcher(event)
}

// Publish will publish logs and events to dispatcher
func (ch *commonHub) Publish(blockEvents data.BlockEvents) {
	subscriptions := ch.subscriptionMapper.Subscriptions()[common.PushLogsAndEvents]
	if len(subscriptions) == 0 {
		return
	}

	type target struct {
		d      dispatcher.EventDispatcher
		events []data.Event
	}

	ch.mutDispatchers.RLock()
	targets := make([]target, 0, len(subscriptions))
	for _, sub := range subscriptions {
		d, ok := ch.dispatchers[sub.DispatcherID]
		if !ok {
			continue
		}

		events := make([]data.Event, 0)
		for _, event := range blockEvents.Events {
			if ch.filter.MatchEvent(sub, event) {
				events = append(events, event)
			}
		}
		targets = append(targets, target{d: d, events: events})
	}
	ch.mutDispatchers.RUnlock()

	// delivery happens outside the lock so a slow or stuck subscriber can
	// never block registration/unregistration of other dispatchers
	for _, t := range targets {
		t.d.PushEvents(t.events)
	}
}

// targetDispatchers returns the (deduplicated) live dispatchers currently
// subscribed to eventType. The hub lock is held only while reading the
// dispatchers map, never while delivering events to them.
func (ch *commonHub) targetDispatchers(eventType string) []dispatcher.EventDispatcher {
	subs := ch.subscriptionMapper.Subscriptions()[eventType]
	if len(subs) == 0 {
		return nil
	}

	ch.mutDispatchers.RLock()
	defer ch.mutDispatchers.RUnlock()

	seen := make(map[uuid.UUID]struct{}, len(subs))
	targets := make([]dispatcher.EventDispatcher, 0, len(subs))
	for _, sub := range subs {
		if _, duplicate := seen[sub.DispatcherID]; duplicate {
			continue
		}
		seen[sub.DispatcherID] = struct{}{}

		if d, ok := ch.dispatchers[sub.DispatcherID]; ok {
			targets = append(targets, d)
		}
	}

	return targets
}

// PublishRevert will publish revert event to dispatcher
func (ch *commonHub) PublishRevert(revertBlock data.RevertBlock) {
	for _, d := range ch.targetDispatchers(common.RevertBlockEvents) {
		d.RevertEvent(revertBlock)
	}
}

// PublishFinalized will publish finalized event to dispatcher
func (ch *commonHub) PublishFinalized(finalizedBlock data.FinalizedBlock) {
	for _, d := range ch.targetDispatchers(common.FinalizedBlockEvents) {
		d.FinalizedEvent(finalizedBlock)
	}
}

// PublishTxs will publish txs event to dispatcher
func (ch *commonHub) PublishTxs(blockTxs data.BlockTxs) {
	for _, d := range ch.targetDispatchers(common.BlockTxs) {
		d.TxsEvent(blockTxs)
	}
}

// PublishBlockEventsWithOrder will publish block events with order to dispatcher
func (ch *commonHub) PublishBlockEventsWithOrder(blockTxs data.BlockEventsWithOrder) {
	for _, d := range ch.targetDispatchers(common.BlockEvents) {
		d.BlockEvents(blockTxs)
	}
}

// PublishScrs will publish scrs events to dispatcher
func (ch *commonHub) PublishScrs(blockScrs data.BlockScrs) {
	for _, d := range ch.targetDispatchers(common.BlockScrs) {
		d.ScrsEvent(blockScrs)
	}
}

// PublishStateAccesses will publish state accesses to dispatcher
func (ch *commonHub) PublishStateAccesses(stateAccesses data.BlockStateAccesses) {
	for _, d := range ch.targetDispatchers(common.BlockStateAccesses) {
		d.StateAccessesEvent(stateAccesses)
	}
}

func (ch *commonHub) registerDispatcher(d dispatcher.EventDispatcher) {
	ch.mutDispatchers.Lock()
	defer ch.mutDispatchers.Unlock()

	if _, ok := ch.dispatchers[d.GetID()]; ok {
		return
	}

	ch.dispatchers[d.GetID()] = d

	log.Info("registered new dispatcher", "dispatcherID", d.GetID())
}

func (ch *commonHub) unregisterDispatcher(d dispatcher.EventDispatcher) {
	ch.mutDispatchers.Lock()
	defer ch.mutDispatchers.Unlock()

	if _, ok := ch.dispatchers[d.GetID()]; ok {
		delete(ch.dispatchers, d.GetID())
	}

	log.Info("unregistered dispatcher", "dispatcherID", d.GetID())

	ch.subscriptionMapper.RemoveSubscriptions(d.GetID())
}

// Close will close the goroutine and channels
func (ch *commonHub) Close() error {
	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (ch *commonHub) IsInterfaceNil() bool {
	return ch == nil
}
