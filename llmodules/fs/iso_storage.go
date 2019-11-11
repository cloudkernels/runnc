package fs

import (
	"os"
	"path/filepath"
	"io/ioutil"
	//"fmt"
	"strings"
	"github.com/nabla-containers/runnc/libcontainer/configs"
	ll "github.com/nabla-containers/runnc/llif"
	"github.com/nabla-containers/runnc/nabla-lib/storage"
	"github.com/nabla-containers/runnc/utils"
	"github.com/pkg/errors"
)

type iSOFsHandler struct{}

func NewISOFsHandler() (ll.FsHandler, error) {
	return &iSOFsHandler{}, nil
}

func (h *iSOFsHandler) FsCreateFunc(i *ll.FsCreateInput) (*ll.LLState, error) {
	fsPath, err := createRootfsISO(i.Config, i.ContainerRoot)
	isoPaths, err := findAllISOs(i.Config)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create rootfs ISO")
	}

	ret := &ll.LLState{}
	ret.Options = map[string]string{
		"FsPath": fsPath,
		"FsStoragePath": isoPaths,
	}

	return ret, nil
}

func (h *iSOFsHandler) FsRunFunc(i *ll.FsRunInput) (*ll.LLState, error) {
	return i.FsState, nil
}

func (h *iSOFsHandler) FsDestroyFunc(i *ll.FsDestroyInput) (*ll.LLState, error) {
	if err := os.RemoveAll(i.ContainerRoot); err != nil {
		return nil, err
	}
	return i.FsState, nil
}

func findAllISOs(config *configs.Config) (string, error) {
	isoPath := config.IsoPaths
	var isos []string
	files, err := ioutil.ReadDir(isoPath)
	if err != nil {
		return "",err
	}
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".ffs" {
			isoFilePath := filepath.Join(isoPath, f.Name())
			isos = append(isos, isoFilePath)
		}
	}
	return strings.Join(isos, ","),nil
}

func createRootfsISO(config *configs.Config, containerRoot string) (string, error) {
	// edw ftiaxnetai to rootfs.iso, apo ta periexomena tou rootfs directory pou
	// exei diabasei apo to config
	// TODO: na doume pws tha erxetai apo to docker to directory me ta isos 
	// gia thn wra tha ftiaksw ena direcotry opou tha periexei kapoia isos kai 
	// tha xrhsimopoiw auta
	//rootfsPath := config.Rootfs
	//rootfsPath, _ = filepath.Abs("/tmp");
        tmp_dir, _ := ioutil.TempDir(config.Rootfs, "nabla-rootfs")
        //tmp_dir, _ := ioutil.TempDir("/tmp/", "nabla-rootfs")
        //rootfsPath := config.Rootfs
        rootfsPath := tmp_dir
	targetISOPath := filepath.Join(containerRoot, "rootfs.ffs")
	//rootfsPath = filepath.Join(rootfsPath, "/etc")
	//targetISOPath := filepath.Join(rootfsPath, "rootfs.ffs")
	//fmt.Fprintln(os.Stderr, containerRoot)
	//fmt.Fprintln(os.Stderr, rootfsPath)
	if err := os.MkdirAll(filepath.Join(rootfsPath, "/etc"), 0755); err != nil {
		return "", errors.Wrap(err, "Unable to create "+filepath.Join(rootfsPath, "/etc"))
	}
	for _, mount := range config.Mounts {
		if (mount.Destination == "/etc/resolv.conf") ||
			(mount.Destination == "/etc/hosts") ||
			(mount.Destination == "/etc/hostname") {
			dest := filepath.Join(rootfsPath, mount.Destination)
			source := mount.Source
			if err := utils.Copy(dest, source); err != nil {
				return "", errors.Wrap(err, "Unable to copy "+source+" to "+dest)
			}
		}
	}
	//_, err := storage.CreateIso(rootfsPath, &targetISOPath)
        _, err := storage.CreateFfs(rootfsPath, &targetISOPath)
	if err != nil {
		return "", errors.Wrap(err, "Error creating iso from rootfs")
	}
	return targetISOPath, nil
}
