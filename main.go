package main

import (
	"dockup/dockup"
	"log"
	"strings"

	"github.com/charmbracelet/huh"
)

func main() {

	var confirm bool
	var selectedList []string
	var bindings string

	var digestImages []string

	imageList := dockup.ImageList()

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Choice a image to update").
				OptionsFunc(func() []huh.Option[string] {
					var options []huh.Option[string]
					for _, l := range imageList {
						options = append(options, huh.NewOption(l.Name, l.Digest))
					}
					return options
				}, bindings).
				Value(&selectedList),
		),
		huh.NewGroup(
			huh.NewConfirm().
				Title("Do you want to restart the containers that use the selected images?").
				Value(&confirm),
		),
	)

	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}

	for _, selected := range selectedList {
		SelectedPrefixRemoved := strings.TrimPrefix(selected, "sha256:")
		for _, image := range imageList {
			ImagePrefixRemoved := strings.TrimPrefix(image.Digest, "sha256:")
			if SelectedPrefixRemoved == ImagePrefixRemoved {
				digestImages = append(digestImages, ImagePrefixRemoved)
			}
		}
	}

	if len(selectedList) > 0 {
		dockup.UpdateImages(digestImages)

		if confirm {
			dockup.RestartContainers(digestImages)
		}
	}

}
