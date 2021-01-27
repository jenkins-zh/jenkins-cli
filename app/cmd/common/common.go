package common

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/google/go-github/v29/github"
	"github.com/jenkins-zh/jenkins-cli/app/config"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/jenkins-zh/jenkins-cli/util"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"net/http"
	"os"
)

const (
	// Since indicate when the feature war added
	Since = "since"
)

// Option contains the common options
type Option struct {
	ExecContext     util.ExecContext
	SystemCallExec  util.SystemCallExec
	LookPathContext util.LookPathContext
	RoundTripper    http.RoundTripper
	Logger          *zap.Logger

	GitHubClient *github.Client

	Stdio terminal.Stdio

	// EditFileName allow editor has a better performance base on this
	EditFileName string
}

// BatchOption represent the options for a batch operation
type BatchOption struct {
	Batch bool

	Stdio terminal.Stdio
}

// MsgConfirm is the interface for confirming a message
type MsgConfirm interface {
	Confirm(message string) bool
}

// Confirm promote user if they really want to do this
func (b *BatchOption) Confirm(message string) bool {
	if !b.Batch {
		confirm := false
		var prompt survey.Prompt
		prompt = &survey.Confirm{
			Message: message,
		}
		_ = survey.AskOne(prompt, &confirm, survey.WithStdio(b.Stdio.In, b.Stdio.Out, b.Stdio.Err))
		return confirm
	}

	return true
}

// GetSystemStdio returns the stdio from system
func GetSystemStdio() terminal.Stdio {
	return terminal.Stdio{
		In:  os.Stdin,
		Out: os.Stdout,
		Err: os.Stderr,
	}
}

// EditContent is the interface for editing content from a file
type EditContent interface {
	Editor(defaultContent, message string) (content string, err error)
}

// Selector is the interface for selecting an option
type Selector interface {
	Select(options []string, message, defaultOpt string) (target string, err error)
}

// Editor edit a file than return the content
func (o *Option) Editor(defaultContent, message string) (content string, err error) {
	var fileName string
	if o.EditFileName != "" {
		fileName = o.EditFileName
	} else {
		fileName = "*.sh"
	}

	prompt := &survey.Editor{
		Message:       message,
		FileName:      fileName,
		Default:       defaultContent,
		HideDefault:   true,
		AppendDefault: true,
	}

	err = survey.AskOne(prompt, &content, survey.WithStdio(o.Stdio.In, o.Stdio.Out, o.Stdio.Err))
	return
}

// Select return a target
func (o *Option) Select(options []string, message, defaultOpt string) (target string, err error) {
	prompt := &survey.Select{
		Message: message,
		Options: options,
		Default: defaultOpt,
	}
	err = survey.AskOne(prompt, &target, survey.WithStdio(o.Stdio.In, o.Stdio.Out, o.Stdio.Err))
	return
}

// SetFlag the flag for batch option
func (b *BatchOption) SetFlag(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&b.Batch, "batch", "b", false, "Batch mode, no need confirm")
}

// WatchOption for the resources which can be watched
type WatchOption struct {
	Watch    bool
	Interval int
	Count    int
}

// SetFlag for WatchOption
func (o *WatchOption) SetFlag(cmd *cobra.Command) {
	cmd.Flags().IntVarP(&o.Interval, "interval", "i", 1, "Interval of watch")
	cmd.Flags().IntVarP(&o.Count, "count", "", 9999, "Count of watch")
}

// InteractiveOption allow user to choose whether the mode is interactive
type InteractiveOption struct {
	Interactive bool
}

// SetFlag set the option flag to this cmd
func (b *InteractiveOption) SetFlag(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&b.Interactive, "interactive", "i", false,
		i18n.T("Interactive mode"))
}

// HookOption is the option whether skip command hook
type HookOption struct {
	SkipPreHook  bool
	SkipPostHook bool
}

// GetAliasesDel returns the aliases for delete command
func GetAliasesDel() []string {
	return []string{"remove", "del"}
}

// GetEditorHelpText returns the help text related a text editor
func GetEditorHelpText() string {
	return `notepad is the default editor of Windows, vim is the default editor of unix.
But if the environment variable "VISUAL" or "EDITOR" exists, jcli will take it.
For example, you can set it under unix like this: export VISUAL=vi`
}

// JenkinsClient is the interface of get Jenkins client
type JenkinsClient interface {
	GetCurrentJenkinsFromOptions() (jenkinsServer *config.JenkinsServer)
	GetCurrentJenkinsAndClient(jClient *client.JenkinsCore) *config.JenkinsServer
}

// JenkinsConfigMgr is the interface of getting configuration
type JenkinsConfigMgr interface {
	GetMirror(string) string

	GetGitHubClient() *github.Client
	SetGitHubClient(gitHubClient *github.Client)
}
