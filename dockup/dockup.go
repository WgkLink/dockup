package dockup

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/fatih/color"
)

type Images struct {
	Id     int
	Name   string
	Digest string
}

func ImageList() []Images {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		panic(err)
	}
	defer cli.Close()

	imageList := []Images{}

	images, err := cli.ImageList(ctx, image.ListOptions{})

	if err != nil {
		panic(err)
	}

	for index, l := range images {

		if len(l.RepoTags) > 0 {
			imageList = append(imageList, Images{
				Id:     index,
				Name:   l.RepoTags[0],
				Digest: l.ID},
			)
		} else {
			imageList = append(imageList, Images{
				Id:     index,
				Name:   l.ID,
				Digest: l.ID},
			)
		}
	}
	return imageList
}

func UpdateImages(imageList []string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		panic(err)
	}
	defer cli.Close()

	for _, l := range imageList {
		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
		s.Prefix = l + " "
		s.Start()
		_, err := cli.ImagePull(ctx, l, image.PullOptions{})

		s.Stop()

		if err != nil {
			red := color.New(color.FgRed).SprintFunc()
			fmt.Printf("%s %s \n", l, red("FAILED"))

		} else {
			green := color.New(color.FgGreen).SprintFunc()
			fmt.Printf("%s %s \n", l, green("UPDATED"))
		}
	}
}

func RestartContainers(imageList []string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, c := range containers {
		for _, image := range imageList {

			ImageIdPrefixRemoved := strings.TrimPrefix(c.ImageID, "sha256:")

			if ImageIdPrefixRemoved == image {
				cli.ContainerRestart(ctx, c.ID, container.StopOptions{})

				cleanedName := strings.TrimPrefix(c.Names[0], "/")

				blue := color.New(color.FgBlue).SprintFunc()
				fmt.Printf("%s %s ", blue(cleanedName), blue("HAS BEEN RESTARTED\n"))
			}
		}
	}
}
