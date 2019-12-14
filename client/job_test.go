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
			name := "fake"
			kind := "fake"

			PrepareOneItem(roundTripper, jobClient.URL, name, kind, "", "")

			result, err := jobClient.Search(name, kind)
			Expect(err).To(BeNil())
			Expect(result).NotTo(BeNil())
			Expect(len(result)).To(Equal(1))
			Expect(result[0].Name).To(Equal("fake"))
		})

		It("basic case without any result items", func() {
			name := "fake"
			kind := "fake"

			PrepareEmptyItems(roundTripper, jobClient.URL, name, kind, "", "")

			result, err := jobClient.Search(name, kind)
			Expect(err).To(BeNil())
			Expect(result).NotTo(BeNil())
			Expect(len(result)).To(Equal(0))
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
			PrepareForGetBuild(roundTripper, jobClient.URL, jobName, 2, "", "")

			result, err := jobClient.GetBuild(jobName, buildID)
			Expect(err).To(BeNil())
			Expect(result).NotTo(BeNil())
		})
	})

	Context("BuildWithParams", func() {
		It("no params", func() {
			jobName := "fake"

			PrepareForBuildWithNoParams(roundTripper, jobClient.URL, jobName, "", "")

			err := jobClient.BuildWithParams(jobName, []ParameterDefinition{})
			Expect(err).To(BeNil())
		})

		It("with params", func() {
			jobName := "fake"

			PrepareForBuildWithParams(roundTripper, jobClient.URL, jobName, "", "")

			err := jobClient.BuildWithParams(jobName, []ParameterDefinition{ParameterDefinition{
				Name:  "name",
				Value: "value",
			}})
			Expect(err).To(BeNil())
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

		It("stop the last job build without a folder", func() {
			jobName := "fakeJob"
			buildID := -1
			request, _ := http.NewRequest("POST", fmt.Sprintf("%s/job/%s/lastBuild/stop", jobClient.URL, jobName), nil)
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

	Context("GetPipeline", func() {
		It("simple case, should success", func() {
			PrepareForPipelineJob(roundTripper, jobClient.URL, "", "")
			job, err := jobClient.GetPipeline("test")
			Expect(err).To(BeNil())
			Expect(job.Script).To(Equal("script"))
		})
	})

	Context("UpdatePipeline", func() {
		It("simple case, should success", func() {
			PrepareForUpdatePipelineJob(roundTripper, jobClient.URL, "", "", "")
			err := jobClient.UpdatePipeline("test", "")
			Expect(err).To(BeNil())
		})
	})

	Context("Create", func() {
		var (
			jobPayload CreateJobPayload
		)

		BeforeEach(func() {
			jobPayload = CreateJobPayload{
				Name: "jobName",
				Mode: "jobType",
			}
		})

		It("create a normal job, should success", func() {
			PrepareForCreatePipelineJob(roundTripper, jobClient.URL, "", "", jobPayload)
			err := jobClient.Create(jobPayload)
			Expect(err).To(BeNil())
		})

		It("create a job by copy mode", func() {
			jobPayload.From = "another-one"
			jobPayload.Mode = "copy"
			PrepareForCreatePipelineJob(roundTripper, jobClient.URL, "", "", jobPayload)
			err := jobClient.Create(jobPayload)
			Expect(err).To(BeNil())
		})
	})

	Context("Delete", func() {
		It("delete a job", func() {
			jobName := "fakeJob"
			request, _ := http.NewRequest("POST", fmt.Sprintf("%s/job/%s/doDelete", jobClient.URL, jobName), nil)
			request.Header.Add("CrumbRequestField", "Crumb")
			request.Header.Add(util.ContentType, util.ApplicationForm)
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

	Context("GetJobInputActions", func() {
		It("simple case, should success", func() {
			PrepareForGetJobInputActions(roundTripper, jobClient.URL, "", "", "jobName", 1)
			actions, err := jobClient.GetJobInputActions("jobName", 1)
			Expect(err).To(BeNil())
			Expect(len(actions)).To(Equal(1))
			Expect(actions[0].Message).To(Equal("message"))
		})
	})

	Context("JobInputSubmit", func() {
		It("simple case, should success", func() {
			PrepareForSubmitInput(roundTripper, jobClient.URL, "/job/jobName", "", "")
			err := jobClient.JobInputSubmit("jobName", "Eff7d5dba32b4da32d9a67a519434d3f", 1, true, nil)
			Expect(err).To(BeNil())
		})
	})

	Context("GetHistory", func() {
		It("simple case, should success", func() {
			jobName := "fakeJob"

			PrepareForGetJob(roundTripper, jobClient.URL, jobName, "", "")
			PrepareForGetBuild(roundTripper, jobClient.URL, jobName, 1, "", "")
			PrepareForGetBuild(roundTripper, jobClient.URL, jobName, 2, "", "")

			builds, err := jobClient.GetHistory(jobName)
			Expect(err).To(BeNil())
			Expect(builds).NotTo(BeNil())
			Expect(len(builds)).To(Equal(2))
		})
	})

	Context("Log", func() {
		It("with a specific build number", func() {
			jobName := "fakeJob"

			PrepareForJobLog(roundTripper, jobClient.URL, jobName, 1, "", "")

			log, err := jobClient.Log(jobName, 1, 0)
			Expect(err).To(BeNil())
			Expect(log.Text).To(Equal("fake log"))
		})

		It("get the last job's log", func() {
			jobName := "fakeJob"

			PrepareForJobLog(roundTripper, jobClient.URL, jobName, -1, "", "")

			log, err := jobClient.Log(jobName, -1, 0)
			Expect(err).To(BeNil())
			Expect(log.Text).To(Equal("fake log"))
		})
	})
})

var _ = Describe("test function ParseJobPath", func() {
	var (
		path    string
		jobName string
	)

	JustBeforeEach(func() {
		path = ParseJobPath(jobName)
	})

	It("job name is empty", func() {
		Expect(path).To(BeEmpty())
	})

	Context("job name is not empty", func() {
		BeforeEach(func() {
			jobName = "abc"
		})

		It("job name separate with whitespaces", func() {
			Expect(path).To(Equal(fmt.Sprintf("/job/%s", jobName)))
		})

		Context("multi level of job name", func() {
			BeforeEach(func() {
				jobName = "abc def"
			})

			It("should success", func() {
				Expect(path).To(Equal("/job/abc/job/def"))
			})
		})
	})

	Context("job name with URL path", func() {
		BeforeEach(func() {
			jobName = "/job/abc/job/def"
		})

		It("should success", func() {
			Expect(path).To(Equal(jobName))
		})
	})
})
