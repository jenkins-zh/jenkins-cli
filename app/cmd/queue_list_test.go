package cmd

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
)

var _ = Describe("queue list command", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
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

	Context("with http requests", func() {
		BeforeEach(func() {
			roundTripper = mhttp.NewMockRoundTripper(ctrl)
			queueListOption.RoundTripper = roundTripper
		})

		It("should success", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			client.PrepareGetQueue(roundTripper, "http://localhost:8080/jenkins", "admin", "111e3a2f0231198855dceaff96f20540a9")

			rootCmd.SetArgs([]string{"queue", "list", "-o", "json"})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())
			Expect(buf.String()).To(Equal(`[
  {
    "Blocked": false,
    "Buildable": true,
    "ID": 62,
    "Params": "",
    "Pending": false,
    "Stuck": true,
    "URL": "queue/item/62/",
    "Why": "等待下一个可用的执行器",
    "BuildableStartMilliseconds": 1567753826770,
    "InQueueSince": 1567753826770,
    "Actions": []
  }
]`))
		})

		It("output with table format", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			client.PrepareGetQueue(roundTripper, "http://localhost:8080/jenkins", "admin", "111e3a2f0231198855dceaff96f20540a9")

			rootCmd.SetArgs([]string{"queue", "list", "-o", "table"})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal(`ID Why                    URL
62 等待下一个可用的执行器 queue/item/62/
`))
		})
	})
})
