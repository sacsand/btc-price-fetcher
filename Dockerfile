# Building the binary of the App
FROM golang:latest
# set work directory
WORKDIR /app
# copy go.mod and go.sum first
COPY go.mod go.sum ./

# Downloads all the dependencies in advance (could be left out, but it's more clear this way)
RUN go mod download
# then copy all files from root directory
COPY .env .
COPY . .

# build the go main 
RUN go build -o main . 
# run test cases
RUN go test 

# Exposes port 80
EXPOSE 80

CMD ["./main"] 