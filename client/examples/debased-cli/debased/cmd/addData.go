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

// addDataCmd represents the addData command
var addDataCmd = &cobra.Command{
	Use:   "addData INTO <table_name> COLUMNS <column_name(s)>... VALUES <value(s)>... PAY <max_payment_allowed>",
	Short: "Adds data to a table",
	Args: func(cmd *cobra.Command, args []string) error {
		// Identify all the args that are <column_names> in an arg
		// Checks the args to ensure that the args are in the correct position and are correct
		fmt.Println("Args validation called.")
		minRequiredArguments := 7

		if len(args) < minRequiredArguments {
			return fmt.Errorf("Requires a minimum amount of %d arguments", minRequiredArguments)
		}

		tableNameArg := args[1]
		maxPayment := args[len(args)-1]

		if !hasValidKeywords(args) {
			return fmt.Errorf("Missing required argument(s): 'INTO', 'COLUMNS', or 'VALUES'")
		}

		columnNames := args[pos(args, tableNameArg)+2 : pos(args, "VALUES")]
		// values := args[pos(args, "VALUES")+1 : pos(args, maxPayment)]

		if args[0] != "INTO" {
			return fmt.Errorf("Requires the 'INTO' (string) argument as the first argument instead of: '%s'", args[0])
		}

		// Check to see if table_name exists in metadata

		if args[2] != "COLUMNS" {
			return fmt.Errorf("Requires the 'COLUMNS' (string) argument as the next required argument instead of: '%s'", args[2])
		}

		// Check to see if the columns in <columnNames> are in the metadata

		if args[3+len(columnNames)] != "VALUES" {
			return fmt.Errorf("Requires the 'VALUES' (string) argument as the next required argument instead of: '%s'", args[3+len(columnNames)])
		}

		if !areColumnsValuesSameSize(args) {
			return fmt.Errorf("Requires the same number of column and value names after 'COLUMNS' and 'VALUES'")
		}

		// Check if the values in the <values> slice (commented above) are located in the metadata

		if args[len(args)-2] != "PAY" {
			return fmt.Errorf("Requires the 'PAY' (string) argument as the next required argument instead of: '%s'", args[len(args)-2])
		}

		if _, err := strconv.Atoi(maxPayment); err != nil {
			return fmt.Errorf("Requires <max_payment_allowed> (int) argument instead of: '%s'", args[len(args)-1])
		}

		// Check if can give that amount of payment (account balance is not enough or something)

		return nil
	},
	// 	Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:

	PreRun: func(cmd *cobra.Command, args []string) {
		// addData INTO <table_name> COLUMNS <column_name(s)>... VALUES <value(s)>... PAY <max_payment_allowed>
		fmt.Println("Prerun addData called")
		// requiredArgs := &RequiredArgs{[]interface{}{&TableNameArg{""}, &ColumnArgs{[]string{}}, &ValueArgs{[]string{}}, &PayArg{0}}}
		addDataArgs := &AddDataRequiredArgs{
			tableName:   "",
			columnNames: []string{},
			valueNames:  []string{},
			payment:     0,
		}
		fmt.Println(args)
		i := 0

		if args[i] == "INTO" {
			i++
			addDataArgs.tableName = args[i]
			fmt.Println("addDataArgs.tableName: " + addDataArgs.tableName)
			i++
		}
		if args[i] == "COLUMNS" {
			i++
			for args[i] != "VALUES" {
				addDataArgs.columnNames = append(addDataArgs.columnNames, args[i])
				i++
			}
			fmt.Printf("addDataArgs.columnNames: %v \n", addDataArgs.columnNames)
		}
		if args[i] == "VALUES" {
			i++
			for args[i] != "PAY" {
				addDataArgs.valueNames = append(addDataArgs.valueNames, args[i])
				i++
			}
			fmt.Printf("addDataArgs.valueNames: %v \n", addDataArgs.valueNames)
		}
		if args[i] == "PAY" {
			i++
			var err error
			addDataArgs.payment, err = strconv.Atoi(args[i])
			if err != nil {
				fmt.Errorf("Payment argument <string> not properly casted to an int")
			}
			fmt.Printf("addDataArgs.payment: %d \n", addDataArgs.payment)
		}

		fmt.Println(addDataArgs)
	},

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("addData called")
		// call a function to send in addData Args
		// DebasedSystem.PendingTransactions
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
