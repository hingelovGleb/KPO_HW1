package report

import (
	"fmt"

	"vyshkat/internal/domain"
)

// Правило для новичков может меняться поэтому оно за интерфейсом
type BeginnerPolicy interface {
	Allows(v domain.Vehicle) bool
	Description() string
}

// базовое правило из задания
type SimplicityPolicy struct {
	MinSimplicity int
}

func (p SimplicityPolicy) Allows(v domain.Vehicle) bool {
	return v.Simplicity() >= p.MinSimplicity
}

func (p SimplicityPolicy) Description() string {
	return fmt.Sprintf("простота от %d", p.MinSimplicity)
}

// усложнённое правило: ещё смотрим на скорость и техсостояние
type SafeBeginnerPolicy struct {
	Base         SimplicityPolicy
	MaxSpeedKmh  int
	MinCondition int
}

func (p SafeBeginnerPolicy) Allows(v domain.Vehicle) bool {
	return p.Base.Allows(v) &&
		v.MaxSpeedKmh() <= p.MaxSpeedKmh &&
		v.Condition() >= p.MinCondition
}

func (p SafeBeginnerPolicy) Description() string {
	return fmt.Sprintf("%s, скорость до %d км/ч, состояние от %d%%",
		p.Base.Description(), p.MaxSpeedKmh, p.MinCondition)
}

func DefaultBeginnerPolicy() BeginnerPolicy {
	return SafeBeginnerPolicy{
		Base:         SimplicityPolicy{MinSimplicity: 6},
		MaxSpeedKmh:  25,
		MinCondition: 80,
	}
}
