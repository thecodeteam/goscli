package commands

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v1"
)

var scsiinitiatorCmdV *cobra.Command

func init() {
	addCommandsScsiInitiator()
	// scsiinitiatorCmd.Flags().StringVar(&scsiinitiatorname, "scsiinitiatorname", "", "GOSCALEIO_TEMP")
	scsiinitiatorCmd.Flags().StringVar(&systemid, "systemid", "", "GOSCALEIO_SYSTEMID")
	scsiinitiatorgetCmd.Flags().StringVar(&systemid, "systemid", "", "GOSCALEIO_SYSTEMID")

	scsiinitiatorCmdV = scsiinitiatorCmd

	// initConfig(scsiinitiatorCmd, "goscli", true, map[string]FlagValue{
	// 	"endpoint": {endpoint, true, false, ""},
	// 	"insecure": {insecure, false, false, ""},
	// })

	scsiinitiatorCmd.Run = func(cmd *cobra.Command, args []string) {
		setGobValues(cmd, "goscli", "")
		cmd.Usage()
	}
}

func addCommandsScsiInitiator() {
	scsiinitiatorCmd.AddCommand(scsiinitiatorgetCmd)
}

var scsiinitiatorCmd = &cobra.Command{
	Use:   "scsiinitiator",
	Short: "scsiinitiator",
	Long:  `scsiinitiator`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

var scsiinitiatorgetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a scsiinitiator",
	Long:  `Get a scsiinitiator`,
	Run:   cmdGetScsiInitiator,
}

func cmdGetScsiInitiator(cmd *cobra.Command, args []string) {
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

	scsiinitiators, err := system.GetScsiInitiator()
	if err != nil {
		log.Fatalf("error getting statistics: %v", err)
	}

	yamlOutput, err := yaml.Marshal(&scsiinitiators)
	if err != nil {
		log.Fatalf("error marshaling: %s", err)
	}
	fmt.Println(string(yamlOutput))

}
