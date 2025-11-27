package ui

import (
	"github.com/charmbracelet/huh"
)

// ShowForm displays the options and returns the selected value
func ShowForm(options []huh.Option[int], title string) (int, error) {
	var selected int
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title(title).
				Options(options...).
				Value(&selected),
		),
	)

	if err := form.Run(); err != nil {
		return 0, err
	}

	return selected, nil
}
