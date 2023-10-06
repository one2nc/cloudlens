package gcp

import (
	"context"

	"github.com/rs/zerolog/log"
	"google.golang.org/api/iterator"

	resourcemanager "cloud.google.com/go/resourcemanager/apiv3"
	"cloud.google.com/go/resourcemanager/apiv3/resourcemanagerpb"
)


func  FetchProjects() {
	ctx := context.Background()
	c, err := resourcemanager.NewProjectsClient(ctx)
	if err != nil {
		log.Print(err)
	}

	defer c.Close()

	req := &resourcemanagerpb.ListProjectsRequest{
			Parent: "organizations/{org_id}",
	}
	it := c.ListProjects(ctx, req)
	projetNames := []string{}
	for {
		resp, err := it.Next()
		log.Print(err)
		if err == iterator.Done {
			break
		}
		if err != nil {
			// TODO: Handle error.
			break
		}
		// TODO: Use resp.

		projetNames = append(projetNames, resp.Name)

	}

	log.Print(projetNames)

}
