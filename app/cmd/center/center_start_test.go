package center

import (
	"github.com/jenkins-zh/jenkins-cli/app/cmd"
	"github.com/jenkins-zh/jenkins-cli/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
)

var _ = Describe("center start command", func() {
	var (
		configFile string
	)
	BeforeEach(func() {
		file, err := ioutil.TempFile(".", "test.yaml")
		Expect(err).NotTo(HaveOccurred())

		configFile = file.Name()
		data, err := cmd.GenerateSampleConfig()
		Expect(err).To(BeNil())
		err = ioutil.WriteFile(configFile, data, 0664)
		Expect(err).To(BeNil())

		cmd.rootOptions.ConfigFile = configFile
	})
	AfterEach(func() {
		os.RemoveAll(configFile)
	})
	It("enable mirror site", func() {
		centerStartOption.SystemCallExec = util.FakeSystemCallExecSuccess
		centerStartOption.LookPathContext = util.FakeLookPath
		cmd.rootCmd.SetArgs([]string{"center", "start", "--dry-run", "--env", "a=b", "--concurrent-indexing=12", "--https-enable"})
		_, err := cmd.rootCmd.ExecuteC()
		Expect(err).To(BeNil())
	})
})
