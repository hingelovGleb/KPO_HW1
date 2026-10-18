package fleet

import (
	"errors"
	"fmt"
	"strings"

	"vyshkat/internal/domain"
)

// Интерфейс объявлен тут а не в inspection. Так fleet не знает про ServiceCenter
// и в тестах его легко подменить.
type Inspector interface {
	Inspect(v domain.Vehicle) domain.InspectionResult
}

var (
	ErrEmptyInventoryNumber = errors.New("нет инвентарного номера")
	ErrDuplicateInventory   = errors.New("такой номер уже есть")
)

type Rejection struct {
	Vehicle domain.Vehicle
	Reason  string
}

type Fleet struct {
	inspector  Inspector
	vehicles   []domain.Vehicle
	things     []domain.Thing
	rejections []Rejection
	numbers    map[string]struct{}
}

func New(inspector Inspector) *Fleet {
	return &Fleet{inspector: inspector, numbers: make(map[string]struct{})}
}

// Отказ техосмотра это нормальная ситуация поэтому он не error.
// error только если с номером что-то не так.
func (f *Fleet) AddVehicle(v domain.Vehicle) (domain.InspectionResult, error) {
	if err := f.checkNumber(v); err != nil {
		return domain.InspectionResult{}, err
	}
	result := f.inspector.Inspect(v)
	if !result.Approved {
		f.rejections = append(f.rejections, Rejection{Vehicle: v, Reason: result.Reason})
		return result, nil
	}
	f.vehicles = append(f.vehicles, v)
	f.numbers[v.InventoryNumber()] = struct{}{}
	return result, nil
}

// вещи ставим без техосмотра
func (f *Fleet) AddThing(t domain.Thing) error {
	if err := f.checkNumber(t); err != nil {
		return err
	}
	f.things = append(f.things, t)
	f.numbers[t.InventoryNumber()] = struct{}{}
	return nil
}

func (f *Fleet) checkNumber(item domain.Inventoried) error {
	num := strings.TrimSpace(item.InventoryNumber())
	if num == "" {
		return ErrEmptyInventoryNumber
	}
	if _, exists := f.numbers[num]; exists {
		return fmt.Errorf("%w (%s)", ErrDuplicateInventory, num)
	}
	return nil
}

// отдаём копии чтобы снаружи не могли поменять парк

func (f *Fleet) Vehicles() []domain.Vehicle {
	return append([]domain.Vehicle(nil), f.vehicles...)
}

func (f *Fleet) Things() []domain.Thing {
	return append([]domain.Thing(nil), f.things...)
}

func (f *Fleet) Rejections() []Rejection {
	return append([]Rejection(nil), f.rejections...)
}
