package create

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	appName  string
	dir      string
	force    bool
	StartCmd = &cobra.Command{
		Use:     "create",
		Short:   "Create a new mod",
		Example: "jframe create -n users",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: re-implement with new template system
			fmt.Println("Module scaffold generator is being refactored. Please create modules manually.")
			os.Exit(1)
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&appName, "name", "n", "", "create a new mod with provided name")
	StartCmd.PersistentFlags().StringVarP(&dir, "path", "p", "internal/mod", "new file will generate under provided path")
	StartCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "Force generate the mod")
}
