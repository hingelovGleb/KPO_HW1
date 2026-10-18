package domain

type baseThing struct {
	inv  string
	name string
}

func (b baseThing) InventoryNumber() string { return b.inv }
func (b baseThing) Name() string            { return b.name }

type Helmet struct {
	baseThing
	Size string
}

func NewHelmet(inv, name, size string) *Helmet {
	return &Helmet{baseThing{inv, name}, size}
}

func (*Helmet) Category() string { return "Шлем" }

type DockStation struct {
	baseThing
	electricMotor
	Slots int
}

func NewDockStation(inv, name string, slots int, dailyKWh float64) *DockStation {
	return &DockStation{baseThing{inv, name}, electricMotor{dailyKWh}, slots}
}

func (*DockStation) Category() string { return "Док-станция" }

type ChargingCabinet struct {
	baseThing
	electricMotor
}

func NewChargingCabinet(inv, name string, dailyKWh float64) *ChargingCabinet {
	return &ChargingCabinet{baseThing{inv, name}, electricMotor{dailyKWh}}
}

func (*ChargingCabinet) Category() string { return "Зарядный шкаф" }
