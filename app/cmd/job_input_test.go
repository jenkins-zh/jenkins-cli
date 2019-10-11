package cmd

import (
	"bytes"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"

	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
	// "github.com/AlecAivazis/survey/v2/core"
	// "github.com/AlecAivazis/survey/v2/terminal"
)

var _ = Describe("job input command", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
		//jenkinsRoot string
		//username string
		//token string
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		jobInputOption.RoundTripper = roundTripper
		rootCmd.SetArgs([]string{})
		rootOptions.Jenkins = ""
		rootOptions.ConfigFile = "test.yaml"

		//jenkinsRoot = "http://localhost:8080/jenkins"
		//username = "admin"
		//token = "111e3a2f0231198855dceaff96f20540a9"
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

		//It("should success, abort without inputs", func() {
		//	data, err := generateSampleConfig()
		//	Expect(err).To(BeNil())
		//	err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
		//	Expect(err).To(BeNil())
		//
		//	jobName := "test"
		//	buildID := 1
		//
		//	client.PrepareForGetJobInputActions(roundTripper, jenkinsRoot, username, token, jobName, buildID)
		//	client.PrepareForSubmitInput(roundTripper, jenkinsRoot, fmt.Sprintf("/job/%s", jobName) , username, token)
		//
		//	// no idea how to let it works, just leave this here
		//	// _, w, err := os.Pipe()
		//
		//	// c, err := expect.NewConsole(expect.WithStdout(w))
		//	// Expect(err).To(BeNil())
		//	// jobInputOption.Stdio = terminal.Stdio{
		//	// 	In:c.Tty(),
		//	// 	Out:c.Tty(),
		//	// 	Err:c.Tty(),
		//	// }
		//	// defer c.Close()
		//
		//	// go func() {
		//	// 	c.ExpectString("Are you going to process or abort this input: message?")
		//	// 	c.SendLine("abort\n")
		//	// 	c.ExpectEOF()
		//	// }()
		//
		//	rootCmd.SetArgs([]string{"job", "input", jobName, "1", "--action", "abort"})
		//
		//	buf := new(bytes.Buffer)
		//	rootCmd.SetOutput(buf)
		//	_, err = rootCmd.ExecuteC()
		//	Expect(err).To(BeNil())
		//
		//	Expect(buf.String()).To(Equal(""))
		//})

		//It("should success, process without inputs", func() {
		//	data, err := generateSampleConfig()
		//	Expect(err).To(BeNil())
		//	err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
		//	Expect(err).To(BeNil())
		//
		//	jobName := "test"
		//	buildID := 1
		//
		//	client.PrepareForGetJobInputActions(roundTripper, jenkinsRoot, username, token, jobName, buildID)
		//	client.PrepareForSubmitProcessInput(roundTripper, jenkinsRoot, fmt.Sprintf("/job/%s", jobName) , username, token)
		//
		//	rootCmd.SetArgs([]string{"job", "input", jobName, "1", "--action", "process"})
		//
		//	buf := new(bytes.Buffer)
		//	rootCmd.SetOutput(buf)
		//	_, err = rootCmd.ExecuteC()
		//	Expect(err).To(BeNil())
		//
		//	Expect(buf.String()).To(Equal(""))
		//})
	})
})
