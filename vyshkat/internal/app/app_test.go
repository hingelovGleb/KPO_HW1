package app

import (
	"bytes"
	"strings"
	"testing"

	"vyshkat/internal/di"
	"vyshkat/internal/domain"
	"vyshkat/internal/fleet"
)

// fake техосмотр который всем отказывает
type rejectAll struct{ calls int }

func (r *rejectAll) Inspect(domain.Vehicle) domain.InspectionResult {
	r.calls++
	return domain.Rejected("fake: техосмотр закрыт")
}

// Главный тест на DI. Меняем техосмотр в контейнере и Fleet берёт подмену
func TestInspectorReplacedInContainer(t *testing.T) {
	c := BuildContainer()
	fake := &rejectAll{}
	di.RegisterInstance[fleet.Inspector](c, fake)

	f := di.Resolve[*fleet.Fleet](c)
	good := domain.NewBicycle(domain.VehicleSpec{InventoryNumber: "VK-1", Name: "Идеальный",
		Simplicity: 10, MaxSpeedKmh: 20, Condition: 100})

	res, _ := f.AddVehicle(good)
	if res.Approved || fake.calls != 1 || len(f.Vehicles()) != 0 {
		t.Error("Fleet должен использовать подменённый техосмотр")
	}
}

func run(input string) string {
	out := &bytes.Buffer{}
	NewConsoleFromContainer(BuildContainer(), strings.NewReader(input), out).Run()
	return out.String()
}

func TestConsoleAddScooter(t *testing.T) {
	// "abc" неправильный ввод. Потом самокат и список номеров
	text := run("1\nabc\n1\nNinebot\n7\n25\n95\n0,6\n6\n0\n")
	for _, want := range []string{"Нужно число от 1 до 3", "Принят: Электросамокат «Ninebot». Номер VK-1", "VK-1  Ninebot"} {
		if !strings.Contains(text, want) {
			t.Errorf("в выводе нет %q", want)
		}
	}
}

func TestDemo(t *testing.T) {
	text := run("8\n0\n")
	for _, want := range []string{"Не принят:", "такой номер уже есть", "Всего: 8", "Всего 10.30 кВт·ч в день"} {
		if !strings.Contains(text, want) {
			t.Errorf("в выводе нет %q", want)
		}
	}
}
