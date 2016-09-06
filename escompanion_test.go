package main

import (
	"testing"

	"github.com/joho/godotenv"
)

func createTestManager() *Manager {
	m := &Manager{
		Version:    "2.3.3",
		GetCommand: DefaultCommandProvider,
	}
	return m
}

func listPlugins() {
	godotenv.Load()

}

func TestCommandProvider(t *testing.T) {

	command := DefaultCommandProvider("2.4.1", "testing")

	if command[1] != "--batch" {
		t.Fatal("command doesnt match")
	}

	command = DefaultCommandProvider("2.3", "testing")

	if command[1] != "--batch" {
		t.Fatal("command doesnt match")
	}
	command = DefaultCommandProvider("2.2", "testing")
	if command[1] != "--batch" {
		t.Fatal("command doesnt match")
	}
	command = DefaultCommandProvider("2.1", "testing")

	command = DefaultCommandProvider("1.7.2", "testing")

	command = DefaultCommandProvider("2.3.2", "testing")

}
