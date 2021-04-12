package example

import (
	"errors"

	"github.com/fabiobentoluiz/eventsourcing"
)

type ProductionOrder struct {
	*eventsourcing.AggregateBase
	BagsToProduce int
	Activated     bool
}

func NewProductionOrder(id string) *ProductionOrder {
	order := ProductionOrder{
		AggregateBase: eventsourcing.NewAggregateBase(id),
	}

	return &order
}

func (order *ProductionOrder) Create(name string) error {
	if name == "" {
		return errors.New("the name cannot be empty")
	}

	created := ProductionOrderCreated{
		ID:            order.AggregateID(),
		Name:          name,
		BagsToProduce: order.BagsToProduce,
	}

	em := eventsourcing.NewEventMessage(order.AggregateID(), &created, eventsourcing.Uint64(uint64(order.CurrentVersion())))
	order.Apply(em, true)

	return nil
}

func (order *ProductionOrder) Apply(evtMessage eventsourcing.EventMessage, isNew bool) {
	if isNew {
		order.TrackChange(evtMessage)
	}

	switch e := evtMessage.Event().(type) {
	case *ProductionOrderCreated:
		order.Activated = true
		order.BagsToProduce = e.BagsToProduce
	}
}
