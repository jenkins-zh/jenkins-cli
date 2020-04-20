package cmd

import (
	"bytes"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"

	"github.com/golang/mock/gomock"
	. "github.com/jenkins-zh/jenkins-cli/app/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("config remove command", func() {
	var (
		ctrl         *gomock.Controller
		buf          *bytes.Buffer
		configPath   string
		err          error
		otherJenkins string
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		rootCmd.SetArgs([]string{})
		buf = new(bytes.Buffer)
		rootCmd.SetOutput(buf)
		rootOptions.Jenkins = ""
		configPath = path.Join(os.TempDir(), "fake.yaml")

		otherJenkins = "other-jenkins"

		var data []byte
		sampleConfig := getSampleConfig()
		sampleConfig.JenkinsServers = append(sampleConfig.JenkinsServers, JenkinsServer{
			Name: otherJenkins,
		})
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

		It("remove a not exist jenkins", func() {
			rootCmd.SetArgs([]string{"config", "remove", "fake", "--configFile", configPath})
			_, err = rootCmd.ExecuteC()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("cannot found by name fake"))
		})

		It("remove the current jenkins", func() {
			rootCmd.SetArgs([]string{"config", "remove", "yourServer", "--configFile", configPath})
			_, err = rootCmd.ExecuteC()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("you cannot remove current Jenkins, if you want to remove it, can select other items before"))
		})

		It("should success", func() {
			rootCmd.SetArgs([]string{"config", "remove", otherJenkins, "--configFile", configPath})
			_, err = rootCmd.ExecuteC()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
