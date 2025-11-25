package servicepb

import (
	"context"

	"github.com/smart-core-os/sc-bos/pkg/gen"
)

// RenameApi returns a gen.ServicesApiServer that changes the names associated with requests before calling client.
func RenameApi(client gen.ServicesApiClient, namer func(n string) string) gen.ServicesApiServer {
	return &rename{client: client, namer: namer}
}

type rename struct {
	gen.UnimplementedServicesApiServer
	client gen.ServicesApiClient
	namer  func(n string) string
}

func (r *rename) GetService(ctx context.Context, request *gen.GetServiceRequest) (*gen.Service, error) {
	request.Name = r.namer(request.Name)
	return r.client.GetService(ctx, request)
}

func (r *rename) PullService(request *gen.PullServiceRequest, server gen.ServicesApi_PullServiceServer) error {
	request.Name = r.namer(request.Name)
	stream, err := r.client.PullService(server.Context(), request)
	if err != nil {
		return err
	}
	for {
		msg, err := stream.Recv()
		if err != nil {
			return err
		}
		for _, change := range msg.Changes {
			change.Name = request.Name
		}
		err = server.Send(msg)
		if err != nil {
			return err
		}
	}
}

func (r *rename) CreateService(ctx context.Context, request *gen.CreateServiceRequest) (*gen.Service, error) {
	request.Name = r.namer(request.Name)
	return r.client.CreateService(ctx, request)
}

func (r *rename) DeleteService(ctx context.Context, request *gen.DeleteServiceRequest) (*gen.Service, error) {
	request.Name = r.namer(request.Name)
	return r.client.DeleteService(ctx, request)
}

func (r *rename) ListServices(ctx context.Context, request *gen.ListServicesRequest) (*gen.ListServicesResponse, error) {
	request.Name = r.namer(request.Name)
	return r.client.ListServices(ctx, request)
}

func (r *rename) PullServices(request *gen.PullServicesRequest, server gen.ServicesApi_PullServicesServer) error {
	request.Name = r.namer(request.Name)
	stream, err := r.client.PullServices(server.Context(), request)
	if err != nil {
		return err
	}
	for {
		msg, err := stream.Recv()
		if err != nil {
			return err
		}
		for _, change := range msg.Changes {
			change.Name = request.Name
		}
		err = server.Send(msg)
		if err != nil {
			return err
		}
	}
}

func (r *rename) StartService(ctx context.Context, request *gen.StartServiceRequest) (*gen.Service, error) {
	request.Name = r.namer(request.Name)
	return r.client.StartService(ctx, request)
}

func (r *rename) ConfigureService(ctx context.Context, request *gen.ConfigureServiceRequest) (*gen.Service, error) {
	request.Name = r.namer(request.Name)
	return r.client.ConfigureService(ctx, request)
}

func (r *rename) StopService(ctx context.Context, request *gen.StopServiceRequest) (*gen.Service, error) {
	request.Name = r.namer(request.Name)
	return r.client.StopService(ctx, request)
}

func (r *rename) GetServiceMetadata(ctx context.Context, request *gen.GetServiceMetadataRequest) (*gen.ServiceMetadata, error) {
	request.Name = r.namer(request.Name)
	return r.client.GetServiceMetadata(ctx, request)
}

func (r *rename) PullServiceMetadata(request *gen.PullServiceMetadataRequest, server gen.ServicesApi_PullServiceMetadataServer) error {
	request.Name = r.namer(request.Name)
	stream, err := r.client.PullServiceMetadata(server.Context(), request)
	if err != nil {
		return err
	}
	for {
		msg, err := stream.Recv()
		if err != nil {
			return err
		}
		for _, change := range msg.Changes {
			change.Name = request.Name
		}
		err = server.Send(msg)
		if err != nil {
			return err
		}
	}
}
