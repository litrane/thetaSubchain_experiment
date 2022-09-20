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
		targetChainIDInt, success := big.NewInt(0).SetString(targetChainIDForChannelRegister, 10)
		if !success {
			panic(fmt.Sprintf("Failed to read amount: %v", amount))
		}
		tools.SubchainChannelRegister(targetChainIDInt, targetChainEthRpcClientURLForChannelRegister, sourceChainEthRpcClientURLForChannelRegister)
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
		targetChainIDInt, success := big.NewInt(0).SetString(targetChainID, 10)
		if !success {
			panic(fmt.Sprintf("Failed to read amount: %v", amount))
		}
		tools.GetChannelStatusFromRegistrar(targetChainIDInt, targetChainEthRpcClientURL)
	},
}

var startVerifyChannelCmd = &cobra.Command{
	Use: "VerifyChannel",
	Run: func(cmd *cobra.Command, args []string) {
		// 加chainID和IP就行
		targetChainIDInt, success := big.NewInt(0).SetString(targetChainIDForVerify, 10)
		if !success {
			panic(fmt.Sprintf("Failed to read amount: %v", amount))
		}
		tools.VerifyChannel(targetChainIDInt, targetChainEthRpcClientURLForVerify)
	},
}

func init() {
	rootCmd.AddCommand(startRegisterChannelCmd)
	rootCmd.AddCommand(startGetMaxProcessedNonceCmd)
	rootCmd.AddCommand(startGetCrossChainFeeFromRegistrarCmd)
	rootCmd.AddCommand(startCheckChannelStatusCmd)
	rootCmd.AddCommand(startVerifyChannelCmd)

	startCheckChannelStatusCmd.PersistentFlags().StringVar(&targetChainID, "tid", "360888", "targetChainID")
	startCheckChannelStatusCmd.PersistentFlags().StringVar(&targetChainEthRpcClientURL, "turl", "http://localhost:19988/rpc", "targetChainEthRpcClientURL")

	startVerifyChannelCmd.PersistentFlags().StringVar(&targetChainIDForVerify, "tid", "360888", "targetChainID")
	startVerifyChannelCmd.PersistentFlags().StringVar(&targetChainEthRpcClientURLForVerify, "turl", "http://localhost:19988/rpc", "targetChainEthRpcClientURL")

	startRegisterChannelCmd.PersistentFlags().StringVar(&targetChainIDForChannelRegister, "tid", "360888", "targetChainID")
	startRegisterChannelCmd.PersistentFlags().StringVar(&targetChainEthRpcClientURLForChannelRegister, "turl", "http://localhost:19988/rpc", "targetChainEthRpcClientURL")
	startRegisterChannelCmd.PersistentFlags().StringVar(&sourceChainEthRpcClientURLForChannelRegister, "surl", "http://localhost:19888/rpc", "sourceChainEthRpcClientURL")

}
