package example

import "github.com/fabiobentoluiz/eventsourcing"

type ProductionOrderRepository interface {
	Load(string, string) (ProductionOrder, error)
	Save(eventsourcing.AggregateRoot, *int) error
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

}
