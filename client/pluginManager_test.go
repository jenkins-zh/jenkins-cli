package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/golang/mock/gomock"
	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PluginManager test", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
		pluginMgr    PluginManager
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		pluginMgr = PluginManager{}
		pluginMgr.RoundTripper = roundTripper
		pluginMgr.URL = "http://localhost"
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("basic function test", func() {
		It("get install plugin query string", func() {
			names := make([]string, 0)
			Expect(getPluginsInstallQuery(names)).To(Equal(""))

			names = append(names, "abc")
			Expect(getPluginsInstallQuery(names)).To(Equal("plugin.abc="))

			names = append(names, "def")
			Expect(getPluginsInstallQuery(names)).To(Equal("plugin.abc=&plugin.def="))

			names = append(names, "")
			Expect(getPluginsInstallQuery(names)).To(Equal("plugin.abc=&plugin.def="))
		})
	})

	Context("GetAvailablePlugins", func() {
		It("no plugin in the list", func() {
			request, _ := http.NewRequest("GET", fmt.Sprintf("%s/pluginManager/plugins", pluginMgr.URL), nil)
			response := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    request,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
					"status": "ok",
					"data": []
				}`)),
			}
			roundTripper.EXPECT().
				RoundTrip(request).Return(response, nil)

			pluginList, err := pluginMgr.GetAvailablePlugins()
			Expect(err).To(BeNil())
			Expect(pluginList).NotTo(BeNil())
			Expect(len(pluginList.Data)).To(Equal(0))
		})

		It("one plugin in the list", func() {
			request, _ := http.NewRequest("GET", fmt.Sprintf("%s/pluginManager/plugins", pluginMgr.URL), nil)
			response := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    request,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
					"status": "ok",
					"data": [{
						"name": "fake"
					}]
				}`)),
			}
			roundTripper.EXPECT().
				RoundTrip(request).Return(response, nil)

			pluginList, err := pluginMgr.GetAvailablePlugins()
			Expect(err).To(BeNil())
			Expect(pluginList).NotTo(BeNil())
			Expect(len(pluginList.Data)).To(Equal(1))
			Expect(pluginList.Data[0].Name).To(Equal("fake"))
		})

		It("response with 500", func() {
			request, _ := http.NewRequest("GET", fmt.Sprintf("%s/pluginManager/plugins", pluginMgr.URL), nil)
			response := &http.Response{
				StatusCode: 500,
				Proto:      "HTTP/1.1",
				Request:    request,
				Body:       ioutil.NopCloser(bytes.NewBufferString("")),
			}
			roundTripper.EXPECT().
				RoundTrip(request).Return(response, nil)

			_, err := pluginMgr.GetAvailablePlugins()
			Expect(err).To(HaveOccurred())
		})
	})

	Context("GetPlugins", func() {
		var (
			api string
		)

		BeforeEach(func() {
			api = fmt.Sprintf("%s/pluginManager/api/json?depth=1", pluginMgr.URL)
		})

		It("no plugin in the list", func() {
			request, _ := http.NewRequest("GET", api, nil)
			response := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    request,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
					"plugins": []
				}`)),
			}
			roundTripper.EXPECT().
				RoundTrip(request).Return(response, nil)

			pluginList, err := pluginMgr.GetPlugins()
			Expect(err).To(BeNil())
			Expect(pluginList).NotTo(BeNil())
			Expect(len(pluginList.Plugins)).To(Equal(0))
		})

		It("one plugin in the list", func() {
			request, _ := http.NewRequest("GET", api, nil)
			response := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    request,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
					"plugins": [{
						"shortName": "fake"
					}]
				}`)),
			}
			roundTripper.EXPECT().
				RoundTrip(request).Return(response, nil)

			pluginList, err := pluginMgr.GetPlugins()
			Expect(err).To(BeNil())
			Expect(pluginList).NotTo(BeNil())
			Expect(len(pluginList.Plugins)).To(Equal(1))
			Expect(pluginList.Plugins[0].ShortName).To(Equal("fake"))
		})

		It("response with 500", func() {
			request, _ := http.NewRequest("GET", api, nil)
			response := &http.Response{
				StatusCode: 500,
				Proto:      "HTTP/1.1",
				Request:    request,
				Body:       ioutil.NopCloser(bytes.NewBufferString("")),
			}
			roundTripper.EXPECT().
				RoundTrip(request).Return(response, nil)

			_, err := pluginMgr.GetPlugins()
			Expect(err).To(HaveOccurred())
		})
	})

	Context("UninstallPlugin", func() {
		var (
			api        string
			pluginName string
		)

		BeforeEach(func() {
			pluginName = "fake"
			api = fmt.Sprintf("%s/pluginManager/plugin/%s/uninstall", pluginMgr.URL, pluginName)
		})

		It("normal case, should success", func() {
			request, _ := http.NewRequest("POST", api, nil)
			request.Header.Add("CrumbRequestField", "Crumb")
			response := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    request,
				Body:       ioutil.NopCloser(bytes.NewBufferString("")),
			}
			roundTripper.EXPECT().
				RoundTrip(request).Return(response, nil)

			// common crumb request
			RequestCrumb(roundTripper, pluginMgr.URL)

			err := pluginMgr.UninstallPlugin(pluginName)
			Expect(err).To(BeNil())
		})

		It("response with 500", func() {
			request, _ := http.NewRequest("POST", api, nil)
			request.Header.Add("CrumbRequestField", "Crumb")
			response := &http.Response{
				StatusCode: 500,
				Proto:      "HTTP/1.1",
				Request:    request,
				Body:       ioutil.NopCloser(bytes.NewBufferString("")),
			}
			roundTripper.EXPECT().
				RoundTrip(request).Return(response, nil)

			// common crumb request
			RequestCrumb(roundTripper, pluginMgr.URL)

			err := pluginMgr.UninstallPlugin(pluginName)
			Expect(err).To(HaveOccurred())
		})
	})
})
