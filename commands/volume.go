package commands

import (
	"fmt"
	"log"

	"github.com/emccode/goscaleio"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v1"
)

var volumeCmdV *cobra.Command

func init() {
	addCommandsVolume()
	// volumeCmd.Flags().StringVar(&volumename, "volumename", "", "GOSCALEIO_TEMP")
	volumeCmd.Flags().StringVar(&systemid, "systemid", "", "GOSCALEIO_SYSTEMID")
	volumegetCmd.Flags().StringVar(&systemid, "systemid", "", "GOSCALEIO_SYSTEMID")
	volumeuseCmd.Flags().StringVar(&volumename, "volumename", "", "GOSCALEIO_VOLUMENAME")
	volumeuseCmd.Flags().StringVar(&volumeid, "volumeid", "", "GOSCALEIO_VOLUMEID")

	volumeCmdV = volumeCmd

	volumeCmd.Run = func(cmd *cobra.Command, args []string) {
		setGobValues(cmd, "goscli", "")
		cmd.Usage()
	}
}

func addCommandsVolume() {
	volumeCmd.AddCommand(volumegetCmd)
	volumeCmd.AddCommand(volumeuseCmd)
	volumeCmd.AddCommand(volumelocalCmd)
}

var volumeCmd = &cobra.Command{
	Use:   "volume",
	Short: "volume",
	Long:  `volume`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

var volumegetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a volume",
	Long:  `Get a volume`,
	Run:   cmdGetVolume,
}

var volumeuseCmd = &cobra.Command{
	Use:   "use",
	Short: "Use a volume",
	Long:  `Use a volume`,
	Run:   cmdUseVolume,
}

var volumelocalCmd = &cobra.Command{
	Use:   "local",
	Short: "Get local volumes",
	Long:  `Get local volumes`,
	Run:   cmdGetVolumeLocal,
}

func cmdGetVolume(cmd *cobra.Command, args []string) {
	client, err := authenticate()
	if err != nil {
		log.Fatalf("error authenticating: %v", err)
	}

	initConfig(cmd, "goscli_system", true, map[string]FlagValue{
		"storagepoolhref": {&storagepoolhref, true, false, ""},
	})

	storagepoolhref = viper.GetString("storagepoolhref")

	protectionDomain := goscaleio.NewProtectionDomain(client)
	targetStoragePool, err := protectionDomain.FindStoragePool("", "", storagepoolhref)
	if err != nil {
		log.Fatalf("err: problem getting system %v", err)
	}

	storagePool := goscaleio.NewStoragePool(client)
	storagePool.StoragePool = targetStoragePool

	volumes, err := storagePool.GetVolume("")
	if err != nil {
		log.Fatalf("error getting volumes: %v", err)
	}

	yamlOutput, err := yaml.Marshal(&volumes)
	if err != nil {
		log.Fatalf("error marshaling: %s", err)
	}
	fmt.Println(string(yamlOutput))

}

func cmdUseVolume(cmd *cobra.Command, args []string) {
}

func cmdGetVolumeLocal(cmd *cobra.Command, args []string) {
	volumeMaps, err := goscaleio.GetLocalVolumeMap()
	if err != nil {
		log.Fatalf("Error getting local volume maps: %s", err)
	}

	yamlOutput, err := yaml.Marshal(&volumeMaps)
	if err != nil {
		log.Fatalf("error marshaling: %s", err)
	}
	fmt.Println(string(yamlOutput))

}
