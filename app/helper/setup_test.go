package helper_test

import (
	"testing"

	"github.com/onsi/ginkgo/reporters"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestApp(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("helper.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "app/helper", []Reporter{junitReporter})
}
