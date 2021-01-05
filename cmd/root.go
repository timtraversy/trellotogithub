package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "trellotogithub",
	Short: "Transfer cards from Trello boards to GitHub Projects",
	Run: func(cmd *cobra.Command, args []string) {
		x := viper.Get("trelloToken")
		fmt.Println(x)
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Prefix = "Fetching projects "
		s.Start()
		s.Stop()
		fmt.Println("Done")
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

type configuration struct {
	TrelloToken     string
	githubToken     string
	trelloBoardID   string
	githubProjectID string
}

// type Configuration struct {
// 	TrelloAuth AuthConfiguration    `yaml:"trello"`
// 	GithubAuth AuthConfiguration    `yaml:"github"`
// 	Mapping    MappingConfiguration `yaml:"mapping"`
// }

// type AuthConfiguration struct {
// 	Token    string `yaml:"token"`
// 	Username string `yaml:"username"`
// 	// Username yaml.Node `yaml:"username"`
// }

// type MappingConfiguration struct {
// 	TrelloBoard          string `yaml:"trello_board"`
// 	TrelloBoardName      string `yaml:"trello_board_name"`
// 	GithubRepository     string `yaml:"github_repository"`
// 	GithubRepositoryName string `yaml:"github_repository_name"`
// 	GithubProject        string `yaml:"github_project"`
// 	GithubProjectName    string `yaml:"github_project_name"`
// }

var config configuration

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.trellotogithub.yaml)")

	rootCmd.PersistentFlags().String("trelloToken", "123", "token received from Trello during authentication process")
	viper.BindPFlag("trelloToken", rootCmd.PersistentFlags().Lookup("trelloToken"))

	err := viper.Unmarshal(&config)
	fmt.Println(viper.Get("trelloToken"))
	fmt.Println(err, config)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".trellotogithub" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".trellotogithub")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
