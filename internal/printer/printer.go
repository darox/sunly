package printer

import (
	"fmt"
	"time"

	"github.com/darox/sunly/pkg/swissmeteo"
	"github.com/jedib0t/go-pretty/v6/table"
)

func PrintCurrentTemperature(w swissmeteo.CurrentWeather, zip string, city string) {
	t := table.NewWriter()

	t.AppendHeader(table.Row{"Zip", "City", "Temperature", "Updated at"})

	c := fmt.Sprintf("%.1f °C", w.Temperature)
	t.AppendRows([]table.Row{
		{zip, city, c, convertTime(w.Time)},
	})

	fmt.Print(t.Render())
}

func convertTime(u int64) (updatedAt string) {
	// Convert time to a human readable format
	h := time.Unix(u/1000, 0)
	return h.Format("15:04 02.01.2006")
}
