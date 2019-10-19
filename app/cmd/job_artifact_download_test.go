package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/golang/mock/gomock"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("job artifact download command", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
		buildID      int
		jobName      string
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		rootCmd.SetArgs([]string{})
		rootOptions.Jenkins = ""
		rootOptions.ConfigFile = "test.yaml"
		buildID = 1
		jobName = "fakeJob"

		jobArtifactDownloadOption.RoundTripper = roundTripper
	})

	AfterEach(func() {
		rootCmd.SetArgs([]string{})
		os.Remove(rootOptions.ConfigFile)
		rootOptions.ConfigFile = ""
		ctrl.Finish()
	})

	Context("basic cases", func() {
		It("invalid build id", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)

			rootCmd.SetArgs([]string{"job", "artifact", "download", jobName, "fakeid"})
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())
			Expect(buf.String()).To(Equal("strconv.Atoi: parsing \"fakeid\": invalid syntax\n"))
		})

		It("should success", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			client.PrepareGetArtifacts(roundTripper, "http://localhost:8080/jenkins", "admin", "111e3a2f0231198855dceaff96f20540a9", jobName, buildID)

			request, _ := http.NewRequest("GET", "http://localhost:8080/jenkins/job/pipeline/1/artifact/a.log", nil)
			request.SetBasicAuth("admin", "111e3a2f0231198855dceaff96f20540a9")
			response := &http.Response{
				StatusCode: 200,
				Request:    request,
				Body:       ioutil.NopCloser(bytes.NewBufferString("")),
			}
			roundTripper.EXPECT().
				RoundTrip(request).Return(response, nil)

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)

			// store artifact files
			tmpdir, err := ioutil.TempDir("", "test-gen-cmd-tree")
			Expect(err).To(BeNil())
			defer os.RemoveAll(tmpdir)

			rootCmd.SetArgs([]string{"job", "artifact", "download", jobName, fmt.Sprintf("%d", buildID),
				"--progress=false", "--download-dir", tmpdir})
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			_, err = os.Stat(filepath.Join(tmpdir, "a.log"))
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal(""))
		})

		It("should success, fake artifact id", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			client.PrepareGetArtifacts(roundTripper, "http://localhost:8080/jenkins", "admin", "111e3a2f0231198855dceaff96f20540a9", jobName, buildID)

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)

			rootCmd.SetArgs([]string{"job", "artifact", "download", jobName, fmt.Sprintf("%d", buildID),
				"--progress=false", "--id", "fakeid"})
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal(""))
		})
	})
})
