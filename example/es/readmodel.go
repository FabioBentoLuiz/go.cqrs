package example

import (
	"log"

	"github.com/fabiobentoluiz/eventsourcing"
)

var fakeDatabase *FakeDatabase

// NewReadModel constructs a new read model
func NewReadModel() *ReadModel {
	if fakeDatabase == nil {
		fakeDatabase = NewFakeDatabase()
	}

	return &ReadModel{}
}

// ProductionOrderListDto provides a lightweight lookup view of a production order
type ProductionOrderListDto struct {
	ID            string
	Name          string
	BagsToProduce int
	Version       uint64
}

type PalletListDto struct {
	ID      string
	OrderID string
	Bags    int
	Version uint64
}

// FakeDatabase is a simple in memory repository
type FakeDatabase struct {
	Pallets []*PalletListDto
	Orders  []*ProductionOrderListDto
}

// NewFakeDatabase constructs a new FakeDatabase
func NewFakeDatabase() *FakeDatabase {
	return &FakeDatabase{
		Pallets: make([]*PalletListDto, 0),
	}
}

// ReadModelFacade is an interface for the readmodel facade
type ReadModelFacade interface {
	GetProductionOrders() []*ProductionOrderListDto
	GetPallets() []*PalletListDto
}

// ReadModel is an implementation of the ReadModelFacade interface.
// ReadModel provides an in memory read model.
type ReadModel struct {
}

// GetProductionOrders returns a slice of all production orders
func (m *ReadModel) GetProductionOrders() []*ProductionOrderListDto {
	return fakeDatabase.Orders
}

// GetPallets returns a slice of all pallets
func (m *ReadModel) GetPallets() []*PalletListDto {
	return fakeDatabase.Pallets
}

// ProductionOrderListView handles messages related to orders and builds an
// in memory read model of order summaries in a list.
type ProductionOrderListView struct {
}

type PalletListView struct {
}

// NewProductionOrderListView constructs a new ProductionOrderListView
func NewProductionOrderListView() *ProductionOrderListView {
	if fakeDatabase == nil {
		fakeDatabase = NewFakeDatabase()
	}

	return &ProductionOrderListView{}
}

// NewPalletListView constructs a new PalletListView
func NewPalletListView() *PalletListView {
	if fakeDatabase == nil {
		fakeDatabase = NewFakeDatabase()
	}

	return &PalletListView{}
}

// Handle processes events related to order and builds an in memory read model
func (v *ProductionOrderListView) Handle(message eventsourcing.EventMessage) {

	switch event := message.Event().(type) {

	case *ProductionOrderCreated:

		fakeDatabase.Orders = append(fakeDatabase.Orders, &ProductionOrderListDto{
			ID:            message.AggregateID(),
			Name:          event.Name,
			BagsToProduce: event.BagsToProduce,
			Version:       0,
		})

	default:
		log.Printf("there is no handler for the event %s", event)

		/*case *InventoryItemRenamed:

			for _, v := range bullShitDatabase.List {
				if v.ID == message.AggregateID() {
					v.Name = event.NewName
					break
				}
			}

		case *InventoryItemDeactivated:
			i := -1
			for k, v := range bullShitDatabase.List {
				if v.ID == message.AggregateID() {
					i = k
					break
				}
			}

			if i >= 0 {
				bullShitDatabase.List = append(
					bullShitDatabase.List[:i],
					bullShitDatabase.List[i+1:]...,
				)
			}*/
	}
}

// Handle processes events related to pallets and builds an in memory read model
func (v *PalletListView) Handle(message eventsourcing.EventMessage) {

	switch event := message.Event().(type) {

	case *PalletCreated:

		fakeDatabase.Pallets = append(fakeDatabase.Pallets, &PalletListDto{
			ID:      message.AggregateID(),
			Bags:    event.Bags,
			OrderID: event.OrderID,
			Version: 0,
		})

	default:
		log.Printf("there is no handler for the event %s", event)
	}
}
