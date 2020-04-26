package client

import (
	"bytes"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"os"
)

var _ = Describe("update center test", func() {
	var (
		ctrl         *gomock.Controller
		manager      *UpdateCenterManager
		roundTripper *mhttp.MockRoundTripper
		responseBody string
		donwloadFile string
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		manager = &UpdateCenterManager{}
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		responseBody = "fake response"
		donwloadFile = "downloadfile.log"
	})

	AfterEach(func() {
		os.Remove(donwloadFile)
		ctrl.Finish()
	})

	Context("DownloadJenkins", func() {
		It("should success with basic cases", func() {
			manager.RoundTripper = roundTripper
			manager.MirrorSite = "http://mirrors.jenkins.io/"
			manager.LTS = false
			manager.ShowProgress = false
			manager.Output = donwloadFile

			request, _ := http.NewRequest("GET", "http://mirrors.jenkins.io/war/latest/jenkins.war", nil)
			response := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Header:     http.Header{},
				Request:    request,
				Body:       ioutil.NopCloser(bytes.NewBufferString(responseBody)),
			}
			roundTripper.EXPECT().
				RoundTrip(NewRequestMatcher(request)).Return(response, nil)
			err := manager.DownloadJenkins()
			Expect(err).To(BeNil())

			_, err = os.Stat(donwloadFile)
			Expect(err).To(BeNil())

			content, readErr := ioutil.ReadFile(donwloadFile)
			Expect(readErr).To(BeNil())
			Expect(string(content)).To(Equal(responseBody))
		})
	})

	Context("Upgrade", func() {
		It("basic cases", func() {
			manager.RoundTripper = roundTripper
			manager.URL = ""

			requestCrumb, _ := http.NewRequest("GET", fmt.Sprintf("%s/crumbIssuer/api/json", ""), nil)
			responseCrumb := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    requestCrumb,
				Body: ioutil.NopCloser(bytes.NewBufferString(`
				{"crumbRequestField":"CrumbRequestField","crumb":"Crumb"}
				`)),
			}
			roundTripper.EXPECT().
				RoundTrip(NewRequestMatcher(requestCrumb)).Return(responseCrumb, nil)

			request, _ := http.NewRequest("POST", "/updateCenter/upgrade", nil)
			request.Header.Add("CrumbRequestField", "Crumb")
			response := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    request,
				Body:       ioutil.NopCloser(bytes.NewBufferString("")),
			}
			roundTripper.EXPECT().
				RoundTrip(NewRequestMatcher(request)).Return(response, nil)

			err := manager.Upgrade()
			Expect(err).To(BeNil())
		})
	})

	Context("Status", func() {
		It("should success", func() {
			manager.RoundTripper = roundTripper
			manager.URL = ""

			request, _ := http.NewRequest("GET", "/updateCenter/api/json?pretty=false&depth=1", nil)
			response := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    request,
				Body: ioutil.NopCloser(bytes.NewBufferString(`
			{"RestartRequiredForCompletion": true}
			`)),
			}
			roundTripper.EXPECT().
				RoundTrip(NewRequestMatcher(request)).Return(response, nil)

			status, err := manager.Status()
			Expect(err).To(BeNil())
			Expect(status).NotTo(BeNil())
			Expect(status.RestartRequiredForCompletion).Should(BeTrue())
		})
	})

	Context("GetUpdateCenterPlugin", func() {
		It("basic cases", func() {
			manager.RoundTripper = roundTripper
			manager.URL = ""

			requestCenter, _ := http.NewRequest("GET", "/updateCenter/site/default/api/json?pretty=true&depth=2", nil)
			responseCenter := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    requestCenter,
				Body: ioutil.NopCloser(bytes.NewBufferString(`
				{
					"_class": "hudson.model.UpdateSite",
					"availables": 
					[
						{
							"name": "absint-a3",
							"sourceId": "default",
							"url": "http://updates.jenkins-ci.org/download/plugins/absint-a3/1.1.0/absint-a3.hpi",
							"version": "1.1.0",
							"categories": [
							  "buildwrapper"
							],
							"compatibleSinceVersion": null,
							"compatibleWithInstalledVersion": true,
							"dependencies": {
							  "command-launcher": "1.0",
							  "jdk-tool": "1.0",
							  "bouncycastle-api": "2.16.0"
							},
							"excerpt": "Provides Jenkins integration for the AbsInt Advanced Analyzer (a³) tools.",
							"installed": null,
							"minimumJavaVersion": null,
							"requiredCore": "1.625.3",
							"title": "AbsInt a³",
							"wiki": "https://plugins.jenkins.io/absint-a3"
						  }
					],
					"connectionCheckUrl": "http://www.google.com/",
					"dataTimestamp": 1567952107517,
					"hasUpdates": true,
					"id": "default",
					"updates":
					[
						{
							"name" : "blueocean-commons",
							"sourceId" : "default",
							"url" : "http://updates.jenkins-ci.org/download/plugins/blueocean-commons/1.19.0/blueocean-commons.hpi",
							"version" : "1.19.0",
							"categories" : [
							  "external",
							  "ui"
							],
							"compatibleSinceVersion" : null,
							"compatibleWithInstalledVersion" : true,
							"dependencies" : {
							  "jackson2-api" : "2.9.8"
							},
							"excerpt" : "This plugin is a part of Blue Ocean UI",
							"installed" : {
							  "active" : true,
							  "backupVersion" : "1.18.0",
							  "bundled" : false,
							  "deleted" : false,
							  "dependencies" : [
								{
								  
								}
							  ],
							  "downgradable" : true,
							  "enabled" : true,
							  "hasUpdate" : true,
							  "longName" : "Common API for Blue Ocean",
							  "minimumJavaVersion" : "1.8",
							  "pinned" : false,
							  "requiredCoreVersion" : "2.138.4",
							  "shortName" : "blueocean-commons",
							  "supportsDynamicLoad" : "MAYBE",
							  "url" : "https://wiki.jenkins-ci.org/display/JENKINS/Blue+Ocean+Plugin",
							  "version" : "1.18.1"
							},
							"minimumJavaVersion" : "1.8",
							"neededDependencies" : [
							  
							],
							"optionalDependencies" : {
							  
							},
							"requiredCore" : "2.138.4",
							"title" : "Common API for Blue Ocean",
							"wiki" : "https://plugins.jenkins.io/blueocean-commons"
						  }
					]
				}
				`)),
			}
			roundTripper.EXPECT().
				RoundTrip(NewRequestMatcher(requestCenter)).Return(responseCenter, nil)

			plugins, err := manager.GetSite()
			Expect(err).To(BeNil())
			Expect(plugins.UpdatePlugins[0].Name).To(Equal("blueocean-commons"))
		})
	})

	Context("SetMirrorCertificate", func() {
		var (
			enableMirror bool
			err          error
		)

		BeforeEach(func() {
			manager.URL = "http://localhost"
			manager.RoundTripper = roundTripper
		})

		JustBeforeEach(func() {
			PrepareForSetMirrorCertificate(roundTripper, manager.URL, "", "", enableMirror)
			err = manager.SetMirrorCertificate(enableMirror)
		})

		Context("enable mirror site", func() {
			BeforeEach(func() {
				enableMirror = true
			})

			It("should success", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("disable mirror site", func() {
			BeforeEach(func() {
				enableMirror = false
			})

			It("should success", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Context("ChangeUpdateCenterSite", func() {
		var (
			name      string
			mirrorURL string
			err       error
		)

		BeforeEach(func() {
			manager.URL = "http://localhost"
			manager.RoundTripper = roundTripper
			name = "default"
		})

		JustBeforeEach(func() {
			PrepareForChangeUpdateCenterSite(roundTripper, manager.URL, "", "", name, mirrorURL)
			err = manager.ChangeUpdateCenterSite(name, mirrorURL)
		})

		Context("enable mirror site", func() {
			BeforeEach(func() {
				mirrorURL = "http://fake.com"
			})

			It("should success", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})

var _ = Describe("GetJenkinsWarURL", func() {
	var (
		mgr    *UpdateCenterManager
		warURL string
	)

	BeforeEach(func() {
		mgr = &UpdateCenterManager{}
	})

	JustBeforeEach(func() {
		warURL = mgr.GetJenkinsWarURL()
	})

	It("default behaviour", func() {
		Expect(warURL).To(Equal("/war/latest/jenkins.war"))
	})

	Context("with lts", func() {
		BeforeEach(func() {
			mgr.LTS = true
		})

		It("with lts", func() {
			Expect(warURL).To(Equal("/war-stable/latest/jenkins.war"))
		})

		Context("with specific version", func() {
			BeforeEach(func() {
				mgr.Version = "fake"
			})

			It("with lts", func() {
				Expect(warURL).To(Equal("/war-stable/fake/jenkins.war"))
			})
		})

		Context("with mirror site", func() {
			BeforeEach(func() {
				mgr.MirrorSite = "http://baidu.com"
			})

			It("with mirror site", func() {
				Expect(warURL).To(Equal("http://baidu.com/war-stable/latest/jenkins.war"))
			})
		})
	})

	Context("with weekly", func() {
		BeforeEach(func() {
			mgr.LTS = false
		})

		It("with weekly", func() {
			Expect(warURL).To(Equal("/war/latest/jenkins.war"))
		})

		Context("with specific version", func() {
			BeforeEach(func() {
				mgr.Version = "fake"
			})

			It("with lts", func() {
				Expect(warURL).To(Equal("/war/fake/jenkins.war"))
			})
		})

		Context("with mirror site", func() {
			BeforeEach(func() {
				mgr.MirrorSite = "http://baidu.com"
			})

			It("with mirror site", func() {
				Expect(warURL).To(Equal("http://baidu.com/war/latest/jenkins.war"))
			})
		})
	})
})
