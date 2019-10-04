package client

import (
	"github.com/golang/mock/gomock"
	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("artifacts test", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
		artifactClient ArtifactClient

		username string
		password string
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		artifactClient = ArtifactClient{}
		artifactClient.RoundTripper = roundTripper
		artifactClient.URL = "http://localhost"

		username = "admin"
		password = "token"
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("List", func() {
		It("should success", func() {
			artifactClient.UserName = username
			artifactClient.Token = password

			jobName := "fakename"
			PrepareGetArtifacts(roundTripper, artifactClient.URL, username, password, jobName, 1)

			artifacts, err := artifactClient.List(jobName, 1)
			Expect(err).To(BeNil())
			Expect(len(artifacts)).To(Equal(1))
		})

		It("should success, with empty artifacts", func() {
			artifactClient.UserName = username
			artifactClient.Token = password

			jobName := "fakename"
			PrepareGetEmptyArtifacts(roundTripper, artifactClient.URL, username, password, jobName, 1)

			artifacts, err := artifactClient.List(jobName, 1)
			Expect(err).To(BeNil())
			Expect(len(artifacts)).To(Equal(0))
		})
	})
})
