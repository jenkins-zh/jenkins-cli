package i18n_test

import (
	"testing"

	"github.com/onsi/ginkgo/reporters"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestI18n(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("test-i18n.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "app", []Reporter{junitReporter})
}
