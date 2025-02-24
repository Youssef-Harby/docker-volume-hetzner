package main // import "github.com/costela/docker-volume-hetzner"

import (
	"fmt"
	"os"
	"strconv"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/sirupsen/logrus"
	"github.com/Youssef-Harby/docker-volume-hetzner/types"
)

const socketAddress = "/run/docker/plugins/hetzner.sock"
const propagatedMountPath = "/mnt"

func main() {
	// Support command-line resizing
	if len(os.Args) > 1 && os.Args[1] == "resize" {
		if len(os.Args) != 4 {
			fmt.Printf("Usage: %s resize <volumeName> <newSizeGB>\n", os.Args[0])
			os.Exit(1)
		}
		volName := os.Args[2]
		newSize, err := strconv.Atoi(os.Args[3])
		if err != nil {
			logrus.Fatalf("invalid new size: %v", err)
		}
		hd := newHetznerDriver()
		resizeReq := &types.ResizeRequest{
			Name:    volName,
			Options: map[string]string{"size": strconv.Itoa(newSize)},
		}
		if err := hd.Resize(resizeReq); err != nil {
			logrus.Fatalf("resize failed: %v", err)
		}
		os.Exit(0)
	}

	// Default behavior: start the plugin
	logrus.SetFormatter(&bareFormatter{})

	var (
		defaultDir  = volume.DefaultDockerRootDirectory
		defaultSize = 10 // GB
		defaultFS   = "ext4"
		defaultProt = true
		logLevel    = "info" // Changed from warn to info for better debugging
	)

	logLevelEnv, err := logrus.ParseLevel(os.Getenv("loglevel"))
	if err != nil {
		logLevelEnv = logrus.InfoLevel
	}

	logrus.SetLevel(logLevelEnv)

	hd := newHetznerDriver()
	h := volume.NewHandler(hd)
	logrus.Infof("listening on %s", socketAddress)
	if err := h.ServeUnix(socketAddress, 0); err != nil {
		logrus.Fatalf("error serving docker socket: %v", err)
	}
}

type bareFormatter struct{}

func (bareFormatter) Format(e *logrus.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("%s\n", e.Message)), nil
}
