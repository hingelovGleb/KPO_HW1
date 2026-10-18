package domain

type VehicleKind struct {
	Title    string
	Electric bool
	Create   func(spec VehicleSpec, dailyKWh float64) Vehicle
}

// Чтобы добавить новый вид транспорта достаточно сделать тип и зарегистрировать его тут.
// Остальной код менять не надо.
type VehicleCatalog struct {
	kinds []VehicleKind
}

func (c *VehicleCatalog) Register(k VehicleKind) {
	c.kinds = append(c.kinds, k)
}

func (c *VehicleCatalog) Kinds() []VehicleKind {
	return c.kinds
}

func DefaultCatalog() *VehicleCatalog {
	c := &VehicleCatalog{}
	c.Register(VehicleKind{
		Title: "Электросамокат", Electric: true,
		Create: func(s VehicleSpec, kwh float64) Vehicle { return NewScooter(s, kwh) },
	})
	c.Register(VehicleKind{
		Title: "Велосипед", Electric: false,
		Create: func(s VehicleSpec, _ float64) Vehicle { return NewBicycle(s) },
	})
	c.Register(VehicleKind{
		Title: "Электровелосипед", Electric: true,
		Create: func(s VehicleSpec, kwh float64) Vehicle { return NewEBike(s, kwh) },
	})
	return c
}
