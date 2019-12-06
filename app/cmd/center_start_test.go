package cmd

import (
	"github.com/jenkins-zh/jenkins-cli/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("center start command", func() {
	It("enable mirror site", func() {
		centerStartOption.SystemCallExec = util.FakeSystemCallExecSuccess
		centerStartOption.LookPathContext = util.FakeLookPath
		rootCmd.SetArgs([]string{"center", "start", "--dry-run", "--env", "a=b", "--concurrent-indexing=12", "--https-enable"})
		_, err := rootCmd.ExecuteC()
		Expect(err).To(BeNil())
	})
})
