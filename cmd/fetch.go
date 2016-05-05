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
	"github.com/nii236/forex/fetcher"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch will download ticker data from GAIN Capital's ratedata website",
	Long: `Fetch will download ticker data from GAIN Capital's ratedata website.

Select which years you would like to fetch from, and a string slice of pairs you wish to retrieve.	
`,
	Run: func(cmd *cobra.Command, args []string) {

		from := viper.GetInt("from")
		to := viper.GetInt("to")
		pairs := viper.GetStringSlice("pairs")

		fetcher.Entry(from, to, pairs)
	},
}

func init() {
	RootCmd.AddCommand(fetchCmd)
}
