# Build the application from source
FROM golang:1.22 AS build

WORKDIR /app

COPY . ./
RUN go mod download

RUN CGO_ENABLED=0 go build -o /bin/app

# Deploy the application binary into a lean image
FROM gcr.io/distroless/static-debian11

COPY --from=build /bin/app /bin
COPY .env.prod /bin

EXPOSE 3000

USER nonroot:nonroot

ENTRYPOINT [ "/bin/app", "/bin/.env.prod" ]