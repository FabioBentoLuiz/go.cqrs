package example

type ProductionOrderCreated struct {
	ID            string
	Name          string
	BagsToProduce int
}

type PalletCreated struct {
	ID      string
	Bags    int
	OrderID string
}
