package cmd

import (
	"fmt"
	"github.com/morikuni/aec"
	"github.com/spf13/cobra"
	"os"
)

var (
	Version   string
	GitCommit string
)

const WelcomeMessage = "Welcome to !"

func init() {
	snCmd.AddCommand(versionCmd)
}

// snCmd represents the base command when called without any sub commands.
var snCmd = &cobra.Command{
	Use:   "iniCmd",
	Short: "Expose your local endpoints to the Internet.",
	Long: `
snCmd combines a reverse proxy and websocket tunnels to expose your internal 
and development endpoints on another network, or to the public Internet via 
an exit-server.`,
	Run: runInlets,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the clients version information.",
	Run:   parseBaseCommand,
}

func getVersion() string {
	if len(Version) != 0 {
		return Version
	}
	return "dev"
}

func parseBaseCommand(_ *cobra.Command, _ []string) {
	printLogo()

	fmt.Println("Version:", getVersion())
	fmt.Println("Git Commit:", GitCommit)
	os.Exit(0)
}

// Execute adds all child commands to the root command(InCmd) and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the InCmd.
func Execute(version, gitCommit string) error {

	// Get Version and GitCommit values from main.go.
	Version = version
	GitCommit = gitCommit

	if err := snCmd.Execute(); err != nil {
		return err
	}
	return nil
}

func runInlets(cmd *cobra.Command, args []string) {
	printLogo()
	cmd.Help()
}

func printLogo() {
	inletsLogo := aec.WhiteF.Apply(inFigletStr)
	fmt.Println(inletsLogo)
}

const inFigletStr = ` _       _      _            _
(_)_ __ | | ___| |_ ___   __| | _____   __
| | '_ \| |/ _ \ __/ __| / _` + "`" + ` |/ _ \ \ / /
| | | | | |  __/ |_\__ \| (_| |  __/\ V /
|_|_| |_|_|\___|\__|___(_)__,_|\___| \_/
`

