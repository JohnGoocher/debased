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
	Use:   "createTable <table name> COLUMNS {<column_name> <data_type>}... PAY <max payment allowed>",
	Short: "Creates a table on the debased network",
	Args: func(cmd *cobra.Command, args []string) error {
		// example usage: 'debased createTable pets cats 40'
		minRequiredArguments := 4
		maxPayment := args[len(args)-1]
		var columnNames []string
		var dataTypeValues []string

		if len(args) < minRequiredArguments {
			return fmt.Errorf("Requires a minimum amount of %d arguments", minRequiredArguments)
		}

		if _, err := strconv.Atoi(maxPayment); err != nil {
			return fmt.Errorf("Requires <max_payment_allowed> (int) as last argument instead of: '%s'", args[len(args)-1])
		}

		for i := 2; i < len(args)-2; i += 2 {
			columnNames = append(columnNames, args[i])
		}

		for i := 3; i < len(args)-2; i += 2 {
			if args[i] == "varchar" || args[i] == "int" {
				dataTypeValues = append(dataTypeValues, args[i])
			} else {
				return fmt.Errorf("Requires data types be of type 'varchar' or 'int' and not: '%s'", args[i])
			}
		}

		if len(columnNames) != len(dataTypeValues) {
			return fmt.Errorf("Requires same number of column names and data types in form {<column_name> <data_type>}")
		}

		if args[len(args)-2] != "PAY" {
			return fmt.Errorf("Requires the 'PAY' (string) argument as the next required argument instead of: '%s'", args[len(args)-2])
		}

		// Check if can give that amount of payment (account balance is not enough or something)

		return nil
	},
	// 	Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		// "createTable <table name> COLUMNS {<column_name> <data_type>}... PAY <max payment allowed>",
		// requiredArgs := &RequiredArgs{[]interface{}{&TableNameArg{""}, &ColumnArgs{[]string{}}, &ValueArgs{[]string{}}, &PayArg{0}}}
		createTableArgs := &CreateTableArgs{
			tableName:       "",
			columnNames:     []string{},
			columnDataTypes: []string{},
			payment:         0,
		}
		i := 0

		createTableArgs.tableName = args[i]
		i++
		if args[i] == "COLUMNS" {
			i++
			for args[i] != "PAY" {
				createTableArgs.columnNames = append(createTableArgs.columnNames, args[i])
				createTableArgs.columnDataTypes = append(createTableArgs.columnDataTypes, args[i+1])
				i += 2
			}
		}
		if args[i] == "PAY" {
			i++
			var err error
			createTableArgs.payment, err = strconv.Atoi(args[i])
			if err != nil {
				fmt.Errorf("Payment argument <string> not properly casted to an int")
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("createTable called.")
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
