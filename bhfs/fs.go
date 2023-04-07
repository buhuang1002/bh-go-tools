package bhfs

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

func FileIsExisted(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func GetFilesFromDir(dir string) ([]string, int64, error) {
	var fileList []string
	var size int64
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		size += info.Size()
		fileList = append(fileList, path)
		return nil
	})
	if err != nil {
		return nil, 0, err
	}
	return fileList, size, nil
}

func MV(src, dst string) error {
	var cmd *exec.Cmd
	cmd = exec.Command("mv", src, dst)
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}

func IsDir(s string) bool {
	stat, err := os.Stat(s)
	if err != nil {
		return false
	}
	return stat.IsDir()
}
