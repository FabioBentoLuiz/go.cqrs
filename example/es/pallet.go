package example

import (
	"github.com/fabiobentoluiz/eventsourcing"
)

type Pallet struct {
	*eventsourcing.AggregateBase
	ID      string
	Bags    int
	OrderID string
}

func NewPallet(id string) *Pallet {
	pallet := Pallet{
		AggregateBase: eventsourcing.NewAggregateBase(id),
	}

	return &pallet
}

func (pallet *Pallet) Create(cmd *CreatePallet) error {

	created := PalletCreated{
		ID:      pallet.AggregateID(),
		Bags:    cmd.Bags,
		OrderID: cmd.OrderID,
	}

	em := eventsourcing.NewEventMessage(pallet.AggregateID(), &created, eventsourcing.Uint64(uint64(pallet.CurrentVersion())))
	pallet.Apply(em, true)

	return nil
}

func (pallet *Pallet) Apply(evtMessage eventsourcing.EventMessage, isNew bool) {
	if isNew {
		pallet.TrackChange(evtMessage)
	}

	switch e := evtMessage.Event().(type) {
	case *PalletCreated:
		pallet.Bags = e.Bags
		pallet.ID = e.ID
		pallet.OrderID = e.OrderID
	}
}
