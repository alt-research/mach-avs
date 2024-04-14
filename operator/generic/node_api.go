package generic_operator

import (
	"context"

	"github.com/Layr-Labs/eigensdk-go/nodeapi"
)

const (
	ServiceOperator           string = "ServiceOperator"
	ServiceOperatorAggregator string = "ServiceOperatorAggregator"
	ServiceOperatorVerifier   string = "ServiceOperatorVerifier"
)

func (o *Operator) StartNodeApi() <-chan error {
	if o.config.EnableNodeApi {
		o.logger.Info("Registering node api")

		o.nodeApi.RegisterNewService(
			ServiceOperator,
			ServiceOperator,
			"operator to commit alert to aggregator",
			nodeapi.ServiceStatusInitializing,
		)

		o.nodeApi.RegisterNewService(
			ServiceOperatorAggregator,
			ServiceOperatorAggregator,
			"operator aggregator work",
			nodeapi.ServiceStatusInitializing,
		)

		o.nodeApi.RegisterNewService(
			ServiceOperatorVerifier,
			ServiceOperatorVerifier,
			"operator verifier work",
			nodeapi.ServiceStatusInitializing,
		)

		o.nodeApi.UpdateHealth(nodeapi.Healthy)
		return o.nodeApi.Start()
	}

	return make(chan error, 1)
}

func (o *Operator) StartMetrics(ctx context.Context) <-chan error {
	var metricsErrChan <-chan error
	if o.config.EnableMetrics {
		metricsErrChan = o.metrics.Start(ctx, o.metricsReg)
	} else {
		metricsErrChan = make(chan error, 1)
	}

	return metricsErrChan
}
