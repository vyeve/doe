package endpoint

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"doe/source/mocks"
	"doe/source/models"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
)

func TestService_readFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.TODO()
	logMock := mocks.NewMockLogger(ctrl)
	clientMock := mocks.NewMockPortServiceClient(ctrl)
	streamMock := mocks.NewMockPortService_SetClient(ctrl)
	srv := service{
		grpcClient: clientMock,
		log:        logMock,
	}

	testCases := []struct {
		name    string
		needErr bool
		msg     string
	}{
		{
			name:    "test with array",
			needErr: true,
			msg:     `[]`,
		},
		{
			name:    "test with bool key",
			needErr: true,
			msg:     `{true:"bar"}`,
		},
		{
			name:    "test with wrong closing bracket",
			needErr: true,
			msg:     `{"foo":"bar"]`,
		},
		{
			name:    "test with wrong closing bracket",
			needErr: true,
			msg:     `{"foo":{}]`,
		},
		{
			name:    "test with OK",
			needErr: false,
			msg:     `{"foo":{}}`,
		},
	}

	streamMock.EXPECT().Send(gomock.Any()).Return(nil).AnyTimes()
	streamMock.EXPECT().CloseSend().Return(nil).AnyTimes()
	logMock.EXPECT().Warn(gomock.Any()).AnyTimes()
	logMock.EXPECT().Warnf(gomock.Any(), gomock.Any()).AnyTimes()
	for _, tc := range testCases {
		clientMock.EXPECT().Set(gomock.Any()).Return(streamMock, nil)
		t.Run(tc.name, func(t *testing.T) {
			err := srv.readFile(ctx, strings.NewReader(tc.msg))
			if tc.needErr {
				if err == nil {
					t.Error("expected not <nil> error")
				}
			} else {
				if err != nil {
					t.Error(err)
				}
			}
		})
	}
}

func TestService_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logMock := mocks.NewMockLogger(ctrl)
	clientMock := mocks.NewMockPortServiceClient(ctrl)
	srv := service{
		grpcClient: clientMock,
		log:        logMock,
	}

	handler := http.HandlerFunc(srv.GetAll)
	testCases := []struct {
		name           string
		expectedStatus int
		getErr         error
		recvErr        error
		url            string
	}{
		{
			name:           "test with GetAll error",
			expectedStatus: http.StatusInternalServerError,
			getErr:         errors.New("test error"),
			url:            "/foobar?limit=10",
		},
		{
			name:           "test without error",
			expectedStatus: http.StatusOK,
			recvErr:        io.EOF,
			url:            "/foobar?hello=10",
		},
	}
	logMock.EXPECT().Warnf(gomock.Any(), gomock.Any()).AnyTimes()
	for _, tc := range testCases {
		recorder := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodGet, tc.url, nil)
		if err != nil {
			t.Fatal(err)
		}
		clientMock.EXPECT().GetAll(gomock.Any(), gomock.Any()).Return(&models.Ports{}, tc.getErr)
		t.Run(tc.name, func(t *testing.T) {
			handler.ServeHTTP(recorder, request)
			if code := recorder.Code; code != tc.expectedStatus {
				t.Errorf("expected code %d, but got %d", tc.expectedStatus, code)
			}
		})
	}
}

func TestService_GetOne(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logMock := mocks.NewMockLogger(ctrl)
	clientMock := mocks.NewMockPortServiceClient(ctrl)
	srv := service{
		grpcClient: clientMock,
		log:        logMock,
	}

	handler := http.HandlerFunc(srv.GetOne)
	testCases := []struct {
		name           string
		expectedStatus int
		portID         string
		getErr         error
	}{
		{
			name:           "test with empty port-id",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "test with get grpc error",
			expectedStatus: http.StatusNotFound,
			portID:         "PID",
			getErr:         errors.New("test error"),
		},
		{
			name:           "test without error",
			expectedStatus: http.StatusOK,
			portID:         "PID",
		},
	}
	logMock.EXPECT().Warn(gomock.Any()).AnyTimes()
	logMock.EXPECT().Warnf(gomock.Any(), gomock.Any()).AnyTimes()
	for _, tc := range testCases {
		recorder := httptest.NewRecorder()
		request := &http.Request{}
		request = mux.SetURLVars(request, map[string]string{"port-id": tc.portID})
		if tc.portID != "" {
			clientMock.EXPECT().GetOne(gomock.Any(), gomock.Any()).Return(new(models.Port), tc.getErr)
		}
		t.Run(tc.name, func(t *testing.T) {
			handler.ServeHTTP(recorder, request)
			if code := recorder.Code; code != tc.expectedStatus {
				t.Errorf("expected code %d, but got %d", tc.expectedStatus, code)
			}
		})
	}
}
