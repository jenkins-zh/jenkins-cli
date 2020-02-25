package cmd

import (
	"bytes"
	"github.com/golang/mock/gomock"
	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
	"github.com/jenkins-zh/jenkins-cli/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"testing"
)

var _ = Describe("cwp command test", func() {
	var (
		ctrl       *gomock.Controller
		localCache string
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		localCache = os.TempDir()
		roundTripper := mhttp.NewMockRoundTripper(ctrl)
		cwpOpts := CWPOptions{
			CommonOption: CommonOption{RoundTripper: roundTripper},
			MetadataURL:  "http://localhost/maven-metadata.xml",
			LocalCache:   localCache,
		}
		prepareMavenMetadataRequest(roundTripper)

		fakeContent := "hello"
		prepareDownloadFileRequest(cwpOpts.GetCWPURL("2.0-alpha-2"), fakeContent, roundTripper)

		cwpOptions.SystemCallExec = util.FakeSystemCallExecSuccess
		cwpOptions.LookPathContext = util.FakeLookPath
	})

	AfterEach(func() {
		os.RemoveAll(localCache)
		ctrl.Finish()
	})

	Context("basic test", func() {
		It("should success", func() {
			rootCmd.SetArgs([]string{"cwp"})
			_, err := rootCmd.ExecuteC()
			Expect(err).To(BeNil())
		})
	})
})

func TestDownload(t *testing.T) {
	ctrl := gomock.NewController(t)

	tmpDir := os.TempDir()
	defer os.RemoveAll(tmpDir)

	roundTripper := mhttp.NewMockRoundTripper(ctrl)
	cwpOpts := CWPOptions{
		CommonOption: CommonOption{RoundTripper: roundTripper},
		MetadataURL:  "http://localhost/maven-metadata.xml",
		LocalCache:   tmpDir,
	}
	prepareMavenMetadataRequest(roundTripper)

	fakeContent := "hello"
	prepareDownloadFileRequest(cwpOpts.GetCWPURL("2.0-alpha-2"), fakeContent, roundTripper)

	err := cwpOpts.Download()
	assert.Nil(t, err)

	var data []byte
	data, err = ioutil.ReadFile(path.Join(tmpDir, "cwp-cli.jar"))
	assert.Nil(t, err)
	assert.Equal(t, fakeContent, string(data))
}

func TestGetLatest(t *testing.T) {
	ctrl := gomock.NewController(t)

	roundTripper := mhttp.NewMockRoundTripper(ctrl)
	cwpOpts := CWPOptions{
		CommonOption: CommonOption{RoundTripper: roundTripper},
		MetadataURL:  "http://localhost/maven-metadata.xml",
	}
	prepareMavenMetadataRequest(roundTripper)

	ver, err := cwpOpts.GetLatest()
	assert.Nil(t, err)
	assert.Equal(t, "2.0-alpha-2", ver)
}

func prepareDownloadFileRequest(url, content string, roundTripper *mhttp.MockRoundTripper) {
	request, _ := http.NewRequest("GET", url, nil)
	response := &http.Response{
		StatusCode: 200,
		Request:    request,
		Body:       ioutil.NopCloser(bytes.NewBufferString(content)),
	}
	roundTripper.EXPECT().
		RoundTrip(request).Return(response, nil)
}

func prepareMavenMetadataRequest(roundTripper *mhttp.MockRoundTripper) {
	request, _ := http.NewRequest("GET", "http://localhost/maven-metadata.xml", nil)
	response := &http.Response{
		StatusCode: 200,
		Request:    request,
		Body:       ioutil.NopCloser(bytes.NewBufferString(getMavenMetadataSample())),
	}
	roundTripper.EXPECT().
		RoundTrip(request).Return(response, nil)
}

func getMavenMetadataSample() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<metadata>
  <groupId>io.jenkins.tools.custom-war-packager</groupId>
  <artifactId>custom-war-packager-cli</artifactId>
  <versioning>
    <latest>2.0-alpha-2</latest>
    <release>2.0-alpha-2</release>
    <versions>
      <version>2.0-alpha-2</version>
	</versions>
    <lastUpdated>20190815083928</lastUpdated>
  </versioning>
</metadata>`
}
