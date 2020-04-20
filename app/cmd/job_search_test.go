package cmd

import (
	"bytes"
	"github.com/jenkins-zh/jenkins-cli/client"
	"io/ioutil"
	"os"

	"github.com/golang/mock/gomock"
	. "github.com/jenkins-zh/jenkins-cli/app/config"
	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("job search command", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		jobSearchOption.RoundTripper = roundTripper
		rootCmd.SetArgs([]string{})
		rootOptions.Jenkins = ""
		rootOptions.ConfigFile = "test.yaml"
	})

	AfterEach(func() {
		rootCmd.SetArgs([]string{})
		os.Remove(rootOptions.ConfigFile)
		rootOptions.ConfigFile = ""
		ctrl.Finish()
	})

	Context("basic cases, need roundtrippper", func() {
		BeforeEach(func() {
			roundTripper = mhttp.NewMockRoundTripper(ctrl)
			jobSearchOption.RoundTripper = roundTripper
		})

		It("should success, search with one result item", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			name := "fake"
			kind := "fake"

			client.PrepareOneItem(roundTripper, "http://localhost:8080/jenkins", name, kind,
				"admin", "111e3a2f0231198855dceaff96f20540a9")

			rootCmd.SetArgs([]string{"job", "search", name, "--type", kind})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal(`Name DisplayName Type        URL
fake fake        WorkflowJob job/fake/
`))
		})

		It("should success, search without keyword", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			client.PrepareOneItem(roundTripper, "http://localhost:8080/jenkins", "", "",
				"admin", "111e3a2f0231198855dceaff96f20540a9")

			rootCmd.SetArgs([]string{"job", "search", "--type=", "--name="})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal(`Name DisplayName Type        URL
fake fake        WorkflowJob job/fake/
`))
		})
	})
})

var _ = Describe("job search command check", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
		rootURL      string
		user         string
		password     string
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		jobSearchOption.RoundTripper = roundTripper
		rootOptions.Jenkins = ""
		rootOptions.ConfigFile = "test.yaml"
		rootURL = "http://localhost:8080/jenkins"
		user = "admin"
		password = "111e3a2f0231198855dceaff96f20540a9"

		config = &Config{
			Current: "fake",
			JenkinsServers: []JenkinsServer{{
				Name:     "fake",
				URL:      rootURL,
				UserName: user,
				Token:    password,
			}},
		}
	})

	AfterEach(func() {
		config = nil
		rootCmd.SetArgs([]string{})
		os.Remove(rootOptions.ConfigFile)
		rootOptions.ConfigFile = ""
		ctrl.Finish()
	})

	Context("basic cases", func() {
		It("without pipeline-restful-api plugin", func() {
			req, _ := client.PrepareForOneInstalledPlugin(roundTripper, rootURL)
			req.SetBasicAuth(user, password)

			err := jobSearchOption.Check()
			Expect(err).To(HaveOccurred())
		})

		It("with pipeline-restful-api 0.2+ plugin", func() {
			req, _ := client.PrepareForOneInstalledPluginWithPluginNameAndVer(roundTripper, rootURL,
				"pipeline-restful-api", "1.0")
			req.SetBasicAuth(user, password)

			err := jobSearchOption.Check()
			Expect(err).NotTo(HaveOccurred())
		})

		It("with pipeline-restful-api 0.2- plugin", func() {
			req, _ := client.PrepareForOneInstalledPluginWithPluginNameAndVer(roundTripper, rootURL,
				"pipeline-restful-api", "0.1")
			req.SetBasicAuth(user, password)

			err := jobSearchOption.Check()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("should be"))
		})
	})
})
