FROM ubuntu:18.04

RUN apt-get update \
        && apt-get install -y python2.7 ffmpeg curl git \
        && ln -s /usr/bin/python2.7 /usr/bin/python

RUN curl -fsSL "https://dl.google.com/go/go1.12.linux-amd64.tar.gz" -o golang.tar.gz \
        && tar -C /usr/local -xzf golang.tar.gz \
        && rm golang.tar.gz

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

RUN mkdir -p /go/src/github.com/bancek/youtube-to-koofr
COPY . /go/src/github.com/bancek/youtube-to-koofr
RUN cd /go/src/github.com/bancek/youtube-to-koofr && dep ensure -vendor-only
RUN go get github.com/revel/cmd/revel && cd /go/src/github.com/revel/cmd && git checkout v0.20.0 && go get github.com/revel/cmd/revel
RUN cd /go && revel build github.com/bancek/youtube-to-koofr /youtube-to-koofr

RUN curl -L https://yt-dl.org/downloads/2019.03.09/youtube-dl -o /usr/local/bin/youtube-dl && chmod a+rx /usr/local/bin/youtube-dl

CMD /youtube-to-koofr/run.sh
