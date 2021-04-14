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
	orderListView := example.NewProductionOrderListView()

	// Create an EventBus
	eventBus := eventsourcing.NewInternalEventBus()
	// Register the listView as an event handler on the event bus
	// for the events specified.
	eventBus.AddHandler(orderListView,
		&example.ProductionOrderCreated{})

	// Create a PalletListView
	palletListView := example.NewPalletListView()
	eventBus.AddHandler(palletListView,
		&example.PalletCreated{})

	// Here we use an in memory event repository.
	orderRepo := example.NewInMemoryOrderRepo(eventBus)

	// Create an ProductionOrderCommandHandler instance
	productionOrderCommandHandler := example.NewProductionOrderCommandHandler(orderRepo)

	// Create a dispatcher
	dispatcher = eventsourcing.NewInMemoryDispatcher()
	// Register the production order command handlers instance as a command handler
	// for the events specified.
	err := dispatcher.RegisterHandler(productionOrderCommandHandler,
		&example.CreateProductionOrder{})
	if err != nil {
		log.Fatal(err)
	}

	palletRepo := example.NewInMemoryPalletRepo(eventBus)

	palletCommandHandler := example.NewPalletCommandHandler(palletRepo)

	err = dispatcher.RegisterHandler(palletCommandHandler,
		&example.CreatePallet{})
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	bags := 400
	orderId := createOneOrder(bags)

	createPallets(orderId, bags)

	orders := readModel.GetProductionOrders()
	for _, o := range orders {
		fmt.Printf("order %v\n", o)
	}

	pallets := readModel.GetPallets()
	for _, p := range pallets {
		fmt.Printf("pallet %v\n", p)
	}
}

func createOneOrder(bags int) string {
	id := eventsourcing.NewUUID()
	command := eventsourcing.NewCommandMessage(
		id,
		&example.CreateProductionOrder{
			Name:          "test-order",
			BagsToProduce: bags,
		})

	err := dispatcher.Dispatch(command)
	if err != nil {
		log.Println(err)
	}

	return id
}

// createPallets simulate the pallet creation
// each pallet has 100 bags and they are always full
func createPallets(orderId string, bags int) {
	bagsPerPallet := 100
	pallets := bags / bagsPerPallet

	for i := 0; i < pallets; i++ {
		id := eventsourcing.NewUUID()
		command := eventsourcing.NewCommandMessage(
			id,
			&example.CreatePallet{
				OrderID: orderId,
				Bags:    bagsPerPallet,
			})

		err := dispatcher.Dispatch(command)
		if err != nil {
			log.Println(err)
		}
	}
}
