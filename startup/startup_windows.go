//go:build windows
// +build windows

package startup

import (
	"errors"
	"fmt"
	"golang.org/x/sys/windows/registry"
	"log/slog"
	"os"
	"strings"
)

const (
	appRegistryName     = "D2Tool"
	startupRegistryPath = "Software\\Microsoft\\Windows\\CurrentVersion\\Run"
)

func StartupRegister(args []string) error {
	// Path to your application's executable
	appExecutable, err := os.Executable()
	if err != nil {
		slog.Warn(fmt.Sprintf("Error getting executable path: %v", err))
		return err
	}

	// Open "Run" registry key
	key, _, err := registry.CreateKey(registry.CURRENT_USER, startupRegistryPath, registry.ALL_ACCESS)
	if err != nil {
		slog.Warn(fmt.Sprintf("Error opening registry key: %v", err))
		return err
	}
	defer key.Close()

	// Write the path of the executable to the registry
	err = key.SetStringValue(appRegistryName, fmt.Sprintf("\"%s\" %s", appExecutable, strings.Join(args, " ")))
	if err != nil {
		slog.Warn(fmt.Sprintf("Error setting registry value: %v", err))
		return err
	}

	slog.Info("Application added to Windows Startup successfully.")
	return nil
}

func StartupRemove() error {
	// Open "Run" registry key
	key, _, err := registry.CreateKey(registry.CURRENT_USER, startupRegistryPath, registry.ALL_ACCESS)
	if err != nil {
		slog.Warn(fmt.Sprintf("Error opening registry key: %v", err))
		return err
	}
	defer key.Close()

	// Delete the value
	err = key.DeleteValue(appRegistryName)
	if err != nil {
		slog.Warn(fmt.Sprintf("Error deleting registry value: %v", err))
		return err
	}

	slog.Info("Application removed from Windows Startup successfully.")
	return nil
}

func IsStartupRegistered() (bool, error) {
	key, _, err := registry.CreateKey(registry.CURRENT_USER, startupRegistryPath, registry.ALL_ACCESS)
	if err != nil {
		slog.Warn(fmt.Sprintf("Error opening registry key: %v", err))
		return false, err
	}
	defer key.Close()

	_, _, err = key.GetStringValue(appRegistryName)
	if err != nil {
		if errors.Is(err, registry.ErrNotExist) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func SupportsStartup() bool {
	return true
}
