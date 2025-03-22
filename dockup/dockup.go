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
	id     int
	name   string
	digest string
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
			tmp[index] = Images{id: index, name: l.RepoTags[0], digest: l.ID}
		} else {
			tmp = append(tmp, Images{id: index, name: l.RepoTags[0], digest: l.ID})
		}
	}
	return tmp
}

func UpdateImages(listImage []Images) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		panic(err)
	}
	defer cli.Close()

	for _, l := range listImage {
		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
		s.Prefix = l.name + " "
		s.Start()
		out, err := cli.ImagePull(ctx, l.name, image.PullOptions{})

		s.Stop()

		if err != nil {
			red := color.New(color.FgRed).SprintFunc()
			fmt.Printf("%s %s \n", l.name, red("FAILED"))

		} else {
			green := color.New(color.FgGreen).SprintFunc()
			fmt.Printf("%s %s \n", l.name, green("UPDATED"))
		}

		defer out.Close()
	}
}
