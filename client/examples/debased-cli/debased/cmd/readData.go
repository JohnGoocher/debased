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
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// readDataCmd represents the readData command
var readDataCmd = &cobra.Command{
	Use:   "readData COLUMNS <column_name(s)>... FROM <table_name> [WHERE] [<condition>] {[AND|OR] [<condition>]}...",
	Short: "Reads data from a table",
	Args: func(cmd *cobra.Command, args []string) error {
		minRequiredArguments := 4

		if len(args) < minRequiredArguments {
			return fmt.Errorf("Requires a minimum amount of %d arguments", minRequiredArguments)
		}

		columnNames := args[pos(args, "COLUMNS")+1 : pos(args, "FROM")]
		// Check to see if the <table_name> is in the metadata
		// tableNameArg := args[pos(args, "FROM")+1]
		// conditionArg := args[pos(args, "WHERE")+1]

		if len(args) < minRequiredArguments {
			return fmt.Errorf("Requires a minimum amount of %d arguments", minRequiredArguments)
		}

		if args[0] != "COLUMNS" {
			return fmt.Errorf("Reqcuires the 'COLUMNS' argument as the first argument instead of: %s", args[0])
		}

		// Check to see if the columns in <columnNames> are in the metadata

		if args[1+len(columnNames)] != "FROM" {
			return fmt.Errorf("Requires the 'FROM' argument as the next argument instead of: %s", args[1+len(columnNames)])
		}

		// tableName := args[pos(args, "FROM")+1]

		// check to see if the tableName is in the metadata

		// Checks to see if the args include the WHERE statement

		if len(args) > 3+len(columnNames) {
			maxOptionalArgs := 2
			fmt.Println("The length of the args is into the optional length")

			if args[3+len(columnNames)] != "WHERE" {
				return fmt.Errorf("Requires the 'WHERE' argument as the next optional argument instead of: %s", args[3+len(columnNames)])
			}

			if len(args) < 3+len(columnNames)+maxOptionalArgs {
				return fmt.Errorf("Requires a condition argument as the next optional argument instead of: %s", args[3+len(columnNames)])
			}

			condition := args[pos(args, "WHERE")+1]

			fmt.Println(condition)

			// Checks to see if a valid condition:
			if !strings.Contains(condition, "=") || !strings.Contains(condition, ">") || !strings.Contains(condition, "<") {
				return errors.New("Requires an '=', '>', or '<' operator in the condition")
			}
		}

		return nil
	},
	// 	Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("readData called")
	},
}

func init() {
	rootCmd.AddCommand(readDataCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// readDataCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// readDataCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
