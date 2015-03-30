package commands

import (
	"fmt"
	"log"

	"github.com/emccode/clue"
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
	protectiondomainuseCmd.Flags().StringVar(&protectiondomainname, "protectiondomainname", "", "GOSCALEIO_PROTECTIONDOMAINNAME")
	protectiondomainuseCmd.Flags().StringVar(&protectiondomainid, "protectiondomainid", "", "GOSCALEIO_PROTECTIONDOMAINID")

	protectiondomainCmdV = protectiondomainCmd

	// initConfig(protectiondomainCmd, "goscli", true, map[string]FlagValue{
	// 	"endpoint": {endpoint, true, false, ""},
	// 	"insecure": {insecure, false, false, ""},
	// })

	protectiondomainCmd.Run = func(cmd *cobra.Command, args []string) {
		setGobValues(cmd, "goscli", "")
		cmd.Usage()
	}
}

func addCommandsProtectionDomain() {
	protectiondomainCmd.AddCommand(protectiondomaingetCmd)
	protectiondomainCmd.AddCommand(protectiondomainuseCmd)
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

var protectiondomainuseCmd = &cobra.Command{
	Use:   "use",
	Short: "Use a protectiondomain",
	Long:  `Use a protectiondomain`,
	Run:   cmdUseProtectionDomain,
}

func cmdGetProtectionDomain(cmd *cobra.Command, args []string) {
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

func cmdUseProtectionDomain(cmd *cobra.Command, args []string) {
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

	protectionDomain, err := system.FindProtectionDomain(protectiondomainid, protectiondomainname)
	if err != nil {
		log.Fatalf("error getting protection domain: %s", err)
	}

	err = clue.EncodeGobFile("goscli_system", clue.UseValue{
		VarMap: map[string]string{
			"protectiondomainid": protectionDomain.ID,
		},
	})
	if err != nil {
		log.Fatalf("error encoding gob file %v", err)
	}

	yamlOutput, err := yaml.Marshal(&protectionDomain)
	if err != nil {
		log.Fatalf("error marshaling: %s", err)
	}
	fmt.Println(string(yamlOutput))

}
