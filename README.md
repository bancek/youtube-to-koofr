# YouTube to Koofr

Convert YouTube video to MP3 and upload it to Koofr.

## Getting started

```
go get -u github.com/revel/cmd/revel

go get github.com/bancek/youtube-to-koofr

export KOOFR_CLIENT_ID="CLIENTID"
export KOOFR_CLIENT_SECRET="CLIENTSECRET"
export KOOFR_REDIRECT_URL="http://localhost:9000/App/Auth"
export APP_SECRET="APPSECRET"

revel run github.com/bancek/youtube-to-koofr
```

Now go to http://localhost:9000/

## Deploy to Docker

```
cd docker

docker build -t youtube-to-koofr .

docker run -d -e KOOFR_CLIENT_ID="CLIENTID" -e KOOFR_CLIENT_SECRET="CLIENTSECRET" -e KOOFR_REDIRECT_URL="http://localhost:8000/App/Auth" -e APP_SECRET="APPSECRET" -p 8000:9000 youtube-to-koofr
```
