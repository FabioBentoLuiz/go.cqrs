package example

import "github.com/fabiobentoluiz/eventsourcing"

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
		if err := order.Create(cmd.Name); err != nil {
			return &eventsourcing.ErrCommandExecution{Command: cmdMessage, Reason: err.Error()}
		}
		return handler.repo.Save(order, eventsourcing.Uint64(uint64(order.OriginalVersion())))
	}

	return nil
}
