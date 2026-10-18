package inventory

import "fmt"

// Выдаёт номера по порядку: VK-1, VK-2, VK-3 и так далее
type Generator struct {
	last int
}

func (g *Generator) Next() string {
	g.last++
	return fmt.Sprintf("VK-%d", g.last)
}
