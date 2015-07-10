package main

import (
	"log"
	"os"
	"github.com/spf13/cobra"
)

func main() {
	if os.Getuid() != 0 {
		log.Fatal("Please run this command as root")
	}

	processArgs();
}

func discoverAndRunBuildType(path string, args BuildArgs) {
 	if cnt, err := OpenCnt(path, args); err == nil {
		cnt.Build()
	} else if pod, err := OpenPod(path, args); err == nil {
		pod.Build()
	} else {
		log.Fatal("Cannot found cnt-manifest.yml")
	}
}

func discoverAndRunPushType(path string, args BuildArgs) {
	if cnt, err := OpenCnt(path, args); err == nil {
		cnt.Push()
	} else if _, err := OpenPod(path, args); err == nil {
//		pod.Build()
	} else {
		log.Fatal("Cannot found cnt-manifest.yml")
	}
}

type BuildArgs struct {
	Zip bool
}

func processArgs() {
	buildArgs := BuildArgs{}

	var cmdBuild = &cobra.Command{
		Use:   "build",
		Short: "build aci or pod",
		Long:  `build an aci or a pod`,
		Run: func(cmd *cobra.Command, args []string) {
			discoverAndRunBuildType(".", buildArgs)
		},
	}
	cmdBuild.Flags().BoolVarP(&buildArgs.Zip, "nozip", "z", true, "zip final image or not")

	var cmdClean = &cobra.Command{
		Use:   "clean",
		Short: "clean build",
		Long:  `clean build, including rootfs`,
		Run: func(cmd *cobra.Command, args []string) {
			os.RemoveAll("target/");
		},
	}

	var push = &cobra.Command{
		Use:   "push",
		Short: "push image(s)",
		Long:  `push images to repository`,
		Run: func(cmd *cobra.Command, args []string) {
			discoverAndRunPushType(".", buildArgs)
		},
	}

	var rootCmd = &cobra.Command{Use: "cnt"}
	rootCmd.AddCommand(cmdBuild, cmdClean, push)

	cntConfig.Load()
	rootCmd.Execute()

	println("Victory !")
}
