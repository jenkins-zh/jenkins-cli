package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/golang/mock/gomock"
	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
	"github.com/jenkins-zh/jenkins-cli/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("job test", func() {
	var (
		ctrl         *gomock.Controller
		jobClient    JobClient
		roundTripper *mhttp.MockRoundTripper
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		jobClient = JobClient{}
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		jobClient.RoundTripper = roundTripper
		jobClient.URL = "http://localhost"
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("Search", func() {
		It("basic case with one result item", func() {
			keyword := "fake"

			request, _ := http.NewRequest("GET", fmt.Sprintf("%s%s%s&max=1", jobClient.URL, "/search/suggest?query=", keyword), nil)
			response := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    request,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"suggestions": [{"name": "fake"}]}`)),
			}
			roundTripper.EXPECT().
				RoundTrip(request).Return(response, nil)

			result, err := jobClient.Search(keyword, 1)
			Expect(err).To(BeNil())
			Expect(result).NotTo(BeNil())
			Expect(len(result.Suggestions)).To(Equal(1))
			Expect(result.Suggestions[0].Name).To(Equal("fake"))
		})

		It("basic case without any result items", func() {
			keyword := "fake"

			request, _ := http.NewRequest("GET", fmt.Sprintf("%s%s%s&max=1", jobClient.URL, "/search/suggest?query=", keyword), nil)
			response := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    request,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"suggestions":[]}`)),
			}
			roundTripper.EXPECT().
				RoundTrip(request).Return(response, nil)

			result, err := jobClient.Search(keyword, 1)
			Expect(err).To(BeNil())
			Expect(result).NotTo(BeNil())
			Expect(len(result.Suggestions)).To(Equal(0))
		})
	})

	Context("Build", func() {
		It("trigger a simple job without a folder", func() {
			jobName := "fakeJob"
			request, _ := http.NewRequest("POST", fmt.Sprintf("%s/job/%s/build", jobClient.URL, jobName), nil)
			request.Header.Add("CrumbRequestField", "Crumb")
			response := &http.Response{
				StatusCode: 201,
				Proto:      "HTTP/1.1",
				Request:    request,
				Body:       ioutil.NopCloser(bytes.NewBufferString("")),
			}
			roundTripper.EXPECT().
				RoundTrip(request).Return(response, nil)

			requestCrumb, _ := http.NewRequest("GET", fmt.Sprintf("%s%s", jobClient.URL, "/crumbIssuer/api/json"), nil)
			responseCrumb := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    requestCrumb,
				Body: ioutil.NopCloser(bytes.NewBufferString(`
				{"crumbRequestField":"CrumbRequestField","crumb":"Crumb"}
				`)),
			}
			roundTripper.EXPECT().
				RoundTrip(requestCrumb).Return(responseCrumb, nil)

			err := jobClient.Build(jobName)
			Expect(err).To(BeNil())
		})

		It("trigger a simple job with an error", func() {
			jobName := "fakeJob"
			request, _ := http.NewRequest("POST", fmt.Sprintf("%s/job/%s/build", jobClient.URL, jobName), nil)
			request.Header.Add("CrumbRequestField", "Crumb")
			response := &http.Response{
				StatusCode: 500,
				Proto:      "HTTP/1.1",
				Request:    request,
				Body:       ioutil.NopCloser(bytes.NewBufferString("")),
			}
			roundTripper.EXPECT().
				RoundTrip(request).Return(response, nil)

			requestCrumb, _ := http.NewRequest("GET", fmt.Sprintf("%s%s", jobClient.URL, "/crumbIssuer/api/json"), nil)
			responseCrumb := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    requestCrumb,
				Body: ioutil.NopCloser(bytes.NewBufferString(`
				{"crumbRequestField":"CrumbRequestField","crumb":"Crumb"}
				`)),
			}
			roundTripper.EXPECT().
				RoundTrip(requestCrumb).Return(responseCrumb, nil)

			err := jobClient.Build(jobName)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("GetBuild", func() {
		It("basic case with the last build", func() {
			jobName := "fake"
			buildID := -1

			request, _ := http.NewRequest("GET", fmt.Sprintf("%s/job/%s/lastBuild/api/json", jobClient.URL, jobName), nil)
			response := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    request,
				Body: ioutil.NopCloser(bytes.NewBufferString(`
				{"displayName":"fake"}
				`)),
			}
			roundTripper.EXPECT().
				RoundTrip(request).Return(response, nil)

			result, err := jobClient.GetBuild(jobName, buildID)
			Expect(err).To(BeNil())
			Expect(result).NotTo(BeNil())
		})

		It("basic case with one build", func() {
			jobName := "fake"
			buildID := 2

			request, _ := http.NewRequest("GET", fmt.Sprintf("%s/job/%s/%d/api/json", jobClient.URL, jobName, buildID), nil)
			response := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    request,
				Body: ioutil.NopCloser(bytes.NewBufferString(`
				{"displayName":"fake"}
				`)),
			}
			roundTripper.EXPECT().
				RoundTrip(request).Return(response, nil)

			result, err := jobClient.GetBuild(jobName, buildID)
			Expect(err).To(BeNil())
			Expect(result).NotTo(BeNil())
		})
	})

	Context("StopJob", func() {
		It("stop a job build without a folder", func() {
			jobName := "fakeJob"
			buildID := 1
			request, _ := http.NewRequest("POST", fmt.Sprintf("%s/job/%s/%d/stop", jobClient.URL, jobName, buildID), nil)
			request.Header.Add("CrumbRequestField", "Crumb")
			response := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    request,
				Body:       ioutil.NopCloser(bytes.NewBufferString("")),
			}
			roundTripper.EXPECT().
				RoundTrip(request).Return(response, nil)

			requestCrumb, _ := http.NewRequest("GET", fmt.Sprintf("%s%s", jobClient.URL, "/crumbIssuer/api/json"), nil)
			responseCrumb := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    requestCrumb,
				Body: ioutil.NopCloser(bytes.NewBufferString(`
				{"crumbRequestField":"CrumbRequestField","crumb":"Crumb"}
				`)),
			}
			roundTripper.EXPECT().
				RoundTrip(requestCrumb).Return(responseCrumb, nil)

			err := jobClient.StopJob(jobName, buildID)
			Expect(err).To(BeNil())
		})
	})

	Context("GetJob", func() {
		It("get a job without in a folder", func() {
			jobName := "fake"

			request, _ := http.NewRequest("GET", fmt.Sprintf("%s/job/%s/api/json", jobClient.URL, jobName), nil)
			response := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    request,
				Body: ioutil.NopCloser(bytes.NewBufferString(`
				{"name":"fake"}
				`)),
			}
			roundTripper.EXPECT().
				RoundTrip(request).Return(response, nil)

			result, err := jobClient.GetJob(jobName)
			Expect(err).To(BeNil())
			Expect(result).NotTo(BeNil())
			Expect(result.Name).To(Equal(jobName))
		})
	})

	Context("GetJobTypeCategories", func() {
		It("simple case, should success", func() {
			request, _ := http.NewRequest("GET", fmt.Sprintf("%s/view/all/itemCategories?depth=3", jobClient.URL), nil)
			response := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    request,
				Body:       ioutil.NopCloser(bytes.NewBufferString("{}")),
			}
			roundTripper.EXPECT().
				RoundTrip(request).Return(response, nil)

			_, err := jobClient.GetJobTypeCategories()
			Expect(err).To(BeNil())
		})
	})

	Context("Delete", func() {
		It("delete a job", func() {
			jobName := "fakeJob"
			request, _ := http.NewRequest("POST", fmt.Sprintf("%s/job/%s/doDelete", jobClient.URL, jobName), nil)
			request.Header.Add("CrumbRequestField", "Crumb")
			request.Header.Add(util.CONTENT_TYPE, util.APP_FORM)
			response := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    request,
				Body:       ioutil.NopCloser(bytes.NewBufferString("")),
			}
			roundTripper.EXPECT().
				RoundTrip(request).Return(response, nil)

			requestCrumb, _ := http.NewRequest("GET", fmt.Sprintf("%s%s", jobClient.URL, "/crumbIssuer/api/json"), nil)
			responseCrumb := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    requestCrumb,
				Body: ioutil.NopCloser(bytes.NewBufferString(`
				{"crumbRequestField":"CrumbRequestField","crumb":"Crumb"}
				`)),
			}
			roundTripper.EXPECT().
				RoundTrip(requestCrumb).Return(responseCrumb, nil)

			err := jobClient.Delete(jobName)
			Expect(err).To(BeNil())
		})
	})
})
