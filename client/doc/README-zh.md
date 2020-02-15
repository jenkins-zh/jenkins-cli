# jcli-client 文档

## 如何安装

go.mod
```
require github.com/jenkins-zh/jenkins-cli
```

## 配置

jcli配置详情
```
- name: dev
  url: http://192.168.1.10
  username: admin
  token: 11132c9ae4b20edbe56ac3e09cb5a3c8c2
  proxy: http://192.168.10.10:47586
  proxyAuth: username:password
```
所以当你想使用jcli-client，你需要配置一下jenkins.yml，并解析
```
type Config struct {
	URL      string `yaml:"url"`
	UserName string `yaml:"username"`
	Token    string `yaml:"token"`
}

func GetJenkinsCore() (core client.JenkinsCore, e error) {
	// 读取配置文件
	jenkinsConfigPath := "./jenkins.yml"
	yamlFile, e := ioutil.ReadFile(jenkinsConfigPath)
	if e != nil {
		return
	}
	// yaml解析
	var config Config
	e = yaml.Unmarshal(yamlFile, &config)
	if e != nil {
		return
	}
	// 封装JenkinsCore
	core = client.JenkinsCore{
		URL:      config.URL,
		UserName: config.UserName,
		Token:    config.Token,
	}
	crumbIssuer, e := core.GetCrumb()
	if e != nil {
		return
	} else if crumbIssuer != nil {
		core.JenkinsCrumb = *crumbIssuer
	}
	return
}
```

## Job API
常用API
```
// GetBuild get build information of a job
func (q *JobClient) GetBuild(jobName string, id int) (job *JobBuild, err error){...}

// Build trigger a job
func (q *JobClient) Build(jobName string) (err error) {...}

// BuildWithParams build a job which has params
func (q *JobClient) BuildWithParams(jobName string, parameters []ParameterDefinition) (err error) {...}

// GetHistory returns the build history of a job
func (q *JobClient) GetHistory(name string) (builds []*JobBuild, err error) {...}

// Log get the log of a job
func (q *JobClient) Log(jobName string, history int, start int64) (jobLog JobLog, err error){...}
```

## User API
常用API
```
// Create will create a user in Jenkins
func (q *UserClient) Create(username, password string) (user *UserForCreate, err error) {...}

// Delete will remove a user from Jenkins
func (q *UserClient) Delete(username string) (err error) {...}
```

## 调用案例
```
type JobBuildOptions struct {
	Env       string 
	JobName   string 
	BranchTag string 
}


func BuildJob(jobBuild JobBuildOptions) (e error) {
	core, e := GetJenkinsCore()
	jobClient = client.JobClient{core}
	if e != nil {
		return
	}
	param1 := client.ParameterDefinition{
		Description:           "pre 和 prd 环境使用 tag 发布",
		Name:                  "BRANCH_TAG",
		Type:                  "Branch or Tag",
		Value:                 jobBuild.BranchTag,
		DefaultParameterValue: client.DefaultParameterValue{Value: "origin/master"},
	}
	param2 := client.ParameterDefinition{
		Description:           "环境选择",
		Name:                  "ENV",
		Type:                  "Choice Parameter",
		Value:                 jobBuild.Env,
		DefaultParameterValue: client.DefaultParameterValue{Value: "qa"},
	}
	param3 := client.ParameterDefinition{
		Description:           "部署Type",
		Name:                  "DEPLOY_TYPE",
		Type:                  "Choice Parameter",
		Value:                 "publish",
		DefaultParameterValue: client.DefaultParameterValue{Value: "publish"},
	}
	params := []client.ParameterDefinition{param1, param2, param3}
	e = jobClient.BuildWithParams(jobBuild.JobName, params)
	return
}
```

## 插件依赖

本 Jenkins 客户端的部分 API 需要依赖特定的插件，请参考下面的列表：

| API | Plugin |
|---|---|
| Search Job | [pipeline-restful-api](https://github.com/jenkinsci/pipeline-restful-api-plugin) |
