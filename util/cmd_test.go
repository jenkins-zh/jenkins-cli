package util

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"testing"
)

var _ = Describe("Test open browser", func() {
	It("should success", func() {
		err := Open("fake://url", FakeExecCommandSuccess)
		Expect(err).To(BeNil())
	})
})

// TestShellProcessSuccess only for test
func TestShellProcessSuccess(t *testing.T) {
	if os.Getenv("GO_TEST_PROCESS") != "1" {
		return
	}
	//os.Exit(0)
}
