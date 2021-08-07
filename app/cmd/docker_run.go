package cmd

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
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
	dockerRunCmd.Flags().StringVarP(&dockerRunOptions.ImageName, "image-name", "", "",
		i18n.T("Name of the image in docker hub which contains upgraded jenkins and plugins"))
	dockerRunCmd.Flags().StringVarP(&dockerRunOptions.IP, "ip", "", "127.0.0.1",
		i18n.T("The ip address of the computer you want to use"))
	dockerRunCmd.Flags().StringVarP(&dockerRunOptions.Tag, "tag", "", "latest",
		i18n.T("The tag of the images"))
	dockerRunCmd.Flags().IntVarP(&dockerRunOptions.DockerPort, "docker-port", "", 2375,
		i18n.T("The port to connect to docker"))
	dockerRunCmd.Flags().StringVarP(&dockerRunOptions.DockerfilePath, "dockerfile-path", "", "./tmp/output/Dockerfile",
		i18n.T("where you want the dockerfile to be placed"))
	dockerRunCmd.Flags().IntVarP(&dockerRunOptions.jenkinsPort, "jenkins-port", "", 8081,
		i18n.T("The port to connect to jenkins"))
}

type DockerRunOptions struct {
	ImageName      string
	Tag            string
	IP             string
	DockerPort     int
	DockerfilePath string
	jenkinsPort    int
}

var dockerRunOptions DockerRunOptions

var dockerRunCmd = &cobra.Command{
	Use:     "docker run",
	Short:   i18n.T("Start a container in docker where all upgraded plugins and jenkins run in order to test their eligibility"),
	Long:    i18n.T("Start a container, where all upgraded plugins and jenkins run, using a image built by Jenkins WAR packager in order to test their eligibility"),
	RunE:    dockerRunOptions.PullImageAndRunContainer,
	Example: `jcli docker run`,
}

func (o *DockerRunOptions) GetDockerIPAndPort() string {
	ip := o.IP
	port := o.DockerPort
	return fmt.Sprintf("tcp://%s:%d", ip, port)
}

func (o *DockerRunOptions) ConnectToDocker() (cli *client.Client, err error) {
	tcp := fmt.Sprintf("tcp://%s:%d", o.IP, o.DockerPort)
	cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation(), client.WithHost(tcp))
	return cli, err
}

func (o *DockerRunOptions) PullImageAndRunContainer(cmd *cobra.Command, args []string) (err error) {
	ctx := context.Background()
	cli, err := o.ConnectToDocker()
	if err != nil {
		cmd.Print("1. ")
		cmd.Println(err)
		return err
	}
	imageName := o.ImageName + ":" + o.Tag
	if o.CheckImageExistsInDockerHub(cmd) {
		reader, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
		if err != nil {
			cmd.Print("2. ")
			cmd.Println(err)
		}
		cmd.Print(reader)
	} else {
		err := o.BuildImage(cmd)
		if err != nil {
			cmd.Print("3. ")
			cmd.Println(err)
		}
	}
	newPort, err := nat.NewPort("tcp", strconv.Itoa(o.jenkinsPort))
	if err != nil {
		cmd.Print("4. ")
		cmd.Println(err)
		return err
	}
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Cmd:   []string{"java -jar", "/usr/share/jenkins/jenkins.war", "--server.port=8080"},
	}, &container.HostConfig{
		PortBindings: nat.PortMap{
			newPort: []nat.PortBinding{
				{
					HostIP:   "127.0.0.1",
					HostPort: strconv.Itoa(o.jenkinsPort),
				},
			},
		},
	}, nil, nil, "jenkinsTest")
	if err != nil {
		cmd.Print("5. ")
		fmt.Println(err)
	}
	fmt.Println(resp.ID)
	return nil
}
func (o *DockerRunOptions) CheckImageExistsInDockerHub(cmd *cobra.Command) bool {
	ip := fmt.Sprintf("https://index.docker.io/v1/repositories/%s/tags/%s", o.ImageName, o.Tag)
	resp, err := http.Get(ip)
	if err != nil {
		cmd.Print("6. ")
		cmd.Println(err)
		return false
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if string(bytes) != "" {
		cmd.Print("7. ")
		cmd.Println(err)
		return false
	}
	return true
}
func (o *DockerRunOptions) BuildImage(cmd *cobra.Command) error {
	ctx := context.Background()
	cli, err := o.ConnectToDocker()
	if err != nil {
		cmd.Print("8. ")
		cmd.Println(err)
		return err
	}
	file, err := os.Open(o.DockerfilePath)
	if err != nil {
		cmd.Print("10. ")
		cmd.Println(err)
		return err
	}
	readDockerfile, err := ioutil.ReadAll(file)
	if err != nil {
		cmd.Print("11. ")
		cmd.Println(err)
		return err
	}
	tarHeader := &tar.Header{
		Name: o.DockerfilePath,
		Size: int64(len(readDockerfile)),
	}
	buffer := new(bytes.Buffer)
	tw := tar.NewWriter(buffer)
	defer tw.Close()
	err = tw.WriteHeader(tarHeader)
	if err != nil {
		cmd.Print("12. ")
		cmd.Println(err)
		return err
	}
	_, err = tw.Write(readDockerfile)
	if err != nil {
		cmd.Print("13. ")
		cmd.Println(err)
		return err
	}
	dockerFileTarReader := bytes.NewReader(buffer.Bytes())
	opts := types.ImageBuildOptions{
		Context:    dockerFileTarReader,
		Dockerfile: o.DockerfilePath,
		Tags:       []string{o.ImageName, ":" + o.Tag},
		Remove:     true,
	}
	resp, err := cli.ImageBuild(ctx, dockerFileTarReader, opts)
	if err != nil {
		cmd.Print("9. ")
		cmd.Println(err)
		return err
	}
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	cmd.Println(buf.String())
	return nil
}
