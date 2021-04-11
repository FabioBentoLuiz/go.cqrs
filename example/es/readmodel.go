package example

import "github.com/fabiobentoluiz/eventsourcing"

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
	ID   string
	Name string
}

// FakeDatabase is a simple in memory repository
type FakeDatabase struct {
	//Details map[string]*InventoryItemDetailsDto
	List []*ProductionOrderListDto
}

// NewFakeDatabase constructs a new FakeDatabase
func NewFakeDatabase() *FakeDatabase {
	return &FakeDatabase{
		//Details: make(map[string]*InventoryItemDetailsDto),
	}
}

// ReadModelFacade is an interface for the readmodel facade
type ReadModelFacade interface {
	GetProductionOrders() []*ProductionOrderListDto
}

// ReadModel is an implementation of the ReadModelFacade interface.
// ReadModel provides an in memory read model.
type ReadModel struct {
}

// GetProductionOrders returns a slice of all inventory items
func (m *ReadModel) GetProductionOrders() []*ProductionOrderListDto {
	return fakeDatabase.List
}

// ProductionOrderListView handles messages related to inventory and builds an
// in memory read model of inventory item summaries in a list.
type ProductionOrderListView struct {
}

// NewProductionOrderListView constructs a new ProductionOrderListView
func NewProductionOrderListView() *ProductionOrderListView {
	if fakeDatabase == nil {
		fakeDatabase = NewFakeDatabase()
	}

	return &ProductionOrderListView{}
}

// Handle processes events related to inventory and builds an in memory read model
func (v *ProductionOrderListView) Handle(message eventsourcing.EventMessage) {

	switch event := message.Event().(type) {

	case *ProductionOrderCreated:

		fakeDatabase.List = append(fakeDatabase.List, &ProductionOrderListDto{
			ID:   message.AggregateID(),
			Name: event.Name,
		})

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
