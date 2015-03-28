package commands

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v1"
)

var protectiondomainCmdV *cobra.Command

func init() {
	addCommandsProtectionDomain()
	// protectiondomainCmd.Flags().StringVar(&protectiondomainname, "protectiondomainname", "", "GOSCALEIO_TEMP")
	protectiondomainCmd.Flags().StringVar(&systemid, "systemid", "", "GOSCALEIO_SYSTEMID")
	protectiondomaingetCmd.Flags().StringVar(&systemid, "systemid", "", "GOSCALEIO_SYSTEMID")

	protectiondomainCmdV = protectiondomainCmd

	initConfig(protectiondomainCmd, "goscli", true, map[string]FlagValue{
		"endpoint": {endpoint, true, false, ""},
		"insecure": {insecure, false, false, ""},
	})

	protectiondomainCmd.Run = func(cmd *cobra.Command, args []string) {
		setGobValues(cmd, "goscli", "")
		cmd.Usage()
	}
}

func addCommandsProtectionDomain() {
	protectiondomainCmd.AddCommand(protectiondomaingetCmd)
}

var protectiondomainCmd = &cobra.Command{
	Use:   "protectiondomain",
	Short: "protectiondomain",
	Long:  `protectiondomain`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

var protectiondomaingetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a protectiondomain",
	Long:  `Get a protectiondomain`,
	Run:   cmdGetProtectionDomain,
}

func cmdGetProtectionDomain(cmd *cobra.Command, args []string) {
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

	protectiondomains, err := system.GetProtectionDomain()
	if err != nil {
		log.Fatalf("error getting protection domains: %v", err)
	}

	yamlOutput, err := yaml.Marshal(&protectiondomains)
	if err != nil {
		log.Fatalf("error marshaling: %s", err)
	}
	fmt.Println(string(yamlOutput))

}
