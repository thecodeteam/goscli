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

var storagepoolCmdV *cobra.Command

func init() {
	addCommandsStoragePool()
	// storagepoolCmd.Flags().StringVar(&storagepoolname, "storagepoolname", "", "GOSCALEIO_TEMP")
	storagepoolCmd.Flags().StringVar(&systemid, "systemid", "", "GOSCALEIO_SYSTEMID")
	storagepoolgetCmd.Flags().StringVar(&systemid, "systemid", "", "GOSCALEIO_SYSTEMID")
	storagepooluseCmd.Flags().StringVar(&storagepoolname, "storagepoolname", "", "GOSCALEIO_STORAGEPOOLNAME")
	storagepooluseCmd.Flags().StringVar(&storagepoolid, "storagepoolid", "", "GOSCALEIO_STORAGEPOOLID")

	storagepoolCmdV = storagepoolCmd

	storagepoolCmd.Run = func(cmd *cobra.Command, args []string) {
		setGobValues(cmd, "goscli", "")
		cmd.Usage()
	}
}

func addCommandsStoragePool() {
	storagepoolCmd.AddCommand(storagepoolgetCmd)
	storagepoolCmd.AddCommand(storagepooluseCmd)
}

var storagepoolCmd = &cobra.Command{
	Use:   "storagepool",
	Short: "storagepool",
	Long:  `storagepool`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

var storagepoolgetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a storagepool",
	Long:  `Get a storagepool`,
	Run:   cmdGetStoragePool,
}

var storagepooluseCmd = &cobra.Command{
	Use:   "use",
	Short: "Use a storagepool",
	Long:  `Use a storagepool`,
	Run:   cmdUseStoragePool,
}

func cmdGetStoragePool(cmd *cobra.Command, args []string) {
	client, err := authenticate()
	if err != nil {
		log.Fatalf("error authenticating: %v", err)
	}

	initConfig(cmd, "goscli_system", true, map[string]FlagValue{
		"protectiondomainhref": {&protectiondomainhref, true, false, ""},
	})

	protectiondomainhref = viper.GetString("protectiondomainhref")

	system := goscaleio.NewSystem(client)
	targetProtectionDomain, err := system.FindProtectionDomain("", "", protectiondomainhref)
	if err != nil {
		log.Fatalf("err: problem getting system %v", err)
	}

	protectionDomain := goscaleio.NewProtectionDomain(client)
	protectionDomain.ProtectionDomain = targetProtectionDomain

	storagePools, err := protectionDomain.GetStoragePool("")
	if err != nil {
		log.Fatalf("error getting protection domains: %v", err)
	}

	yamlOutput, err := yaml.Marshal(&storagePools)
	if err != nil {
		log.Fatalf("error marshaling: %s", err)
	}
	fmt.Println(string(yamlOutput))

}

func cmdUseStoragePool(cmd *cobra.Command, args []string) {
	client, err := authenticate()
	if err != nil {
		log.Fatalf("error authenticating: %v", err)
	}

	initConfig(cmd, "goscli_system", true, map[string]FlagValue{
		"protectiondomainhref": {&protectiondomainhref, true, false, ""},
	})

	protectiondomainhref = viper.GetString("protectiondomainhref")

	system := goscaleio.NewSystem(client)
	targetProtectionDomain, err := system.FindProtectionDomain("", "", protectiondomainhref)
	if err != nil {
		log.Fatalf("error getting protection domain: %s", err)
	}

	protectionDomain := goscaleio.NewProtectionDomain(client)
	protectionDomain.ProtectionDomain = targetProtectionDomain
	storagePool, err := protectionDomain.FindStoragePool(storagepoolid, storagepoolname, "")
	if err != nil {
		log.Fatalf("error getting storage pool: %s", err)
	}

	link, err := goscaleio.GetLink(storagePool.Links, "self")
	if err != nil {
		log.Fatalf("Err: problem getting self link")
	}

	storagepoolhref = link.HREF

	err = clue.EncodeGobFile("goscli_system", clue.UseValue{
		VarMap: map[string]string{
			"protectiondomainid":   protectionDomain.ProtectionDomain.ID,
			"protectiondomainhref": protectiondomainhref,
			"protectiondomainname": protectionDomain.ProtectionDomain.Name,
			"storagepoolid":        storagePool.ID,
			"storagepoolhref":      storagepoolhref,
			"storagepoolname":      storagePool.Name,
		},
	})
	if err != nil {
		log.Fatalf("error encoding gob file %v", err)
	}

	yamlOutput, err := yaml.Marshal(&storagePool)
	if err != nil {
		log.Fatalf("error marshaling: %s", err)
	}
	fmt.Println(string(yamlOutput))

}
