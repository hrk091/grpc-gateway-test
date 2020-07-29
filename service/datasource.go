package service

import (
	"context"
	"grpc-gateway-test/driver"
	pb "grpc-gateway-test/model"
	"github.com/mitchellh/mapstructure"
	"log"
	"github.com/golang/protobuf/ptypes"
	"time"
)

var (
	sourceStore = getDatasourceStore()
)

// getStore returns store driver for runtime environment
func getDatasourceStore() driver.DatastoreDriver {
	// Env varsに基づくDIはここで実装
	return &driver.FirestoreDriver{
		Doc1stKey: "",
		Col2ndKey: "",
	}
}

type DatasourceRpcService struct{}

func (s *DatasourceRpcService) GetList(ctx context.Context, req *pb.Empty) (*pb.DatasourceList, error) {
	log.Printf("GetList called")
	docs, err := sourceStore.GetAll(ctx, "datasources", "")
	if err != nil {
		return nil, err
	}
	sourceList := &pb.DatasourceList{}
	for _, doc := range docs {
		source := new(pb.Datasource)
		err = mapstructure.Decode(doc, source)
		if err != nil {
			return nil, err
		}
		// time.Time型で渡されるが、JSON mapのためにはprotobuf.Timestampへの変換が必要
		source.Timestamp, err = ptypes.TimestampProto(doc["timestamp"].(time.Time))
		if err != nil {
			continue
		}
		sourceList.Datasource = append(sourceList.Datasource, source)
	}
	return sourceList, nil
}
