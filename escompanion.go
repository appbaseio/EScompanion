package main

import (
	"bytes"
	"errors"
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

func (m *Manager) getPluginCommand() string {
	if match, _ := regexp.Match("5+", []byte(m.Version)); match {
		return "elasticsearch-plugin"
	} else {
		return "plugin"
	}

}

func (m *Manager) run(plugin string) (string, error) {
	var cmd *exec.Cmd

	log.Println(m.GetCommand(m.Version, plugin))
	cmd = exec.Command(m.getPluginCommand(), m.GetCommand(m.Version, plugin)...)
	var stderr, stdout bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return "", errors.New(err.Error() + ":" + stdout.String())
	}
	return stdout.String(), nil
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

	} else if match, _ := regexp.Match("2.1+", []byte(version)); match {
		return []string{"install", plugin}
	} else if match, _ := regexp.Match("2.0+", []byte(version)); match {
		return []string{"install", plugin}
	} else if match, _ := regexp.Match("2.+", []byte(version)); match {
		return []string{"install", "--batch", plugin}
	} else if match, _ := regexp.Match("1.7+", []byte(version)); match {
		return []string{"--install", plugin}
	} else if match, _ := regexp.Match("1.6+", []byte(version)); match {
		return []string{"--install", plugin}
	} else if match, _ := regexp.Match("5+", []byte(version)); match {
		return []string{"install", "--batch", plugin}
	}
	panic("Invalid Version")

}

func WgetFile(url string, path string) error {

	cmd := exec.Command("wget", url, "-O", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	log.Println(string(output))

	return err

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

func installAwsPlugin(version *string, m *Manager) {

	if match, _ := regexp.Match("5+", []byte(*version)); match {
		log.Println("Installing repository s3")
		m.Install("repository-s3")
	} else if match, _ := regexp.Match("2.+", []byte(*version)); match {
		log.Println("Installing cloud-aws")
		m.Install("cloud-aws")
	} else if match, _ := regexp.Match("1.7+", []byte(*version)); match {
		log.Println("Installing the cloud-aws/2.7.1")
		m.Install("elasticsearch/elasticsearch-cloud-aws/2.7.1")
	} else if match, _ := regexp.Match("1.6+", []byte(*version)); match {
		log.Println("Installing the cloud-aws/2.6.1")
		m.Install("elasticsearch/elasticsearch-cloud-aws/2.6.1")
	}
}
func installAzurePlugin(version *string, m *Manager) {
	if match, _ := regexp.Match("5+", []byte(*version)); match {
		log.Println("Installing repository azure")
		m.Install("repository-azure")
	} else if match, _ := regexp.Match("2.+", []byte(*version)); match {
		log.Println("Installing cloud-azure")
		m.Install("cloud-azure")
	} else if match, _ := regexp.Match("1.7+", []byte(*version)); match {
		log.Println("Installing the cloud-azure/2.8.3")
		m.Install("elasticsearch/elasticsearch-cloud-azure/2.8.3")
	} else if match, _ := regexp.Match("1.6+", []byte(*version)); match {
		log.Println("Installing the cloud-azure/2.7.1")
		m.Install("elasticsearch/elasticsearch-cloud-azure/2.7.1")
	}
}

func main() {
	err := godotenv.Load()

	version := flag.String("version", "2.3", "The elasticsearch version installed on the system")
	esPath := flag.String("path", "", "The elasticsearch yml file location")
	configUrl := flag.String("url", "", "Location of the elasticsearch url")
	backup := flag.Bool("backup", false, "Should the s3 backup plugin be installed?")
	provider := flag.String("provider", "aws", "the remote storage provider (supported values : aws , azure )")

	flag.Parse()

	if *configUrl != "" && *esPath != "" {
		log.Println("Detecting new elasticsearch yml file url ", *configUrl)
		WgetFile(*configUrl, *esPath)
	}
	plugins := flag.Args()

	log.Println("Using elasticsearch version ", *version)

	m := &Manager{
		Version:    *version,
		GetCommand: DefaultCommandProvider,
	}

	if *backup {
		switch *provider {
		case "aws":
			installAwsPlugin(version, m)
		case "azure":
			installAzurePlugin(version, m)
		}

	}

	if len(plugins) < 1 {
		log.Fatalln("You havent provided any plugins")
		os.Exit(0)
	}
	log.Println("The plugins being installed are ")
	for i := 0; i < len(plugins); i++ {
		log.Println("+ ", plugins[i])
	}

	log.Println("Starting installation ..")

	// Main Execution

	m.Install("license")

	message, err := m.Install(plugins...)

	if err != nil {
		log.Fatal("There was an error while installing the plugins ", err)
	} else {
		log.Println(message)
	}
}
