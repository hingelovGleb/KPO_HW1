package domain

// Inventoried это всё что стоит на балансе и имеет инвентарный номер
type Inventoried interface {
	InventoryNumber() string
	Name() string
}

// EnergyConsumer реализуют только те кто реально тратит энергию
type EnergyConsumer interface {
	DailyEnergyKWh() float64
}

type Vehicle interface {
	Inventoried
	Kind() string
	Simplicity() int // от 1 до 10
	MaxSpeedKmh() int
	Condition() int // техсостояние в процентах
}

type Thing interface {
	Inventoried
	Category() string
}
