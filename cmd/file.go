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
	"embed"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// Post is a struct that holds information to be added into an html file
type Post struct {
	Title    string
	Contents template.HTML
}

var tmpl embed.FS

// fileCmd represents the file command
var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "makesite file transforms an individual txt file into an html file",
	Long: `makesite file transforms an individual txt file into an html file

	The makesite file <filename>.txt command will create a <filename>.html file`,
	Run: func(cmd *cobra.Command, args []string) {
		for argNum := range args {
			arg := args[argNum]
			outputFile, err := save(arg)
			if err != nil {
				fmt.Printf("Error transforming %s to .html file!\n", arg)
			} else {
				fmt.Printf("Successfully created %s based on %s\n", outputFile, arg)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(fileCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func save(fileName string) (outputFileName string, err error) {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory")
		return "", err
	}

	postContents, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	newPost := new(Post)
	contentLines := strings.Split(string(postContents), "\n")
	if len(contentLines) > 0 {
		newPost.Title = contentLines[0]
	}
	for line := range contentLines {
		if line != 0 && contentLines[line] != "\n" {
			newPost.Contents += template.HTML("<p>" + contentLines[line] + "</p>\n")
		}
	}

	t, err := template.ParseFS(tmpl, "template.tmpl")
	if err != nil {
		return "", err
	}
	newFileName := path + "/" + fileName[0:len(fileName)-4] + ".html"
	newFile, err := os.Create(newFileName)
	if err != nil {
		return "", err
	}
	err = t.Execute(newFile, newPost)
	if err != nil {
		return "", err
	}
	return newFileName, nil
}
