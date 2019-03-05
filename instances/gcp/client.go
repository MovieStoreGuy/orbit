package gcp

import (
	"context"
	"fmt"

	"github.com/MovieStoreGuy/orbit/types"

	"google.golang.org/api/compute/v1"
	htransport "google.golang.org/api/transport/http"
)

type Client struct {
	service *compute.Service
}

func NewClient(ctx context.Context) (*Client, error) {
	hc, _, err := htransport.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	service, err := compute.New(hc)
	if err != nil {
		return nil, err
	}
	return &Client{service}, nil

}

func (c *Client) GetInstances(ctx context.Context, project string, filters ...string) ([]*types.Machine, error) {
	zones, errors := c.getAllZones(ctx, project)
	machines := make([]*types.Machine, 0)
	for {
		select {
		case zone, open := <-zones:
			if !open {
				return machines, nil
			}
			fmt.Println("Checking instances in zone: ", zone)
			call := c.service.Instances.List(project, zone)
			for _, filter := range filters {
				call = call.Filter(filter)
			}
			err := call.Pages(ctx, func(instances *compute.InstanceList) error {
				for _, instance := range instances.Items {
					address := make([]string, 0)
					for _, device := range instance.NetworkInterfaces {
						for _, ip := range device.AccessConfigs {
							address = append(address, ip.NatIP)
						}
					}
					machines = append(machines, &types.Machine{
						Name:       instance.Name,
						Labels:     instance.Labels,
						ExternalIP: address,
						Zone:       zone,
						Status:     instance.Status,
					})
				}
				return nil
			})
			if err != nil {
				return nil, err
			}
		case err, open := <-errors:
			if !open {
				return machines, nil
			}
			if err != nil {
				return nil, err
			}
		}
	}
	return machines, nil
}

// getAllZones returns processes every zone accessible by the authorised client
// Since this process can take some time, it returns a channel of zones and errors
// that are to be handled by the callee
func (c *Client) getAllZones(ctx context.Context, project string) (<-chan string, <-chan error) {
	zones := make(chan string)
	errch := make(chan error)
	go func() {
		defer func() {
			close(zones)
			close(errch)
		}()
		err := c.service.Zones.List(project).Pages(ctx, func(list *compute.ZoneList) error {
			for _, items := range list.Items {
				zones <- items.Name
			}
			return nil
		})
		if err != nil {
			errch <- err
		}
	}()
	return zones, errch
}
