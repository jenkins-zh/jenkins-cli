package cmd

import (
	"bytes"
	"fmt"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/Netflix/go-expect"
	"github.com/golang/mock/gomock"
	"github.com/hinshun/vt10x"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
	"github.com/jenkins-zh/jenkins-cli/util"
)

var _ = Describe("job delete command", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
		err          error
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		jobDeleteOption.RoundTripper = roundTripper
		rootCmd.SetArgs([]string{})
		rootOptions.Jenkins = ""
		rootOptions.ConfigFile = "test.yaml"
	})

	AfterEach(func() {
		rootCmd.SetArgs([]string{})
		err = os.Remove(rootOptions.ConfigFile)
		rootOptions.ConfigFile = ""
		ctrl.Finish()
	})

	Context("basic cases", func() {
		It("should not error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("should success, with batch mode", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			jobName := "fakeJob"
			request, _ := http.NewRequest("POST", fmt.Sprintf("http://localhost:8080/jenkins/job/%s/doDelete", jobName), nil)
			request.Header.Add("CrumbRequestField", "Crumb")
			request.SetBasicAuth("admin", "111e3a2f0231198855dceaff96f20540a9")
			request.Header.Add(util.ContentType, util.ApplicationForm)
			response := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    request,
				Body:       ioutil.NopCloser(bytes.NewBufferString("")),
			}
			roundTripper.EXPECT().
				RoundTrip(request).Return(response, nil)

			requestCrumb, _ := http.NewRequest("GET", "http://localhost:8080/jenkins/crumbIssuer/api/json", nil)
			requestCrumb.SetBasicAuth("admin", "111e3a2f0231198855dceaff96f20540a9")
			responseCrumb := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    requestCrumb,
				Body: ioutil.NopCloser(bytes.NewBufferString(`
				{"crumbRequestField":"CrumbRequestField","crumb":"Crumb"}
				`)),
			}
			roundTripper.EXPECT().
				RoundTrip(requestCrumb).Return(responseCrumb, nil)

			rootCmd.SetArgs([]string{"job", "delete", jobName, "-b", "true"})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal(""))
		})
	})
})

func TestDeleteJob(t *testing.T) {
	RunPromptCommandTest(t, PromptCommandTest{
		Args: []string{"job", "delete", "fake", "-b=false"},
		Procedure: func(c *expect.Console) {
			c.ExpectString("Are you sure to delete job fake ?")
			c.SendLine("n")
			c.ExpectEOF()
		},
		BatchOption: &jobDeleteOption.BatchOption,
		Expected:    nil,
	})
}

type PromptCommandTest struct {
	Message     string
	MsgConfirm  MsgConfirm
	BatchOption *BatchOption
	Procedure   func(*expect.Console)
	Args        []string
	Expected    interface{}
}

type PromptTest struct {
	Message    string
	MsgConfirm MsgConfirm
	Procedure  func(*expect.Console)
	Expected   interface{}
}

func RunPromptCommandTest(t *testing.T, test PromptCommandTest) {
	RunTest(t, test.Procedure, func(stdio terminal.Stdio) (err error) {
		var data []byte
		data, err = generateSampleConfig()
		err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)

		test.BatchOption.Stdio = stdio
		rootCmd.SetArgs(test.Args)
		_, err = rootCmd.ExecuteC()
		return
	})
}

func RunPromptTest(t *testing.T, test PromptTest) {
	var answer interface{}
	RunTest(t, test.Procedure, func(stdio terminal.Stdio) error {
		batch := &BatchOption{
			Batch: false,
			Stdio: stdio,
		}
		answer = batch.Confirm(test.Message)
		return nil
	})
	require.Equal(t, test.Expected, answer)
}

func Stdio(c *expect.Console) terminal.Stdio {
	return terminal.Stdio{In: c.Tty(), Out: c.Tty(), Err: c.Tty()}
}

func RunTest(t *testing.T, procedure func(*expect.Console), test func(terminal.Stdio) error) {
	t.Parallel()

	// Multiplex output to a buffer as well for the raw bytes.
	buf := new(bytes.Buffer)

	//c, err := expect.NewConsole(expect.WithStdout(buf))
	//c, err := expect.NewConsole(expect.WithStdout(os.Stdout))
	c, _, err := vt10x.NewVT10XConsole(expect.WithStdout(buf))

	require.Nil(t, err)
	defer c.Close()

	donec := make(chan struct{})
	go func() {
		defer close(donec)
		procedure(c)
	}()

	err = test(Stdio(c))
	fmt.Println("Raw output: ", buf.String())
	require.Nil(t, err)

	// Close the slave end of the pty, and read the remaining bytes from the master end.
	c.Tty().Close()
	<-donec

	// Dump the terminal's screen.
	//fmt.Sprintf("\n%s", expect.StripTrailingEmptyLines(state.String()))
}
