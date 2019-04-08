// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

// createTableCmd represents the createTable command
var createTableCmd = &cobra.Command{
	Use:   "createTable <table name> {<column_name> <data_type>}... <max payment allowed>",
	Short: "Creates a table on the debased network",
	Args: func(cmd *cobra.Command, args []string) error {
		// example usage: 'debased createTable pets cats 40'
		minNArguments := 4

		if len(args) < minNArguments {
			return fmt.Errorf("Requires a minimum amount of %d arguments", minNArguments)
		}

		// tableName := args[1]
		maxPayment := args[len(args)-1]

		if _, err := strconv.Atoi(maxPayment); err != nil {
			return fmt.Errorf("Requires <max_payment_allowed> (int) argument instead of: '%s'", args[len(args)-1])
		}

		// columnNames := args[1:len(args):2]

		return nil
	},
	// 	Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("createTable called")
	},
}

func init() {
	rootCmd.AddCommand(createTableCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createTableCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createTableCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
