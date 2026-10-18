package app

import "vyshkat/internal/domain"

// Показывает все сценарии на готовых данных
func (c *Console) RunDemo() {
	spec := func(name string, simplicity, speed, condition int) domain.VehicleSpec {
		return domain.VehicleSpec{InventoryNumber: c.numbers.Next(), Name: name,
			Simplicity: simplicity, MaxSpeedKmh: speed, Condition: condition}
	}

	c.printf("\n=== Добавляем транспорт ===\n")
	scooter := domain.NewScooter(spec("Ninebot Max G30", 7, 25, 95), 0.6)
	c.acceptVehicle(scooter)
	c.acceptVehicle(domain.NewScooter(spec("Kugoo Kirin G4", 4, 40, 90), 1.1))
	c.acceptVehicle(domain.NewBicycle(spec("Stels Navigator 300", 9, 20, 88)))
	c.acceptVehicle(domain.NewEBike(spec("Eltreco XT 750", 6, 25, 75), 0.9))
	c.acceptVehicle(domain.NewScooter(spec("Старый самокат", 8, 25, 40), 0.5))
	c.acceptVehicle(domain.NewEBike(spec("Очень быстрый велик", 5, 60, 90), 7.5))

	c.printf("\n=== Пробуем добавить второй раз тот же номер ===\n")
	dup := domain.VehicleSpec{InventoryNumber: scooter.InventoryNumber(),
		Name: "Копия", Simplicity: 7, MaxSpeedKmh: 25, Condition: 99}
	c.acceptVehicle(domain.NewScooter(dup, 0.5))

	c.printf("\n=== Добавляем вещи ===\n")
	c.acceptThing(domain.NewHelmet(c.numbers.Next(), "Синий шлем", "M"))
	c.acceptThing(domain.NewHelmet(c.numbers.Next(), "Детский шлем", "S"))
	c.acceptThing(domain.NewDockStation(c.numbers.Next(), "Док у Покровки", 10, 3.2))
	c.acceptThing(domain.NewChargingCabinet(c.numbers.Next(), "Шкаф на складе", 4.5))

	c.printf("\n=== Отчёты ===\n")
	c.PrintSummary()
	c.PrintEnergy()
	c.PrintBeginners()
	c.PrintInventory()
	c.PrintRejections()
}
