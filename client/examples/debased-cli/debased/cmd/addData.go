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

// finds the position of the string in the slice of strings.
// Returns -1 if the string is not found.

func pos(s []string, value string) int {
	for i, v := range s {
		if value == v {
			return i
		}
	}
	return -1
}

func contains(s []string, value string) bool {
	for _, v := range s {
		if value == v {
			return true
		}
	}
	return false
}

// func isValidKeyword(value string) bool {
// 	if value == "INTO" ||  {
// 		return fmt.Errorf("Valid keyword ", value)
// 	}
// }

func hasValidKeywords(s []string) bool {
	if !contains(s, "INTO") {
		return false
	}
	if !contains(s, "COLUMNS") {
		return false
	}
	if !contains(s, "VALUES") {
		return false
	}
	return true
}

func areColumnsValuesSameSize(s []string) bool {
	columnSize := pos(s, "VALUES") - pos(s, "COLUMNS") + 1
	valueSize := len(s) - 1 - pos(s, "VALUES") + 1
	return columnSize == valueSize
}

// addDataCmd represents the addData command
var addDataCmd = &cobra.Command{
	Use:   "addData INTO <table_name> COLUMNS <column_name(s)>... VALUES <value(s)>... <max_payment_allowed>",
	Short: "Adds data to a table",
	Args: func(cmd *cobra.Command, args []string) error {
		// Identify all the args that are <column_names> in an arg
		minNArguments := 7

		if len(args) < minNArguments {
			return fmt.Errorf("Requires a minimum amount of %d arguments", minNArguments)
		}

		tableNameArg := args[1]
		maxPayment := args[len(args)-1]

		if !hasValidKeywords(args) {
			return fmt.Errorf("Missing required argument(s): 'INTO', 'COLUMNS', or 'VALUES'")
		}

		columnNames := args[pos(args, tableNameArg)+1 : pos(args, "VALUES")]
		// values := args[pos(args, "VALUES")+1 : pos(args, maxPayment)]

		if args[0] != "INTO" {
			return fmt.Errorf("Requires the 'INTO' (string) argument as the first argument instead of: '%s'", args[0])
		}

		// Check to see if table_name exists in metadata

		if args[2] != "COLUMNS" {
			return fmt.Errorf("Requires the 'COLUMNS' (string) argument as the next required argument instead of: '%s'", args[2])
		}

		// Check to see if the columns in <columnNames> are in the metadata

		if args[3+len(columnNames)] == "VALUES" {
			return fmt.Errorf("Requires the 'VALUES' (string) argument as the next required argument instead of: '%s'", args[3+len(columnNames)])
		}

		if !areColumnsValuesSameSize(args) {
			return fmt.Errorf("Requires the same number of column and value names after 'COLUMNS' and 'VALUES'")
		}

		// Check if the values in the <values> slice (commented above) are located in the metadata

		if _, err := strconv.Atoi(maxPayment); err != nil {
			return fmt.Errorf("Requires <max_payment_allowed> (int) argument instead of: '%s'", args[len(args)-1])
		}

		// Check if can give that amount of payment (account balance is not enough or something)

		// Check that for every column after COLUMNS there should be a value after VALUES

		return nil
	},
	// 	Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("addData called")
	},
}

func init() {
	rootCmd.AddCommand(addDataCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addDataCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addDataCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
