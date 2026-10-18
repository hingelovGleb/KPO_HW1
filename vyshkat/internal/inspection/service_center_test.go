package inspection

import (
	"testing"

	"vyshkat/internal/domain"
)

func spec(condition, speed int) domain.VehicleSpec {
	return domain.VehicleSpec{InventoryNumber: "VK-1", Name: "Тест",
		Simplicity: 7, MaxSpeedKmh: speed, Condition: condition}
}

func TestServiceCenterDecisions(t *testing.T) {
	sc := NewServiceCenter(DefaultConfig())

	tests := []struct {
		name     string
		vehicle  domain.Vehicle
		approved bool
	}{
		{"исправный самокат", domain.NewScooter(spec(90, 25), 0.6), true},
		{"ровно 70% проходит", domain.NewBicycle(spec(70, 20)), true},
		{"плохое состояние", domain.NewScooter(spec(50, 25), 0.6), false},
		{"слишком быстрый", domain.NewEBike(spec(90, 80), 1), false},
		{"слишком большой расход", domain.NewEBike(spec(90, 25), 9), false},
	}
	for _, tt := range tests {
		if got := sc.Inspect(tt.vehicle).Approved; got != tt.approved {
			t.Errorf("%s: получили %v а ждали %v", tt.name, got, tt.approved)
		}
	}
}
