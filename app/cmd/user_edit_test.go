package cmd

import (
	"github.com/Netflix/go-expect"
	"io/ioutil"
	"testing"
	"time"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/golang/mock/gomock"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
)

func TestEditUser(t *testing.T) {
	RunEditCommandTest(t, EditCommandTest{
		Procedure: func(c *expect.Console) {
			c.ExpectString("Edit user description")
			c.SendLine("")
			go c.ExpectEOF()
			time.Sleep(time.Millisecond)
			c.Send("\x1b")
			c.SendLine(":wq!")
		},
		Test: func(stdio terminal.Stdio) (err error) {
			data, err := generateSampleConfig()
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)

			var (
				description = "fake-description"
			)

			ctrl := gomock.NewController(t)
			roundTripper := mhttp.NewMockRoundTripper(ctrl)

			client.PrepareForEditUserDesc(roundTripper, "http://localhost:8080/jenkins",
				"admin", "admin", description, "111e3a2f0231198855dceaff96f20540a9")

			rootCmd.SetArgs([]string{"user", "edit", "--desc", description})

			userEditOption.CommonOption.Stdio = stdio
			_, err = rootCmd.ExecuteC()
			return
		},
	})
}
