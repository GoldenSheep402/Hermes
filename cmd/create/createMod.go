package create

import (
	"bytes"
	"fmt"
	"github.com/juanjiTech/jframe/mod/example"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io/fs"
	"log"
	"os"
	"path"
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
			err := load()
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			fmt.Println("Module " + appName + " generate success under " + dir)
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&appName, "name", "n", "", "create a new mod with provided name")
	StartCmd.PersistentFlags().StringVarP(&dir, "path", "p", "internal/mod", "new file will generate under provided path")
	StartCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "Force generate the mod")
}

func load() error {
	if appName == "" {
		return errors.New("mod name should not be empty, use -n")
	}

	dirEntries, err := example.FS.ReadDir(".")
	if err != nil {
		return err
	}

	err = cloneDir("", dirEntries)
	if err != nil {
		return err
	}

	return nil
}

// 克隆目录以及目录下文件
func cloneDir(subPath string, entry []fs.DirEntry) error {
	for _, dirEntry := range entry {
		if dirEntry.IsDir() {
			if _, err := os.Stat(path.Join(dir, appName, subPath, dirEntry.Name())); os.IsNotExist(err) {
				if err = os.MkdirAll(path.Join(dir, appName, subPath, dirEntry.Name()), os.ModePerm); err != nil {
					return err
				}
			} else if err != nil {
				return err
			}
			subDir, _ := example.FS.ReadDir(path.Join(subPath, dirEntry.Name()))
			if err := cloneDir(path.Join(subPath, dirEntry.Name()), subDir); err != nil {
				return err
			}
			continue
		}
		if dirEntry.Name() == "embed.go" {
			continue
		}
		file, err := example.FS.ReadFile(path.Join(subPath, dirEntry.Name()))
		if err != nil {
			return err
		}
		file = bytes.ReplaceAll(file, []byte("example"), []byte(appName))

		// check if file exist
		fp := path.Join(dir, appName, subPath, dirEntry.Name())
		_, err = os.Stat(fp)
		if err == nil && !force {
			return errors.New("file " + fp + " is existed, use -f to force generate")
		}

		f, err := os.Create(fp)
		if err != nil {
			log.Println(err)
		}
		_, err = f.WriteString(bytes.NewBuffer(file).String())
		if err != nil {
			log.Println(err)
		}
		_ = f.Close()
	}
	return nil
}
