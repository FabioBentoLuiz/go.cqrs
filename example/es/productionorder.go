package example

import (
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

func (order *ProductionOrder) Create(cmd *CreateProductionOrder) error {
	created := ProductionOrderCreated{
		ID:            order.AggregateID(),
		Name:          cmd.Name,
		BagsToProduce: cmd.BagsToProduce,
	}

	em := eventsourcing.NewEventMessage(order.AggregateID(), &created, eventsourcing.Int64(order.CurrentVersion()))
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
