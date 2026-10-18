package domain

import "testing"

// новый вид добавляется без правки старого кода
func TestCatalogNewKind(t *testing.T) {
	cat := DefaultCatalog()
	cat.Register(VehicleKind{Title: "Моноколесо", Electric: true,
		Create: func(s VehicleSpec, kwh float64) Vehicle { return NewScooter(s, kwh) }})

	if len(cat.Kinds()) != 4 || cat.Kinds()[3].Title != "Моноколесо" {
		t.Error("новый вид не добавился")
	}
}
