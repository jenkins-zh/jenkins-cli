package client

import (
	"github.com/golang/mock/gomock"
	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("core test", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
		coreClient   CoreClient

		username string
		password string
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		coreClient = CoreClient{}
		coreClient.RoundTripper = roundTripper
		coreClient.URL = "http://localhost"

		username = "admin"
		password = "token"

		coreClient.UserName = username
		coreClient.Token = password
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("Get", func() {
		It("should success", func() {
			PrepareRestart(roundTripper, coreClient.URL, username, password, 503)

			err := coreClient.Restart()
			Expect(err).To(BeNil())
		})

		It("should error, 400", func() {
			PrepareRestart(roundTripper, coreClient.URL, username, password, 400)

			err := coreClient.Restart()
			Expect(err).To(HaveOccurred())
		})

		It("should success", func() {
			PrepareRestartDirectly(roundTripper, coreClient.URL, username, password, 503)

			err := coreClient.RestartDirectly()
			Expect(err).To(BeNil())
		})

		It("GetIdentity", func() {
			PrepareForGetIdentity(roundTripper, coreClient.URL, username, password)

			identity, err := coreClient.GetIdentity()
			Expect(err).NotTo(HaveOccurred())
			Expect(identity).To(Equal(JenkinsIdentity{
				Fingerprint:   "fingerprint",
				PublicKey:     "publicKey",
				SystemMessage: "systemMessage",
			}))
		})
	})

	Context("shutdown", func() {
		var (
			err  error
			safe bool
		)

		JustBeforeEach(func() {
			PrepareForShutdown(roundTripper, coreClient.URL, username, password, safe)
			err = coreClient.Shutdown(safe)
		})

		Context("shutdown safely", func() {
			BeforeEach(func() {
				safe = true
			})
			It("should success", func() {
				Expect(err).To(BeNil())
			})
		})

		Context("shutdown not safely", func() {
			BeforeEach(func() {
				safe = false
			})
			It("should success", func() {
				Expect(err).To(BeNil())
			})
		})
	})

	Context("prepare shutdown", func() {
		var (
			err    error
			cancel bool
		)

		JustBeforeEach(func() {
			PrepareForCancelShutdown(roundTripper, coreClient.URL, username, password, cancel)
			err = coreClient.PrepareShutdown(cancel)
		})

		Context("cancelQuietDown", func() {
			BeforeEach(func() {
				cancel = true
			})
			It("should success", func() {
				Expect(err).To(BeNil())
			})
		})

		Context("quietDown", func() {
			BeforeEach(func() {
				cancel = false
			})
			It("should success", func() {
				Expect(err).To(BeNil())
			})
		})
	})
})
