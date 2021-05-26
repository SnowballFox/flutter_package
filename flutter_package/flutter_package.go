package flutter_package

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

var (
	source        string
	dest          string
	versionNo     string
	versionNumber string
)

type Item struct {
	Version string `yaml:"version"`
}

func Start() {
	sep := string(os.PathSeparator)
	pwd, _ := os.Getwd()
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "s", Value: pwd + sep, Usage: "Source directory"},
			&cli.StringFlag{Name: "d", Value: pwd + sep + "dist" + sep, Usage: "Destination directory"},
		},
		Action: func(c *cli.Context) error {
			source = c.String("s")
			dest = c.String("d")

			if !AskForConfirmation("[Source directory]: " + source + "\n[Destination directory]: " + dest + "\nconfirm!(y/n)") {
				source = AskForInformation("Please input source directory: ")
				if !strings.HasSuffix(source, sep) {
					source = source + sep
				}

				dest = AskForInformation("Please input destination directory: ")
				if !strings.HasSuffix(dest, sep) {
					dest = dest + sep
				}
			}

			if !Exists(source + "pubspec.yaml") {
				fmt.Println("File 'pubspec.yaml' not found, maybe it is not a flutter project ?")
				return nil
			}

			vs, _ := readVersionFromConfig()
			versionNo = vs[0]
			versionNumber = vs[1]

			if !AskForConfirmation("[Version no]: " + versionNo + "\n[Version number]: " + versionNumber + "\nconfirm(y/n)") {
				versionNo = AskForInformation("Please input version no: ")
				versionNumber = AskForInformation("Please input version number: ")
			}

			ExecuteCommand(source, "flutter", "pub", "get")

			ExecuteCommand(source, "flutter", "build", "apk", "--target-platform", "android-arm64")
			if !Exists(dest) {
				err := os.Mkdir(dest, os.ModePerm)
				if err != nil {
					log.Fatal(err)
					return nil
				}
			}

			appPath := dest + "app_" + versionNo + "_" + versionNumber + ".apk"
			if Exists(appPath) {
				err := os.Remove(appPath)
				if err != nil {
					log.Fatal(err)
					return nil
				}
			}
			sourceAppPath := source + "build" + sep + "app" + sep + "outputs" + sep + "flutter-apk" + sep + "app-release.apk"

			time.Sleep(2 * time.Second)
			err := MoveAndCopyContent(sourceAppPath, appPath)

			if err != nil {
				log.Fatal(err)
				return nil
			}

			fmt.Println("Package apk success !!!")
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func readVersionFromConfig() ([]string, error) {
	var item Item
	yamlFile, ex := ioutil.ReadFile(source + "pubspec.yaml")
	if ex != nil {
		return nil, ex
	}
	ex = yaml.Unmarshal(yamlFile, &item)
	if ex != nil {
		return nil, ex
	}
	return strings.Split(item.Version, "+"), nil
}
