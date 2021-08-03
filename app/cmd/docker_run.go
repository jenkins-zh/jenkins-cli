package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
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
	dockerRunCmd.Flags().IntVarP(&dockerRunOptions.Port, "port", "", 2375,
		i18n.T("The port to connect"))
}

type DockerRunOptions struct {
	ImageName string
	Tag       string
	IP        string
	Port      int
}

var dockerRunOptions DockerRunOptions

var dockerRunCmd = &cobra.Command{
	Use:     "docker run",
	Short:   i18n.T("Start a container in docker where all upgraded plugins and jenkins run in order to test their eligibility"),
	Long:    i18n.T("Start a container, where all upgraded plugins and jenkins run, using a image built by Jenkins WAR packager in order to test their eligibility"),
	PreRunE: dockerRunOptions.CheckImageExistsOrNot,
	Example: `jcli docker run`,
	RunE:    dockerRunOptions.PullImageAndRunContainer,
}

func (o *DockerRunOptions) GetDockerIPAndPort() string {
	ip := o.IP
	port := o.Port
	return fmt.Sprintf("tcp://%s:%d", ip, port)
}

func (o *DockerRunOptions) PullImageAndRunContainer(cmd *cobra.Command, args []string) (err error) {
	tcp := o.GetDockerIPAndPort()
	ctx := context.Background()
	cmd.Println(ctx)
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation(), client.WithHost(tcp))
	if err != nil {
		cmd.Print("1 ")
		cmd.Println(err)
	}

	imageName := o.ImageName + ":" + o.Tag

	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		cmd.Print("2 " + imageName)
		cmd.Println(err)
	}
	cmd.Println(out)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
	}, nil, nil, nil, "")
	if err != nil {
		cmd.Print("3 ")
		cmd.Println(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		cmd.Print("4 ")
		cmd.Println(err)
	}

	cmd.Println("The container ID is " + resp.ID)
	return nil
}
func (o *DockerRunOptions) CheckImageExistsOrNot(cmd *cobra.Command, args []string) (err error) {
	ip := fmt.Sprintf("https://index.docker.io/v1/repositories/%s/tags/%s", o.ImageName, o.Tag)
	resp, err := http.Get(ip)
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if string(bytes) != "" {
		return err
	}
	return nil
}
