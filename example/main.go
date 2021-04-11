package main

import (
	"fmt"
	"log"

	"github.com/fabiobentoluiz/eventsourcing"
	example "github.com/fabiobentoluiz/eventsourcing/example/es"
)

var (
	readModel  example.ReadModelFacade
	dispatcher eventsourcing.Dispatcher
)

func init() {
	// CQRS Infrastructure configuration

	// Configure the read model

	// Create a readModel instance
	readModel = example.NewReadModel()

	// Create a ProductionOrderListView
	listView := example.NewProductionOrderListView()

	// Create an EventBus
	eventBus := eventsourcing.NewInternalEventBus()
	// Register the listView as an event handler on the event bus
	// for the events specified.
	eventBus.AddHandler(listView,
		&example.ProductionOrderCreated{})

	// Here we use an in memory event repository.
	repo := example.NewInMemoryRepo(eventBus)

	// Create an ProductionOrderCommandHandler instance
	productionOrderCommandHandler := example.NewProductionOrderCommandHandler(repo)

	// Create a dispatcher
	dispatcher = eventsourcing.NewInMemoryDispatcher()
	// Register the inventory command handlers instance as a command handler
	// for the events specified.
	err := dispatcher.RegisterHandler(productionOrderCommandHandler,
		&example.CreateProductionOrder{})
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	id := eventsourcing.NewUUID()
	em := eventsourcing.NewCommandMessage(
		id,
		&example.CreateProductionOrder{
			Name: "test-order",
		})

	err := dispatcher.Dispatch(em)
	if err != nil {
		log.Println(err)
	}

	orders := readModel.GetProductionOrders()

	for _, o := range orders {
		fmt.Printf("order %v", o)
	}
}

func init() {

}
