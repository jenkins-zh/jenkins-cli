package main

import (
	"fmt"
	"github.com/cucumber/godog"
	messages "github.com/cucumber/messages-go/v10"
	"os/exec"
	"strings"
)

func thereAreGodogs(available int) error {
	Godogs = available
	return nil
}

func iEat(num int) error {
	if Godogs < num {
		return fmt.Errorf("you cannot eat %d godogs, there are %d available", num, Godogs)
	}
	Godogs -= num
	return nil
}

func thereShouldBeRemaining(remaining int) error {
	if Godogs != remaining {
		return fmt.Errorf("expected %d godogs to be remaining, but there is %d", remaining, Godogs)
	}
	return nil
}

func showVersion(subStr string) (err error) {
	cmd := exec.Command("jcli", "version")
	var data []byte
	if data, err = cmd.CombinedOutput(); err == nil {
		if !strings.Contains(string(data), subStr) {
			err = fmt.Errorf("do not contain %s", subStr)
		}
	}
	return
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^there are (\d+) godogs$`, thereAreGodogs)
	s.Step(`^I eat (\d+)$`, iEat)
	s.Step(`^there should be (\d+) remaining$`, thereShouldBeRemaining)
	s.Step(`^Show version contains (\w+)$`, showVersion)

	s.BeforeScenario(func(*messages.Pickle) {
		Godogs = 0 // clean the state before every scenario
	})
}
