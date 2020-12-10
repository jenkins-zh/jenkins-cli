package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	ext "github.com/linuxsuren/cobra-extension/version"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// DocOption is the option for doc generating
type DocOption struct {
	DocType string
}

const (
	// DocTypeMarkdown represents markdown type of doc
	DocTypeMarkdown string = "Markdown"
	// DocTypeManPage represents man page type of doc
	DocTypeManPage string = "ManPage"
)

var docOption DocOption

func init() {
	rootCmd.AddCommand(docCmd)
	docCmd.Flags().StringVarP(&docOption.DocType, "doc-type", "", DocTypeMarkdown,
		"Which type of document will generate")

	err := docCmd.RegisterFlagCompletionFunc("doc-type", func(cmd *cobra.Command, args []string, toComplete string) (
		i []string, directive cobra.ShellCompDirective) {
		return []string{DocTypeMarkdown, DocTypeManPage}, cobra.ShellCompDirectiveDefault
	})
	if err != nil {
		docCmd.PrintErrf("register flag doc-type for sub-command doc failed %#v\n", err)
	}
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
	Use: "doc",
	Example: `jcli doc tmp
jcli doc --doc-type ManPage /usr/local/share/man/man1`,
	Short:  i18n.T("Generate document for all jcl commands"),
	Long:   i18n.T("Generate document for all jcl commands"),
	Args:   cobra.MinimumNArgs(1),
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		outputDir := args[0]
		if err = os.MkdirAll(outputDir, os.FileMode(0755)); err != nil {
			return
		}

		switch docOption.DocType {
		case DocTypeMarkdown:
			now := time.Now().Format(time.RFC3339)
			prepender := func(filename string) string {
				name := filepath.Base(filename)
				base := strings.TrimSuffix(name, path.Ext(name))
				return fmt.Sprintf(gendocFrontmatterTemplate, now,
					strings.Replace(base, "_", " ", -1),
					ext.GetVersion())
			}

			linkHandler := func(name string) string {
				base := strings.TrimSuffix(name, path.Ext(name))
				return "/commands/" + strings.ToLower(base) + "/"
			}

			rootCmd.DisableAutoGenTag = true
			err = doc.GenMarkdownTreeCustom(rootCmd, outputDir, prepender, linkHandler)
		case DocTypeManPage:
			header := &doc.GenManHeader{
				Title:   "Jenkins CLI",
				Section: "1",
				Source:  "Jenkins Chinese Community",
			}
			err = doc.GenManTree(rootCmd, header, outputDir)
		}
		return
	},
}
