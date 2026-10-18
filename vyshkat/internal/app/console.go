package app

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"vyshkat/internal/domain"
	"vyshkat/internal/fleet"
	"vyshkat/internal/inventory"
	"vyshkat/internal/report"
)

// Console отвечает только за ввод и вывод. Логика лежит в fleet и report.
type Console struct {
	fleet    *fleet.Fleet
	reporter *report.Reporter
	catalog  *domain.VehicleCatalog
	numbers  *inventory.Generator
	in       *bufio.Scanner
	out      io.Writer
}

func NewConsole(f *fleet.Fleet, r *report.Reporter, cat *domain.VehicleCatalog,
	gen *inventory.Generator, in io.Reader, out io.Writer) *Console {
	return &Console{fleet: f, reporter: r, catalog: cat, numbers: gen,
		in: bufio.NewScanner(in), out: out}
}

const menu = `
Что сделать?
  1. Добавить транспорт
  2. Добавить вещь
  3. Сколько всего техники и вещей
  4. Сколько энергии тратит парк за день
  5. Что можно дать новичку
  6. Список всех номеров
  7. Кто не прошёл техосмотр
  8. Демо (заполнить парк примерами)
  0. Выход
`

func (c *Console) Run() {
	c.printf("ВышКат. Учёт самокатов и велосипедов\n")
	for {
		c.printf(menu)
		choice, ok := c.ask("Номер пункта")
		if !ok || choice == "0" {
			c.printf("Пока!\n")
			return
		}
		ok = true
		switch choice {
		case "1":
			ok = c.addVehicle()
		case "2":
			ok = c.addThing()
		case "3":
			c.PrintSummary()
		case "4":
			c.PrintEnergy()
		case "5":
			c.PrintBeginners()
		case "6":
			c.PrintInventory()
		case "7":
			c.PrintRejections()
		case "8":
			c.RunDemo()
		default:
			c.printf("Такого пункта нет\n")
		}
		// ok=false значит ввод закончился посреди вопросов
		if !ok {
			c.printf("Пока!\n")
			return
		}
	}
}

func (c *Console) addVehicle() bool {
	kinds := c.catalog.Kinds()
	for i, k := range kinds {
		c.printf("  %d. %s\n", i+1, k.Title)
	}
	idx, ok := c.askInt("Какой транспорт", 1, len(kinds))
	if !ok {
		return false
	}
	kind := kinds[idx-1]

	var s domain.VehicleSpec
	if s.Name, ok = c.askText("Название"); !ok {
		return false
	}
	// тут проверяем только что ввели число. Нормы проверяет техосмотр
	if s.Simplicity, ok = c.askInt("Насколько просто ездить (1-10)", 0, 100); !ok {
		return false
	}
	if s.MaxSpeedKmh, ok = c.askInt("Макс. скорость в км/ч", 0, 200); !ok {
		return false
	}
	if s.Condition, ok = c.askInt("Состояние в процентах", 0, 100); !ok {
		return false
	}
	kwh := 0.0
	if kind.Electric {
		if kwh, ok = c.askFloat("Расход энергии за день в кВт·ч"); !ok {
			return false
		}
	}

	s.InventoryNumber = c.numbers.Next()
	c.acceptVehicle(kind.Create(s, kwh))
	return true
}

func (c *Console) acceptVehicle(v domain.Vehicle) {
	res, err := c.fleet.AddVehicle(v)
	switch {
	case err != nil:
		c.printf("Ошибка: %s «%s»: %v\n", v.Kind(), v.Name(), err)
	case res.Approved:
		c.printf("Принят: %s «%s». Номер %s\n", v.Kind(), v.Name(), v.InventoryNumber())
	default:
		c.printf("Не принят: %s «%s». Причина: %s\n", v.Kind(), v.Name(), res.Reason)
	}
}

