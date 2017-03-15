package queues

import (
	"fmt"
	"github.com/kubilus1/rabbitmqbeat/module/rabbitmq"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/metricbeat/mb"
	"github.com/pkg/errors"
	"time"
)

const (
	beatType = "rabbitmq.queues"
)

// init registers the MetricSet with the central registry.
// The New method will be called after the setup of the module and before starting to fetch data
func init() {
	if err := mb.Registry.AddMetricSet("rabbitmq", "queues", New); err != nil {
		panic(err)
	}
}

// MetricSet type defines all fields of the MetricSet
// As a minimum it must inherit the mb.BaseMetricSet fields, but can be extended with
// additional entries. These variables can be used to persist data or configuration between
// multiple fetch calls.
type MetricSet struct {
	mb.BaseMetricSet
	api rabbitmq.Api
}

// New create a new instance of the MetricSet
// Part of new is also setting up the configuration by processing additional
// configuration entries if needed.
func New(base mb.BaseMetricSet) (mb.MetricSet, error) {

	// Unpack additional configuration options.
	config := struct {
		Hosts    []string `config:"hosts" validate:"nonzero,required"`
		Username string   `config:"username"`
		Password string   `config:"password"`
	}{
		Username: "",
		Password: "",
	}
	err := base.Module().UnpackConfig(&config)
	if err != nil {
		return nil, err
	}

	client := rabbitmq.NewClient(base.Host(), config.Username, config.Password)

	return &MetricSet{
		BaseMetricSet: base,
		api:           client,
	}, nil
}

// Fetch methods implements the data gathering and data conversion to the right format
// It returns the event which is then forward to the output. In case of an error, a
// descriptive error must be returned.
func (m *MetricSet) Fetch() ([]common.MapStr, error) {
	queues, err := m.api.Queues()
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("%v api fetch failed.", beatType))
	}
	for _, queue := range queues {
		queue["@timestamp"] = common.Time(time.Now())
		queue["type"] = beatType
		queue["backing_queue_status"] = nil
	}
	return queues, nil
}
