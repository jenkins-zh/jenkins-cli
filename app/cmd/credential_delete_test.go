package cmd

import (
	"bytes"
	"github.com/Netflix/go-expect"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/jenkins-zh/jenkins-cli/client"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
)

var _ = Describe("credential delete command", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
		buf          *bytes.Buffer
		store        string
		id           string
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		rootCmd.SetArgs([]string{})
		buf = new(bytes.Buffer)
		rootCmd.SetOutput(buf)
		rootOptions.Jenkins = ""
		rootOptions.ConfigFile = "test.yaml"

		credentialDeleteOption.RoundTripper = roundTripper

		store = "system"
		id = "fake-id"
	})

	AfterEach(func() {
		rootCmd.SetArgs([]string{})
		os.Remove(rootOptions.ConfigFile)
		rootOptions.ConfigFile = ""
		ctrl.Finish()
	})

	Context("basic cases", func() {
		var (
			err error
		)

		BeforeEach(func() {
			var data []byte
			data, err = generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())
		})

		It("lack of the necessary parameters", func() {
			rootCmd.SetArgs([]string{"credential", "delete", "--store=", "--id="})
			_, err = rootCmd.ExecuteC()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("the store or id of target credential is empty"))
		})

		It("should success", func() {
			client.PrepareForDeleteCredential(roundTripper, "http://localhost:8080/jenkins",
				"admin", "111e3a2f0231198855dceaff96f20540a9", store, id)

			rootCmd.SetArgs([]string{"credential", "delete", store, id, "-b"})
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())
		})
	})
})

func TestConfirmCommands(t *testing.T) {
	RunPromptCommandTest(t, PromptCommandTest{
		Args: []string{"credential", "delete", "fake-store", "fake-id", "-b=false"},
		Procedure: func(c *expect.Console) {
			c.ExpectString("Are you sure to delete credential fake-id")
			c.SendLine("n")
			c.ExpectEOF()
		},
		BatchOption: &credentialDeleteOption.BatchOption,
	})

	RunPromptCommandTest(t, PromptCommandTest{
		Args: []string{"job", "stop", "fake", "-b=false"},
		Procedure: func(c *expect.Console) {
			c.ExpectString("Are you sure to stop job fake ?")
			c.SendLine("n")
			c.ExpectEOF()
		},
		BatchOption: &jobStopOption.BatchOption,
	})

	RunPromptCommandTest(t, PromptCommandTest{
		Args: []string{"job", "build", "fake", "-b=false"},
		Procedure: func(c *expect.Console) {
			c.ExpectString("Are you sure to build job fake")
			c.SendLine("n")
			c.ExpectEOF()
		},
		BatchOption: &jobBuildOption.BatchOption,
	})

	RunPromptCommandTest(t, PromptCommandTest{
		Args: []string{"job", "delete", "fake", "-b=false"},
		Procedure: func(c *expect.Console) {
			c.ExpectString("Are you sure to delete job fake ?")
			c.SendLine("n")
			c.ExpectEOF()
		},
		BatchOption: &jobDeleteOption.BatchOption,
	})

	RunPromptCommandTest(t, PromptCommandTest{
		Args: []string{"user", "delete", "fake-user", "-b=false"},
		Procedure: func(c *expect.Console) {
			c.ExpectString("Are you sure to delete user fake-user ?")
			c.SendLine("n")
			c.ExpectEOF()
		},
		BatchOption: &userDeleteOption.BatchOption,
	})

	RunPromptCommandTest(t, PromptCommandTest{
		Args: []string{"restart", "-b=false"},
		Procedure: func(c *expect.Console) {
			c.ExpectString("Are you sure to restart Jenkins http://localhost:8080/jenkins?")
			c.SendLine("n")
			c.ExpectEOF()
		},
		BatchOption: &restartOption.BatchOption,
	})

	RunPromptTest(t, PromptTest{
		Message:    "message",
		MsgConfirm: &common.BatchOption{},
		Procedure: func(c *expect.Console) {
			c.ExpectString("message")
			c.SendLine("y")
			c.ExpectEOF()
		},
		Expected: true,
	})

	RunEditorTest(t, EditorTest{
		Message:        "message",
		DefaultContent: "hello",
		EditContent:    &common.CommonOption{},
		Procedure: func(c *expect.Console) {
			c.ExpectString("message")
			c.SendLine("")
			go c.ExpectEOF()
			time.Sleep(time.Millisecond)
			c.Send("ddigood\x1b")
			c.SendLine(":wq!")
		},
		Expected: "good\n",
	})
}
