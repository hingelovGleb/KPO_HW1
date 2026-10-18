package report

import "vyshkat/internal/domain"

// Отчётам нужно только читать. Добавлять технику им не надо.
type Source interface {
	Vehicles() []domain.Vehicle
	Things() []domain.Thing
}

type Summary struct {
	Vehicles int
	Things   int
	ByKind   map[string]int
}

func (s Summary) Total() int { return s.Vehicles + s.Things }

type Energy struct {
	VehiclesKWh float64
	ThingsKWh   float64
	Consumers   int
}

func (e Energy) TotalKWh() float64 { return e.VehiclesKWh + e.ThingsKWh }

type Reporter struct {
	src    Source
	policy BeginnerPolicy
}

func NewReporter(src Source, policy BeginnerPolicy) *Reporter {
	return &Reporter{src: src, policy: policy}
}

func (r *Reporter) Summary() Summary {
	s := Summary{ByKind: make(map[string]int)}
	for _, v := range r.src.Vehicles() {
		s.Vehicles++
		s.ByKind[v.Kind()]++
	}
	for _, t := range r.src.Things() {
		s.Things++
		s.ByKind[t.Category()]++
	}
	return s
}

// считаем только тех кто реализует EnergyConsumer
func (r *Reporter) Energy() Energy {
	var e Energy
	for _, v := range r.src.Vehicles() {
		if c, ok := v.(domain.EnergyConsumer); ok {
			e.VehiclesKWh += c.DailyEnergyKWh()
			e.Consumers++
		}
	}
	for _, t := range r.src.Things() {
		if c, ok := t.(domain.EnergyConsumer); ok {
			e.ThingsKWh += c.DailyEnergyKWh()
			e.Consumers++
		}
	}
	return e
}

func (r *Reporter) BeginnerVehicles() []domain.Vehicle {
	var result []domain.Vehicle
	for _, v := range r.src.Vehicles() {
		if r.policy.Allows(v) {
			result = append(result, v)
		}
	}
	return result
}

func (r *Reporter) BeginnerRule() string { return r.policy.Description() }

// сначала вся техника потом все вещи
func (r *Reporter) Inventory() []domain.Inventoried {
	var items []domain.Inventoried
	for _, v := range r.src.Vehicles() {
		items = append(items, v)
	}
	for _, t := range r.src.Things() {
		items = append(items, t)
	}
	return items
}
