package example

import "github.com/fabiobentoluiz/eventsourcing"

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
func (r *InMemoryOrderRepo) Save(aggregate eventsourcing.AggregateRoot, _ *uint64) error {

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
func (r *InMemoryPalletRepo) Save(aggregate eventsourcing.AggregateRoot, _ *uint64) error {

	//TODO: Look at the expected version
	for _, v := range aggregate.GetChanges() {
		r.current[aggregate.AggregateID()] = append(r.current[aggregate.AggregateID()], v)
		r.publisher.PublishEvent(v)
	}

	return nil
}
