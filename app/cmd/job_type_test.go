package cmd

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"

	. "github.com/jenkins-zh/jenkins-cli/app/config"
	"github.com/jenkins-zh/jenkins-cli/client"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
)

var _ = Describe("job type command", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		jobTypeOption.RoundTripper = roundTripper
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
		It("GetCategories", func() {
			data, err := GenerateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			request, _ := http.NewRequest("GET", "http://localhost:8080/jenkins/view/all/itemCategories?depth=3", nil)
			request.SetBasicAuth("admin", "111e3a2f0231198855dceaff96f20540a9")
			response := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    request,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"_class":"jenkins.model.item_category.Categories","categories":[{"description":"description","id":"standalone-projects","items":[{"displayName":"Freestyle project","iconFilePathPattern":"static/da605e5f/images/:size/freestyleproject.png","description":"description","iconClassName":"icon-freestyle-project","class":"hudson.model.FreeStyleProject","order":1}],"minToShow":1,"name":"Nested Projects","order":1}]}`)),
			}
			roundTripper.EXPECT().
				RoundTrip(client.NewRequestMatcher(request)).Return(response, nil)

			config = &Config{
				Current: "fake",
				JenkinsServers: []JenkinsServer{{
					Name:     "fake",
					URL:      "http://localhost:8080/jenkins",
					UserName: "admin",
					Token:    "111e3a2f0231198855dceaff96f20540a9",
				}},
			}
			jclient := &client.JobClient{
				JenkinsCore: client.JenkinsCore{
					RoundTripper: jobTypeOption.RoundTripper,
				},
			}
			getCurrentJenkinsAndClient(&(jclient.JenkinsCore))

			typeMap, types, err := GetCategories(jclient)
			Expect(err).To(BeNil())
			Expect(len(typeMap)).To(Equal(1))
			Expect(len(types)).To(Equal(1))
		})

		It("should success, empty list", func() {
			data, err := GenerateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			request, _ := http.NewRequest("GET", "http://localhost:8080/jenkins/view/all/itemCategories?depth=3", nil)
			request.SetBasicAuth("admin", "111e3a2f0231198855dceaff96f20540a9")
			response := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    request,
				Body:       ioutil.NopCloser(bytes.NewBufferString("{}")),
			}
			roundTripper.EXPECT().
				RoundTrip(client.NewRequestMatcher(request)).Return(response, nil)

			rootCmd.SetArgs([]string{"job", "type"})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal("DisplayName Class\n"))
		})

		It("should success, empty list", func() {
			data, err := GenerateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			request, _ := http.NewRequest("GET", "http://localhost:8080/jenkins/view/all/itemCategories?depth=3", nil)
			request.SetBasicAuth("admin", "111e3a2f0231198855dceaff96f20540a9")
			response := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    request,
				Body:       ioutil.NopCloser(bytes.NewBufferString("{}")),
			}
			roundTripper.EXPECT().
				RoundTrip(client.NewRequestMatcher(request)).Return(response, nil)

			rootCmd.SetArgs([]string{"job", "type"})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal("DisplayName Class\n"))
		})

		It("should success, one item", func() {
			data, err := GenerateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			request, _ := http.NewRequest("GET", "http://localhost:8080/jenkins/view/all/itemCategories?depth=3", nil)
			request.SetBasicAuth("admin", "111e3a2f0231198855dceaff96f20540a9")
			response := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    request,
				Body: ioutil.NopCloser(bytes.NewBufferString(`
				{"categories":[
					{
					  "description" : "description",
					  "id" : "standalone-projects",
					  "items" : [
						{
						  "displayName" : "displayName",
						  "iconFilePathPattern" : "iconFilePathPattern",
						  "description" : "description",
						  "iconClassName" : "iconClassName",
						  "class" : "class",
						  "order" : 1
						}],
						"minToShow" : 1,
						"name" : "Nested Projects",
						"order" : 1
					}]
				}`)),
			}
			roundTripper.EXPECT().
				RoundTrip(client.NewRequestMatcher(request)).Return(response, nil)

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			rootCmd.SetArgs([]string{"job", "type"})

			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal(`DisplayName Class
displayName class
`))
		})
	})
})
