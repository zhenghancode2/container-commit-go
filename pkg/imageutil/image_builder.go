package imageutil

import (
	"os"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/tarball"

	"container-commit-go/pkg/logger"
	"container-commit-go/pkg/utils"
)

func CreateLayerFromRootfs(rootfs string) (v1.Layer, error) {
	// 创建一个新的tar文件
	tarPath, err := utils.TarDir(rootfs, os.TempDir(), "layer-*.tar")
	if err != nil {
		return nil, err
	}
	logger.Debug("Created layer tar file", logger.WithString("tar path", tarPath))
	// 生成layer
	return tarball.LayerFromFile(tarPath)
}

func SaveImage(rootfs string) (v1.Image, error) {
	// 创建一个新的layer
	layer, err := CreateLayerFromRootfs(rootfs)
	if err != nil {
		return nil, err
	}
	img, err := mutate.AppendLayers(empty.Image, layer)
	if err != nil {
		return nil, err
	}
	logger.Debug("Image created")
	return img, nil
}

func SaveImageAsTar(rootfs, tarPath, imageRef string) error {
	// 创建一个新的layer
	layer, err := CreateLayerFromRootfs(rootfs)
	if err != nil {
		return err
	}
	img, err := mutate.AppendLayers(empty.Image, layer)
	if err != nil {
		return err
	}
	ref, err := name.ParseReference(imageRef)
	if err != nil {
		return err
	}

	logger.Info("Saving image", logger.WithString("image ref", imageRef))
	// 将镜像保存到本地tar文件
	if err := os.Remove(tarPath); err != nil {
		logger.Warn("Failed to remove existing tar file", logger.WithError(err))
	}
	return tarball.WriteToFile(tarPath, ref, img)
}
