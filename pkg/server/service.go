package server

import (
	"github.com/hashicorp/go-hclog"
	hznpb "github.com/hashicorp/horizon/pkg/pb"
	"github.com/jinzhu/gorm"

	"github.com/itsopenmiso/openmiso-hzn/pkg/pb"
)

// service implements pb.WaypointHznServer.
type service struct {
	DB         *gorm.DB
	Domain     string
	Namespace  string
	HznControl hznpb.ControlManagementClient

	Logger hclog.Logger

	// Token public key is derived from the HznControl client on startup
	tokenPub []byte
}

var _ pb.WaypointHznServer = (*service)(nil)
