package main

import (
	"fmt"

	"github.com/jenkins-zh/jenkins-cli/app/cmd"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
)

func main() {
	if err := i18n.LoadTranslations("jcli", nil); err != nil {
		fmt.Println(err)
	}

	cmd.Execute()
}
