package report

import (
	"math"
	"testing"

	"vyshkat/internal/domain"
)

// подмена парка для отчётов
type fakeSource struct {
	vehicles []domain.Vehicle
	things   []domain.Thing
}

func (f fakeSource) Vehicles() []domain.Vehicle { return f.vehicles }
func (f fakeSource) Things() []domain.Thing     { return f.things }

func vspec(simplicity, speed, condition int) domain.VehicleSpec {
	return domain.VehicleSpec{InventoryNumber: "VK", Name: "Тест",
		Simplicity: simplicity, MaxSpeedKmh: speed, Condition: condition}
}

func sampleSource() fakeSource {
	return fakeSource{
		vehicles: []domain.Vehicle{
			domain.NewScooter(vspec(7, 25, 95), 0.6),
			domain.NewScooter(vspec(4, 40, 90), 1.0),
			domain.NewBicycle(vspec(9, 20, 85)),
			domain.NewEBike(vspec(6, 25, 75), 0.9),
		},
		things: []domain.Thing{
			domain.NewHelmet("VK", "Шлем", "M"),
			domain.NewDockStation("VK", "Док", 10, 3.2),
		},
	}
}

// велосипед и шлем энергию не тратят
func TestEnergy(t *testing.T) {
	e := NewReporter(sampleSource(), DefaultBeginnerPolicy()).Energy()
	if math.Abs(e.TotalKWh()-5.7) > 1e-9 || e.Consumers != 4 {
		t.Errorf("неверный расход: %+v", e)
	}
}

func TestSafeBeginnerPolicy(t *testing.T) {
	p := DefaultBeginnerPolicy()
	tests := []struct {
		name string
		v    domain.Vehicle
		want bool
	}{
		{"подходит", domain.NewBicycle(vspec(6, 25, 80)), true},
		{"сложный", domain.NewBicycle(vspec(5, 20, 90)), false},
		{"быстрый", domain.NewScooter(vspec(8, 30, 90), 1), false},
		{"потрёпанный", domain.NewEBike(vspec(8, 25, 79), 1), false},
	}
	for _, tt := range tests {
		if got := p.Allows(tt.v); got != tt.want {
			t.Errorf("%s: получили %v а ждали %v", tt.name, got, tt.want)
		}
	}
}

// один и тот же парк с разными правилами даёт разные списки
func TestBeginnerVehiclesUsesInjectedPolicy(t *testing.T) {
	src := sampleSource()
	safe := NewReporter(src, DefaultBeginnerPolicy()).BeginnerVehicles()
	basic := NewReporter(src, SimplicityPolicy{MinSimplicity: 6}).BeginnerVehicles()

	if len(safe) != 2 || len(basic) != 3 {
		t.Errorf("получили %d и %d а ждали 2 и 3", len(safe), len(basic))
	}
}
