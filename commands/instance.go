package commands

import (
	"fmt"
	"log"

	"github.com/emccode/clue"
	"github.com/emccode/goscaleio"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v1"
)

var instanceCmdV *cobra.Command

func init() {
	addCommandsInstance()
	// instanceCmd.Flags().StringVar(&instancename, "instancename", "", "GOSCALEIO_TEMP")
	instancegetCmd.Flags().StringVar(&systemid, "systemid", "", "GOSCALEIO_SYSTEMID")

	instanceCmdV = instanceCmd

	// initConfig(instanceCmd, "goscli", true, map[string]FlagValue{
	// 	"endpoint": {endpoint, true, false, ""},
	// 	"insecure": {insecure, false, false, ""},
	// })

	instanceCmd.Run = func(cmd *cobra.Command, args []string) {
		setGobValues(cmd, "goscli", "")
		cmd.Usage()
	}
}

func addCommandsInstance() {
	instanceCmd.AddCommand(instancegetCmd)
}

var instanceCmd = &cobra.Command{
	Use:   "instance",
	Short: "instance",
	Long:  `instance`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

var instancegetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a instance",
	Long:  `Get a instance`,
	Run:   cmdGetInstance,
}

func cmdGetInstance(cmd *cobra.Command, args []string) {

	getValue := clue.GetValue{}
	if err := clue.DecodeGobFile("goscli", &getValue); err != nil {
		log.Fatalf("Problem with client DecodeGobFile: %v", err)
	}

	client, err := goscaleio.NewClient()
	if err != nil {
		log.Fatalf("error with NewClient: %s", err)
	}

	client.Token = *getValue.VarMap["token"]

	if len(args) == 0 {
		if systemid == "" {
			systems, err := client.GetInstance("")
			if err != nil {
				log.Fatalf("err: problem getting instance %v", err)
			}

			yamlOutput, err := yaml.Marshal(&systems)
			if err != nil {
				log.Fatalf("error marshaling: %s", err)
			}
			fmt.Println(string(yamlOutput))
			return

		}

		system, err := client.FindSystem(systemid, "")
		if err != nil {
			log.Fatalf("err: problem getting instance %v", err)
		}

		yamlOutput, err := yaml.Marshal(&system)
		if err != nil {
			log.Fatalf("error marshaling: %s", err)
		}
		fmt.Println(string(yamlOutput))
		return

	}

	if systemid == "" {
		log.Fatalf("missing systemid when doing detailed get")
	}

	if len(args) > 1 {
		log.Fatalf("too many arguments specified")
	}

	system, err := client.FindSystem(systemid, "")
	if err != nil {
		log.Fatalf("err: problem getting instance %v", err)
	}

	switch args[0] {
	case "statistics":
	case "user":
	case "scsiinitiator":
	case "protectiondomain":
	case "sdc":
	default:
		log.Fatalf("err: must specify statistics|user|scsiinitiator|protectiondomain|sdc")
	}

	var href string
	for _, link := range system.System.Links {
		switch {
		case args[0] == "statistics" && link.Rel == "/api/System/relationship/Statistics":
			href = link.HREF
			break
		case args[0] == "user" && link.Rel == "/api/System/relationship/User":
			href = link.HREF
			break
		case args[0] == "scsiinitiator" && link.Rel == "/api/System/relationship/ScsiInitiator":
			href = link.HREF
			break
		case args[0] == "protectiondomain" && link.Rel == "/api/System/relationship/ProtectionDomain":
			href = link.HREF
			break
		case args[0] == "sdc" && link.Rel == "/api/System/relationship/Sdc":
			href = link.HREF
			break
		}
	}

	if href == "" {
		log.Fatalf("couldn't find link")
	}

	fmt.Println(href)

}
