FROM golang:1.10

ENV APP_DIR $GOPATH/src/github.com/gregbiv/sandbox

COPY . ${APP_DIR}
WORKDIR ${APP_DIR}

ARG GITHUB_TOKEN=${GITHUB_TOKEN}

# Install go dependecy manager
RUN curl https://glide.sh/get | sh

RUN make deps-dev
RUN make run
