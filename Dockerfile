FROM golang:1.16
WORKDIR /go/src/app
COPY . .

RUN make build 

FROM alpine:latest  
RUN apk --no-cache add ca-certificates ipvsadm
WORKDIR /
COPY --from=0 /go/src/app/dist/kpng-backend-ipvs .
CMD ["/kpng-backend-ipvs"]  
