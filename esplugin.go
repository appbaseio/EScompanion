package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/hypersleep/easyssh"
	"github.com/joho/godotenv"
)

type Executor interface {
	run(command string, args ...string) (string, error)
}

type RemoteExecutor struct{}
type LocalExecutor struct{}

type CommandProvider func(version string) string

func (remote RemoteExecutor) run(command string, args ...string) (string, error) {
	log.Println(os.Getenv("key"))
	ssh := &easyssh.MakeConfig{
		User:   os.Getenv("USER"),
		Server: os.Getenv("SERVER"),
		Port:   os.Getenv("PORT"),
		// optionally key as well
		//  Key: "~/.ssh/id_rsa"
	}
	if password := os.Getenv("PASSWORD"); password != "" {
		ssh.Password = password
	} else {
		ssh.Key = os.Getenv("KEY")
	}

	var resp string
	var err error

	if len(args) > 0 {
		resp, err = ssh.Run(fmt.Sprintf("%s %s", command, strings.Join(args, " ")))
	} else {
		resp, err = ssh.Run(command)
	}
	return resp, err
}

func (local LocalExecutor) run(command string, args ...string) (string, error) {
	var cmd *exec.Cmd

	if len(args) > 0 {
		cmd = exec.Command(fmt.Sprintf("%s %s", command, strings.Join(args, " ")))
	} else {
		cmd = exec.Command(command)

	}
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

type Manager struct {
	Version    string
	Exec       Executor
	GetCommand CommandProvider
}

func DefaultCommandProvider(version string) string {

	if match, _ := regexp.Match("2.+", []byte(version)); match {
		return "/usr/share/elasticsearch/bin/plugin install"
	} else if match, _ := regexp.Match("1.7*", []byte(version)); match {
		return "/usr/share/elasticsearch/bin/plugin --install"
	} else {
		panic("Invalid Version")
	}

}

func (m *Manager) Install(plugins ...string) (string, error) {
	command := m.GetCommand(m.Version)
	for i := 0; i < len(plugins); i++ {
		log.Println("Installing plugin ", plugins[i])
		resp, err := m.Exec.run(command, plugins[i])
		log.Println(resp)
		if err != nil {
			return "", err
		}
	}

	return "All plugins installed successfully", nil
}

func main() {
	err := godotenv.Load()

	var executor Executor
	if err != nil {
		log.Println(".env file not found")
		log.Println("installing plugins on the local system")
		executor = LocalExecutor{}
	} else {
		log.Println(".env file found")
		log.Println("installing plugins on the system with ip ", os.Getenv("Server"))
		executor = RemoteExecutor{}
	}

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
		Exec:       executor,
		GetCommand: DefaultCommandProvider,
	}

	message, err := m.Install(plugins...)

	if err != nil {
		log.Fatal("There was an error while installing the plugins ", err)
	} else {
		log.Println(message)
	}
}
