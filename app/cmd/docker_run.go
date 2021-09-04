package cmd

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(dockerRunCmd)
	dockerRunCmd.Flags().StringVarP(&dockerRunOptions.ImageName, "image-name", "", "test",
		i18n.T("name of the image which contains upgraded jenkins and plugins"))
	dockerRunCmd.Flags().StringVarP(&dockerRunOptions.IP, "docker-ip", "", "127.0.0.1",
		i18n.T("the ip address of the computer you want to use"))
	dockerRunCmd.Flags().StringVarP(&dockerRunOptions.Tag, "tag", "", "latest",
		i18n.T("the tag of the images"))
	dockerRunCmd.Flags().IntVarP(&dockerRunOptions.DockerPort, "docker-port", "", 2375,
		i18n.T("the port to connect to docker"))
	dockerRunCmd.Flags().StringVarP(&dockerRunOptions.DockerfilePath, "dockerfile-path", "", "/tmp/Dockerfile",
		i18n.T("where you want the dockerfile to be placed"))
	dockerRunCmd.Flags().IntVarP(&dockerRunOptions.JenkinsPort, "jenkins-port", "", 8081,
		i18n.T("The port you want to used to connect to jenkins"))
	dockerRunCmd.Flags().StringVarP(&dockerRunOptions.WarPath, "war-path", "", "",
		i18n.T("where your war is placed"))
}

//DockerRunOptions contains some of the options used to create a docker image and run a container
type DockerRunOptions struct {
	ImageName      string
	Tag            string
	IP             string
	DockerPort     int
	DockerfilePath string
	JenkinsPort    int
	WarPath        string
}

var dockerRunOptions DockerRunOptions

var dockerRunCmd = &cobra.Command{
	Use:   "docker run",
	Short: i18n.T("Create a image and start a container in docker where all upgraded plugins and jenkins run in order to test their eligibility"),
	Long: i18n.T("This command is used to create a image and start a container, where all upgraded plugins and jenkins run in. This command relies on the war file provided by Custom War Packager.\n" +
		"The war can be made by calling 'jcli cwp --config-path <yamlfilename>' and the yaml file will be created using the command 'jcli create yaml'."),
	PreRunE: dockerRunOptions.createDockerfile,
	RunE:    dockerRunOptions.createImageAndRunContainer,
	Example: `jcli docker run --docker-ip <ip> --docker-port <port> --war-path <pathToWar>`,
}

//GetDockerIPAndPort returns a string contains IP and port of a local or remote host
func (o *DockerRunOptions) GetDockerIPAndPort() string {
	ip := o.IP
	port := o.DockerPort
	return fmt.Sprintf("tcp://%s:%d", ip, port)
}

//ConnectToDocker returns a client which is used to connect to a local or remote docker host
func (o *DockerRunOptions) ConnectToDocker() (cli *client.Client, err error) {
	tcp := fmt.Sprintf("tcp://%s:%d", o.IP, o.DockerPort)
	cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation(), client.WithHost(tcp))
	return cli, err
}

func (o *DockerRunOptions) createImageAndRunContainer(cmd *cobra.Command, args []string) (err error) {
	ctx := context.Background()
	cli, err := o.ConnectToDocker()
	if err != nil {
		cmd.Println(err)
		return err
	}
	imageName := o.ImageName + ":" + o.Tag
	err = o.buildImage(cmd)
	if err != nil {
		cmd.Println(err)
	}
	jenkinsPort, err := nat.NewPort("tcp", "8080")
	if err != nil {
		cmd.Println(err)
	}
	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			jenkinsPort: []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: strconv.Itoa(o.JenkinsPort),
				},
			},
		},
	}
	exposedPorts := map[nat.Port]struct{}{
		jenkinsPort: struct{}{},
	}
	config := &container.Config{
		Image:        imageName,
		ExposedPorts: exposedPorts,
	}

	resp, err := cli.ContainerCreate(ctx, config, hostConfig, nil, nil, "jenkinstest2")
	if err != nil {
		fmt.Println(err)
	}
	cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	fmt.Println(resp.ID)
	return nil
}

func (o *DockerRunOptions) buildImage(cmd *cobra.Command) error {
	ctx := context.Background()
	cli, _ := o.ConnectToDocker()
	dockerFileTarReader, _ := o.getTarReader(cmd)
	buildOptions := types.ImageBuildOptions{
		Context:    dockerFileTarReader,
		Dockerfile: o.DockerfilePath,
		Remove:     true,
		Tags:       []string{o.ImageName},
	}
	imageBuildResponse, err := cli.ImageBuild(
		ctx,
		dockerFileTarReader,
		buildOptions,
	)

	if err != nil {
		return err
	}
	defer imageBuildResponse.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(imageBuildResponse.Body)
	cmd.Println(buf.String())
	return nil
}
func (o *DockerRunOptions) getTarReader(cmd *cobra.Command) (*bytes.Reader, error) {
	src := []string{o.DockerfilePath, o.WarPath, "/tmp/jenkins-cli-docker.sh"}
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	for _, fileName := range src {
		dockerFileReader, err := os.Open(fileName)
		if err != nil {
			return nil, err
		}
		readDockerFile, err := ioutil.ReadAll(dockerFileReader)
		if err != nil {
			return nil, err
		}
		tarHeader := &tar.Header{
			Name: fileName,
			Size: int64(len(readDockerFile)),
		}
		err = tw.WriteHeader(tarHeader)
		if err != nil {
			return nil, err
		}
		_, err = tw.Write(readDockerFile)
		if err != nil {
			return nil, err
		}
	}
	dockerFileTarReader := bytes.NewReader(buf.Bytes())
	return dockerFileTarReader, nil
}
func (o *DockerRunOptions) createDockerfile(cmd *cobra.Command, args []string) (err error) {
	sh := "java -jar /usr/share/jenkins/jenkins.war\n" +
		"tail -f /dev/null"
	err = ioutil.WriteFile("/tmp/jenkins-cli-docker.sh", []byte(sh), 0)
	if err != nil {
		panic(err)
	}
	dockerfileString := fmt.Sprintf("FROM adoptopenjdk/openjdk11\n"+
		"LABEL Version=\"1.0-SNAPSHOT\"Description=\"Jenkins formula generated by jcli\"Vendor=\"Chinese Jenkins Community\"\n"+
		"ADD %s /usr/share/jenkins/jenkins.war\n"+
		"ADD /tmp/jenkins-cli-docker.sh /usr/local/bin/jenkins-cli-docker.sh\n"+
		"ENTRYPOINT [\"sh\", \"/usr/local/bin/jenkins-cli-docker.sh\"]", o.WarPath)
	err = ioutil.WriteFile("/tmp/Dockerfile", []byte(dockerfileString), 0)
	if err != nil {
		panic(err)
	}
	return nil
}
