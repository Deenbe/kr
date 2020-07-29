/*
Copyright © 2020 kr contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"kr/lib"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var cfgFile string
var config *lib.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kr",
	Short: "Kinesis reset (kr) is a utility reset KCL consumer state to a known point in time",
	Long:  `kr reads the target stream to find a record created at the specified time. If a record is not created at that point it discovers the first one created after that point. Once the record is discovered, it updates the KCL state table in DynamoDB to the sequence number of that record.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := config.CalculatePointInTime(time.Now)
		if err != nil {
			return err
		}
		return lib.Reset(t, config)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	config = &lib.Config{}

	rootCmd.Flags().StringVar(&config.StreamName, "stream-name", "", "kinesis stream name")
	rootCmd.Flags().StringVar(&config.ConsumerName, "consumer-name", "", "kcl consumer name")
	rootCmd.Flags().StringVar(&config.Rewind, "rewind", "", "time window to rewind the stream")
	rootCmd.Flags().StringVar(&config.Since, "since", "", "date and time to rewind the stream")
	rootCmd.Flags().BoolVar(&config.Update, "update", false, "update sequence number in dynamodb")

	rootCmd.MarkFlagRequired("stream-name")
	rootCmd.MarkFlagRequired("consumer-name")
}
