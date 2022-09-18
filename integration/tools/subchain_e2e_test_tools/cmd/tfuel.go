package cmd

import (
	"fmt"
	"math/big"

	"github.com/spf13/cobra"
	"github.com/thetatoken/thetasubchain/integration/tools/subchain_e2e_test_tools/tools"
)

var startMainchainTFuelLockCmd = &cobra.Command{
	Use: "MainchainTFuelLock",
	Run: func(cmd *cobra.Command, args []string) {
		amountInt, success := big.NewInt(0).SetString(amount, 10)
		if !success {
			panic(fmt.Sprintf("Failed to read amount: %v", amount))
		}
		targetChainIDInt, success := big.NewInt(0).SetString(targetChainID, 10)
		if !success {
			panic(fmt.Sprintf("Failed to read amount: %v", amount))
		}
		tools.MainchainTFuelLock(amountInt, targetChainIDInt, targetChainEthRpcClientURL)
	},
}

var startSubchainTFuelBurnCmd = &cobra.Command{
	Use: "SubchainTFuelBurn",
	Run: func(cmd *cobra.Command, args []string) {
		amountInt, success := big.NewInt(0).SetString(amount, 10)
		if !success {
			panic(fmt.Sprintf("Failed to read amount: %v", amount))
		}
		tools.SubchainTFuelBurn(amountInt)
	},
}

func init() {
	rootCmd.AddCommand(startMainchainTFuelLockCmd)
	rootCmd.AddCommand(startSubchainTFuelBurnCmd)
	startMainchainTFuelLockCmd.PersistentFlags().StringVar(&amount, "amount", "10", "amount")
	startSubchainTFuelBurnCmd.PersistentFlags().StringVar(&amount, "amount", "10", "amount")

	startMainchainTFuelLockCmd.PersistentFlags().StringVar(&targetChainID, "targetChainID", "360777", "targetChainID")
	startMainchainTFuelLockCmd.PersistentFlags().StringVar(&targetChainEthRpcClientURL, "targetChainEthRpcClientURL", "http://localhost:19888/rpc", "targetChainEthRpcClientURL")
}
