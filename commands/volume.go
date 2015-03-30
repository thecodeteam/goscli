package commands

import (
	"fmt"
	"log"

	"github.com/emccode/goscaleio"
	types "github.com/emccode/goscaleio/types/v1"
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
	volumegetCmd.Flags().StringVar(&volumeid, "volumeid", "", "GOSCALEIO_VOLUMEID")
	volumeuseCmd.Flags().StringVar(&volumename, "volumename", "", "GOSCALEIO_VOLUMENAME")
	volumeuseCmd.Flags().StringVar(&volumeid, "volumeid", "", "GOSCALEIO_VOLUMEID")
	volumecreateCmd.Flags().StringVar(&volumename, "volumename", "", "GOSCALEIO_VOLUMENAME")
	volumecreateCmd.Flags().StringVar(&volumeid, "volumeid", "", "GOSCALEIO_VOLUMEID")
	volumecreateCmd.Flags().StringVar(&volumeusermcache, "volumeusermcache", "", "GOSCALEIO_VOLUMEUSERMCACHE")
	volumecreateCmd.Flags().StringVar(&volumetype, "volumetype", "", "GOSCALEIO_VOLUMETYPE")
	volumecreateCmd.Flags().StringVar(&volumesizeinkb, "volumesizeinkb", "", "GOSCALEIO_VOLUMESIZEINKB")
	volumemapsdcCmd.Flags().StringVar(&volumeid, "volumeid", "", "GOSCALEIO_VOLUMEID")
	volumemapsdcCmd.Flags().StringVar(&sdcid, "sdcid", "", "GOSCALEIO_SDCID")
	volumemapsdcCmd.Flags().StringVar(&allowmultiplemappings, "allowmultiplemappings", "", "GOSCALEIO_ALLOWMULTIPLEMAPPINGS")
	volumemapsdcCmd.Flags().StringVar(&allsdcs, "allsdcs", "", "GOSCALEIO_ALLSDCS")

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
	volumeCmd.AddCommand(volumecreateCmd)
	volumeCmd.AddCommand(volumemapsdcCmd)
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

var volumecreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create volume",
	Long:  `Create volume`,
	Run:   cmdCreateVolume,
}

var volumemapsdcCmd = &cobra.Command{
	Use:   "map",
	Short: "Map volume",
	Long:  `Map volume`,
	Run:   cmdMapVolumeSdc,
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

	initConfig(cmd, "goscli_system", true, map[string]FlagValue{
		"volumeid": {&volumeid, false, false, ""},
	})

	volumes, err := storagePool.GetVolume("", volumeid)
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

func cmdCreateVolume(cmd *cobra.Command, args []string) {
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

	initConfig(cmd, "goscli_system", true, map[string]FlagValue{
		"volumename":       {&volumename, true, false, ""},
		"volumesizeinkb":   {&volumesizeinkb, true, false, ""},
		"volumetype":       {&volumetype, false, false, ""},
		"volumeusermcache": {&volumeusermcache, false, false, ""},
	})

	volume := &types.VolumeParam{
		Name:           volumename,
		VolumeSizeInKb: volumesizeinkb,
		VolumeType:     volumetype,
		UseRmCache:     volumeusermcache,
	}

	volumeResp, err := storagePool.CreateVolume(volume)
	if err != nil {
		log.Fatalf("err: problem creating volume: %s", err)
	}

	fmt.Println("Successfuly created volume with ID of", volumeResp.ID)

}

func cmdMapVolumeSdc(cmd *cobra.Command, args []string) {
	client, err := authenticate()
	if err != nil {
		log.Fatalf("error authenticating: %v", err)
	}

	initConfig(cmd, "goscli_system", true, map[string]FlagValue{
		"volumeid":              {&volumeid, true, false, ""},
		"sdcid":                 {&sdcid, true, false, ""},
		"allowmultiplemappings": {&allowmultiplemappings, false, false, ""},
		"allsdcs":               {&allsdcs, false, false, ""},
	})

	storagePool := goscaleio.NewStoragePool(client)
	targetVolumes, err := storagePool.GetVolume("", volumeid)
	if err != nil {
		log.Fatalf("error getting volume: %s", err)
	}

	volume := goscaleio.NewVolume(client)
	volume.Volume = targetVolumes[0]

	mapVolumeSdcParam := &types.MapVolumeSdcParam{
		SdcID: sdcid,
		AllowMultipleMappings: allowmultiplemappings,
		AllSdcs:               allsdcs,
	}

	err = volume.MapVolumeSdc(mapVolumeSdcParam)
	if err != nil {
		log.Fatalf("err: problem creating volume: %s", err)
	}

	fmt.Println(fmt.Sprintf("Successfuly mapped volume %s to %s", volume.Volume.ID, sdcid))

}
