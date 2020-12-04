package cmd

import (
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	cobra_ext "github.com/linuxsuren/cobra-extension"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// JobTypeOption is the job type cmd option
type JobTypeOption struct {
	cobra_ext.OutputOption
	common.Option
}

var jobTypeOption JobTypeOption

func init() {
	jobCmd.AddCommand(jobTypeCmd)
	jobTypeOption.SetFlagWithHeaders(jobTypeCmd, "DisplayName,Class")
}

var jobTypeCmd = &cobra.Command{
	Use:   "type",
	Short: i18n.T("Print the types of job which in your Jenkins"),
	Long:  i18n.T("Print the types of job which in your Jenkins"),
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		jclient := &client.JobClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: jobTypeOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClientOrDie(&(jclient.JenkinsCore))

		var jobCategories []client.JobCategory
		jobCategories, err = jclient.GetJobTypeCategories()
		if err == nil {
			var jobCategoryItems []client.JobCategoryItem
			for _, jobCategory := range jobCategories {
				for _, item := range jobCategory.Items {
					jobCategoryItems = append(jobCategoryItems, item)
				}
			}
			jobTypeOption.Writer = cmd.OutOrStdout()
			err = jobTypeOption.OutputV2(jobCategoryItems)
		}
		return
	},
}

// GetCategories returns the categories of current Jenkins
func GetCategories(jclient *client.JobClient) (
	typeMap map[string]string, types []string, err error) {
	typeMap = make(map[string]string)
	var categories []client.JobCategory
	if categories, err = jclient.GetJobTypeCategories(); err == nil {
		for _, category := range categories {
			for _, item := range category.Items {
				typeMap[item.DisplayName] = item.Class
			}
		}

		types = make([]string, len(typeMap))
		i := 0
		for tp := range typeMap {
			types[i] = tp
			i++
		}
	}
	return
}
