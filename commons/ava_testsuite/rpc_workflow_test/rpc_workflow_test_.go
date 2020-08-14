package rpc_workflow_test

import (
	"time"

	"github.com/ava-labs/avalanche-e2e-tests/commons/ava_networks"
	"github.com/ava-labs/avalanche-e2e-tests/commons/ava_services"
	"github.com/kurtosis-tech/kurtosis/commons/networks"
	"github.com/kurtosis-tech/kurtosis/commons/testsuite"
	"github.com/palantir/stacktrace"
)

const (
	stakerUsername    = "staker"
	stakerPassword    = "test34test!23"
	delegatorUsername = "delegator"
	delegatorPassword = "test34test!23"
	seedAmount        = uint64(50000000000000)
	stakeAmount       = uint64(30000000000000)
	delegatorAmount   = uint64(30000000000000)

	regularNodeServiceID   networks.ServiceID = "validator-node"
	delegatorNodeServiceID networks.ServiceID = "delegator-node"

	networkAcceptanceTimeoutRatio                          = 0.3
	normalNodeConfigID            networks.ConfigurationID = "normal-config"
)

// StakingNetworkRPCWorkflowTest ...
type StakingNetworkRPCWorkflowTest struct {
	ImageName string
}

// Run implements the Kurtosis Test interface
func (test StakingNetworkRPCWorkflowTest) Run(network networks.Network, context testsuite.TestContext) {
	// =============================== SETUP GECKO CLIENTS ======================================
	castedNetwork := network.(ava_networks.TestGeckoNetwork)
	networkAcceptanceTimeout := time.Duration(networkAcceptanceTimeoutRatio * float64(test.GetExecutionTimeout().Nanoseconds()))
	stakerClient, err := castedNetwork.GetGeckoClient(regularNodeServiceID)
	if err != nil {
		context.Fatal(stacktrace.Propagate(err, "Could not get staker client"))
	}

	delegatorClient, err := castedNetwork.GetGeckoClient(delegatorNodeServiceID)
	if err != nil {
		context.Fatal(stacktrace.Propagate(err, "Could not get delegator client"))
	}

	executor := NewRPCWorkflowTestExecutor(stakerClient, delegatorClient, networkAcceptanceTimeout)

	if err := executor.ExecuteTest(); err != nil {
		context.Fatal(stacktrace.Propagate(err, "RPCWorkflow Test failed."))
	}
}

// GetNetworkLoader implements the Kurtosis Test interface
func (test StakingNetworkRPCWorkflowTest) GetNetworkLoader() (networks.NetworkLoader, error) {
	// Define possible service configurations.
	serviceConfigs := map[networks.ConfigurationID]ava_networks.TestGeckoNetworkServiceConfig{
		normalNodeConfigID: *ava_networks.NewTestGeckoNetworkServiceConfig(true, ava_services.LOG_LEVEL_DEBUG, test.ImageName, 2, 2, make(map[string]string)),
	}
	// Define which services use which configurations.
	desiredServices := map[networks.ServiceID]networks.ConfigurationID{
		regularNodeServiceID:   normalNodeConfigID,
		delegatorNodeServiceID: normalNodeConfigID,
	}
	// Return a Gecko test net with this service:configuration mapping.
	return ava_networks.NewTestGeckoNetworkLoader(
		true,
		test.ImageName,
		ava_services.LOG_LEVEL_DEBUG,
		2,
		2,
		serviceConfigs,
		desiredServices)
}

// GetExecutionTimeout implements the Kurtosis Test interface
func (test StakingNetworkRPCWorkflowTest) GetExecutionTimeout() time.Duration {
	return 5 * time.Minute
}

// GetSetupBuffer implements the Kurtosis Test interface
func (test StakingNetworkRPCWorkflowTest) GetSetupBuffer() time.Duration {
	// TODO drop this down when the availability checker doesn't have a sleep (becuase we spin up a bunch of nodes before the test starts executing)
	return 6 * time.Minute
}
