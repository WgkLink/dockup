package main

import (
	"dockup/dockup"
	"log"

	"github.com/charmbracelet/huh"
)

func main() {

	//var selectedList []string
	var confirm bool
	var selectedList []string
	var bindings string

	listImage := dockup.ListImages()

	form := huh.NewForm(
		huh.NewGroup(

			huh.NewMultiSelect[string]().
				Title("Choice a image to updade").
				OptionsFunc(func() []huh.Option[string] {
					var options []huh.Option[string]
					for _, l := range listImage {
						options = append(options, huh.NewOption(l.Name, l.Name))
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
	dockup.UpdateImages(selectedList)
}
