FROM golang:1.14-buster as build

WORKDIR /go/src/app
ADD . /go/src/app

RUN go build -mod=vendor cmd/treactor/main.go

# Now copy it into our base image.
FROM gcr.io/distroless/base-debian10

WORKDIR /app
COPY --from=build /go/src/app/main /app/treactor
COPY --from=build /go/src/app/elements.yaml /app/elements.yaml

CMD ["/app/treactor"]
