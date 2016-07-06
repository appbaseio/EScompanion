package main

import (
	"log"
	"os/user"
	"regexp"
	"strings"
	"testing"

	"github.com/joho/godotenv"
)

func createTestManager() *Manager {
	m := &Manager{
		Version:    "2.3.3",
		Exec:       RemoteExecutor{},
		GetCommand: DefaultCommandProvider,
	}
	return m
}

func listPlugins() {
	godotenv.Load()
	resp, _ := RemoteExecutor{}.run("/usr/share/elasticsearch/bin/plugin list")
	log.Println(resp)

}

func TestLocalExecutor(t *testing.T) {

	user, _ := user.Current()
	resp, err := LocalExecutor{}.run("whoami")
	resp = strings.TrimSpace(resp)
	if err != nil {
		t.Fatal(err)
	} else {
		if resp != user.Username {
			t.Fatal("the server has different user : ", resp)
		}
	}
}

func TestManager(t *testing.T) {
	m := createTestManager()

	if command := m.GetCommand(m.Version); command != "plugin install" {
		t.Fatal("wrong command ", command)
	}

}

func TestOldESManager(t *testing.T) {
	m := createTestManager()
	m.Version = "1.7.2"

	if command := m.GetCommand(m.Version); command != "plugin --install" {
		t.Fatal("wrong command ", command)
	}

}

func testCommandProvider(version string) string {
	if match, _ := regexp.Match("2.+", []byte(version)); match {
		return "/usr/share/elasticsearch/bin/plugin install"
	} else if match, _ := regexp.Match("1.7*", []byte(version)); match {
		return "/usr/share/elasticsearch/bin/plugin --install"
	} else {
		panic("Invalid Version")
	}
}

func TestSinglePluginInstall(t *testing.T) {
	godotenv.Load()

	m := createTestManager()
	m.GetCommand = testCommandProvider
	resp, err := m.Install("mobz/elasticsearch-head")
	if err != nil {
		t.Fatal("There was an error installing plugins :", err)
	} else {
		t.Log(resp)
	}

	listPlugins()

}
func TestMultiplePluginInstall(t *testing.T) {
	godotenv.Load()

	m := createTestManager()
	m.GetCommand = testCommandProvider
	resp, err := m.Install("mobz/elasticsearch-head", "appbaseio/dejaVu")
	if err != nil {
		t.Fatal("There was an error installing plugins :", err)
	} else {
		t.Log(resp)
	}

	listPlugins()

}
