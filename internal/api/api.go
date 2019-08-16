package api

import (
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/mdevilliers/ports-ahoy/internal/store"
	"github.com/rs/zerolog"
)

type apiService struct {
	logger zerolog.Logger
	client remoteClient
}

type remoteClient interface {
	Get(key string) (*store.Port, error)
}

// New returns an endpoint for returning port information
func New(logger zerolog.Logger, client remoteClient) *apiService { // nolint
	return &apiService{
		logger: logger,
		client: client,
	}
}

// WEbservice configures the JSON endpoints
func (a *apiService) WebService() *restful.WebService {

	service := new(restful.WebService)
	service.
		Path("/api/v1/ports").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	service.Route(service.GET("/{port-id}").To(a.getPortByID))

	a.logger.Debug().Msg("creating `/api/v1/ports` endpoint")

	return service
}

func (a *apiService) getPortByID(request *restful.Request, response *restful.Response) {

	portID := request.PathParameter("port-id")
	port, err := a.client.Get(portID)

	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		return
	}

	// nolint: errcheck
	response.WriteHeaderAndEntity(http.StatusOK, port)
}
