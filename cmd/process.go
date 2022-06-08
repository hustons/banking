/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
  "log"
  "time"
  "github.com/spf13/cobra"
  "github.com/shuston/banking/transactions"
)

var revolutFile string
var aibFile string
var start string

// processCmd represents the process command
var processCmd = &cobra.Command{
  Use:   "process",
  Short: "Process transaction data into formatted banking data",
  Long: `Process transaction data into formatted banking data.

Transaction data can be supplied as exported data from financial sources.
Supported transaction data exports: Revolut, Allied Irish Bank

Multiple classification types can be specified in a classifications.json file, mapping transaction descriptions to defined classifications.
Formatted data output includes these supplied classifications by default.
Any transactions which cannot be automatically classified based on the classifications.json file will be classified as "UNKNOWN"
Foratted data can be used elsewhere to understand spending types and remaining budget per classification.`,
  Run: func(cmd *cobra.Command, args []string) {
    startDate := getStartDate(start)
    transactions.Process(revolutFile, aibFile, startDate)
  },
}

func init() {
  rootCmd.AddCommand(processCmd)

  processCmd.Flags().StringVarP(&revolutFile, "revolutFile", "r", "", "The path to the exported Revolut statement file")
  processCmd.Flags().StringVarP(&aibFile, "aibFile", "a", "", "The path to the exported AIB statement file")
  processCmd.Flags().StringVarP(&start, "startDate", "s", "", "The earliest date to be processed (format yyyymmdd)")
}

func getStartDate(start string) time.Time {
  if start != "" {
    startDate, err := time.Parse("20060102", start)
    if err != nil {
      log.Fatal("Could not parse start date: ", err)
    }
    return startDate.Add(time.Hour * -24)
  }

  return time.Unix(0, 0)
}
