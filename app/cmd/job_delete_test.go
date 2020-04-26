package cmd

import (
	"bytes"
	"fmt"
	"github.com/Netflix/go-expect"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/jenkins-zh/jenkins-cli/client"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/golang/mock/gomock"
	"github.com/hinshun/vt10x"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/require"

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
				RoundTrip(client.NewRequestMatcher(request)).Return(response, nil)

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
				RoundTrip(client.NewRequestMatcher(requestCrumb)).Return(responseCrumb, nil)

			rootCmd.SetArgs([]string{"job", "delete", jobName, "-b", "true"})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal(""))
		})
	})
})

type PromptCommandTest struct {
	Message     string
	MsgConfirm  common.MsgConfirm
	BatchOption *common.BatchOption
	Procedure   func(*expect.Console)
	Args        []string
}

type EditCommandTest struct {
	Message          string
	DefaultContent   string
	EditContent      common.EditContent
	CommonOption     *common.CommonOption
	BatchOption      *common.BatchOption
	ConfirmProcedure func(*expect.Console)
	Procedure        func(*expect.Console)
	Test             func(stdio terminal.Stdio) (err error)
	Expected         string
	Args             []string
}

type PromptTest struct {
	Message    string
	MsgConfirm common.MsgConfirm
	Procedure  func(*expect.Console)
	Expected   interface{}
}

type EditorTest struct {
	Message        string
	DefaultContent string
	EditContent    common.EditContent
	Procedure      func(*expect.Console)
	Expected       string
}

//func RunPromptCommandTest(t *testing.T, test PromptCommandTest) {
//	RunTest(t, func(stdio terminal.Stdio) (err error) {
//		var data []byte
//		data, err = generateSampleConfig()
//		err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
//
//		test.BatchOption.Stdio = stdio
//		rootOptions.ConfigFile = "test.yaml"
//		rootCmd.SetArgs(test.Args)
//		_, err = rootCmd.ExecuteC()
//		return
//	}, test.Procedure)
//}

func RunPromptTest(t *testing.T, test PromptTest) {
	var answer interface{}
	RunTest(t, func(stdio terminal.Stdio) error {
		batch := &common.BatchOption{
			Batch: false,
			Stdio: stdio,
		}
		answer = batch.Confirm(test.Message)
		return nil
	}, test.Procedure)
	require.Equal(t, test.Expected, answer)
}

func RunEditorTest(t *testing.T, test EditorTest) {
	var content string
	RunTest(t, func(stdio terminal.Stdio) (err error) {
		editor := &common.CommonOption{
			Stdio: stdio,
		}
		content, err = editor.Editor(test.DefaultContent, test.Message)
		return nil
	}, test.Procedure)
	require.Equal(t, test.Expected, content)
}

func Stdio(c *expect.Console) terminal.Stdio {
	return terminal.Stdio{In: c.Tty(), Out: c.Tty(), Err: c.Tty()}
}

func RunTest(t *testing.T, test func(terminal.Stdio) error, procedures ...func(*expect.Console)) {
	//t.Parallel()

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
		for _, procedure := range procedures {
			if procedure != nil {
				procedure(c)
			}
		}
	}()

	err = test(Stdio(c))
	//fmt.Println("Raw output: ", buf.String())
	require.Nil(t, err)

	// Close the slave end of the pty, and read the remaining bytes from the master end.
	c.Tty().Close()
	<-donec

	// Dump the terminal's screen.
	//fmt.Sprintf("\n%s", expect.StripTrailingEmptyLines(state.String()))
}
