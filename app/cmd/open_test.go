package cmd

import (
	"bytes"
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"path"
)

var _ = Describe("test open", func() {
	var (
		err         error
		jenkinsName string
		configFile  string
		cmdArgs     []string
	)

	BeforeEach(func() {
		configFile = path.Join(os.TempDir(), "fake.yaml")

		data, err := generateSampleConfig()
		Expect(err).To(BeNil())
		err = ioutil.WriteFile(configFile, data, 0664)
		Expect(err).To(BeNil())
		openOption.ExecContext = util.FakeExecCommandSuccess
		jenkinsName = "fake"

		cmdArgs = []string{"open", jenkinsName, "--configFile", configFile}
	})

	AfterEach(func() {
		os.Remove(configFile)
	})

	JustBeforeEach(func() {
		buf := new(bytes.Buffer)
		rootCmd.SetOut(buf)
		rootCmd.SetArgs(cmdArgs)
		_, err = rootCmd.ExecuteC()
	})

	It("open a not exists Jenkins", func() {
		Expect(err).To(HaveOccurred())
		Expect(fmt.Sprint(err)).To(ContainSubstring("no URL found with Jenkins " + jenkinsName))
	})

	Context("give a right config", func() {
		BeforeEach(func() {
			cmdArgs = []string{"open", "yourServer", "--configFile", configFile}
		})

		It("should success", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			os.Remove(rootOptions.ConfigFile)
		})
	})

	Context("open the current jenkins with config page", func() {
		BeforeEach(func() {
			cmdArgs = []string{"open", "--interactive=false", "--config", "--configFile", configFile}
		})

		It("should success", func() {
			Expect(err).NotTo(HaveOccurred())
		})
	})
})

//func TestOpenJenkins(t *testing.T) {
//	RunEditCommandTest(t, EditCommandTest{
//		ConfirmProcedure: func(c *expect.Console) {
//			c.ExpectString("Choose a Jenkins which you want to open:")
//			// filter away everything
//			c.SendLine("z")
//			// send enter (should get ignored since there are no answers)
//			c.SendLine(string(terminal.KeyEnter))
//
//			// remove the filter we just applied
//			c.SendLine(string(terminal.KeyBackspace))
//
//			// press enter
//			c.SendLine(string(terminal.KeyEnter))
//		},
//		Test: func(stdio terminal.Stdio) (err error) {
//			configFile := path.Join(os.TempDir(), "fake.yaml")
//			defer os.Remove(configFile)
//
//			var data []byte
//			data, err = generateSampleConfig()
//			err = ioutil.WriteFile(configFile, data, 0664)
//
//			openOption.ExecContext = util.FakeExecCommandSuccess
//			openOption.CommonOption.Stdio = stdio
//			rootCmd.SetArgs([]string{"open", "--interactive", "--configFile", configFile})
//			_, err = rootCmd.ExecuteC()
//			return
//		},
//	})
//}
