package example

import "github.com/fabiobentoluiz/eventsourcing"

// InMemoryRepo provides an in memory repository implementation.
type InMemoryRepo struct {
	current   map[string][]eventsourcing.EventMessage
	publisher eventsourcing.EventBus
}

// NewInMemoryRepo constructs an InMemoryRepo instance.
func NewInMemoryRepo(eventBus eventsourcing.EventBus) *InMemoryRepo {
	return &InMemoryRepo{
		current:   make(map[string][]eventsourcing.EventMessage),
		publisher: eventBus,
	}
}

// Load loads an aggregate of the specified type.
func (r *InMemoryRepo) Load(aggregateType string, id string) (*ProductionOrder, error) {

	events, ok := r.current[id]
	if !ok {
		return nil, &eventsourcing.ErrAggregateNotFound{}
	}

	inventoryItem := NewProductionOrder(id)

	for _, v := range events {
		inventoryItem.Apply(v, false)
		inventoryItem.IncrementVersion()
	}

	return inventoryItem, nil
}

// Save persists an aggregate.
func (r *InMemoryRepo) Save(aggregate eventsourcing.AggregateRoot, _ *uint64) error {

	//TODO: Look at the expected version
	for _, v := range aggregate.GetChanges() {
		r.current[aggregate.AggregateID()] = append(r.current[aggregate.AggregateID()], v)
		r.publisher.PublishEvent(v)
	}

	return nil

}
