package example

import (
	"log"

	"github.com/fabiobentoluiz/eventsourcing"
)

//########## Orders

type ProductionOrderRepository interface {
	Load(string, string) (*ProductionOrder, error)
	Save(eventsourcing.AggregateRoot, *uint64) error
}

type ProductionOrderCommandHandler struct {
	repo ProductionOrderRepository
}

func NewProductionOrderCommandHandler(repo ProductionOrderRepository) *ProductionOrderCommandHandler {
	handler := ProductionOrderCommandHandler{
		repo: repo,
	}

	return &handler
}

func (handler *ProductionOrderCommandHandler) Handle(cmdMessage eventsourcing.CommandMessage) error {

	switch cmd := cmdMessage.Command().(type) {
	case *CreateProductionOrder:
		order := NewProductionOrder(cmdMessage.AggregateID())
		order.BagsToProduce = cmd.BagsToProduce
		if err := order.Create(cmd.Name); err != nil {
			return &eventsourcing.ErrCommandExecution{Command: cmdMessage, Reason: err.Error()}
		}
		return handler.repo.Save(order, eventsourcing.Uint64(uint64(order.OriginalVersion())))

		// case *DeactivateInventoryItem:

		// item, _ = h.repo.Load(reflect.TypeOf(&InventoryItem{}).Elem().Name(), message.AggregateID())
		// if err := item.Deactivate(); err != nil {
		// 	return &ycq.ErrCommandExecution{Command: message, Reason: err.Error()}
		// }
		// return h.repo.Save(item, ycq.Int(item.OriginalVersion()))
	}

	return nil
}

//########## Pallets

type PalletRepository interface {
	Load(string, string) (*Pallet, error)
	Save(eventsourcing.AggregateRoot, *uint64) error
}

type PalletCommandHandler struct {
	repo PalletRepository
}

func NewPalletCommandHandler(repo PalletRepository) *PalletCommandHandler {
	handler := PalletCommandHandler{
		repo: repo,
	}

	return &handler
}

func (handler *PalletCommandHandler) Handle(cmdMessage eventsourcing.CommandMessage) error {

	switch cmd := cmdMessage.Command().(type) {
	case *CreatePallet:
		pallet := NewPallet(cmdMessage.AggregateID())
		if err := pallet.Create(cmd); err != nil {
			return &eventsourcing.ErrCommandExecution{Command: cmdMessage, Reason: err.Error()}
		}
		return handler.repo.Save(pallet, eventsourcing.Uint64(uint64(pallet.OriginalVersion())))

	default:
		log.Printf("there is no handler for the command %s", cmd)
	}

	return nil
}
