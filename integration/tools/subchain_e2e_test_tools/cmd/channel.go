package cmd

import (
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
		tools.GetMaxProcessedNonceFromRegistrar()
	},
}

func init() {
	rootCmd.AddCommand(startRegisterChannelCmd)
	rootCmd.AddCommand(startGetMaxProcessedNonceCmd)
	rootCmd.AddCommand(startGetCrossChainFeeFromRegistrarCmd)
}
