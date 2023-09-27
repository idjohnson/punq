package cmd

import (
	"fmt"

	"github.com/mogenius/punq/kubernetes"
	"github.com/mogenius/punq/operator"
	"github.com/mogenius/punq/services"

	"github.com/mogenius/punq/utils"

	"github.com/spf13/cobra"
)

var operatorCmd = &cobra.Command{
	Use:   "operator",
	Short: "Run the operator inside the cluster!",
	Long:  `Run the operator inside the cluster!`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.PrintLogo()
		println("\n###############################################")
		utils.IsNewReleaseAvailable()
		println("###############################################\n")
		utils.PrintSettings()

		contexts := services.ListContexts()
		PrintInfo(fmt.Sprintf("Initialized operator with %d contexts.", len(contexts)))
		kubernetes.ContextUpdateLocalCache(contexts)

		go operator.InitBackend()
		operator.InitFrontend()
	},
}

func init() {
	rootCmd.AddCommand(operatorCmd)
}
