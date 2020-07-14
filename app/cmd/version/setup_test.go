package version_test

import (
	"testing"

	"github.com/onsi/ginkgo/reporters"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("test-cmd-version.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "app/cmd/version", []Reporter{junitReporter})
}
