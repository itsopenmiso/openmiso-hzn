package server

import (
	"context"
	"fmt"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/horizon/pkg/dbx"
	hznpb "github.com/hashicorp/horizon/pkg/pb"

	"github.com/hashicorp/waypoint-hzn/pkg/models"
	"github.com/hashicorp/waypoint-hzn/pkg/pb"
)

func (s *service) RegisterHostname(
	ctx context.Context,
	req *pb.RegisterHostnameRequest,
) (*pb.RegisterHostnameResponse, error) {
	L := hclog.FromContext(ctx)

	// Determine the full hostname
	var hostname string
	for {
		if req.Hostname == "" {
			hostname = petname.Generate(3, "-")
		} else {
			hostname = req.Hostname
		}

		hostname = hostname + "." + s.Domain

		var host models.Hostname
		//host.RegistrationId = reg.Id
		host.Hostname = hostname
		//host.Labels = req.Labels

		if err := dbx.Check(s.DB.Create(&host)); err != nil {
			// For now, assume the failure is because of failing the unique
			// constraint. If we autogenerated the name, retry, otherwise return
			// an error.
			if req.Hostname == "" {
				continue
			}

			L.Error("error creating hostname", "error", err)
			return nil, fmt.Errorf("requested hostname is not available")
		}

		break
	}

	L.Debug("adding label link", "hostname", hostname, "target", req.Labels)
	_, err := s.HznControl.AddLabelLink(ctx, &hznpb.AddLabelLinkRequest{
		Labels: hznpb.MakeLabels(":hostname", hostname),
		// TODO
		//Account: account,
		// Target: nil,
	})
	if err != nil {
		return nil, err
	}
	L.Info("added label link", "hostname", hostname, "target", req.Labels)

	return &pb.RegisterHostnameResponse{
		Fqdn: hostname,
	}, nil
}
