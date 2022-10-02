/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	greetingBanner = `
 █████╗ ██████╗ ██████╗ ██╗     ███████╗    ██████╗  █████╗ ████████╗ █████╗ 
██╔══██╗██╔══██╗██╔══██╗██║     ██╔════╝    ██╔══██╗██╔══██╗╚══██╔══╝██╔══██╗
███████║██████╔╝██████╔╝██║     █████╗█████╗██║  ██║███████║   ██║   ███████║
██╔══██║██╔═══╝ ██╔═══╝ ██║     ██╔══╝╚════╝██║  ██║██╔══██║   ██║   ██╔══██║
██║  ██║██║     ██║     ███████╗███████╗    ██████╔╝██║  ██║   ██║   ██║  ██║
╚═╝  ╚═╝╚═╝     ╚═╝     ╚══════╝╚══════╝    ╚═════╝ ╚═╝  ╚═╝   ╚═╝   ╚═╝  ╚═╝
																			 
	`
)

var (
	flags struct {
		// configFilePath is the path to the config file that we will read.
		configFilePath string
		// inputFilesPath is the folder to the input file that we will read.
		inputFilesPath string
		// output file path is the path to the output file that we will write.
		ouputFilePath string
		// debug mode is a flag that will enable debug mode.
		debug bool
	}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "apd [command]",
	Short: "Apple Data CLI",
	Long:  `Apple Data CLI is a command line tool to process Apple part data.`,
	Run: func(cmd *cobra.Command, args []string) {
		printGreeting()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.appledata.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().StringVar(&flags.configFilePath, "config", "", "config file path")
	rootCmd.PersistentFlags().StringVar(&flags.inputFilesPath, "input", "", "input files folder")
	rootCmd.PersistentFlags().StringVar(&flags.ouputFilePath, "out", "", "output file path")
	rootCmd.PersistentFlags().BoolVar(&flags.debug, "debug", false, "debug mode")
}

func printGreeting() {
	fmt.Println(greetingBanner)
}
