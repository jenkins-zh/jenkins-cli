package cmd

import (
	"testing"

	"github.com/onsi/ginkgo/reporters"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("test-app.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "app/cmd", []Reporter{junitReporter})
}
