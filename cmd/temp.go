/*
Copyright © 2023 Dario Mader maderdario@gmail.com
*/
package cmd

import (
	"fmt"

	"github.com/darox/sunly/internal/swissmeteo"
	"github.com/spf13/cobra"
)

// tempCmd represents the temp command
var tempCmd = &cobra.Command{
	Use:   "temp",
	Short: "Returns the temperature of a location by providing a postal code",
	Long:  `Returns the temperature of a location by providing a postal code`,
	Run: func(cmd *cobra.Command, args []string) {
		getTemp(args[0])
	},
}

func init() {
	rootCmd.AddCommand(tempCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tempCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tempCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getTemp(zip string) {

	// Get the weather
	t, err := swissmeteo.GetTemp(zip)
	if err != nil {
		fmt.Printf("Something went wrong when fetching the temperature: %s\n", err)
		return
	}

	// Print the temperature to the console
	fmt.Printf("The temperature in %s is %0.1f °C\n", zip, t)
}
