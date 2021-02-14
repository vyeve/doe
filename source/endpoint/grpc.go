package endpoint

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strconv"

	"doe/source/logger"
	"doe/source/models"

	"github.com/gorilla/mux"
	gRPC "google.golang.org/grpc"
)

const defaultLimitForAllPorts = 1000

// ServiceInterface describes methods of clientAPI microservice
type ServiceInterface interface {
	UploadData(w http.ResponseWriter, r *http.Request)
	GetOne(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
}

type service struct {
	grpcClient models.PortServiceClient
	log        logger.Logger
}

// NewService initializes gRPC client
func NewService(params ServerParams) (ServiceInterface, error) {
	conn, err := gRPC.Dial(os.Getenv(CommunicationURLEnvKey), gRPC.WithInsecure())
	if err != nil {
		return nil, err
	}
	cl := models.NewPortServiceClient(conn)
	srv := service{
		grpcClient: cl,
		log:        params.Logger,
	}
	return srv, nil
}

// UploadData redirect upload file to gRPC server
func (s service) UploadData(w http.ResponseWriter, r *http.Request) {
	f, _, err := r.FormFile("file")
	// f, _, err := formFileFunc(r, "file")
	if err != nil {
		s.log.Warnf("Failed to extract file from request. err: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer func() {
		err = f.Close()
		if err != nil {
			s.log.Warn(err)
		}
	}()

	err = s.readFile(r.Context(), f)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			s.log.Warnf("Failed to write data. err: %v", err)
		}
		return
	}
	_, err = w.Write([]byte("file successfully upload"))
	if err != nil {
		s.log.Warnf("Failed to write data. err: %v", err)
	}
}

// GetOne receives portID, and redirect found port from gRPC server to response
func (s service) GetOne(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	id := mux.Vars(r)["port-id"]
	if id == "" {
		s.log.Warn("port-id is required")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	req := models.Request{
		PortID: id,
	}
	port, err := s.grpcClient.GetOne(ctx, &req)
	if err != nil {
		s.log.Warnf("Failed to get port by ID. err=%s", err)
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"message": "port not found"}`)
		return
	}
	p, err := json.Marshal(port)
	if err != nil {
		s.log.Warnf("Failed to Marshal Port[%+v]. err=%s", port, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "{%q: %s}\n", id, p) // nolint: errcheck
}

// GetAll redirects all found ports from gRPC server to response
func (s service) GetAll(w http.ResponseWriter, r *http.Request) {
	limit := parseLimit(r)
	ports, err := s.grpcClient.GetAll(r.Context(), &models.Request{Limit: limit})
	if err != nil {
		s.log.Warnf("Failed to get ports' info. err=%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	p, err := json.Marshal(ports)
	if err != nil {
		s.log.Warnf("Failed to marshal port: %v", err)

	}
	w.Write(p) // nolint: errcheck
}

func (s service) readFile(ctx context.Context, reader io.Reader) error {
	stream, err := s.grpcClient.Set(ctx)
	if err != nil {
		s.log.Warnf("failed to connect to domain server. err=%s", err)
		return err
	}

	err = s.parseJSON(reader, stream)
	if err != nil {
		return err
	}

	if err = stream.CloseSend(); err != nil {
		s.log.Warnf("failed to close stream. err=%s", err)
		return err
	}
	return nil
}

func (s service) parseJSON(reader io.Reader, stream models.PortService_SetClient) error {
	dec := json.NewDecoder(reader)
	// read opening bracket
	t, err := dec.Token()
	if err != nil || !reflect.DeepEqual(t, json.Delim('{')) {
		s.log.Warn("not valid data in file")
		return errors.New("not valid token")
	}
	var key json.Token
	for dec.More() {
		key, err = dec.Token()
		if err != nil {
			s.log.Warn("not valid data in file")
			return err
		}
		id, ok := key.(string)
		if !ok {
			s.log.Warn("not valid data in file")
			return err
		}
		var m models.Port
		if err = dec.Decode(&m); err != nil {
			s.log.Warn("not valid data in file")
			return err
		}
		m.PortID = id
		if err = stream.Send(&m); err != nil {
			s.log.Warnf("failed to send message to domain service. err=%s", err)
			return err
		}
	}
	// read closing bracket
	if t, err = dec.Token(); err != nil || !reflect.DeepEqual(t, json.Delim('}')) {
		s.log.Warnf("not valid data in file")
		return err
	}
	return nil
}

func parseLimit(r *http.Request) int32 {
	parameters := r.URL.Query()
	limits, ok := parameters["limit"]
	if !ok || len(limits) != 1 {
		return defaultLimitForAllPorts
	}
	limit, err := strconv.Atoi(limits[0])
	if err != nil || limit < 1 {
		return defaultLimitForAllPorts
	}
	return int32(limit)
}
