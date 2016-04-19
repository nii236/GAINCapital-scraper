// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
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

	"github.com/nii236/forex/fetcher"
	"github.com/spf13/cobra"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		from, err := strconv.Atoi(cmd.Flag("from").Value.String())
		to, err := strconv.Atoi(cmd.Flag("to").Value.String())
		if err != nil {
			fmt.Println(err)
			return
		}

		fetcher.Do(from, to)
		fmt.Println("fetch called")
	},
}

func init() {
	RootCmd.AddCommand(fetchCmd)
	fetchCmd.Flags().IntP("from", "f", 0, "From which year")
	fetchCmd.Flags().IntP("to", "t", 0, "From which year")
	fetchCmd.Flags().StringP("url", "u", "http://www.google.com", "Which URL to fetch from")
}
