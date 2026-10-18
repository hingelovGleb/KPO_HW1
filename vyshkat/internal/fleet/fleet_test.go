package fleet

import (
	"errors"
	"testing"

	"vyshkat/internal/domain"
)

// Подмена техосмотра. Возвращает заданный ответ и считает вызовы.
type stubInspector struct {
	result domain.InspectionResult
	calls  int
}

func (s *stubInspector) Inspect(domain.Vehicle) domain.InspectionResult {
	s.calls++
	return s.result
}

func scooter(inv string) domain.Vehicle {
	return domain.NewScooter(domain.VehicleSpec{InventoryNumber: inv, Name: "Тест",
		Simplicity: 7, MaxSpeedKmh: 25, Condition: 90}, 0.5)
}

func TestApprovedVehicleIsAdded(t *testing.T) {
	f := New(&stubInspector{result: domain.Approved()})

	res, err := f.AddVehicle(scooter("VK-1"))
	if err != nil || !res.Approved || len(f.Vehicles()) != 1 {
		t.Error("самокат должен попасть в парк")
	}
}

func TestRejectedVehicleIsNotAdded(t *testing.T) {
	f := New(&stubInspector{result: domain.Rejected("тормоза")})

	_, _ = f.AddVehicle(scooter("VK-1"))
	if len(f.Vehicles()) != 0 || len(f.Rejections()) != 1 {
		t.Error("самокат должен попасть только в журнал отказов")
	}
}

func TestDuplicateNumber(t *testing.T) {
	insp := &stubInspector{result: domain.Approved()}
	f := New(insp)
	_, _ = f.AddVehicle(scooter("VK-1"))

	_, err := f.AddVehicle(scooter("VK-1"))
	if !errors.Is(err, ErrDuplicateInventory) || insp.calls != 1 {
		t.Error("дубль номера не должен доходить до техосмотра")
	}
}
