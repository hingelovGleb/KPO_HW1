package inspection

import (
	"fmt"
	"strings"

	"vyshkat/internal/domain"
)

type Config struct {
	MinCondition int
	MaxDailyKWh  float64
	MaxSpeedKmh  int
}

func DefaultConfig() Config {
	return Config{MinCondition: 70, MaxDailyKWh: 5, MaxSpeedKmh: 45}
}

type ServiceCenter struct {
	cfg Config
}

func NewServiceCenter(cfg Config) *ServiceCenter {
	return &ServiceCenter{cfg: cfg}
}

// Собираем сразу все причины отказа а не только первую
func (s *ServiceCenter) Inspect(v domain.Vehicle) domain.InspectionResult {
	var problems []string

	if v.Simplicity() < 1 || v.Simplicity() > 10 {
		problems = append(problems, "простота должна быть от 1 до 10")
	}
	if v.Condition() < s.cfg.MinCondition {
		problems = append(problems, fmt.Sprintf("плохое состояние %d%% (нужно от %d%%)",
			v.Condition(), s.cfg.MinCondition))
	}
	if v.MaxSpeedKmh() <= 0 || v.MaxSpeedKmh() > s.cfg.MaxSpeedKmh {
		problems = append(problems, fmt.Sprintf("скорость %d км/ч (можно до %d)",
			v.MaxSpeedKmh(), s.cfg.MaxSpeedKmh))
	}
	if ec, ok := v.(domain.EnergyConsumer); ok {
		kwh := ec.DailyEnergyKWh()
		if kwh <= 0 || kwh > s.cfg.MaxDailyKWh {
			problems = append(problems, fmt.Sprintf("расход %.1f кВт·ч (можно больше 0 и до %.0f)",
				kwh, s.cfg.MaxDailyKWh))
		}
	}

	if len(problems) > 0 {
		return domain.Rejected(strings.Join(problems, "; "))
	}
	return domain.Approved()
}
