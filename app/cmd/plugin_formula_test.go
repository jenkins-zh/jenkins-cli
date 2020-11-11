package cmd

import (
	"github.com/magiconair/properties/assert"
	"testing"

	jenkinsFormula "github.com/jenkins-zh/jenkins-formulas/pkg/common"
)

func TestSortPlugins(t *testing.T) {
	plugins := []jenkinsFormula.Plugin{{
		GroupId:    "a",
		ArtifactId: "b",
	}, {
		GroupId:    "a",
		ArtifactId: "a",
	}, {
		GroupId:    "b",
		ArtifactId: "c",
	}}

	plugins = SortPlugins(plugins)
	assert.Equal(t, plugins[0].GroupId, "a")
	assert.Equal(t, plugins[0].ArtifactId, "a")
}
