# ./Dockerfile

FROM golang:1.20-alpine AS builder

# Move to working directory (/build).
WORKDIR /build

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

# Copy the code into the container.
COPY . .

# Set necessary environment variables needed for our image 
# and build the API server.
RUN go build -o apiserver .

FROM scratch

# Copy binary and config files from /build 
# to root folder of scratch container.
COPY --from=builder ["/build/apiserver", "/"]

# Export necessary port.
EXPOSE 3001

# Command to run when starting the container.
ENTRYPOINT ["/apiserver"]