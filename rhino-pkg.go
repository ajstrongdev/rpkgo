package main

import (
	"fmt"
	"os"
	"strings"
	"os/exec"
	"regexp"
)

var version = "2.0.0"
var nala = false;

const (
	resetColor = "\033[0m"
	defaultColor = "\033[39m"
)

func main() {
	checkNala() // Check if Nala is installed
	if len(os.Args) < 2 {
		fmt.Println("Usage: rhino-pkg <command>")
		os.Exit(1)
	}
	command := os.Args[1]
	args := os.Args[2:]
	switch command {
	case "search":
		search(args)
	case "install":
		install(args)
	case "remove":
		//remove(args)
	case "update":
		//update(args)
	case "cleanup":
		//cleanup()
	default:
		fmt.Println("Unknown command:", command)
	}
}

func checkNala() bool {
	_, err := exec.LookPath("nala")
	if err == nil {
		nala = true
	}
	return err == nil
}

// Check which managers are installed
func checkManagers(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// Removes colour from pacstall output
func stripAnsi(s string) string {
	re := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	return re.ReplaceAllString(s, "")
}

func search(args []string) {
	query := strings.Join(args, " ")
	fmt.Printf("Searching for packages matching '%s'...\n", query)
	type Package struct {
		Name   string
		Source string
	}
	var packageList []Package
	searchCommands := map[string][]string{
		"apt":      {"apt-cache", "search", query},
		"pacstall": {"pacstall", "-S", query},
		"flatpak":  {"flatpak", "search", query},
		"snap":     {"snap", "find", query},
	}
	// Check which package managers are installed on your system
	availableManagers := make(map[string][]string)
	for source, cmd := range searchCommands {
		if checkManagers(cmd[0]) {
			availableManagers[source] = cmd
		}
	}
	// Query package managers
	for source, cmd := range availableManagers {
		output := captureCommandOutput(cmd[0], cmd[1:]...)
		for _, line := range strings.Split(output, "\n") {
			cleanLine := stripAnsi(line)
			fields := strings.Fields(cleanLine)
			if len(fields) > 0 {
				packageList = append(packageList, Package{fields[0], source})
			}
		}
	}
	for i, pkg := range packageList {
		fmt.Printf("[%d] %s (%s)\n", i, pkg.Name, pkg.Source)
	}
}

func install(args []string) {
	query := strings.Join(args, " ")
	fmt.Printf("Searching for packages matching '%s'...\n", query)
	type Package struct {
		Name   string
		Source string
	}
	var packageList []Package
	searchCommands := map[string][]string{
		"apt":      {"apt-cache", "search", query},
		"pacstall": {"pacstall", "-S", query},
		"flatpak":  {"flatpak", "search", query},
		"snap":     {"snap", "find", query},
	}
	// Check which package managers are installed on your system
	availableManagers := make(map[string][]string)
	for source, cmd := range searchCommands {
		if checkManagers(cmd[0]) {
			availableManagers[source] = cmd
		}
	}
	// Query package managers
	for source, cmd := range availableManagers {
		output := captureCommandOutput(cmd[0], cmd[1:]...)
		for _, line := range strings.Split(output, "\n") {
			cleanLine := stripAnsi(line)
			fields := strings.Fields(cleanLine)
			if len(fields) > 0 {
				packageList = append(packageList, Package{fields[0], source})
			}
		}
	}
	for i, pkg := range packageList {
		fmt.Printf("[%d] %s (%s)\n", i, pkg.Name, pkg.Source)
	}
	listLength := len(packageList)
	fmt.Printf("Enter the number of the package to install [0-%d]: ", listLength-1)
	var selection int
	_, err := fmt.Scanf("%d", &selection)
	if err != nil {
		fmt.Println("Error reading input:", err)
		os.Exit(1)
	}
	if selection < 0 || selection >= len(packageList) {
		fmt.Println("Invalid selection")
		os.Exit(1)
	}
	pkg := packageList[selection]
	fmt.Printf("Would you like to install %s (%s)? [y/N] ", pkg.Name, pkg.Source)
	var confirmation string
	_, err = fmt.Scanf("%s", &confirmation)
	if err != nil {
		fmt.Println("Error reading input:", err)
		os.Exit(1)
	}
	if confirmation == "y" {
		switch pkg.Source {
		case "apt":
			if !nala {
				runCommand("sudo", "apt", "install", "-y", pkg.Name)
			} else {
				runCommand("sudo", "nala", "install", "-y", pkg.Name)
			}
		case "pacstall":
			runCommand("pacstall", "-I", pkg.Name)
		case "flatpak":
			runCommand("flatpak", "install", pkg.Name)
		case "snap":
			runCommand("sudo", "snap", "install", pkg.Name)
		default:
			fmt.Println("Unknown package source:", pkg.Source)
		}
	} else {
		os.Exit(1)
	}
}

func captureCommandOutput(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error running %s: %v\n", name, err)
		return ""
	}
	return string(output)
}

func runCommand(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error running %s: %v\n", name, err)
	}
}