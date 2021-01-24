# DOE

## To start application:
```bash
docker-compose build && docker-compose up
```

## Usage

### Upload file
```bash
curl --location --request POST 'localhost:9099/doe/upload' \
--header 'Content-Type: multipart/form-data' \
--form 'file=@"/PATH/TO/FILE/ports.json"'
```

### Get port info by ID:
```bash
curl --location --request GET 'localhost:9099/doe/ports/PORT_ID'
```

### Get ports' info:
```bash
curl --location --request GET 'localhost:9099/doe/ports?limit=100'
```
If limit not set, default limit [1000] is used