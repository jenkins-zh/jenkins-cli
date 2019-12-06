package cmd

import (
	"github.com/jenkins-zh/jenkins-cli/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("center start command", func() {
	It("enable mirror site", func() {
		centerStartOption.SystemCallExec = util.FakeSystemCallExecSuccess
		rootCmd.SetArgs([]string{"center", "start", "--dry-run"})
		_, err := rootCmd.ExecuteC()
		Expect(err).To(BeNil())
	})
})
