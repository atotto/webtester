FROM debian:latest

ENV DEBIAN_FRONTEND=noninteractive

RUN apt update -qq && apt install -y --no-install-recommends wget ca-certificates gpg-agent gnupg && \
	echo "deb http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google-chrome.list && \
	wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add - && \
    apt update && apt install -y --no-install-recommends \
    git openssh-client google-chrome-stable libnss3 libx11-6 && \
    rm -rf /var/lib/apt/lists/*

# setup golang
RUN curl -sL "https://storage.googleapis.com/golang/$(curl -s https://golang.org/VERSION?m=text).linux-amd64.tar.gz" | tar xz -C /usr/local && ln -sf /usr/local/go/bin/go /usr/local/bin/go && ln -sf /usr/local/go/bin/gofmt /usr/local/bin/gofmt
