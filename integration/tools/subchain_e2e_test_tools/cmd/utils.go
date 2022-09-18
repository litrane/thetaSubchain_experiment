package cmd

import (
	"fmt"
	"math/big"

	"github.com/spf13/cobra"
	"github.com/thetatoken/thetasubchain/integration/tools/subchain_e2e_test_tools/tools"
)

var startOneAccountRegisterCmd = &cobra.Command{
	Use:   "RegisterSubchain",
	Short: "Start Thetasubchain node.",
	Run: func(cmd *cobra.Command, args []string) {
		subchainID, success := big.NewInt(0).SetString(subchainID, 10)
		if !success {
			panic(fmt.Sprintf("Failed to read subchainID: %v", amount))
		}
		tools.OneAccountRegister(subchainID)
	},
}

var accountID int
var validatorAddrStr string

var startOneAccountStakeCmd = &cobra.Command{
	Use:   "AccountStake",
	Short: "Start Thetasubchain node.",
	Run: func(cmd *cobra.Command, args []string) {
		subchainIDInt, success := big.NewInt(0).SetString(subchainID, 10)
		if !success {
			panic(fmt.Sprintf("Failed to read subchainID: %v", subchainIDInt))
		}
		tools.StakeToValidatorFromAccount(accountID, validatorAddrStr, subchainIDInt)
	},
}

var deploySubchainMockTokensCmd = &cobra.Command{
	Use:   "DeployTokens",
	Short: "Deploy mock tokens to both the mainchain and subchain.",
	Run: func(cmd *cobra.Command, args []string) {
		tools.DeployTokens()
	},
}

func init() {
	rootCmd.AddCommand(startOneAccountRegisterCmd)
	rootCmd.AddCommand(startOneAccountStakeCmd)
	rootCmd.AddCommand(deploySubchainMockTokensCmd)
	startOneAccountStakeCmd.PersistentFlags().IntVar(&accountID, "accountID", 1, "accountID")
	startOneAccountStakeCmd.PersistentFlags().StringVar(&validatorAddrStr, "validator", "0x2E833968E5bB786Ae419c4d13189fB081Cc43bab", "validator")
	startOneAccountStakeCmd.PersistentFlags().StringVar(&subchainID, "subchainID", "360777", "subchainID")

	startOneAccountRegisterCmd.PersistentFlags().StringVar(&subchainID, "subchainID", "360777", "subchainID")
}
