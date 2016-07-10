package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"regexp"

	"github.com/joho/godotenv"
)

type CommandProvider func(version string, plugin string) []string

func (m *Manager) run(plugin string) (string, error) {
	var cmd *exec.Cmd

	log.Println(m.GetCommand(m.Version, plugin))
	cmd = exec.Command("plugin", m.GetCommand(m.Version, plugin)...)
	output, err := cmd.CombinedOutput()
	log.Println(string(output))
	if err != nil {
		return "", err
	}
	return string(output), nil
}

type Manager struct {
	Version    string
	GetCommand CommandProvider
}

func DefaultCommandProvider(version string, plugin string) []string {

	if match, _ := regexp.Match("2.+", []byte(version)); match {
		return []string{"install", "--batch", plugin}
	} else if match, _ := regexp.Match("1.7*", []byte(version)); match {
		return []string{"--install", "--batch", plugin}
	} else {
		panic("Invalid Version")
	}

}

func (m *Manager) Install(plugins ...string) (string, error) {
	for i := 0; i < len(plugins); i++ {
		log.Println("Installing plugin ", plugins[i])
		resp, err := m.run(plugins[i])
		log.Println(resp)
		if err != nil {
			return "", err
		}
	}

	return "All plugins installed successfully", nil
}

func main() {
	err := godotenv.Load()

	version := flag.String("version", "2.3", "The elasticsearch version installed on the system")
	flag.Parse()

	plugins := flag.Args()

	log.Println("Using elasticsearch version ", *version)

	if len(plugins) < 1 {
		log.Fatalln("You havent provided any plugins")
		os.Exit(-1)
	}
	log.Println("The plugins being installed are ")
	for i := 0; i < len(plugins); i++ {
		log.Println("+ ", plugins[i])
	}

	log.Println("Starting installation ..")

	// Main Execution

	m := &Manager{
		Version:    *version,
		GetCommand: DefaultCommandProvider,
	}
	m.Install("license")
	message, err := m.Install(plugins...)

	if err != nil {
		log.Fatal("There was an error while installing the plugins ", err)
	} else {
		log.Println(message)
	}
}
