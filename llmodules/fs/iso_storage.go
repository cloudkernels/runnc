package fs

import (
	"os"
	"path/filepath"
	"fmt"
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
	FsStorPaths, err := createRootfsISO(i.Config, i.ContainerRoot)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create rootfs ISO")
	}

	ret := &ll.LLState{}
	ret.Options = map[string]string{
		"FsPath": FsStorPaths[0],
		"StorPath": FsStorPaths[1],
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

func createRootfsISO(config *configs.Config, containerRoot string) ([]string, error) {
	rootfsPath := config.Rootfs
	targetISOPath := filepath.Join(containerRoot, "rootfs.iso")
	storagePath := filepath.Join(containerRoot, "storage.ext2")
	if err := os.MkdirAll(filepath.Join(rootfsPath, "/etc"), 0755); err != nil {
		return []string{"",""}, errors.Wrap(err, "Unable to create "+filepath.Join(rootfsPath, "/etc"))
	}
	//if err := os.MkdirAll(filepath.Join(rootfsPath, "/storage"), 0755); err != nil {
	//	return []string{"",""}, errors.Wrap(err, "Unable to create "+filepath.Join(rootfsPath, "/etc"))
	//}
	fmt.Println(rootfsPath, targetISOPath)
	for _, mount := range config.Mounts {
		fmt.Println("mount.Source: ", mount.Source, "Mount.Destination: ", mount.Destination)
		if (mount.Type == "bind") {
			dest := filepath.Join(rootfsPath, mount.Destination)
			source := mount.Source
			fmt.Println("source: ", source, "dest: ", dest)
			if err := utils.Copy(dest, source); err != nil {
				return []string{"skata2","skata2"}, errors.Wrap(err, "2.Unable to copy "+source+" to "+dest)
			}
			if (mount.Destination == "/storage.ext2") {
				//storagePath = dest
				if err := utils.Copy(storagePath, dest); err != nil {
					return []string{"skata","skata"}, errors.Wrap(err, "1.Unable to copy "+dest+" to "+storagePath)
				}
			}
		}
	}
	_, err := storage.CreateIso(rootfsPath, &targetISOPath)
	if err != nil {
		return []string{"skata3","skata3"}, errors.Wrap(err, "Error creating iso from rootfs")
	}
	FsStorPaths := []string{rootfsPath, storagePath}
	return FsStorPaths, nil
}
