package cmd

import (
	"bytes"
	"github.com/jenkins-zh/jenkins-cli/app"
	"github.com/jenkins-zh/jenkins-cli/client"
	"io/ioutil"
	"os"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("version command", func() {
	var (
		ctrl *gomock.Controller
		buf  *bytes.Buffer
		err  error
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		rootCmd.SetArgs([]string{})
		rootOptions.Jenkins = ""
		rootOptions.ConfigFile = "test.yaml"
		buf = new(bytes.Buffer)
		rootCmd.SetOutput(buf)

		var data []byte
		data, err = generateSampleConfig()
		Expect(err).To(BeNil())
		err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		rootCmd.SetArgs([]string{})
		os.Remove(rootOptions.ConfigFile)
		rootOptions.ConfigFile = ""
		ctrl.Finish()
	})

	Context("normal case", func() {
		It("fake jenkins", func() {
			rootCmd.SetArgs([]string{"version", "--jenkins", "fakeJenkins"})
			_, err = rootCmd.ExecuteC()
			Expect(err).To(HaveOccurred())
			Expect(buf.String()).To(ContainSubstring("cannot found the configuration: fakeJenkins"))
		})

		It("should success", func() {
			rootCmd.SetArgs([]string{"version", "--jenkins", "yourServer"})
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())
			Expect(buf.String()).To(ContainSubstring("Current Jenkins is:"))
			Expect(buf.String()).To(ContainSubstring(`Version: 
Commit: 
`))
		})

		It("Output changelog", func() {
			ghClient, teardown := client.PrepareForGetJCLIAsset("v0.0.1")
			defer teardown()

			app.SetVersion("dev-v0.0.1")

			versionOption.GitHubClient = ghClient
			rootCmd.SetArgs([]string{"version", "--changelog"})
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())
			Expect(buf.String()).To(ContainSubstring("body"))
		})

		It("Show the latest", func() {
			ghClient, teardown := client.PrepareForGetLatestJCLIAsset()
			defer teardown()

			versionOption.GitHubClient = ghClient

			rootCmd.SetArgs([]string{"version", "--changelog=false", "--show-latest"})
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())
			Expect(buf.String()).To(ContainSubstring(`tagName
body`))
		})
	})
})
