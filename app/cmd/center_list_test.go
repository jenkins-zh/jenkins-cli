package cmd

import (
	"bytes"
	_ "embed"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

//go:embed center_list.txt
var resultOneVersionData string

var _ = Describe("center list command", func() {
	var (
		ctrl *gomock.Controller
		// roundTripper *mhttp.MockRoundTripper
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		// roundTripper = mhttp.NewMockRoundTripper(ctrl)
		// centerListOption.RoundTripper = roundTripper
		rootCmd.SetArgs([]string{})
		rootOptions.Jenkins = ""
		rootOptions.ConfigFile = "test.yaml"

	})

	AfterEach(func() {
		rootCmd.SetArgs([]string{})
		os.Remove(rootOptions.ConfigFile)
		rootOptions.ConfigFile = ""
		ctrl.Finish()
	})

	Context("basic cases", func() {
		jenkinsVersions := []string{"Jenkins 2.289.2", "Jenkins 2.289.1"}
		result := []string{
			"You already have the latest version of Jenkins installed!",
			resultOneVersionData,
		}
		It("no jenkins version information in the list", func() {
			data, err := GenerateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())
			temp, _ := getChangelog(LtsURL, jenkinsVersions[0][8:], 10, mockGetVersionData)
			Expect(temp).To(Equal(result[0]))
		})

		It("one jenkins version information in the list", func() {
			data, err := GenerateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())
			temp, _ := getChangelog(LtsURL, jenkinsVersions[1][8:], 10, mockGetVersionData)
			Expect(temp).To(Equal(result[1]))
		})
	})
})

func mockGetVersionData(rss string) ([]Item, string, error) {
	requestVersionData, _ := http.NewRequest(http.MethodGet, rss, nil)
	responseVersionData := &http.Response{
		StatusCode: 200,
		Proto:      "HTTP/1.1",
		Request:    requestVersionData,
		Body:       ioutil.NopCloser(bytes.NewBufferString(versionXML())),
	}
	bytes, err := ioutil.ReadAll(responseVersionData.Body)
	if err != nil {
		return nil, "", err
	}
	var centerListOption CenterListOption
	xml.Unmarshal(bytes, &centerListOption)
	return centerListOption.Channel.Items, centerListOption.Channel.Title, nil
}

func versionXML() string {
	return `<rss>
<channel>
<title> Jenkins LTS Changelog </title>
<item>
<title>Jenkins 2.289.2</title>
<description>&lt;strong&gt;&lt;/strong&gt;&lt;ul&gt;&lt;li&gt; Security: Important security fixes. &lt;/li&gt;
&lt;li&gt;RFE: Winstone 5.18: Update Jetty from 9.4.39.v20210325 to 9.4.41.v20210516 for bug fixes and enhancements. &lt;/li&gt;
 &lt;/ul&gt; </description>
<pubDate> Wed, 30 Jun 2021 00:00:00 +0000 </pubDate>
</item>
<item>
<title>Jenkins 2.289.1</title>
<description> &lt;strong&gt;Changes since 2.289:&lt;/strong&gt; &lt;ul&gt;
&lt;li&gt; Bug: Fix form submission for some specific form validation cases (regression in 2.289). &lt;/li&gt;
&lt;li&gt; Bug: Wrap the build name in the build results list if it is too long. &lt;/li&gt;
 &lt;/ul&gt; </description>
<pubDate> Wed, 2 Jun 2021 00:00:00 +0000 </pubDate>
</item>
</channel>
</rss>`
}
