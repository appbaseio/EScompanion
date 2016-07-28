package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
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

	if match, _ := regexp.Match("2.2+", []byte(version)); match {
		return []string{"install", "--batch", plugin}
	} else if match, _ := regexp.Match("2.3+", []byte(version)); match {
		return []string{"install", "--batch", plugin}

	} else if match, _ := regexp.Match("2.+", []byte(version)); match {
		return []string{"install", plugin}
	} else if match, _ := regexp.Match("1.7+", []byte(version)); match {
		return []string{"--install", plugin}
	} else {
		panic("Invalid Version")
	}

}

func SaveToFile(url string, path string) error {
	res, err := http.Get(url)

	if err != nil {
		return err
	}

	file, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil
	}
	return ioutil.WriteFile(path, file, 0644)
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

	dePath := "./data.zip"
	deConfigURL := "http://localhost:8000/d.zip"
	version := flag.String("version", "2.3", "The elasticsearch version installed on the system")
	esPath := flag.String("path", dePath, "The elasticsearch yml file location")
	configUrl := flag.String("url", deConfigURL, "Location of the elasticsearch url")

	flag.Parse()
	if *configUrl != deConfigURL && *esPath != dePath {
		log.Println("Detecting new elasticsearch yml file url ", *configUrl)
		SaveToFile(*configUrl, *esPath)
	}
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
