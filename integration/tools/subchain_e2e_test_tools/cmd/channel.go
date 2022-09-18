package cmd

import (
	"fmt"
	"math/big"

	"github.com/spf13/cobra"
	"github.com/thetatoken/thetasubchain/integration/tools/subchain_e2e_test_tools/tools"
)

var startRegisterChannelCmd = &cobra.Command{
	Use: "RegisterChannel",
	Run: func(cmd *cobra.Command, args []string) {
		tools.SubchainChannelRegister(big.NewInt(360888), "http://localhost:19988/rpc")
	},
}

var startGetMaxProcessedNonceCmd = &cobra.Command{
	Use: "GetMaxProcessedNonceFromRegistrar",
	Run: func(cmd *cobra.Command, args []string) {
		tools.GetMaxProcessedNonceFromRegistrar()
	},
}

var startGetCrossChainFeeFromRegistrarCmd = &cobra.Command{
	Use: "GetCrossChainFeeFromRegistrar",
	Run: func(cmd *cobra.Command, args []string) {
		tools.GetCrossChainFeeFromRegistrar()
	},
}

var startCheckChannelStatusCmd = &cobra.Command{
	Use: "GetChannelStatusFromRegistrar",
	Run: func(cmd *cobra.Command, args []string) {
		// 加chainID和IP就行
		targetChainIDInt, success := big.NewInt(0).SetString(targetChainID, 10)
		if !success {
			panic(fmt.Sprintf("Failed to read amount: %v", amount))
		}
		tools.GetChannelStatusFromRegistrar(targetChainIDInt, targetChainEthRpcClientURL)
	},
}

func init() {
	rootCmd.AddCommand(startRegisterChannelCmd)
	rootCmd.AddCommand(startGetMaxProcessedNonceCmd)
	rootCmd.AddCommand(startGetCrossChainFeeFromRegistrarCmd)
	rootCmd.AddCommand(startCheckChannelStatusCmd)

	startCheckChannelStatusCmd.PersistentFlags().StringVar(&targetChainID, "targetChainID", "360888", "targetChainID")
	startCheckChannelStatusCmd.PersistentFlags().StringVar(&targetChainEthRpcClientURL, "targetChainEthRpcClientURL", "http://localhost:19988/rpc", "targetChainEthRpcClientURL")
}
