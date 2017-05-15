FROM ubuntu:14.04

RUN apt-get update
RUN apt-get install -y software-properties-common
RUN echo "deb http://us.archive.ubuntu.com/ubuntu/ precise universe" >> /etc/apt/sources.list
RUN apt-get update
RUN apt-get install -y python2.7 libav-tools curl libavcodec-extra-53 git
RUN ln -s /usr/bin/python2.7 /usr/bin/python
RUN ln -s /usr/bin/avconv /usr/bin/ffmpeg

RUN curl -fsSL "https://golang.org/dl/go1.5.2.linux-amd64.tar.gz" -o golang.tar.gz \
        && tar -C /usr/local -xzf golang.tar.gz \
        && rm golang.tar.gz

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

RUN go get -u github.com/revel/cmd/revel
RUN go get github.com/bancek/youtube-to-koofr/app
RUN cd /go/src/github.com/bancek/youtube-to-koofr/app; go get ./...
RUN revel build github.com/bancek/youtube-to-koofr /youtube-to-koofr

RUN sudo curl -L https://yt-dl.org/downloads/2017.05.14/youtube-dl -o /usr/local/bin/youtube-dl && chmod a+rx /usr/local/bin/youtube-dl

CMD /youtube-to-koofr/run.sh
