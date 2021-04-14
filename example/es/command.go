package example

type CreateProductionOrder struct {
	Name          string
	BagsToProduce int
}

type CreatePallet struct {
	OrderID string
	Bags    int
}
