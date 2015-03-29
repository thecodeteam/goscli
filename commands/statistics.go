package commands

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v1"
)

var statisticsCmdV *cobra.Command

func init() {
	addCommandsStatistics()
	// statisticsCmd.Flags().StringVar(&statisticsname, "statisticsname", "", "GOSCALEIO_TEMP")
	statisticsCmd.Flags().StringVar(&systemid, "systemid", "", "GOSCALEIO_SYSTEMID")
	statisticsgetCmd.Flags().StringVar(&systemid, "systemid", "", "GOSCALEIO_SYSTEMID")

	statisticsCmdV = statisticsCmd

	// initConfig(statisticsCmd, "goscli", true, map[string]FlagValue{
	// 	"endpoint": {endpoint, true, false, ""},
	// 	"insecure": {insecure, false, false, ""},
	// })

	statisticsCmd.Run = func(cmd *cobra.Command, args []string) {
		setGobValues(cmd, "goscli", "")
		cmd.Usage()
	}
}

func addCommandsStatistics() {
	statisticsCmd.AddCommand(statisticsgetCmd)
}

var statisticsCmd = &cobra.Command{
	Use:   "statistics",
	Short: "statistics",
	Long:  `statistics`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

var statisticsgetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a statistics",
	Long:  `Get a statistics`,
	Run:   cmdGetStatistics,
}

func cmdGetStatistics(cmd *cobra.Command, args []string) {
	client, err := authenticate()
	if err != nil {
		log.Fatalf("error authenticating: %v", err)
	}

	initConfig(cmd, "goscli", true, map[string]FlagValue{
		"systemid": {&systemid, true, false, ""},
	})

	systemid = viper.GetString("systemid")

	system, err := client.FindSystem(systemid, "")
	if err != nil {
		log.Fatalf("err: problem getting system %v", err)
	}

	statisticss, err := system.GetStatistics()
	if err != nil {
		log.Fatalf("error getting statistics: %v", err)
	}

	yamlOutput, err := yaml.Marshal(&statisticss)
	if err != nil {
		log.Fatalf("error marshaling: %s", err)
	}
	fmt.Println(string(yamlOutput))

}
