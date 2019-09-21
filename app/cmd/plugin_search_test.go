package cmd

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
)

var _ = Describe("plugin search command", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		pluginSearchOption.RoundTripper = roundTripper
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

	Context("basic cases", func() {
		It("should success, empty list", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			request, _ := client.PrepareForEmptyAvaiablePluginList(roundTripper, "http://localhost:8080/jenkins")
			request.SetBasicAuth("admin", "111e3a2f0231198855dceaff96f20540a9")
			rootCmd.SetArgs([]string{"plugin", "search", "fake"})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal(""))
		})

		It("many plugins in the list", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			request, _ := client.PrepareForManyAvaiablePlugin(roundTripper, "http://localhost:8080/jenkins")
			request.SetBasicAuth("admin", "111e3a2f0231198855dceaff96f20540a9")
			request, _ = client.RequestUpdateCenter(roundTripper, "http://localhost:8080/jenkins")
			request.SetBasicAuth("admin", "111e3a2f0231198855dceaff96f20540a9")
			request, _ = client.PrepareForOneInstalledPlugin(roundTripper, "http://localhost:8080/jenkins")
			request.SetBasicAuth("admin", "111e3a2f0231198855dceaff96f20540a9")

			rootCmd.SetArgs([]string{"plugin", "search", "fake"})
			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal(`number name       installed version  installedVersion title
0      fake-ocean true      1.19.011 1.18.111         fake-ocean
1      fake-ln    true      1.19.011 1.18.1           fake-ln
2      fake-is    true      1.19.1   1.18.111         fake-is
3      fake-oa    false     1.13.011                  fake-oa
4      fake-open  false     1.13.0                    fake-open
5      fake       true               1.0              fake
`))
		})

		It("should success, empty updateCenter list", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			request, _ := client.PrepareForManyAvaiablePlugin(roundTripper, "http://localhost:8080/jenkins")
			request.SetBasicAuth("admin", "111e3a2f0231198855dceaff96f20540a9")
			request, _ = client.NoAvailablePlugins(roundTripper, "http://localhost:8080/jenkins")
			request.SetBasicAuth("admin", "111e3a2f0231198855dceaff96f20540a9")
			request, _ = client.PrepareForManyInstalledPlugin(roundTripper, "http://localhost:8080/jenkins")
			request.SetBasicAuth("admin", "111e3a2f0231198855dceaff96f20540a9")
			rootCmd.SetArgs([]string{"plugin", "search", "fake"})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal(`number name       installed version installedVersion title
0      fake-ocean true              1.18.111         fake-ocean
1      fake-ln    true              1.18.1           fake-ln
2      fake-is    true              1.18.111         fake-is
3      fake-oa    false                              fake-oa
4      fake-open  false                              fake-open
5      fake       true              1.0              fake
`))
		})

		It("should success, null updateCenter and 500 installed list", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			request, _ := client.PrepareForManyAvaiablePlugin(roundTripper, "http://localhost:8080/jenkins")
			request.SetBasicAuth("admin", "111e3a2f0231198855dceaff96f20540a9")
			request, _ = client.Request500UpdateCenter(roundTripper, "http://localhost:8080/jenkins")
			request.SetBasicAuth("admin", "111e3a2f0231198855dceaff96f20540a9")
			request, _ = client.PrepareFor500InstalledPluginList(roundTripper, "http://localhost:8080/jenkins")
			request.SetBasicAuth("admin", "111e3a2f0231198855dceaff96f20540a9")
			rootCmd.SetArgs([]string{"plugin", "search", "fake"})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal(`number name       installed version installedVersion title
0      fake-ocean false                              fake-ocean
1      fake-ln    false                              fake-ln
2      fake-is    false                              fake-is
3      fake-oa    false                              fake-oa
4      fake-open  false                              fake-open
5      fake       false                              fake
`))
		})
	})
})
