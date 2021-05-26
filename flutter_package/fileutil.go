package flutter_package

import (
	"io"
	"log"
	"os"
)

func MoveAndCopyContent(src string, dest string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return err
	}

	defer func(in *os.File) {
		err := in.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(in)

	out, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(out)

	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func posString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}

func ContainsString(slice []string, element string) bool {
	return !(posString(slice, element) == -1)
}
