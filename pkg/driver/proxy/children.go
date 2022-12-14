package proxy

import (
	"context"

	"github.com/smart-core-os/sc-api/go/traits"
)

// childrenFetcher implements pull.Fetcher to pull or poll children from a client
type childrenFetcher struct {
	client traits.ParentApiClient
	name   string
}

func (c childrenFetcher) Pull(ctx context.Context, changes chan<- *traits.PullChildrenResponse_Change) error {
	stream, err := c.client.PullChildren(ctx, &traits.PullChildrenRequest{Name: c.name})
	if err != nil {
		return err
	}

	for {
		recv, err := stream.Recv()
		if err != nil {
			return err
		}
		for _, change := range recv.Changes {
			change := change
			select {
			case <-ctx.Done():
				return ctx.Err()
			case changes <- change:
			}
		}
	}
}

func (c childrenFetcher) Poll(ctx context.Context, changes chan<- *traits.PullChildrenResponse_Change) error {
	req := &traits.ListChildrenRequest{
		Name: c.name,
	}
	for {
		children, err := c.client.ListChildren(ctx, req)
		if err != nil {
			return err
		}

		for _, child := range children.Children {
			child := child
			change := &traits.PullChildrenResponse_Change{
				Name:     c.name,
				NewValue: child,
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case changes <- change:
			}
		}

		if children.NextPageToken == "" {
			return nil
		}
		req.PageToken = children.NextPageToken
	}
}
