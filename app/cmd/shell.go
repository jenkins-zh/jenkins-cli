package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/util"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"go.uber.org/zap"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/spf13/cobra"
)

// ShellOptions is the option of shell command
type ShellOptions struct {
	CommonOption
}

var shellOptions ShellOptions

func init() {
	rootCmd.AddCommand(shellCmd)
}

const (
	defaultRcFile = `
if [ -f /etc/bashrc ]; then
    source /etc/bashrc
fi
if [ -f ~/.bashrc ]; then
    source ~/.bashrc
fi
if type -t __start_jcli >/dev/null; then true; else
	source <(jcli completion)
fi
[[ -r "/usr/local/etc/profile.d/bash_completion.sh" ]] && . "/usr/local/etc/profile.d/bash_completion.sh"
`

	zshRcFile = `
if [ -f /etc/zshrc ]; then
    source /etc/zshrc
fi
if [ -f ~/.zshrc ]; then
    source ~/.zshrc
fi
`
)

var shellCmd = &cobra.Command{
	Use:     "shell [<name>]",
	Short:   i18n.T("Create a sub shell so that changes to a specific Jenkins remain local to the shell."),
	Long:    i18n.T("Create a sub shell so that changes to a specific Jenkins remain local to the shell."),
	Aliases: []string{"sh"},
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			jenkinsName := args[0]
			setCurrentJenkins(jenkinsName)
		}
	},
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		var tmpDirName string
		tmpDirName, err = ioutil.TempDir("", ".jcli-shell-")
		if err != nil {
			return err
		}
		tmpConfigFileName := filepath.Join(tmpDirName, "/config")

		var data []byte
		config := getConfig()
		if data, err = yaml.Marshal(&config); err == nil {
			err = ioutil.WriteFile(tmpConfigFileName, data, 0644)
		}
		if err != nil {
			return
		}

		fullShell := os.Getenv("SHELL")
		shell := filepath.Base(fullShell)
		if fullShell == "" && runtime.GOOS == "windows" {
			// SHELL is set by git-bash but not cygwin :-(
			shell = "cmd.exe"
		}

		prompt := createNewBashPrompt(os.Getenv("PS1"))
		rcFile := defaultRcFile + "\nexport PS1=" + prompt + "\nexport JCLI_CONFIG=\"" + tmpConfigFileName + "\"\n"
		tmpRCFileName := tmpDirName + "/.bashrc"

		err = ioutil.WriteFile(tmpRCFileName, []byte(rcFile), 0760)
		if err != nil {
			return err
		}

		logger.Debug("temporary shell profile loaded", zap.String("path", tmpRCFileName))
		//e := exec.Command(shell, "-rcfile", tmpRCFileName, "-i")
		e := util.ExecCommand(shellOptions.ExecContext, shell, "-rcfile", tmpRCFileName, "-i")
		if shell == "zsh" {
			env := os.Environ()
			env = append(env, fmt.Sprintf("ZDOTDIR=%s", tmpDirName))
			//e = exec.Command(shell, "-i")
			e = util.ExecCommand(shellOptions.ExecContext, shell, "-i")
			e.Env = env
		} else if shell == "cmd.exe" {
			env := os.Environ()
			env = append(env, fmt.Sprintf("JCLI_CONFIG=%s", tmpConfigFileName))
			//e = exec.Command(shell)
			e = util.ExecCommand(shellOptions.ExecContext, shell)
			e.Env = env
		}

		e.Stdout = cmd.OutOrStdout()
		e.Stderr = cmd.OutOrStderr()
		e.Stdin = os.Stdin
		err = e.Run()
		if deleteError := os.RemoveAll(tmpDirName); deleteError != nil {
			err = fmt.Errorf("cannot remove dir %s, the shell command error is %#v", tmpDirName, err)
		}
		return err
	},
}

func createNewBashPrompt(prompt string) string {
	if prompt == "" {
		return "'[\\u@\\h \\W jcli> ]\\$ '"
	}
	if prompt[0] == '"' {
		return prompt[0:1] + "jcli> " + prompt[1:]
	}
	if prompt[0] == '\'' {
		return prompt[0:1] + "jcli> " + prompt[1:]
	}
	return "'jcli> " + prompt + "'"
}
