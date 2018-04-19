package handlers_test

import (
	"net/http"
	"testing"

	"github.com/getfider/fider/app/handlers"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/mock"
	. "github.com/onsi/gomega"
)

func TestStatusHandler(t *testing.T) {
	RegisterTestingT(t)

	settings := &models.SystemSettings{
		BuildTime: "today",
	}

	server, _ := mock.NewServer()
	status, query := server.ExecuteAsJSON(handlers.Status(settings))

	Expect(query.String("build")).To(Equal("today"))
	Expect(status).To(Equal(http.StatusOK))
}

func TestCSPReportHandler(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()
	status, _ := server.ExecutePost(
		handlers.CSPReport(),
		"This is the message body.",
	)

	Expect(status).To(Equal(http.StatusOK))
}
