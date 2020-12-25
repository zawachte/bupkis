package util

import (
	"fmt"
	"strings"
	"time"

	"github.com/docker/go-units"
	"github.com/zawachte-msft/bupkis/pkg/registry"
)

func ImagesToNestedArray(images []registry.ImageData) [][]string {
	data := [][]string{}
	for _, image := range images {

		imageName := fmt.Sprintf("%s/%s", image.Hostname, image.Name)
		createdAt := time.Unix(image.Created.Unix(), 0)

		if createdAt.IsZero() {
			continue
		}

		createdAgo := fmt.Sprintf("%s ago", units.HumanDuration(time.Now().UTC().Sub(createdAt)))

		data = append(data, []string{imageName, image.Tag, createdAgo})
	}
	return data
}

func ParseImageName(imageName string) registry.ImageData {
	returnData := registry.ImageData{}
	returnData.Hostname = GetHostnameFromImage(imageName)
	returnData.Name = GetRepositoryFromImage(imageName)
	returnData.Tag = GetTagFromImage(imageName)
	return returnData
}

func GetHostnameFromImage(imageName string) string {
	delimbedImageName := strings.Split(imageName, "/")
	return delimbedImageName[0]
}

func GetTagFromImage(imageName string) string {
	if !strings.Contains(imageName, ":") {
		return ""
	}

	delimbedImageName := strings.Split(imageName, ":")
	return delimbedImageName[len(delimbedImageName)-1]
}

func GetRepositoryFromImage(imageName string) string {
	delimbedImageName := strings.Split(imageName, "/")

	stripTag := strings.Split(delimbedImageName[len(delimbedImageName)-1], ":")
	delimbedImageName[len(delimbedImageName)-1] = stripTag[0]

	return strings.Join(delimbedImageName[1:], "/")
}
