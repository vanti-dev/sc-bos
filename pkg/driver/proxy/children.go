package proxy

import (
	"context"

	"google.golang.org/protobuf/proto"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-api/go/types"
	"github.com/smart-core-os/sc-bos/pkg/util/chans"
)

// childrenFetcher implements pull.Fetcher to pull or poll children from a client
type childrenFetcher struct {
	client traits.ParentApiClient
	name   string
	known  map[string]*traits.Child // in case of polling, this tracks seen children so we correctly send changes
}

func (c *childrenFetcher) Pull(ctx context.Context, changes chan<- *traits.PullChildrenResponse_Change) error {
	stream, err := c.client.PullChildren(ctx, &traits.PullChildrenRequest{Name: c.name})
	if err != nil {
		return err
	}
	for {
		msg, err := stream.Recv()
		if err != nil {
			return err
		}
		for _, change := range msg.Changes {
			if err := chans.SendContext(ctx, changes, change); err != nil {
				return err
			}
		}
	}
}

func (c *childrenFetcher) Poll(ctx context.Context, changes chan<- *traits.PullChildrenResponse_Change) error {
	if c.known == nil {
		c.known = make(map[string]*traits.Child)
	}
	unseen := make(map[string]struct{}, len(c.known))
	for s := range c.known {
		unseen[s] = struct{}{}
	}

	req := &traits.ListChildrenRequest{Name: c.name, PageSize: 1000}
	for {
		res, err := c.client.ListChildren(ctx, req)
		if err != nil {
			return err
		}

		for _, node := range res.Children {
			// we do extra work here to try and send out more accurate changes to make callers lives easier
			change := &traits.PullChildrenResponse_Change{
				Type:     types.ChangeType_ADD,
				NewValue: node,
			}
			if old, ok := c.known[node.Name]; ok {
				change.Type = types.ChangeType_UPDATE
				change.OldValue = old
				delete(unseen, node.Name)
			}
			if proto.Equal(change.OldValue, change.NewValue) {
				continue
			}

			c.known[node.Name] = node
			if err := chans.SendContext(ctx, changes, change); err != nil {
				return err
			}
		}

		req.PageToken = res.NextPageToken
		if req.PageToken == "" {
			break
		}
	}

	for name := range unseen {
		node := c.known[name]
		delete(c.known, name)
		change := &traits.PullChildrenResponse_Change{
			Type:     types.ChangeType_REMOVE,
			OldValue: node,
		}
		if err := chans.SendContext(ctx, changes, change); err != nil {
			return err
		}
	}
	return nil
}
