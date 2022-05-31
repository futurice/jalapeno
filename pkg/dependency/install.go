package dependency

import (
	"fmt"
	"runtime"

	"github.com/futurice/sre-dx-cli-template/pkg/utils"
)

type PackageManager interface {
	InstallPackage(pkg string) error
}

func Install(pkg string) error {
	manager := selectPackageManger()
	err := manager.InstallPackage(pkg)
	return err
}

func selectPackageManger() PackageManager {
	os := runtime.GOOS
	switch os {
	case "windows":
		return &Scoop{}
	case "darwin":
		return &HomeBrew{}
	case "linux":
		return &Apt{}
	default:
		fmt.Printf("%s.\n", os)
		return nil
	}
}

type HomeBrew struct {}

func (h *HomeBrew) InstallPackage(pkg string) error {
	err := utils.RunCommand([]string{"brew", "install", pkg})
	
	return err
}

type Scoop struct {}

func (s *Scoop) InstallPackage(pkg string) error {
	err := utils.RunCommand([]string{"scoop", "install", pkg})
	
	return err
}

type Apt struct {}

func (a *Apt) InstallPackage(pkg string) error {
	err := utils.RunCommand([]string{"apt", "install", pkg})
	
	return err
}