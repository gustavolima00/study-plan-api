package healthcheck

import (
	mockhealthcheck "go-api/.internal/mocks/services/healthcheck"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// HandlerTestSuite ...
type HandlerTestSuite struct {
	suite.Suite

	MockService *mockhealthcheck.Service

	Handler Handler
}

// SetupTest ...
func (s *HandlerTestSuite) SetupTest() {
	t := s.T()
	s.MockService = mockhealthcheck.NewService(t)
	s.Handler = New(Params{
		HealthcheckService: s.MockService,
	})
}

// SetupSubTest ...
func (s *HandlerTestSuite) SetupSubTest() {
	s.SetupTest() // Clean up the mocks
}

// TestHandlerTestSuite ...
func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

// TestGetAPIStatus ...
func (s *HandlerTestSuite) TestGetAPIStatus() {
	tests := map[string]struct {
		MockSetup      func()
		ExpectedStatus int
		ExpectedBody   string
	}{
		"success": {
			MockSetup: func() {
				s.MockService.On("OnlineSince").Return(time.Duration(10*time.Second), nil)
			},
			ExpectedStatus: http.StatusOK,
			ExpectedBody:   "{\"online_time\":\"10s\"}\n",
		},
		"fail to get online time": {
			MockSetup: func() {
				s.MockService.On("OnlineSince").Return(time.Duration(0), assert.AnError)
			},
			ExpectedStatus: http.StatusInternalServerError,
		},
	}

	for name, tc := range tests {
		s.Run(name, func() {
			if tc.MockSetup != nil {
				tc.MockSetup()
			}

			resp, err := runHandler(s.Handler.GetAPIStatus)

			s.NoError(err)
			s.Equal(tc.ExpectedStatus, resp.Code)
			if tc.ExpectedBody != "" {
				s.Equal(tc.ExpectedBody, resp.Body.String())
			}
		})
	}
}

func runHandler(f func(e echo.Context) error) (*httptest.ResponseRecorder, error) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := f(c)
	return rec, err
}
