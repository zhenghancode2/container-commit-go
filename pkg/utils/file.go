package utils

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
)

func TarDir(dir, dstDir, pattern string) (string, error) {
	// 创建一个新的tar文件
	tarFile, err := os.CreateTemp(dstDir, pattern)
	if err != nil {
		return "", err
	}
	defer tarFile.Close()

	tarPath := tarFile.Name()
	tw := tar.NewWriter(tarFile)
	defer tw.Close()

	// 使用filepath.Walk遍历rootfs目录并将其打包为tar文件
	err = filepath.Walk(dir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		if relPath == "." {
			return nil
		}
		hdr, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		hdr.Name = relPath
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}
		if info.Mode().IsRegular() {
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()
			if _, err := io.Copy(tw, f); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return tarPath, err
	}
	return tarPath, nil
}
