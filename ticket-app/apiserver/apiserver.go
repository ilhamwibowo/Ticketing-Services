// ticket-app/apiserver/apiserver.go
package apiserver

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"ticket/storage"
)

var defaultStopTimeout = time.Second * 30

type APIServer struct {
	addr    string
	storage *storage.Storage
}

func NewAPIServer(addr string, storage *storage.Storage) (*APIServer, error) {
	if addr == "" {
		return nil, errors.New("addr cannot be blank")
	}

	return &APIServer{
		addr:    addr,
		storage: storage,
	}, nil
}

func (s *APIServer) Start(stop <-chan struct{}) error {
	router := s.setupRouter()

	server := &http.Server{
		Addr:    s.addr,
		Handler: router,
	}

	go func() {
		logrus.WithField("addr", server.Addr).Info("starting server")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("listen: %s\n", err)
		}
	}()

	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), defaultStopTimeout)
	defer cancel()

	logrus.WithField("timeout", defaultStopTimeout).Info("stopping server")
	return server.Shutdown(ctx)
}

func (s *APIServer) setupRouter() http.Handler {
	router := gin.Default()

	router.GET("/", s.defaultRoute)
	router.POST("/seats", s.createSeat)
	router.GET("/seats", s.listSeat)
	router.GET("/seats/status/:event_id/:seat_number", s.listSeat)

  router.GET("/events", s.getAllEvents)
  router.GET("/events/:event_id/empty-seats", s.getEmptySeats)
  
  router.POST("/book/:event_id/:seat_number", s.holdSeat)
  router.POST("/webhook/payment", s.paymentWebhook)
  
  return router
}