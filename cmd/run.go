/*
Copyright © 2022 shfz

*/
package cmd

import (
	"log"
	"os"

	"github.com/shfz/shfz/run"
	"github.com/spf13/cobra"
)

var (
	file     string
	number   int
	parallel int
	timeout  int
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A scenario execute command",
	Long: `Execute a javascript scenario file written in javascript or compiled with Typescript
by specifying the total number of executions and the number of parallel executions.`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := os.Stat(file)
		if err != nil {
			log.Fatal("[-] No file file found :", file)
		}
		if err := run.Run(file, number, parallel, timeout); err != nil {
			log.Fatal("[-] Failed to run :", err)
		}
		log.Println("[+] Finish")
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	runCmd.Flags().StringVarP(&file, "file", "f", "", "scenario file (required)")
	if err := runCmd.MarkFlagRequired("file"); err != nil {
		panic(err)
	}
	runCmd.Flags().IntVarP(&parallel, "parallel", "p", 1, "number of parallel executions")
	runCmd.Flags().IntVarP(&number, "number", "n", 1, "total number of executions")
	runCmd.Flags().IntVarP(&timeout, "timeout", "t", 30, "scenario execution timeout(seconds)")
}
