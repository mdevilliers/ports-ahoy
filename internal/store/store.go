package store

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mdevilliers/ports-ahoy/internal/env"

	cacheservice "github.com/mdevilliers/cache-service/proto/v1"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// NewCacheServiceClientFromEnvironment returns a connected GRPC client for the cacheservice or an error
// Client defaults use are for a developer environment
func NewCacheServiceClientFromEnvironment(ctx context.Context) (cacheservice.CacheClient, error) {

	host := env.LookUpWithDefaultStr("CACHE_SERVICE_GRPC_HOST", "0.0.0.0")
	port := env.LookUpWithDefaultStr("CACHE_SERVICE_GRPC_PORT", "3000")

	address := fmt.Sprintf("%s:%s", host, port)

	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrap(err, "error connecting to endpoint")
	}

	return cacheservice.NewCacheClient(conn), nil

}

type remoteCacheService struct {
	client cacheservice.CacheClient
}

// Port encapsulates details about a physical coastal port
type Port struct {
	Key         string    `json:"key"`
	Name        string    `json:"name"`
	Coordinates []float64 `json:"coordinates"`
	City        string    `json:"city"`
	Province    string    `json:"province"`
	Country     string    `json:"country"`
	Alias       []string  `json:"alias"`
	Regions     []string  `json:"regions"`
	Timezone    string    `json:"timezone"`
	Unlocs      []string  `json:"unlocs"`
	Code        string    `json:"code"`
}

// New returns an instance of a store for ports
func New(client cacheservice.CacheClient) *remoteCacheService {
	return &remoteCacheService{client: client}
}

// Save saves the store or returns an error
func (s *remoteCacheService) Save(p Port) error {

	bytes, err := json.Marshal(p)

	if err != nil {
		return errors.Wrap(err, "error marshalling port")
	}

	resp, err := s.client.Set(context.Background(), &cacheservice.SetRequest{
		Key:      p.Name,
		Contents: string(bytes),
	})

	if err != nil {
		return errors.Wrap(err, "GRPC error when saving a Port")
	}

	if !resp.GetStatus().GetOk() {
		return errors.New(resp.GetStatus().GetError().GetMessage())
	}

	return nil
}

// Get returns a Port or returns an error
func (s *remoteCacheService) Get(key string) (*Port, error) {

	if key == "" {
		return nil, errors.New("key not specified")
	}

	resp, err := s.client.GetByKey(context.Background(), &cacheservice.GetByKeyRequest{
		Key: key,
	})

	if err != nil {
		return nil, errors.Wrap(err, "GRPC error when getting a Port")
	}

	if !resp.GetStatus().GetOk() {
		return nil, errors.New(resp.GetStatus().GetError().GetMessage())
	}

	p := Port{Key: key}

	err = json.Unmarshal([]byte(resp.GetContents()), &p)

	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling Port")
	}

	return &p, nil
}
