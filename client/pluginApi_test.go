package client

import (
	"github.com/golang/mock/gomock"
	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("plugin api test", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
		pluginAPI    PluginAPI
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		pluginAPI = PluginAPI{
			RoundTripper: roundTripper,
		}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("ShowTrend", func() {
		It("basic case", func() {
			keyword := "fake"

			PrepareShowTrend(roundTripper, keyword)

			trend, err := pluginAPI.ShowTrend(keyword)
			Expect(err).To(BeNil())
			Expect(trend).NotTo(Equal(""))
		})
	})
})
