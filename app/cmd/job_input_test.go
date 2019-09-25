package cmd

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
	"github.com/jenkins-zh/jenkins-cli/client"
	expect "github.com/Netflix/go-expect"
	"github.com/AlecAivazis/survey/v2/core"
	"github.com/AlecAivazis/survey/v2/terminal"
)

func init() {
	// disable color output for all prompts to simplify testing
	core.DisableColor = true
}

var _ = Describe("job input command", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		jobInputOption.RoundTripper = roundTripper
		rootCmd.SetArgs([]string{})
		rootOptions.Jenkins = ""
		rootOptions.ConfigFile = "test.yaml"
	})

	AfterEach(func() {
		rootCmd.SetArgs([]string{})
		os.Remove(rootOptions.ConfigFile)
		rootOptions.ConfigFile = ""
		ctrl.Finish()
	})

	Context("basic cases", func() {
		It("no params, will error",func(){
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			rootCmd.SetArgs([]string{"job", "input"})

			jobInputCmd.SetHelpFunc(func(cmd *cobra.Command, _ []string) {
				cmd.Print("help")
			})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal("help"))
		})

		It("should success, without inputs", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			jobName := "fakeJob"
			buildID := 1

			client.PrepareForGetJobInputActions(roundTripper, "http://localhost:8080/jenkins", "admin", "111e3a2f0231198855dceaff96f20540a9", jobName, buildID)
			
			_, w, err := os.Pipe()

			c, err := expect.NewConsole(expect.WithStdout(w))
			Expect(err).To(BeNil())
			jobInputOption.Stdio = terminal.Stdio{
				In:c.Tty(), 
				Out:c.Tty(),
				Err:c.Tty(),
			}
			defer c.Close()

			go func() {
				c.ExpectString("Are you going to process or abort this input: message?")
				c.SendLine("abort")
				c.ExpectEOF()
			}()

			rootCmd.SetArgs([]string{"job", "input", jobName, "1"})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal("Only process or abort is accepted!\n"))
		})
	})
})
