package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"os"
	"path"

	"golang.org/x/exp/slices"
)

//go:embed init.sh
var initSh []byte

//go:embed source_init.txt
var initIntroduction []byte

//go:embed usage.txt
var UsageText []byte

func printUsage() {
	fmt.Println(string(UsageText))
}

func main() {
	if len(os.Args) == 1 {
		printUsage()
		os.Exit(1)
	}

	if os.Args[1] == "init" {
		if len(os.Args) == 3 && os.Args[2] == "-" {
			fmt.Print(string(initSh))
			return
		}
		fmt.Print(string(initIntroduction))
		return
	}

	installedVersions := getInstalledVersions()

	if os.Args[1] == "global" {
		if len(os.Args) == 2 {
			globalVersion, ok := getGlobalVersion()
			if !ok {
				fmt.Println("Global version is not set")
				return
			}
			if slices.Contains(installedVersions, globalVersion) {
				fmt.Println(globalVersion)
				return
			} else {
				fmt.Println("Error: Global version " + globalVersion + " not installed")
				fmt.Println("Please try specifiying a version with 'xienv global <version>'")
				fmt.Println("Valid versions are:")
				for _, version := range installedVersions {
					fmt.Println("  " + version)
				}
			}
		}
		if len(os.Args) != 3 {
			printUsage()
			os.Exit(1)
		}

		newVersion := os.Args[2]
		if !slices.Contains(installedVersions, newVersion) {
			fmt.Println("Error: Version " + newVersion + " not found")
			fmt.Println("Valid versions are:")
			for _, version := range installedVersions {
				fmt.Println("  " + version)
			}
			os.Exit(1)
		}

		setGlobalVersion(newVersion)
		fmt.Println("Global version set to " + newVersion)
		return
	}

	if os.Args[1] == "local" {
		if len(os.Args) == 2 {
			localVersion, ok, _ := getLocalVersion()
			if !ok {
				fmt.Println("local version is not set")
				return
			}
			if slices.Contains(installedVersions, localVersion) {
				fmt.Println(localVersion)
				return
			} else {
				fmt.Println("Error: local version " + localVersion + " not installed")
				fmt.Println("Please try specifiying a version with 'xienv local <version>'")
				fmt.Println("Valid versions are:")
				for _, version := range installedVersions {
					fmt.Println("  " + version)
				}
			}
		}

		if len(os.Args) != 3 {
			printUsage()
			os.Exit(1)
		}

		newVersion := os.Args[2]
		if !slices.Contains(installedVersions, newVersion) {
			fmt.Println("Error: Version " + newVersion + " not found")
			fmt.Println("Valid versions are:")
			for _, version := range installedVersions {
				fmt.Println("  " + version)
			}
			os.Exit(1)
		}

		setLocalVersion(newVersion)
		fmt.Println("Local version set to " + newVersion)
		return
	}

	currentVersion, isGlobal, isVersionSet, path := getVersion()

	if os.Args[1] == "versions" {
		if len(os.Args) > 2 && os.Args[2] == "--short" {
			for _, version := range installedVersions {
				fmt.Println(version)
			}
			return
		}
		for _, version := range installedVersions {
			if version == currentVersion {
				if isGlobal {
					fmt.Println("* " + version + " (global)")
				} else {
					fmt.Println("* " + version + " (local at " + path + ")")
				}
			} else {
				fmt.Println("  " + version)
			}
		}
		return
	}

	if os.Args[1] == "version" && len(os.Args) > 2 && os.Args[2] == "--no-error" {
		if isVersionSet {
			fmt.Print(currentVersion)
		}
		return
	}

	if os.Args[1] == "version" {
		if !isVersionSet {
			fmt.Println("Error: No version specified")
			fmt.Println("Please try specifiying a version with 'xienv global <version>' or 'xienv local <version>'")
			fmt.Println("Available versions are:")
			for _, version := range installedVersions {
				fmt.Println("  " + version)
			}
			os.Exit(1)

		}
		fmt.Print(currentVersion)
		return
	}

	if os.Args[1] == "check" {
		if !isVersionSet {
			fmt.Println("Error: No version specified")
			fmt.Println("Please try specifiying a version with 'xienv global <version>' or 'xienv local <version>'")
			fmt.Println("Available versions are:")
			for _, version := range installedVersions {
				fmt.Println("  " + version)
			}
			os.Exit(1)
		}
		if !slices.Contains(installedVersions, currentVersion) {
			if isGlobal {
				fmt.Println("Error: Global version " + currentVersion + " not installed")
				fmt.Println("Please try specifiying a version with 'xienv global <version>'")
				fmt.Println("Valid versions are:")
				for _, version := range installedVersions {
					fmt.Println("  " + version)
				}
			} else {
				fmt.Println("Error: Local version " + currentVersion + " not installed")
				fmt.Println("Please try specifiying a version with 'xienv local <version>'")
				fmt.Println("Valid versions are:")
				for _, version := range installedVersions {
					fmt.Println("  " + version)
				}
			}
			os.Exit(1)
		}
		os.Exit(0)
	}

	printUsage()
	os.Exit(1)
}

func getInstalledVersions() []string {
	entries, err := os.ReadDir("/tools/Xilinx/Vivado")
	if err != nil {
		// TODO: error handling
		panic(err)
	}

	var versions []string
	for _, e := range entries {
		versions = append(versions, e.Name())
	}
	return versions
}

// return version, isGlobal, isOk, localPath
func getVersion() (string, bool, bool, string) {
	localVersion, ok, path := getLocalVersion()
	if !ok {
		globalVersion, ok := getGlobalVersion()
		if !ok {
			return "", false, false, ""
		}
		return globalVersion, true, true, ""
	}
	return localVersion, false, true, path

}

func getGlobalVersion() (string, bool) {
	home, err := os.UserHomeDir()
	f, err := os.OpenFile(home+"/.xienv/version", os.O_RDONLY, 0666)
	if err != nil {
		return "", false
	}
	reader := bufio.NewReaderSize(f, 4096)
	line, _, err := reader.ReadLine()
	return string(line), true
}

func setGlobalVersion(ver string) {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	_, err = os.Stat(home + "/.xienv")
	if err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(home+"/.xienv", 0775)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	f, err := os.OpenFile(home+"/.xienv/version", os.O_CREATE+os.O_WRONLY, 0664)
	if err != nil {
		panic(err)
	}
	f.WriteString(ver)
}

func getLocalVersionAt(dir string) (string, bool, string) {
	f, err := os.OpenFile(dir+"/.xilinx_version", os.O_RDONLY, 0664)
	if err != nil {
		if dir == "/" {
			return "", false, ""
		}
		parent := path.Dir(dir)
		version, ok, path := getLocalVersionAt(parent)
		return version, ok, path
	}
	reader := bufio.NewReaderSize(f, 4096)
	line, _, err := reader.ReadLine()
	return string(line), true, dir
}

// return version, isOk, path
func getLocalVersion() (string, bool, string) {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	version, ok, path := getLocalVersionAt(path)
	return version, ok, path
	// TODO: find parent directory with .xilinx_version
}

func setLocalVersion(ver string) {
	f, err := os.OpenFile(".xilinx_version", os.O_WRONLY, 0664)
	if err != nil {
		panic(err)
	}
	f.WriteString(ver)
}
