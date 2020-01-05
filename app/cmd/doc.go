package cmd

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/jenkins-zh/jenkins-cli/app"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func init() {
	rootCmd.AddCommand(docCmd)
}

const (
	gendocFrontmatterTemplate = `---
date: %s
title: "%s"
version: %s
---
`
)

var docCmd = &cobra.Command{
	Use:   "doc <output dir>",
	Short: i18n.T("Generate document for all jcl commands"),
	Long:  i18n.T("Generate document for all jcl commands"),
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		now := time.Now().Format(time.RFC3339)
		prepender := func(filename string) string {
			name := filepath.Base(filename)
			base := strings.TrimSuffix(name, path.Ext(name))
			return fmt.Sprintf(gendocFrontmatterTemplate, now,
				strings.Replace(base, "_", " ", -1),
				app.GetVersion())
		}

		linkHandler := func(name string) string {
			base := strings.TrimSuffix(name, path.Ext(name))
			return "/commands/" + strings.ToLower(base) + "/"
		}

		outputDir := args[0]

		rootCmd.DisableAutoGenTag = true
		err = doc.GenMarkdownTreeCustom(rootCmd, outputDir, prepender, linkHandler)
		return
	},
}
