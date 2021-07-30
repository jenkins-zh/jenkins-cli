package cmd

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("config add command", func() {
	var (
		ctrl       *gomock.Controller
		buf        io.Writer
		err        error
		configPath string
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		rootCmd.SetArgs([]string{})
		buf = new(bytes.Buffer)
		rootCmd.SetOutput(buf)
		rootOptions.Jenkins = ""

		configPath = path.Join(os.TempDir(), "fake.yaml")

		var data []byte
		data, err = GenerateSampleConfig()
		Expect(err).To(BeNil())
		err = ioutil.WriteFile(configPath, data, 0664)
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		rootCmd.SetArgs([]string{})
		os.Remove(configPath)
		rootOptions.ConfigFile = ""
		ctrl.Finish()
	})

	Context("basic cases", func() {
		It("lack of name", func() {
			rootCmd.SetArgs([]string{"config", "add", "--configFile", configPath})
			_, err = rootCmd.ExecuteC()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("name cannot be empty"))
		})

		It("add an exist one", func() {
			rootCmd.SetArgs([]string{"config", "add", "--name", "yourServer", "--configFile", configPath})
			_, err = rootCmd.ExecuteC()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("jenkins yourServer is existed"))
		})

		It("should success", func() {
			rootCmd.SetArgs([]string{"config", "add", "--name", "fake", "--configFile", configPath})
			_, err = rootCmd.ExecuteC()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
