package dockup

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/fatih/color"
)

type Images struct {
	Id     int
	Name   string
	Digest string
}

func ListImages() []Images {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		panic(err)
	}
	defer cli.Close()

	tmp := []Images{{}}

	listImage, err := cli.ImageList(ctx, image.ListOptions{})

	if err != nil {
		panic(err)
	}
	for index, l := range listImage {
		if index == 0 {
			tmp[index] = Images{Id: index, Name: l.RepoTags[0], Digest: l.ID}
		} else {
			tmp = append(tmp, Images{Id: index, Name: l.RepoTags[0], Digest: l.ID})
		}
	}
	return tmp
}

func UpdateImages(listImage []string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		panic(err)
	}
	defer cli.Close()

	for _, l := range listImage {
		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
		s.Prefix = l + " "
		s.Start()
		out, err := cli.ImagePull(ctx, l, image.PullOptions{})

		s.Stop()

		if err != nil {
			red := color.New(color.FgRed).SprintFunc()
			fmt.Printf("%s %s \n", l, red("FAILED"))

		} else {
			green := color.New(color.FgGreen).SprintFunc()
			fmt.Printf("%s %s \n", l, green("UPDATED"))
		}

		defer out.Close()
	}
}
