/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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

// Package cmd holds the commands for the makesite CLI
package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// Green allows Green terminal output
var Green = "\033[32m"

// Bold allows Bold terminal output
var Bold = "\033[1m"

// Reset allows a return to normal terminal output
var Reset = "\033[0m"

// dirCmd represents the dir command
var dirCmd = &cobra.Command{
	Use:   "dir",
	Short: "A brief description of your command",
	Long: `makesite dir transforms your txt files into html files within a directory

	The makesite dir <dirname> command will create these files for each .txt
		file in the directory`,
	Run: func(cmd *cobra.Command, args []string) {
		r_status, _ := cmd.Flags().GetBool("recursive")
		for argNum := range args {
			arg := args[argNum]
			fmt.Printf("Attempting to makesite from files in %s directory...\n", arg)
			directory, err := os.Stat(arg)
			if err != nil {
				fmt.Printf("Issue checking if %s is a directory\n", arg)
				continue
			}
			if !directory.Mode().IsDir() {
				fmt.Printf("%s is not a directory!\n", arg)
			} else {
				numSaved, err := saveDir(arg, r_status)
				if err != nil {
					fmt.Printf("Error generating pages: %s\n", err)
				} else {
					fmt.Printf("%s%sSuccess!%s Generated %s%d%s page(s) from %s directory.\n",
						Green, Bold, Reset, Bold, numSaved, Reset, arg)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(dirCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dirCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	rootCmd.Flags().BoolP("recursive", "r", false, "Creates html files from txt files found in current directory's subdirectories")
}

func saveDir(dirName string, recursive bool) (numSaved int, err error) {
	numSaved = 0
	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		return 0, err
	}
	for _, file := range files {
		if len(file.Name()) > 4 && file.Name()[len(file.Name())-4:] == ".txt" {
			save(dirName + file.Name())
			numSaved += 1
		}
		if file.IsDir() && recursive {
			err = filepath.Walk(dirName+file.Name(), func(path string, info os.FileInfo, err error) error {
				if err == nil && len(info.Name()) > 4 && info.Name()[len(info.Name())-4:] == ".txt" {
					save(path)
					numSaved += 1
				} else {
					return err
				}
				return nil
			})
			if err != nil {
				return numSaved, err
			}
		}
	}
	return numSaved, nil
}
