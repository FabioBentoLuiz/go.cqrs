package example

import (
	"fmt"
	"reflect"

	"github.com/EventStore/EventStore-Client-Go/client"
	"github.com/fabiobentoluiz/eventsourcing"
)

// ###### order

// InMemoryOrderRepo provides an in memory repository implementation.
type InMemoryOrderRepo struct {
	current   map[string][]eventsourcing.EventMessage
	publisher eventsourcing.EventBus
}

// NewInMemoryRepo constructs an InMemoryRepo instance.
func NewInMemoryOrderRepo(eventBus eventsourcing.EventBus) *InMemoryOrderRepo {
	return &InMemoryOrderRepo{
		current:   make(map[string][]eventsourcing.EventMessage),
		publisher: eventBus,
	}
}

// Load loads an aggregate of the specified type.
func (r *InMemoryOrderRepo) Load(aggregateType string, id string) (*ProductionOrder, error) {

	events, ok := r.current[id]
	if !ok {
		return nil, &eventsourcing.ErrAggregateNotFound{}
	}

	order := NewProductionOrder(id)

	for _, v := range events {
		order.Apply(v, false)
		order.IncrementVersion()
	}

	return order, nil
}

// Save persists an aggregate.
func (r *InMemoryOrderRepo) Save(aggregate eventsourcing.AggregateRoot, _ *int64) error {

	//TODO: Look at the expected version
	for _, v := range aggregate.GetChanges() {
		r.current[aggregate.AggregateID()] = append(r.current[aggregate.AggregateID()], v)
		r.publisher.PublishEvent(v)
	}

	return nil
}

// ####### repo

// InMemoryPalletRepo provides an in memory repository implementation.
type InMemoryPalletRepo struct {
	current   map[string][]eventsourcing.EventMessage
	publisher eventsourcing.EventBus
}

// NewInMemoryRepo constructs an InMemoryRepo instance.
func NewInMemoryPalletRepo(eventBus eventsourcing.EventBus) *InMemoryPalletRepo {
	return &InMemoryPalletRepo{
		current:   make(map[string][]eventsourcing.EventMessage),
		publisher: eventBus,
	}
}

// Load loads an aggregate of the specified type.
func (r *InMemoryPalletRepo) Load(aggregateType string, id string) (*Pallet, error) {

	events, ok := r.current[id]
	if !ok {
		return nil, &eventsourcing.ErrAggregateNotFound{}
	}

	pallet := NewPallet(id)

	for _, v := range events {
		pallet.Apply(v, false)
		pallet.IncrementVersion()
	}

	return pallet, nil
}

// Save persists an aggregate.
func (r *InMemoryPalletRepo) Save(aggregate eventsourcing.AggregateRoot, _ *int64) error {

	//TODO: Look at the expected version
	for _, v := range aggregate.GetChanges() {
		r.current[aggregate.AggregateID()] = append(r.current[aggregate.AggregateID()], v)
		r.publisher.PublishEvent(v)
	}

	return nil
}

// ############### EventStoreDB

type ProductionOrderRepo struct {
	repo *eventsourcing.GetEventStoreCommonDomainRepo
}

// NewOrderRepo constructs a new InventoryItemRepository.
func NewProductionOrderRepo(eventStore *client.Client, eventBus eventsourcing.EventBus) (*ProductionOrderRepo, error) {

	r, err := eventsourcing.NewCommonDomainRepository(eventStore, eventBus)
	if err != nil {
		return nil, err
	}

	ret := &ProductionOrderRepo{
		repo: r,
	}

	// An aggregate factory creates an aggregate instance given the name of an aggregate.
	aggregateFactory := eventsourcing.NewDelegateAggregateFactory()
	aggregateFactory.RegisterDelegate(&ProductionOrder{},
		func(id string) eventsourcing.AggregateRoot { return NewProductionOrder(id) })
	ret.repo.SetAggregateFactory(aggregateFactory)

	// A stream name delegate constructs a stream name.
	// A common way to construct a stream name is to use a bounded context and
	// an aggregate id.
	// The interface for a stream name delegate takes a two strings. One may be
	// the aggregate type and the other the aggregate id. In this case the first
	// argument and the second argument are concatenated with a hyphen.
	streamNameDelegate := eventsourcing.NewDelegateStreamNamer()
	streamNameDelegate.RegisterDelegate(func(t string, id string) string {
		return t + "-" + id
	}, &ProductionOrder{})
	ret.repo.SetStreamNameDelegate(streamNameDelegate)

	// An event factory creates an instance of an event given the name of an event
	// as a string.
	eventFactory := eventsourcing.NewDelegateEventFactory()
	eventFactory.RegisterDelegate(&ProductionOrderCreated{},
		func() interface{} { return &ProductionOrderCreated{} })
	ret.repo.SetEventFactory(eventFactory)

	return ret, nil
}

// Load loads events for an aggregate.
//
// Returns an *InventoryAggregate.
func (r *ProductionOrderRepo) Load(aggregateType, id string) (*ProductionOrder, error) {
	ar, err := r.repo.Load(reflect.TypeOf(&ProductionOrder{}).Elem().Name(), id)
	if _, ok := err.(*eventsourcing.ErrAggregateNotFound); ok {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if ret, ok := ar.(*ProductionOrder); ok {
		return ret, nil
	}

	return nil, fmt.Errorf("could not cast aggregate returned to type of %s", reflect.TypeOf(&ProductionOrder{}).Elem().Name())
}

// Save persists an aggregate.
func (r *ProductionOrderRepo) Save(aggregate eventsourcing.AggregateRoot, expectedVersion *int64) error {
	return r.repo.Save(aggregate, expectedVersion)
}
