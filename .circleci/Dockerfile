FROM debian:latest

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update -qq && apt-get install -y --no-install-recommends apt-utils wget curl ca-certificates gpg-agent gnupg && \
	echo "deb http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google-chrome.list && \
	wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add - && \
    apt-get update && apt-get install -y --no-install-recommends \
    git openssh-client google-chrome-stable libnss3 libx11-6 jq && \
    rm -rf /var/lib/apt/lists/*

# setup golang
RUN curl -sL "https://dl.google.com/go/$(curl -s https://go.dev/dl/?mode=json | jq -r '.[0].version').linux-amd64.tar.gz" | tar xz -C /usr/local && ln -sf /usr/local/go/bin/go /usr/local/bin/go && ln -sf /usr/local/go/bin/gofmt /usr/local/bin/gofmt

