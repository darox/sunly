package printer

import (
	"fmt"
	"time"

	"github.com/darox/sunly/pkg/swissmeteo"
	"github.com/jedib0t/go-pretty/v6/table"
)

// Map Icon ID to Emoji
func getEmojiForIcon(iconID string) string {
	switch {
	case iconID == "1":
		return "🔆"
	case iconID == "2":
		return "🌤️"
	case iconID == "3":
		return "⛅"
	case iconID == "4":
		return "⛅"
	case iconID == "5":
		return "🌦️"
	case iconID == "6":
		return "🌦️"
	case iconID == "7":
		return "🌨️"
	case iconID == "8":
		return "🌨️"
	case iconID == "9":
		return "🌧️"
	case iconID == "10":
		return "🌧️"
	case iconID == "11":
		return "🌨️"
	case iconID == "12":
		return "🌩️"
	case iconID == "13":
		return "🌩️"
	case iconID == "14":
		return "🌨️"
	case iconID == "15":
		return "🌨️"
	case iconID == "16":
		return "🌨️"
	case iconID == "17":
		return "🌨️"
	case iconID == "18":
		return "🌨️"
	case iconID == "19":
		return "🌨️"
	case iconID == "20":
		return "🌨️"
	case iconID == "21":
		return "🌨️"
	case iconID == "22":
		return "🌨️"
	case iconID == "23":
		return "🌩️"
	case iconID == "24":
		return "🌩️"
	case iconID == "25":
		return "🌩️"
	case iconID == "26":
		return "🔆💨"
	case iconID == "27":
		return "🔆🌁"
	case iconID == "28":
		return "🌁"
	case iconID == "29":
		return "🌦️"
	case iconID == "30":
		return "🌨️"
	case iconID == "31":
		return "🌦️"
	case iconID == "32":
		return "🌦️"
	}

	return ""
}

func PrintCurrentWeather(w swissmeteo.CurrentWeather) {
	t := table.NewWriter()

	t.AppendHeader(table.Row{"Zip", "City", "Condition", "Temperature", "Updated at"})

	temp := fmt.Sprintf("%s °C", w.Current.Temperature)

	t.AppendRows([]table.Row{
		{w.Zip, w.CityName, getEmojiForIcon(w.Current.WeatherSymbolID), temp, convertTime(w.Timestamp)},
	})

	fmt.Print(t.Render())

}

func convertTime(u int) (updatedAt string) {
	// Convert time to a human readable format
	h := time.Unix(int64(u), 0)
	return h.Format("15:04 02.01.2006")
}
