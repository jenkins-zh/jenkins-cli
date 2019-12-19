package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
)

// PrepareForComputerListRequest only for test
func PrepareForComputerListRequest(roundTripper *mhttp.MockRoundTripper, rootURL, user, password string) {
	request, _ := http.NewRequest("GET", fmt.Sprintf("%s/computer/api/json", rootURL), nil)
	response := &http.Response{
		StatusCode: 200,
		Request:    request,
		Body:       ioutil.NopCloser(bytes.NewBufferString(PrepareForComputerList())),
	}
	roundTripper.EXPECT().
		RoundTrip(request).Return(response, nil)
	if user != "" && password != "" {
		request.SetBasicAuth(user, password)
	}
}

// PrepareForLaunchComputer only for test
func PrepareForLaunchComputer(roundTripper *mhttp.MockRoundTripper, rootURL, user, password, name string) {
	request, _ := http.NewRequest("POST", fmt.Sprintf("%s/computer/%s/launchSlaveAgent", rootURL, name), nil)
	PrepareCommonPost(request, "", roundTripper, user, password, rootURL)
}

// PrepareForComputerLogRequest only for test
func PrepareForComputerLogRequest(roundTripper *mhttp.MockRoundTripper, rootURL, user, password, name string) {
	PrepareForComputerLogRequestWithCode(roundTripper, rootURL, user, password, name, 200)
}

// PrepareForComputerLogRequestWithCode only for test
func PrepareForComputerLogRequestWithCode(roundTripper *mhttp.MockRoundTripper, rootURL, user, password,
	name string, statusCode int) {
	request, _ := http.NewRequest("GET", fmt.Sprintf("%s/computer/%s/logText/progressiveText", rootURL, name), nil)
	response := &http.Response{
		StatusCode: statusCode,
		Request:    request,
		Body:       ioutil.NopCloser(bytes.NewBufferString(`fake-log`)),
	}
	roundTripper.EXPECT().
		RoundTrip(request).Return(response, nil)
	if user != "" && password != "" {
		request.SetBasicAuth(user, password)
	}
}

// PrepareForComputerList only for test
func PrepareForComputerList() string {
	return `
{
  "_class" : "hudson.model.ComputerSet",
  "busyExecutors" : 1,
  "computer" : [
    {"offline" : true},
    {
      "_class" : "hudson.model.Hudson$MasterComputer",
      "actions" : [
        {
          
        },
        {
          
        }
      ],
      "assignedLabels" : [
        {
          "name" : "master"
        }
      ],
      "description" : "Jenkins的master节点",
      "displayName" : "master",
      "executors" : [
        {
          
        },
        {
          
        }
      ],
      "icon" : "computer.png",
      "iconClassName" : "icon-computer",
      "idle" : false,
      "jnlpAgent" : false,
      "launchSupported" : true,
      "loadStatistics" : {
        "_class" : "hudson.model.Label$1"
      },
      "manualLaunchAllowed" : true,
      "monitorData" : {
        "hudson.node_monitors.SwapSpaceMonitor" : {
          "_class" : "hudson.node_monitors.SwapSpaceMonitor$MemoryUsage2",
          "availablePhysicalMemory" : 533508096,
          "availableSwapSpace" : 0,
          "totalPhysicalMemory" : 16656797696,
          "totalSwapSpace" : 0
        },
        "hudson.node_monitors.TemporarySpaceMonitor" : {
          "_class" : "hudson.node_monitors.DiskSpaceMonitorDescriptor$DiskSpace",
          "timestamp" : 1574952137138,
          "path" : "/tmp",
          "size" : 31745785856
        },
        "hudson.node_monitors.DiskSpaceMonitor" : {
          "_class" : "hudson.node_monitors.DiskSpaceMonitorDescriptor$DiskSpace",
          "timestamp" : 1574952136930,
          "path" : "/var/jenkins_home",
          "size" : 31745785856
        },
        "hudson.node_monitors.ArchitectureMonitor" : "Linux (amd64)",
        "hudson.node_monitors.ResponseTimeMonitor" : {
          "_class" : "hudson.node_monitors.ResponseTimeMonitor$Data",
          "timestamp" : 1574952136931,
          "average" : 0
        },
        "hudson.node_monitors.ClockMonitor" : {
          "_class" : "hudson.util.ClockDifference",
          "diff" : 0
        }
      },
      "numExecutors" : 2,
      "offline" : false,
      "offlineCause" : null,
      "offlineCauseReason" : "",
      "oneOffExecutors" : [
        {
          
        }
      ],
      "temporarilyOffline" : false
    }
  ],
  "displayName" : "节点列表",
  "totalExecutors" : 2
}`
}
