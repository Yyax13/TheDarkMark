package misc

import (
	"os"
	"io"
	"path/filepath"
)

func CopyFile(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}

	defer srcFile.Close()
	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destFile.Close()
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	return destFile.Sync()

}

func CopyDir(src, dest string) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err

		}

		targetPath := filepath.Join(dest, relPath)
		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
            
		}

		return CopyFile(path, targetPath)

	})

}