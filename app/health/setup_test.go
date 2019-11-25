package health_test

import (
	"testing"

	"github.com/onsi/ginkgo/reporters"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestI18n(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("test-health.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "command health check", []Reporter{junitReporter})
}
