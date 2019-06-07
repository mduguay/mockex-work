# Having this as an arg allows us to override it at build time
ARG GO_VERSION=1.12.5

# Container for building the app
# Should I use ${GO_VERSION}-alpine ?
FROM golang:${GO_VERSION} as builder

# Principle of least privilege. Create non-root user
RUN useradd -u 1134 mockexuser

# This works with alpine image
# RUN addgroup -S mockexuser && adduser -S -G mockexuser mockexuser

# Set the working directory outside $GOPATH to support modules
WORKDIR /src

# Some deps need to be installed on the alpine image before go mod download will work
# RUN apk add git

# Fetch dependencies. They will be cached, speeding this up next time
COPY ./go.mod ./go.sum ./
RUN go mod download

# Import code
COPY ./ ./

# Not dependent on C libraries
# Mark build as statically linked.
# Build executable to /app. 
RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -o /app .

# Really small image for just the app, doesn't even need an OS!
FROM scratch as final

# Copy the compiled app from the builder
COPY --from=builder /app /app

# Copy the permissionless user
COPY --from=builder /etc/passwd /etc/passwd

# Declare the port on which the server will be exposed
EXPOSE 8080

# Perform actions as unprivileged user
USER mockexuser

# Run the compiled binary
ENTRYPOINT [ "/app" ]