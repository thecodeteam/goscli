package commands

import (
	"fmt"
	"log"

	"github.com/emccode/clue"
	"github.com/emccode/goscaleio"
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
	// initConfig(systemCmd, "goscli", true, map[string]FlagValue{
	// 	"endpoint": {endpoint, true, false, ""},
	// 	"insecure": {insecure, false, false, ""},
	// })

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

	initConfig(cmd, "goscli", true, map[string]FlagValue{
		"systemhref": {&systemhref, true, false, ""},
	})

	systemhref = viper.GetString("systemhref")

	system, err := client.FindSystem("", systemhref)
	if err != nil {
		log.Fatalf("err: problem getting system: %v", err)
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
		protectionDomains, err := system.GetProtectionDomain("")
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

	system, err := client.FindSystem(systemid, "")
	if err != nil {
		log.Fatalf("error getting system: %s", err)
	}

	link, err := goscaleio.GetLink(system.System.Links, "self")
	if err != nil {
		log.Fatalf("Err: problem getting self link")
	}

	systemhref = link.HREF

	err = clue.EncodeGobFile("goscli", clue.UseValue{
		VarMap: map[string]string{
			"token":      client.Token,
			"endpoint":   client.SIOEndpoint.String(),
			"insecure":   viper.GetString("insecure"),
			"systemid":   systemid,
			"systemhref": systemhref,
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
