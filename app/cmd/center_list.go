package cmd

import (
	"encoding/xml"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"strings"
)

const LTSURL = "https://www.jenkins.io/changelog-stable/rss.xml"
const WEEKLYURL = "https://www.jenkins.io/changelog/rss.xml"

//the width of the description column
const WIDTH_OF_DESCRIPTION = 60
const NUMBER_OF_LINES_OF_DESCRIPTION = 10
const ASCII_OF_LINE_FEED = 10
const ASCII_OF_SPACE = 32

type CenterListOption struct {
	Channel      Channel `xml:"channel"`
	RoundTripper http.RoundTripper
}
type Channel struct {
	Title string `xml:"title"`
	Items []Item `xml:"item"`
}
type Item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

var centerListOption CenterListOption

func init() {
	centerCmd.AddCommand(centerListCmd)
}

var centerListCmd = &cobra.Command{
	Use:   "list",
	Short: i18n.T("Print the information of recent-released Jenkins"),
	Long:  i18n.T("Print the information of recent-released Jenkins"),
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		jenkins := getCurrentJenkinsFromOptionsOrDie()
		jclient := &client.JenkinsStatusClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: centerOption.RoundTripper,
			},
		}
		jclient.URL = jenkins.URL
		jclient.UserName = jenkins.UserName
		jclient.Token = jenkins.Token
		jclient.Proxy = jenkins.Proxy
		jclient.ProxyAuth = jenkins.ProxyAuth
		status, error := jclient.Get()
		if error != nil {
			return error
		}
		jenkinsVersion := status.Version
		changeLog, err := printChangelog(LTSURL, jenkinsVersion)
		cmd.Println(changeLog)
		//err = printChangelog(WEEKLYURL,cmd)
		return err
	},
}

func printChangelog(rss string, version string) (changelog string, err error) {
	resp, err := http.Get(rss)
	if err != nil {
		return "", err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var centerListOption CenterListOption
	xml.Unmarshal(bytes, &centerListOption)
	t := table.NewWriter()
	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}
	t.AppendHeader(table.Row{centerListOption.Channel.Title, centerListOption.Channel.Title, centerListOption.Channel.Title, centerListOption.Channel.Title}, rowConfigAutoMerge)
	t.AppendRow(table.Row{"Index", "Title", "Description", "PubDate"})
	var temp string
	var isTheLatestVersion = 1
	for index, item := range centerListOption.Channel.Items {
		if compareVersionTitle(version, item.Title[8:]) >= 0 {
			break
		}
		isTheLatestVersion = 0
		temp = trimXMLSymbols(item.Description)
		temp = regulateWidthAndLines(temp, WIDTH_OF_DESCRIPTION, NUMBER_OF_LINES_OF_DESCRIPTION)
		t.AppendRow([]interface{}{index + 1, item.Title, temp, item.PubDate[:17]})
		t.AppendSeparator()
	}
	resp.Body.Close()
	if isTheLatestVersion == 1 {
		return "You already have the latest version of Jenkins installed!", nil
	} else {
		return t.Render(), nil
	}
}
func regulateWidthAndLines(content string, width int, numberOfLines int) string {
	var count = 0
	myContent := []uint8(content)
	var i int
	for i = 0; i < len(myContent); i++ {
		if myContent[i] == ASCII_OF_LINE_FEED {
			count++
		}
		if count == numberOfLines {
			break
		}
	}
	myContent = myContent[:i]
	indexInLine := 0
	for index, char := range myContent {
		indexInLine++
		if indexInLine%width == 0 {
			if char == ASCII_OF_SPACE {
				myContent[index] = ASCII_OF_LINE_FEED
				indexInLine = 0
			} else {
				indexInLine = 0
				for ; myContent[index] != ASCII_OF_SPACE; index-- {
					indexInLine++
				}
				myContent[index] = ASCII_OF_LINE_FEED
			}
		}
		if char == ASCII_OF_LINE_FEED {
			indexInLine = 0
		}
	}
	return string(myContent)
}

func compareVersionTitle(versionOne string, versionTwo string) int {
	versionOneString := strings.Split(versionOne, ".")
	versionTwoString := strings.Split(versionTwo, ".")
	if strings.Compare(versionOne, versionTwo) == 0 {
		return 0
	}
	if versionOneString[0] > versionTwoString[0] {
		return 1
	} else if versionOneString[0] < versionTwoString[0] {
		return -1
	}
	if versionOneString[1] != "" && versionTwoString[1] != "" && versionOneString[1] > versionTwoString[1] {
		return 1
	} else if versionOneString[1] < versionTwoString[1] {
		return -1
	}
	if versionOneString[2] != "" && versionTwoString[2] != "" && versionOneString[2] > versionTwoString[2] {
		return 1
	}
	return -1
}

func trimXMLSymbols(temp string) string {
	temp = strings.TrimSpace(temp)
	temp = strings.Replace(temp, "<br>", "", -1)
	temp = strings.Replace(temp, "<br/>", "", -1)
	temp = strings.Replace(temp, "<br />", "", -1)
	temp = strings.Replace(temp, "<ul>", "", -1)
	temp = strings.Replace(temp, "<li>", "", -1)
	temp = strings.Replace(temp, "</li>", "\n", -1)
	temp = strings.Replace(temp, "</ul>", "", -1)
	temp = strings.Replace(temp, "<strong>", "", -1)
	temp = strings.Replace(temp, "</strong>", "", -1)
	temp = strings.Replace(temp, "<code>", "", -1)
	temp = strings.Replace(temp, "</code>", "", -1)
	temp = strings.Replace(temp, "<em>", "", -1)
	temp = strings.Replace(temp, "</em>", "", -1)
	return temp
}