func (c *Console) addThing() bool {
	c.printf("  1. Шлем\n  2. Док-станция\n  3. Зарядный шкаф\n")
	cat, ok := c.askInt("Какая вещь", 1, 3)
	if !ok {
		return false
	}
	name, ok := c.askText("Название")
	if !ok {
		return false
	}

	var thing domain.Thing
	switch cat {
	case 1:
		size, ok := c.askText("Размер (S, M или L)")
		if !ok {
			return false
		}
		thing = domain.NewHelmet(c.numbers.Next(), name, size)
	case 2:
		slots, ok := c.askInt("Сколько мест для самокатов", 1, 100)
		if !ok {
			return false
		}
		kwh, ok := c.askFloat("Расход энергии за день в кВт·ч")
		if !ok {
			return false
		}
		thing = domain.NewDockStation(c.numbers.Next(), name, slots, kwh)
	default:
		kwh, ok := c.askFloat("Расход энергии за день в кВт·ч")
		if !ok {
			return false
		}
		thing = domain.NewChargingCabinet(c.numbers.Next(), name, kwh)
	}
	c.acceptThing(thing)
	return true
}

func (c *Console) acceptThing(t domain.Thing) {
	if err := c.fleet.AddThing(t); err != nil {
		c.printf("Ошибка: %s «%s»: %v\n", t.Category(), t.Name(), err)
		return
	}
	c.printf("Добавлено: %s «%s». Номер %s\n", t.Category(), t.Name(), t.InventoryNumber())
}

func (c *Console) PrintSummary() {
	s := c.reporter.Summary()
	c.printf("\nТранспорт: %d\nВещи: %d\nВсего: %d\n", s.Vehicles, s.Things, s.Total())
	kinds := make([]string, 0, len(s.ByKind))
	for k := range s.ByKind {
		kinds = append(kinds, k)
	}
	sort.Strings(kinds)
	for _, k := range kinds {
		c.printf("  %s: %d\n", k, s.ByKind[k])
	}
}

func (c *Console) PrintEnergy() {
	e := c.reporter.Energy()
	c.printf("\nТранспорт тратит %.2f кВт·ч в день\n", e.VehiclesKWh)
	c.printf("Вещи тратят %.2f кВт·ч в день\n", e.ThingsKWh)
	c.printf("Всего %.2f кВт·ч в день\n", e.TotalKWh())
}

func (c *Console) PrintBeginners() {
	c.printf("\nМожно дать новичку (%s):\n", c.reporter.BeginnerRule())
	list := c.reporter.BeginnerVehicles()
	if len(list) == 0 {
		c.printf("  Пока ничего нет\n")
	}
	for _, v := range list {
		c.printf("  %s  %s «%s»\n", v.InventoryNumber(), v.Kind(), v.Name())
	}
}

func (c *Console) PrintInventory() {
	c.printf("\nВсе номера:\n")
	items := c.reporter.Inventory()
	if len(items) == 0 {
		c.printf("  Пока ничего нет\n")
	}
	for _, it := range items {
		c.printf("  %s  %s\n", it.InventoryNumber(), it.Name())
	}
}

func (c *Console) PrintRejections() {
	c.printf("\nНе прошли техосмотр:\n")
	rej := c.fleet.Rejections()
	if len(rej) == 0 {
		c.printf("  Таких нет\n")
	}
	for _, r := range rej {
		c.printf("  %s «%s». Причина: %s\n", r.Vehicle.Kind(), r.Vehicle.Name(), r.Reason)
	}
}

func (c *Console) printf(format string, args ...any) {
	fmt.Fprintf(c.out, format, args...)
}

// ok=false если ввод закончился
func (c *Console) ask(prompt string) (string, bool) {
	c.printf("%s: ", prompt)
	if !c.in.Scan() {
		return "", false
	}
	return strings.TrimSpace(c.in.Text()), true
}

func (c *Console) askText(prompt string) (string, bool) {
	for {
		s, ok := c.ask(prompt)
		if !ok || s != "" {
			return s, ok
		}
		c.printf("Нужно что-то ввести\n")
	}
}

func (c *Console) askInt(prompt string, min, max int) (int, bool) {
	for {
		s, ok := c.ask(prompt)
		if !ok {
			return 0, false
		}
		n, err := strconv.Atoi(s)
		if err == nil && n >= min && n <= max {
			return n, true
		}
		c.printf("Нужно число от %d до %d\n", min, max)
	}
}

func (c *Console) askFloat(prompt string) (float64, bool) {
	for {
		s, ok := c.ask(prompt)
		if !ok {
			return 0, false
		}
		x, err := strconv.ParseFloat(strings.Replace(s, ",", ".", 1), 64)
		if err == nil && x >= 0 {
			return x, true
		}
		c.printf("Нужно число, например 1.5\n")
	}
}
