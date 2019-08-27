package cmd

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Table util test", func() {
	var (
		ctrl *gomock.Controller
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		config = nil
	})

	AfterEach(func() {
		config = nil
		ctrl.Finish()
	})

	Context("basic test", func() {
		It("getJenkinsNames", func() {
			config = &Config{
				JenkinsServers: []JenkinsServer{JenkinsServer{
					Name: "a",
				}, JenkinsServer{
					Name: "b",
				}},
			}

			names := getJenkinsNames()
			Expect(names).To(Equal([]string{"a", "b"}))

			config.JenkinsServers = []JenkinsServer{}
			names = getJenkinsNames()
			Expect(names).To(Equal([]string{}))
		})

		It("getCurrentJenkins", func() {
			config = &Config{}
			current := getCurrentJenkins()
			Expect(current).To(BeNil())

			config.Current = "test"
			config.JenkinsServers = []JenkinsServer{JenkinsServer{
				Name: "test",
			}}
			current = getCurrentJenkins()
			Expect(current).To(Equal(&config.JenkinsServers[0]))
		})

		It("findSuiteByName", func() {
			config = &Config{}
			suite := findSuiteByName("fake")
			Expect(suite).To(BeNil())

			pluginName := "plugin-one"
			config.PluginSuites = []PluginSuite{PluginSuite{
				Name: pluginName,
			}}
			suite = findSuiteByName(pluginName)
			Expect(suite).NotTo(BeNil())
			Expect(suite.Name).To(Equal(pluginName))
		})
	})
})
