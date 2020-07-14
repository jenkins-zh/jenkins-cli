package version_test

import (
	"bytes"
	"github.com/jenkins-zh/jenkins-cli/app"
	"github.com/jenkins-zh/jenkins-cli/app/cmd"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("version command", func() {
	var (
		ctrl        *gomock.Controller
		buf         *bytes.Buffer
		tempFile    *os.File
		err         error
		rootCmd     *cobra.Command
		rootOptions *cmd.RootOptions
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())

		tempFile, err = ioutil.TempFile("", "test.yaml")
		Expect(err).NotTo(HaveOccurred())

		rootCmd = cmd.GetRootCommand()
		rootOptions = cmd.GetRootOptions()

		rootOptions.ConfigFile = tempFile.Name()
		buf = new(bytes.Buffer)
		rootCmd.SetOutput(buf)

		var data []byte
		data, err = cmd.GenerateSampleConfig()
		Expect(err).To(BeNil())
		err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		os.Remove(rootOptions.ConfigFile)
		rootOptions.ConfigFile = ""
		ctrl.Finish()
	})

	Context("normal case", func() {
		It("Output changelog", func() {
			ghClient, teardown := client.PrepareForGetJCLIAsset("v0.0.1")
			defer teardown()

			app.SetVersion("dev-v0.0.1")

			rootOptions.SetGitHubClient(ghClient)

			rootCmd.SetArgs([]string{"version", "--changelog"})
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())
			Expect(buf.String()).To(ContainSubstring("body"))
		})

		It("Show the latest", func() {
			ghClient, teardown := client.PrepareForGetLatestJCLIAsset()
			defer teardown()

			rootOptions.SetGitHubClient(ghClient)

			rootCmd.SetArgs([]string{"version", "--changelog=false", "--show-latest"})
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())
			Expect(buf.String()).To(ContainSubstring(`tagName
body`))
		})
	})
})
