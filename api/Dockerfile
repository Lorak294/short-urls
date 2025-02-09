FROM golang:alpine as builder

# directory setup and build

RUN mkdir /build

ADD . /build/

WORKDIR /build

RUN go build -o main .

# main stage

FROM alpine
# add new user
RUN adduser -S -D -H -h /app appuser

USER appuser

# copy build executable
COPY . /app

COPY --from=builder /build/main /app/

WORKDIR /app

# expose port
EXPOSE 3000

# run the executable
CMD [ "./main" ]