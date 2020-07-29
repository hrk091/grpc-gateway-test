package env

import (
	"grpc-gateway-test/util"
	"log"
	"os"
)

const (
	defaultHttpPort  = "8080"
	defaultGrpcPort  = "15000"
)

var (
	ProjectID  string
	BucketName string
	HttpPort   string
	GrpcPort   string
	NoAuth     bool
)

func init() {
	ProjectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	if ProjectID == "" {
		log.Fatalf("ENV 'GOOGLE_CLOUD_PROJECT' must be set.")
	}
	BucketName = ProjectID + ".appspot.com"
	HttpPort = util.Getenv("PORT", defaultHttpPort)
	GrpcPort = util.Getenv("MY_GRPC_PORT", defaultGrpcPort)
	NoAuth = os.Getenv("NO_AUTH") != ""
}
