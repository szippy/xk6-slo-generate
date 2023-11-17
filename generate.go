package slo_generate

import (
	"fmt"
	"sort"

	"go.k6.io/k6/js/modules"
)

// init is called by the Go runtime at application startup.
func init() {
	modules.Register("k6/x/slo-generate", new(SLO))
}

type (
	RootModule struct{}

	// ModuleInstance represents an instance of the JS module.
	ModuleInstance struct {
		// vu provides methods for accessing internal k6 objects for a VU
		vu modules.VU

		slo *SLO
	}
)

// Check the k6 test

// Ensure the interfaces are implemented correctly.
var (
	_ modules.Instance = &ModuleInstance{}
	_ modules.Module   = &RootModule{}
)

// New returns a pointer to a new RootModule instance.
func New() *RootModule {
	return &RootModule{}
}

// NewModuleInstance implements the modules.Module interface returning a new instance for each VU.
func (*RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return &ModuleInstance{
		vu: vu,
		// comparator: &Compare{vu: vu},
		slo: &SLO{vu: vu},
	}
}

// type result for generate.go
type SLO struct {
	vu         modules.VU
	threashold float64
	percentage string
	slo        float64
}

func (s *SLO) Generate(data []float64, percentage string) (float64, float64) {
	s.percentage = percentage

	thresholdLookup := map[string]float64{
		"99.00": 0.10,
		"99.9":  0.15,
		"99.99": 0.20,
	}

	if len(data) <= 40 {
		println(" You do not have enough data to create an SLO ")
		return -1.0, -1.0
	}

	sort.Float64s(data)

	// Step 2: Calculate percentiles
	n := len(data)
	p95Index := int(float64(n) * 0.95)
	//p99Index := int(float64(n) * 0.99)

	p95 := data[p95Index]
	//p99 := data[p99Index]

	latencyVariance, ok := thresholdLookup[percentage]
	if !ok {
		// Handle the case where the percentage is not in the lookup table
		// error out here instead
		latencyVariance = 0.0
	}

	// Step 4: Calculate the final threshold
	// 100 * .15 = 15 + 100 = 115
	//p99Threashold := p99 + (latencyVariance * p99)
	p95Threashold := p95 + (latencyVariance * p95)

	s.threashold = p95Threashold

	// Calculate what to set SLOs to based on this.
	//p99Slo := p99Threashold + (latencyVariance * p99)
	p95Slo := p95Threashold + (latencyVariance * p95)

	s.slo = p95Slo

	fmt.Printf("Based on current testing, %s%% of requests to this endpoint over a 5-minute window "+
		"should return with a valid response code and correct data within %.2fms measured at the 95th percentile. The alert "+
		"threshold should be set to %.2fms to allow for adequate time to respond. \n\n", percentage, s.slo, s.threashold)

	r, err := s.vu.metrics.NewMetric("suggestedSLO", "String", "String")
	print(r)
	print(err)

	return s.threashold, s.slo

}

//func (r *Registry) NewMetric(name string, typ MetricType, t ...ValueType) (*Metric, error)

// Exports implements the modules.Instance interface and returns the exported types for the JS module.
func (mi *ModuleInstance) Exports() modules.Exports {
	return modules.Exports{
		Default: mi.slo,
	}
}
