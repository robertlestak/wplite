FROM golang:1.22 as cli-builder

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /bin/wplite cmd/wplite/*.go

FROM ubuntu:22.04 as builder

RUN apt update && apt install -y \
    vim \
    git \
    zip \
    unzip \
    wget \
    curl \
    sudo \
    && rm -rf /var/lib/apt/lists/*

# install nodejs 20
RUN curl -sL https://deb.nodesource.com/setup_20.x | bash -

# install nodejs and npm
RUN apt-get install -y nodejs

# until merged into core, we need to use the forked core with sqlite support
RUN git clone https://github.com/aristath/wordpress-develop /usr/src/wordpress-develop && \
    cd /usr/src/wordpress-develop && \
    git checkout sqlite && \
    npm install && \
    npm run build

COPY scripts /wplite/scripts

FROM wordpress:latest as app

RUN apt update && apt install -y \
    sqlite3 \
    curl \
    unzip \
    jq \
    && rm -rf /var/lib/apt/lists/*

COPY --from=cli-builder /bin/wplite /bin/wplite
COPY --from=builder /usr/src/wordpress-develop/build/ /var/www/html/
RUN chown -R www-data:www-data /var/www/html
COPY --from=builder /wplite/scripts/wp-patch/wp-config.php /var/www/html/wp-config.php

COPY --from=builder /wplite /wplite

RUN bash /wplite/scripts/install-wp-cli.sh

ENTRYPOINT [ "bash", "/wplite/scripts/entrypoint.sh" ]