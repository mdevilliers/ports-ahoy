package importer

import (
	"strings"
	"testing"

	"github.com/mdevilliers/ports-ahoy/internal/importer/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func Test_Import(t *testing.T) {

	testCases := []struct {
		name      string
		json      string
		err       error
		importErr error
		storeErr  error
		count     int
	}{
		{
			name:  "ok",
			count: 2,
			json: `{
	     "AEAJM": {
	       "name": "Ajman",
	       "city": "Ajman",
	       "country": "United Arab Emirates",
	       "alias": [],
	       "regions": [],
	       "coordinates": [
	         55.5136433,
	         25.4052165
	       ],
	       "province": "Ajman",
	       "timezone": "Asia/Dubai",
	       "unlocs": [
	         "AEAJM"
	       ],
	       "code": "52000"
	     },
	     "AEAUH": {
	       "name": "Abu Dhabi",
	       "coordinates": [
	         54.37,
	         24.47
	       ],
	       "city": "Abu Dhabi",
	       "province": "Abu Z¸aby [Abu Dhabi]",
	       "country": "United Arab Emirates",
	       "alias": [],
	       "regions": [],
	       "timezone": "Asia/Dubai",
	       "unlocs": [
	         "AEAUH"
	       ],
	       "code": "52001"
			}
	     }`,
		},
		{
			name:  "no records",
			count: 0,
			json:  `{}`,
		},
		{
			name:      "not well json",
			json:      ``,
			importErr: ErrInvalidJSON,
		},
		{
			name:  "store error",
			count: 1,
			json: `{
	     "AEAJM": {
	       "name": "Ajman",
	       "city": "Ajman",
	       "country": "United Arab Emirates",
	       "alias": [],
	       "regions": [],
	       "coordinates": [
	         55.5136433,
	         25.4052165
	       ],
	       "province": "Ajman",
	       "timezone": "Asia/Dubai",
	       "unlocs": [
	         "AEAJM"
	       ],
	       "code": "52000"
	     }
		}`,
			storeErr:  errors.New("booyah"),
			importErr: errors.New("booyah"),
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {

			store := &mocks.FakeStorer{}

			store.SaveReturns(tc.storeErr)

			r := strings.NewReader(tc.json)
			i, err := New(r, store)

			if tc.err != nil {
				require.Equal(t, tc.err, err)
			}

			err = i.Import()

			if tc.importErr != nil {
				require.EqualError(t, tc.importErr, err.Error())
			}

			require.Equal(t, tc.count, store.SaveCallCount())

		})
	}
}
