package ava_services

import (
	"time"

	"github.com/kurtosis-tech/ava-e2e-tests/gecko_client"
	"github.com/kurtosis-tech/kurtosis/commons/services"
	"github.com/palantir/stacktrace"
	"github.com/sirupsen/logrus"
)

// Implements ServiceAvailabilityCheckerCore
type GeckoServiceAvailabilityCheckerCore struct{}

func (g GeckoServiceAvailabilityCheckerCore) IsServiceUp(toCheck services.Service, dependencies []services.Service) bool {
	// NOTE: we don't check the dependencies intentionally, because we don't need to - a Gecko service won't report itself
	//  as up until its bootstrappers are up

	castedService := toCheck.(GeckoService)
	jsonRpcSocket := castedService.GetJsonRpcSocket()
	client := gecko_client.NewGeckoClient(jsonRpcSocket.GetIpAddr(), jsonRpcSocket.GetPort())
	healthInfo, err := client.HealthApi().GetLiveness()
	if err != nil {
		logrus.Trace(stacktrace.Propagate(err, "Error occurred getting liveness info"))
		return false
	}

	if healthInfo.Healthy {
		time.Sleep(time.Second)
	}

	return healthInfo.Healthy
}

func (g GeckoServiceAvailabilityCheckerCore) GetTimeout() time.Duration {
	return 90 * time.Second
}
