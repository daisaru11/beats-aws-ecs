ARG GO_VERSION
ARG BEAT_DOCKER_IMAGE

### Build Beats binary ###
FROM golang:${GO_VERSION} AS beats

LABEL maintainer "Daichi Sakai <daisaru11@gmail.com>"

ARG GO_VERSION
ARG GO_PLATFORM
ARG BEATS_VERSION
ARG BEAT_NAME

RUN git clone --depth=1 -b v${BEATS_VERSION} https://github.com/elastic/beats /go/src/github.com/elastic/beats
WORKDIR /go/src/github.com/elastic/beats/${BEAT_NAME}
RUN make ${BEAT_NAME}
RUN mkdir -p target \ 
 && mv ${BEAT_NAME} "target/${BEAT_NAME}-${BEATS_VERSION}-go${GO_VERSION}-${GO_PLATFORM}"

### Build beats-aws-ecs plugin binary ###
FROM golang:${GO_VERSION} AS beats_aws_ecs

LABEL maintainer "Daichi Sakai <daisaru11@gmail.com>"

ARG GO_VERSION
ARG GO_PLATFORM
ARG BEATS_VERSION
ARG BEATS_AWS_ECS_VERSION

RUN git clone --depth=1 -b v${BEATS_VERSION} https://github.com/elastic/beats /go/src/github.com/elastic/beats
RUN go get github.com/pkg/errors

COPY . /go/src/github.com/daisaru11/beats-aws-ecs
WORKDIR /go/src/github.com/daisaru11/beats-aws-ecs

RUN CGO_ENABLED=1 GOOS=linux go build -buildmode=plugin
RUN mkdir -p target \
 && mv beats-aws-ecs.so "target/beats-aws-ecs-${BEATS_AWS_ECS_VERSION}-${BEATS_VERSION}-go${GO_VERSION}-${GO_PLATFORM}.so"

### Build Beats Image ###
FROM ${BEAT_DOCKER_IMAGE}

LABEL maintainer "Daichi Sakai <daisaru11@gmail.com>"

ARG GO_VERSION
ARG GO_PLATFORM
ARG BEATS_VERSION
ARG BEAT_NAME
ARG BEATS_AWS_ECS_VERSION

COPY --from=beats /go/src/github.com/elastic/beats/${BEAT_NAME}/target/${BEAT_NAME}-${BEATS_VERSION}-go${GO_VERSION}-${GO_PLATFORM} /usr/share/${BEAT_NAME}/${BEAT_NAME}
COPY --from=beats_aws_ecs /go/src/github.com/daisaru11/beats-aws-ecs/target/beats-aws-ecs-${BEATS_AWS_ECS_VERSION}-${BEATS_VERSION}-go${GO_VERSION}-${GO_PLATFORM}.so /usr/share/${BEAT_NAME}/beats-aws-ecs.so
