package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cascCmd)

	healthCheckRegister.Register(cascCmd, &CASCOptions{})
	healthCheckRegister.RegisterPath(getCmdPath(cascCmd)+".*", &CASCOptions{})
}

type CASCOptions struct {
}

func (o *CASCOptions) Check() (err error) {
	fmt.Println("hello")
	err = fmt.Errorf("Sdfsdf fake error")
	return
}

var cascCmd = &cobra.Command{
	Use:   "casc",
	Short: i18n.T("Configuration as Code"),
	Long:  i18n.T("Configuration as Code"),
}
