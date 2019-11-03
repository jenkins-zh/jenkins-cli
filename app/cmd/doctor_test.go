package cmd

import (
	"bytes"

	//"bytes"
	//"fmt"
	"io/ioutil"
	"strings"

	"github.com/jenkins-zh/jenkins-cli/client"
	"gopkg.in/yaml.v2"

	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("doctor command", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		doctorOption.RoundTripper = roundTripper
		rootOptions.Jenkins = ""
		rootOptions.ConfigFile = "test.yaml"
		config = &Config{
			Current: "a",
			JenkinsServers: []JenkinsServer{JenkinsServer{
				Name:     "a",
				URL:      "http://localhost:8080/jenkins",
				UserName: "admin",
				Token:    "111e3a2f0231198855dceaff96f20540a9",
			}, JenkinsServer{
				Name:     "b",
				URL:      "http://localhost:8080/jenkins1",
				UserName: "admin1",
				Token:    "111e3a2f0231198855dceaff96f20540a9",
			}, JenkinsServer{
				Name:     "c",
				URL:      "http://localhost:8080/jenkins2",
				UserName: "admin2",
				Token:    "111e3a2f0231198855dceaff96f20540a9",
			}},
		}
	})

	AfterEach(func() {
		config = nil
		doctorCmd.SetArgs([]string{})
		ctrl.Finish()
	})

	Context("test mode", func() {
		It("test JenkinsServers status", func() {
			config = &Config{
				Current: "a",
				JenkinsServers: []JenkinsServer{JenkinsServer{
					Name:     "a",
					URL:      "http://localhost:8080/jenkins",
					UserName: "admin",
					Token:    "111e3a2f0231198855dceaff96f20540a9",
				}, JenkinsServer{
					Name:     "b",
					URL:      "http://localhost:8080/jenkins1",
					UserName: "admin1",
					Token:    "111e3a2f0231198855dceaff96f20540a9",
				}, JenkinsServer{
					Name:     "a",
					URL:      "http://localhost:8080/jenkins2",
					UserName: "admin2",
					Token:    "111e3a2f0231198855dceaff96f20540a9",
				}},
			}

			data, err := yaml.Marshal(&config)
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			names := getJenkinsNames()
			var outString string
			outString += "Begining checking the name in the configuration file is duplicated：\n"
			duplicateName := ""
			for i := range names {
				for j := range names {
					if i != j && names[i] == names[j] && !strings.Contains(duplicateName, names[i]) {
						duplicateName += names[i] + " "
					}
				}
			}
			if duplicateName == "" {
				outString += "  Checked it sure. no duplicated config Name\n"
			} else {
				outString += "  Duplicate names: " + duplicateName + "\n"
			}
			outString += "Begining checking JenkinsServer status form the configuration files: \n"
			jenkinsServers := config.JenkinsServers
			for k := range jenkinsServers {
				if k == 1 {
					request, _ := client.PrepareFor500InstalledPluginList(roundTripper, jenkinsServers[k].URL)
					request.SetBasicAuth(jenkinsServers[k].UserName, jenkinsServers[k].Token)
					outString += "  checking the No." + string(k) + " - " + jenkinsServers[k].Name + " status: ***unavailable*** unexpected status code: 500\n"
				} else {
					request, _ := client.PrepareForEmptyInstalledPluginList(roundTripper, jenkinsServers[k].URL)
					request.SetBasicAuth(jenkinsServers[k].UserName, jenkinsServers[k].Token)
					outString += "  checking the No." + string(k) + " - " + jenkinsServers[k].Name + " status: ***available***\n"
				}
			}
			current := getCurrentJenkins()
			outString += "Begining checking the current JenkinsServer's plugins status: \n"
			request, _ := client.PrepareFor500InstalledPluginList(roundTripper, current.URL,2)
			request.SetBasicAuth(current.UserName, current.Token)
			outString += "  No plugins have updated...\n"
			outString += "Checked is done."
			rootCmd.SetArgs([]string{"doctor"})
			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())
			//fmt.Println(outString)
			//Expect(buf.String()).To(Equal(outString))
		})
	})
})
