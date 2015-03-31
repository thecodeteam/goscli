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
	volumeunmapsdcCmd.Flags().StringVar(&volumeid, "volumeid", "", "GOSCALEIO_VOLUMEID")
	volumeunmapsdcCmd.Flags().StringVar(&sdcid, "sdcid", "", "GOSCALEIO_SDCID")
	volumeunmapsdcCmd.Flags().StringVar(&ignorescsiinitiators, "ignoreScsiInitiators", "", "GOSCALEIO_IGNORESCSIINITIATORS")
	volumeunmapsdcCmd.Flags().StringVar(&allsdcs, "allsdcs", "", "GOSCALEIO_ALLSDCS")
	volumesnapshotCmd.Flags().StringVar(&volumeid, "volumeid", "", "GOSCALEIO_VOLUMEID")
	volumesnapshotCmd.Flags().StringVar(&snapshotname, "snapshotname", "", "GOSCALEIO_SNAPSHOTNAME")

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
	volumeCmd.AddCommand(volumeunmapsdcCmd)
	volumeCmd.AddCommand(volumesnapshotCmd)
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

var volumeunmapsdcCmd = &cobra.Command{
	Use:   "unmap",
	Short: "Unmap volume",
	Long:  `Unmap volume`,
	Run:   cmdUnmapVolumeSdc,
}

var volumesnapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Snapshot volume",
	Long:  `Snapshot volume`,
	Run:   cmdSnapshotVolume,
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

	if len(args) == 0 {
		yamlOutput, err := yaml.Marshal(&volumes)
		if err != nil {
			log.Fatalf("error marshaling: %s", err)
		}
		fmt.Println(string(yamlOutput))
		return
	}

	if len(args) > 1 {
		log.Fatalf("Too many arguments")
	}

	volume := goscaleio.NewVolume(client)
	volume.Volume = volumes[0]
	var yamlOutput []byte
	switch args[0] {
	case "vtree":
		vtree, err := volume.GetVTree()
		if err != nil {
			log.Fatalf("error getting scsiinitiators: %v", err)
		}

		yamlOutput, err = yaml.Marshal(&vtree)
	default:
		log.Fatalf("need to specify vtree as argument")
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
		log.Fatalf("err: problem mapping volume: %s", err)
	}

	fmt.Println(fmt.Sprintf("Successfuly mapped volume %s to %s", volume.Volume.ID, sdcid))

}

func cmdUnmapVolumeSdc(cmd *cobra.Command, args []string) {
	client, err := authenticate()
	if err != nil {
		log.Fatalf("error authenticating: %v", err)
	}

	initConfig(cmd, "goscli_system", true, map[string]FlagValue{
		"volumeid":             {&volumeid, true, false, ""},
		"sdcid":                {&sdcid, true, false, ""},
		"ignoreScsiInitiators": {&ignorescsiinitiators, false, false, ""},
		"allsdcs":              {&allsdcs, false, false, ""},
	})

	storagePool := goscaleio.NewStoragePool(client)
	targetVolumes, err := storagePool.GetVolume("", volumeid)
	if err != nil {
		log.Fatalf("error getting volume: %s", err)
	}

	volume := goscaleio.NewVolume(client)
	volume.Volume = targetVolumes[0]

	unmapVolumeSdcParam := &types.UnmapVolumeSdcParam{
		SdcID:                sdcid,
		IgnoreScsiInitiators: ignorescsiinitiators,
		AllSdcs:              allsdcs,
	}

	err = volume.UnmapVolumeSdc(unmapVolumeSdcParam)
	if err != nil {
		log.Fatalf("err: problem unmapping volume: %s", err)
	}

	fmt.Println(fmt.Sprintf("Successfuly unmapped volume %s to %s", volume.Volume.ID, sdcid))

}

func cmdSnapshotVolume(cmd *cobra.Command, args []string) {
	client, err := authenticate()
	if err != nil {
		log.Fatalf("error authenticating: %v", err)
	}

	initConfig(cmd, "goscli", true, map[string]FlagValue{
		"systemhref":   {&systemhref, true, false, ""},
		"volumeid":     {&volumeid, true, false, ""},
		"snapshotname": {&snapshotname, false, false, ""},
	})

	systemhref = viper.GetString("systemhref")

	system, err := client.FindSystem("", systemhref)
	if err != nil {
		log.Fatalf("err: problem getting system: %v", err)
	}

	snapshotDef := &types.SnapshotDef{
		VolumeID:     volumeid,
		SnapshotName: snapshotname,
	}

	var snapshotDefs []*types.SnapshotDef
	snapshotDefs = append(snapshotDefs, snapshotDef)
	snapshotVolumesParam := &types.SnapshotVolumesParam{
		SnapshotDefs: snapshotDefs,
	}

	snapshotVolumesResp, err := system.CreateSnapshotConsistencyGroup(snapshotVolumesParam)
	if err != nil {
		log.Fatalf("error creating snapshot: %s", err)
	}

	yamlOutput, err := yaml.Marshal(&snapshotVolumesResp)
	if err != nil {
		log.Fatalf("error marshaling: %s", err)
	}
	fmt.Println(string(yamlOutput))

}
