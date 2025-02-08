package telemetry

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"sync"
)

const (
	_baseMetricName = "bompedido"
)

type metricItem struct {
	name string
	err  error
	tags map[string]string
}

var (
	metrics = make(map[string]prometheus.Counter)
	mutex   = sync.Mutex{}
	ch      = make(chan metricItem)
)

func init() {
	go incrementMetrics()
}

func IncrementMetric(name string) {
	IncrementMetricWithTags(name, nil)
}

func IncrementMetricWithTags(name string, tags map[string]string) {
	ch <- metricItem{
		name: name,
		tags: tags,
	}
}

func IncrementErrorMetric(name string, err error) {
	IncrementErrorMetricWithTags(name, err, nil)
}

func IncrementErrorMetricWithTags(name string, err error, tags map[string]string) {
	ch <- metricItem{
		name: fmt.Sprintf("%s_error", name),
		err:  err,
		tags: tags,
	}
}

func incrementMetrics() {
	for aMetric := range ch {
		mutex.Lock()
		metricID := aMetric.String()
		metric, found := metrics[metricID]

		if !found {
			labels := prometheus.Labels{}
			if aMetric.tags != nil {
				labels = aMetric.tags
			}
			if aMetric.err != nil {
				labels["error"] = aMetric.err.Error()
			}
			metric = promauto.NewCounter(prometheus.CounterOpts{
				Name:        fmt.Sprintf("%s_%s", _baseMetricName, aMetric.name),
				ConstLabels: labels,
			})
			metrics[metricID] = metric
		}
		mutex.Unlock()
		metric.Inc()
	}
}

func (m metricItem) String() string {
	if m.err != nil {
		return fmt.Sprintf("name: %s, error: %s, tags: %v", m.name, m.err, m.tags)
	}
	return fmt.Sprintf("name: %s, tags: %v", m.name, m.tags)
}
