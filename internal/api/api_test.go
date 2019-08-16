package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emicklei/go-restful"
	"github.com/mdevilliers/ports-ahoy/internal/api/mocks"
	"github.com/mdevilliers/ports-ahoy/internal/store"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func Test_GetPortByID(t *testing.T) {

	fake := &mocks.FakeRemoteClient{}
	log := zerolog.New(nil).With().Logger()

	api := New(log, fake)

	restful.Add(api.WebService())

	testCases := []struct {
		name     string
		request  string
		httpCode int
		storeErr error
	}{
		{
			name:     "ok",
			request:  "/api/v1/ports/foo",
			httpCode: http.StatusOK,
		},
		{
			name:     "store errors",
			request:  "/api/v1/ports/foo",
			httpCode: http.StatusNotFound,
			storeErr: errors.New("booyah"),
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {

			fake.GetReturns(&store.Port{}, tc.storeErr)

			httpWriter := httptest.NewRecorder()

			req, _ := http.NewRequest("GET", tc.request, nil)

			restful.DefaultContainer.ServeHTTP(httpWriter, req)
			require.Equal(t, tc.httpCode, httpWriter.Code)
		},
		)
	}

}
