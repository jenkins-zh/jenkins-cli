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
		updateMgr    UpdateCenterManager
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		pluginMgr = PluginManager{}
		pluginMgr.RoundTripper = roundTripper
		pluginMgr.URL = "http://localhost"
		updateMgr = UpdateCenterManager{}
		updateMgr.RoundTripper = roundTripper
		updateMgr.URL = "http://localhost"
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
			PrepareForEmptyAvaiablePluginList(roundTripper, pluginMgr.URL)

			pluginList, err := pluginMgr.GetAvailablePlugins()
			Expect(err).To(BeNil())
			Expect(pluginList).NotTo(BeNil())
			Expect(len(pluginList.Data)).To(Equal(0))
		})

		It("one plugin in the list", func() {
			PrepareForOneAvaiablePlugin(roundTripper, pluginMgr.URL)

			pluginList, err := pluginMgr.GetAvailablePlugins()
			Expect(err).To(BeNil())
			Expect(pluginList).NotTo(BeNil())
			Expect(len(pluginList.Data)).To(Equal(1))
			Expect(pluginList.Data[0].Name).To(Equal("fake"))
		})

		It("many plugins in the list", func() {
			PrepareForManyAvaiablePlugin(roundTripper, pluginMgr.URL)

			pluginList, err := pluginMgr.GetAvailablePlugins()
			Expect(err).To(BeNil())
			Expect(pluginList).NotTo(BeNil())
			Expect(len(pluginList.Data)).To(Equal(6))
			Expect(pluginList.Data[0].Name).To(Equal("fake-ocean"))
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
		It("no plugin in the list", func() {
			PrepareForEmptyInstalledPluginList(roundTripper, pluginMgr.URL)

			pluginList, err := pluginMgr.GetPlugins()
			Expect(err).To(BeNil())
			Expect(pluginList).NotTo(BeNil())
			Expect(len(pluginList.Plugins)).To(Equal(0))
		})

		It("one plugin in the list", func() {
			PrepareForOneInstalledPlugin(roundTripper, pluginMgr.URL)

			pluginList, err := pluginMgr.GetPlugins()
			Expect(err).To(BeNil())
			Expect(pluginList).NotTo(BeNil())
			Expect(len(pluginList.Plugins)).To(Equal(1))
			Expect(pluginList.Plugins[0].ShortName).To(Equal("fake"))
		})

		It("response with 500", func() {
			PrepareFor500InstalledPluginList(roundTripper, pluginMgr.URL)

			_, err := pluginMgr.GetPlugins()
			Expect(err).To(HaveOccurred())
		})

		It("test with parameter", func() {
			PrepareForManyInstalledPlugins(roundTripper, pluginMgr.URL,2)

			pluginList, err := pluginMgr.GetPlugins(2)
			Expect(err).To(BeNil())
			Expect(pluginList).NotTo(BeNil())
			Expect(len(pluginList.Plugins[0].Dependencies)).To(Equal(1))
		})
	})

	Context("InstallPlugin", func() {
		var (
			pluginName string
		)

		BeforeEach(func() {
			pluginName = "fake"
		})

		It("normal case, should success", func() {
			PrepareForInstallPlugin(roundTripper, pluginMgr.URL, pluginName, "", "")

			err := pluginMgr.InstallPlugin([]string{pluginName})
			Expect(err).To(BeNil())
		})

		It("with 400", func() {
			PrepareForInstallPluginWithCode(roundTripper, 400, pluginMgr.URL, pluginName, "", "")

			err := pluginMgr.InstallPlugin([]string{pluginName})
			Expect(err).NotTo(BeNil())
		})

		It("with 400, error message", func() {
			response := PrepareForInstallPluginWithCode(roundTripper, 400, pluginMgr.URL, pluginName, "", "")
			response.Header = map[string][]string{
				"X-Error": []string{"X-Error"},
			}

			err := pluginMgr.InstallPlugin([]string{pluginName})
			Expect(err).To(Equal(fmt.Errorf("X-Error")))
		})
	})

	Context("UninstallPlugin", func() {
		var (
			pluginName string
		)

		BeforeEach(func() {
			pluginName = "fake"
		})

		It("normal case, should success", func() {
			PrepareForUninstallPlugin(roundTripper, pluginMgr.URL, pluginName)

			err := pluginMgr.UninstallPlugin(pluginName)
			Expect(err).To(BeNil())
		})

		It("response with 500", func() {
			PrepareForUninstallPluginWith500(roundTripper, pluginMgr.URL, pluginName)

			err := pluginMgr.UninstallPlugin(pluginName)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("Upload", func() {
		It("normal case, should success", func() {
			tmpfile, err := ioutil.TempFile("", "example")
			Expect(err).To(BeNil())

			PrepareForUploadPlugin(roundTripper, pluginMgr.URL)

			err = pluginMgr.Upload(tmpfile.Name())
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("UpdateCenter", func() {
		It("normal case, should success", func() {
			PrepareForRequestUpdateCenter(roundTripper, pluginMgr.URL)

			site, err := updateMgr.GetSite()
			Expect(err).To(BeNil())
			Expect(site).NotTo(BeNil())
			Expect(site.ID).To(Equal("default"))
		})
	})

	Context("NullUpdateCenter", func() {
		It("normal case, should success", func() {
			PrepareForNoAvailablePlugins(roundTripper, pluginMgr.URL)

			site, err := updateMgr.GetSite()
			Expect(err).To(BeNil())
			Expect(site).NotTo(BeNil())
			Expect(site.ID).To(Equal("default"))
		})
	})

	Context("ManyInstalledPlugins", func() {
		It("normal case, should success", func() {
			PrepareForManyInstalledPlugins(roundTripper, pluginMgr.URL)

			pluginList, err := pluginMgr.GetPlugins()
			Expect(err).To(BeNil())
			Expect(pluginList).NotTo(BeNil())
			Expect(len(pluginList.Plugins)).To(Equal(4))
			Expect(pluginList.Plugins[0].ShortName).To(Equal("fake-ocean"))
		})
	})

	Context("500UpdateCenter", func() {
		It("normal case, should success", func() {
			PrepareForRequest500UpdateCenter(roundTripper, pluginMgr.URL)

			_, err := updateMgr.GetSite()
			Expect(err).To(HaveOccurred())
		})
	})

	Context("CheckUpdate", func() {
		It("normal case, should success", func() {
			PrepareCheckUpdate(roundTripper, pluginMgr.URL, "", "")

			err := pluginMgr.CheckUpdate(func(_ *http.Response) {
				// do nothing
			})
			Expect(err).To(BeNil())
		})
	})
})
