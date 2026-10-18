package app

import (
	"io"

	"vyshkat/internal/di"
	"vyshkat/internal/domain"
	"vyshkat/internal/fleet"
	"vyshkat/internal/inspection"
	"vyshkat/internal/inventory"
	"vyshkat/internal/report"
)

// Тут собираются все зависимости. Больше нигде контейнер не используется.
func BuildContainer() *di.Container {
	c := di.New()

	di.RegisterInstance(c, inspection.DefaultConfig())
	di.Register(c, func(c *di.Container) fleet.Inspector {
		return inspection.NewServiceCenter(di.Resolve[inspection.Config](c))
	})
	di.Register(c, func(c *di.Container) *fleet.Fleet {
		return fleet.New(di.Resolve[fleet.Inspector](c))
	})

	di.RegisterInstance(c, report.DefaultBeginnerPolicy())
	di.Register(c, func(c *di.Container) *report.Reporter {
		return report.NewReporter(di.Resolve[*fleet.Fleet](c), di.Resolve[report.BeginnerPolicy](c))
	})

	di.RegisterInstance(c, domain.DefaultCatalog())
	di.RegisterInstance(c, &inventory.Generator{})
	return c
}

func NewConsoleFromContainer(c *di.Container, in io.Reader, out io.Writer) *Console {
	return NewConsole(
		di.Resolve[*fleet.Fleet](c),
		di.Resolve[*report.Reporter](c),
		di.Resolve[*domain.VehicleCatalog](c),
		di.Resolve[*inventory.Generator](c),
		in, out,
	)
}
