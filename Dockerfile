FROM golang:1.15.7

ENV TARGET_PATH /usr/local/bin
ENV TARGET_NAME_CLIENT client-ms
ENV TARGET_PATH_CLIENT ${TARGET_PATH}/${TARGET_NAME_CLIENT}
ENV TARGET_NAME_SERVER endpoint-ms
ENV TARGET_PATH_SERVER ${TARGET_PATH}/${TARGET_NAME_SERVER}
ENV SOURCE_PATH /app/

ENV PACKAGES   ./source/data \
               ./source/endpoint 

COPY ./source ${SOURCE_PATH}/source
COPY go.mod ${SOURCE_PATH}
COPY go.sum ${SOURCE_PATH}


WORKDIR $SOURCE_PATH
RUN go mod download
RUN for package in $PACKAGES; do go test -cover -covermode=count $PACKAGES; done

WORKDIR $SOURCE_PATH/source/cmd/client
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-w' -o ${TARGET_PATH_CLIENT}

WORKDIR $SOURCE_PATH/source/cmd/endpoint
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-w' -o ${TARGET_PATH_SERVER}