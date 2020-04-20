package cmd

import (
	"bytes"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
)

var _ = Describe("config select command", func() {
	var (
		ctrl       *gomock.Controller
		buf        *bytes.Buffer
		configPath string
		err        error
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		rootCmd.SetArgs([]string{})
		buf = new(bytes.Buffer)
		rootCmd.SetOutput(buf)
		rootOptions.Jenkins = ""
		configPath = path.Join(os.TempDir(), "fake.yaml")

		var data []byte
		sampleConfig := getSampleConfig()
		data, err = yaml.Marshal(&sampleConfig)
		Expect(err).To(BeNil())
		err = ioutil.WriteFile(configPath, data, 0664)
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		rootCmd.SetArgs([]string{})
		os.Remove(configPath)
		ctrl.Finish()
	})

	Context("basic cases", func() {
		var (
			err error
		)

		It("select a config", func() {
			rootCmd.SetArgs([]string{"config", "select", "yourServer", "fake", "--configFile", configPath})
			_, err = rootCmd.ExecuteC()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})

//func TestConfigSelect(t *testing.T) {
//	RunEditCommandTest(t, EditCommandTest{
//		ConfirmProcedure: func(c *expect.Console) {
//			c.ExpectString("Choose a Jenkins as the current one:")
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
//			configSelectOptions.CommonOption.Stdio = stdio
//			rootCmd.SetArgs([]string{"config", "select", "--configFile", configFile})
//			_, err = rootCmd.ExecuteC()
//			return
//		},
//	})
//}
