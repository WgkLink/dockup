package dockup

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
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

func UpdateImages(listImage *[]Images) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		panic(err)
	}
	defer cli.Close()

	for _, l := range *listImage {
		out, err := cli.ImagePull(ctx, l.name, image.PullOptions{})

		if err != nil {
			panic(err)
		}

		defer out.Close()

		io.Copy(os.Stdout, out)
	}
}
