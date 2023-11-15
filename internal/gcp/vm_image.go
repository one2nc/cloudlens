package gcp

import (
	"context"
	"strings"

	"github.com/one2nc/cloudlens/internal"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/compute/v1"
)

func ListImages(ctx context.Context) ([]ImageResp, error) {
	var imagesResp []ImageResp

	client, err := compute.NewService(ctx)
	if err != nil {
		log.Printf("Failed to create Compute Engine client: %v\n", err)
		return imagesResp, err
	}
	project := ctx.Value(internal.KeyActiveProject).(string)
	images, err := client.Images.List(project).Context(ctx).Do()
	if err != nil {
		log.Printf("Failed to fetch images: %v\n", err)
		return imagesResp, err
	}
	for _, image := range images.Items {
		createdAt, err := GetLocalTime(image.CreationTimestamp)
		if err != nil {
			log.Print(err)
			return imagesResp, err
		}
		res := ImageResp{
			Name:      image.Name,
			CreatedAt: createdAt,
			Location:  strings.Join(image.StorageLocations, ","),
			Status:    image.Status,
		}
		imagesResp = append(imagesResp, res)
	}
	return imagesResp, nil
}

func GetImage(ctx context.Context, ImageId string) (*compute.Image, error) {

	client, err := compute.NewService(ctx)
	if err != nil {
		log.Printf("Failed to create Compute Engine client: %v\n", err)
		return nil, err
	}
	project := ctx.Value(internal.KeyActiveProject).(string)
	image, err := client.Images.Get(project, ImageId).Context(ctx).Do()

	if err != nil {
		log.Printf("Failed to fetch snapshots: %v\n", err)
		return nil, err
	}

	return image, nil
}
