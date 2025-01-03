FROM golang:1.23.4-alpine as build
WORKDIR /api
COPY . .
RUN go mod download && go mod verify
RUN GOOS=linux GOARCH=amd64 go build
RUN ls

FROM scratch
EXPOSE 5050
WORKDIR /proxyapi
COPY --from=build /api/congressProxy  .
CMD ["./congressProxy"]