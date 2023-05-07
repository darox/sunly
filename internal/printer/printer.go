package printer

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
)

func PrintCurrentTemperature(zip string, location string, temperature float64, updatedAt string) {
	t := table.NewWriter()

	t.AppendHeader(table.Row{"Zip", "Location", "Temperature", "Updated at"})

	c := fmt.Sprintf("%.1f °C", temperature)
	t.AppendRows([]table.Row{
		{zip, location, c, updatedAt},
	})

	fmt.Print(t.Render())
}
