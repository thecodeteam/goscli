package commands

import (
	"fmt"
	"log"

	"github.com/emccode/clue"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v1"
)

var systemCmdV *cobra.Command

func init() {
	addCommandsSystem()
	// systemCmd.Flags().StringVar(&systemname, "systemname", "", "GOSCALEIO_TEMP")
	systemgetCmd.Flags().StringVar(&systemid, "systemid", "", "GOSCALEIO_SYSTEMID")
	systemuseCmd.Flags().StringVar(&systemid, "systemid", "", "GOSCALEIO_SYSTEMID")

	systemCmdV = systemCmd

	initConfig(systemCmd, "goscli", true, map[string]FlagValue{
		"endpoint": {endpoint, true, false, ""},
		"insecure": {insecure, false, false, ""},
	})

	systemCmd.Run = func(cmd *cobra.Command, args []string) {
		setGobValues(cmd, "goscli", "")
		cmd.Usage()
	}
}

func addCommandsSystem() {
	systemCmd.AddCommand(systemgetCmd)
	systemCmd.AddCommand(systemuseCmd)
}

var systemCmd = &cobra.Command{
	Use:   "system",
	Short: "system",
	Long:  `system`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

var systemgetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a system",
	Long:  `Get a system`,
	Run:   cmdGetSystem,
}

var systemuseCmd = &cobra.Command{
	Use:   "use",
	Short: "Use a system",
	Long:  `Use a system`,
	Run:   cmdUseSystem,
}

func cmdGetSystem(cmd *cobra.Command, args []string) {
	client, err := authenticate()
	if err != nil {
		log.Fatalf("error authenticating: %v", err)
	}

	initConfig(cmd, "goscli_system", true, map[string]FlagValue{
		"systemid": {systemid, true, false, ""},
	})

	systemid = viper.GetString("systemid")

	system, err := client.FindSystem(systemid)
	if err != nil {
		log.Fatalf("err: problem getting system %v", err)
	}

	if len(args) == 0 {

		yamlOutput, err := yaml.Marshal(&system)
		if err != nil {
			log.Fatalf("error marshaling: %s", err)
		}
		fmt.Println(string(yamlOutput))
		return

	}

	if len(args) > 1 {
		log.Fatalf("too many arguments specified")
	}

	var yamlOutput []byte
	switch args[0] {
	case "statistics":
		statistics, err := system.GetStatistics()
		if err != nil {
			log.Fatalf("error getting statistics: %v", err)
		}

		yamlOutput, err = yaml.Marshal(&statistics)
	case "user":
		users, err := system.GetUser()
		if err != nil {
			log.Fatalf("error getting users: %v", err)
		}

		yamlOutput, err = yaml.Marshal(&users)
	case "scsiinitiator":
		scsiInitiators, err := system.GetScsiInitiator()
		if err != nil {
			log.Fatalf("error getting scsiinitiators: %v", err)
		}

		yamlOutput, err = yaml.Marshal(&scsiInitiators)
	case "protectiondomain":
		protectionDomains, err := system.GetProtectionDomain()
		if err != nil {
			log.Fatalf("error getting protectiondomains: %v", err)
		}

		yamlOutput, err = yaml.Marshal(&protectionDomains)
	case "sdc":
		sdcs, err := system.GetSdc()
		if err != nil {
			log.Fatalf("error getting sdcs: %v", err)
		}

		yamlOutput, err = yaml.Marshal(&sdcs)
	default:
		log.Fatalf("err: must specify statistics|user|scsiinitiator|protectiondomain|sdc")
	}

	if err != nil {
		log.Fatalf("error marshaling: %s", err)
	}
	fmt.Println(string(yamlOutput))

}

func cmdUseSystem(cmd *cobra.Command, args []string) {
	client, err := authenticate()
	if err != nil {
		log.Fatalf("error authenticating: %v", err)
	}

	system, err := client.FindSystem(systemid)
	if err != nil {
		log.Fatalf("error getting system: %s", err)
	}

	err = clue.EncodeGobFile("goscli_system", clue.UseValue{
		VarMap: map[string]string{
			"systemid": systemid,
		},
	})

	if err != nil {
		log.Fatalf("error encoding gob file %v", err)
	}

	yamlOutput, err := yaml.Marshal(&system)
	if err != nil {
		log.Fatalf("error marshaling: %s", err)
	}
	fmt.Println(string(yamlOutput))

}
