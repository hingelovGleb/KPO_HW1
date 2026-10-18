package domain

type VehicleSpec struct {
	InventoryNumber string
	Name            string
	Simplicity      int
	MaxSpeedKmh     int
	Condition       int
}

// общие поля для всего транспорта
type baseVehicle struct {
	spec VehicleSpec
}

func (b baseVehicle) InventoryNumber() string { return b.spec.InventoryNumber }
func (b baseVehicle) Name() string            { return b.spec.Name }
func (b baseVehicle) Simplicity() int         { return b.spec.Simplicity }
func (b baseVehicle) MaxSpeedKmh() int        { return b.spec.MaxSpeedKmh }
func (b baseVehicle) Condition() int          { return b.spec.Condition }

// встраиваем в то что работает от электричества
type electricMotor struct {
	dailyKWh float64
}

func (m electricMotor) DailyEnergyKWh() float64 { return m.dailyKWh }

type Bicycle struct {
	baseVehicle
}

func NewBicycle(spec VehicleSpec) *Bicycle {
	return &Bicycle{baseVehicle{spec}}
}

func (*Bicycle) Kind() string { return "Велосипед" }

type Scooter struct {
	baseVehicle
	electricMotor
}

func NewScooter(spec VehicleSpec, dailyKWh float64) *Scooter {
	return &Scooter{baseVehicle{spec}, electricMotor{dailyKWh}}
}

func (*Scooter) Kind() string { return "Электросамокат" }

type EBike struct {
	baseVehicle
	electricMotor
}

func NewEBike(spec VehicleSpec, dailyKWh float64) *EBike {
	return &EBike{baseVehicle{spec}, electricMotor{dailyKWh}}
}

func (*EBike) Kind() string { return "Электровелосипед" }
