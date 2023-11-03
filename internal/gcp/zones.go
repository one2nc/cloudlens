package gcp

import (
	"context"

	"github.com/one2nc/cloudlens/internal"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/compute/v1"
)

func FecthZones(ctx context.Context) ([]string, error) {

	zonesResp := []string{}

	// Create a client with default credentials.
	client, err := compute.NewService(ctx)
	if err != nil {
		log.Printf("Failed to create compute service client: %v", err)
		return zonesResp, err

	}

	// List available zones in your project.
	projectID := ctx.Value(internal.KeyActiveProject).(string) // Replace with your GCP project ID.
	zones, err := client.Zones.List(projectID).Context(ctx).Do()
	if err != nil {
		log.Printf("Failed to list zones: %v", err)
		return zonesResp, err
	}

	for _, zone := range zones.Items {
		zonesResp = append(zonesResp, zone.Name)
	}

	return zonesResp, nil
}
