package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func init() {
	rootCmd.AddCommand(docCmd)
}

const (
	gendocFrontmatterTemplate = `---
date: %s
title: "%s"
anchor: %s
url: %s
---
`
)

var docCmd = &cobra.Command{
	Use:   "doc <output dir>",
	Short: i18n.T("Generate document for all jcl commands"),
	Long:  i18n.T("Generate document for all jcl commands"),
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		now := time.Now().Format(time.RFC3339)
		prepender := func(filename string) string {
			name := filepath.Base(filename)
			base := strings.TrimSuffix(name, path.Ext(name))
			url := "/commands/" + strings.ToLower(base) + "/"
			return fmt.Sprintf(gendocFrontmatterTemplate, now, strings.Replace(base, "_", " ", -1), base, url)
		}

		linkHandler := func(name string) string {
			base := strings.TrimSuffix(name, path.Ext(name))
			return "/commands/" + strings.ToLower(base) + "/"
		}

		outputDir := args[0]

		err := doc.GenMarkdownTreeCustom(rootCmd, outputDir, prepender, linkHandler)
		if err != nil {
			cmd.PrintErr(err)
		}
	},
}
