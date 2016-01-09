package models

import (
	"fmt"
	"github.com/koofr/go-httpclient"
	"github.com/koofr/go-koofrclient"
	"github.com/revel/revel"
	"golang.org/x/oauth2"
	"net/http"
	"os"
)

var KoofrBaseUrl = "https://app.koofr.net"
var KoofrOAuthConfig *oauth2.Config

func init() {
	revel.OnAppStart(InitKoofrOAuthConfig)
}

func InitKoofrOAuthConfig() {
	clientId := revel.Config.StringDefault("koofr.client_id", "")
	clientSecret := revel.Config.StringDefault("koofr.client_secret", "")
	redirectUrl := revel.Config.StringDefault("koofr.redirect_url", "")

	if clientId == "" || clientSecret == "" || redirectUrl == "" {
		panic("Missing koofr.client_id, koofr.client_secret, koofr.redirect_url")
	}

	KoofrOAuthConfig = &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Scopes:       []string{"public"},
		RedirectURL:  redirectUrl,
		Endpoint: oauth2.Endpoint{
			AuthURL:  KoofrBaseUrl + "/oauth2/auth",
			TokenURL: KoofrBaseUrl + "/oauth2/token",
		},
	}
}

func GetKoofrClient(accessToken string) *koofrclient.KoofrClient {
	client := koofrclient.NewKoofrClient(KoofrBaseUrl, false)
	client.HTTPClient.Headers.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	return client
}

func KoofrCreateShortUrl(koofr *koofrclient.KoofrClient, mountId string, path string) (shortUrl string, err error) {
	link := map[string]interface{}{}

	_, err = koofr.Request(&httpclient.RequestData{
		Method:         "POST",
		Path:           "/api/v2/mounts/" + mountId + "/links",
		ExpectedStatus: []int{http.StatusCreated},
		ReqEncoding:    httpclient.EncodingJSON,
		ReqValue: map[string]string{
			"path": path,
		},
		RespEncoding: httpclient.EncodingJSON,
		RespValue:    &link,
	})

	if err != nil {
		return "", err
	}

	shortUrl = link["shortUrl"].(string)

	return shortUrl, nil
}

func KoofrUpload(koofr *koofrclient.KoofrClient, filePath string, name string) (shortUrl string, err error) {
	mounts, err := koofr.Mounts()
	if err != nil {
		return "", err
	}

	primaryMountId := ""

	for _, mount := range mounts {
		if mount.IsPrimary {
			primaryMountId = mount.Id
			break
		}
	}

	if primaryMountId == "" {
		return "", fmt.Errorf("Primary mount id not found")
	}

	reader, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	defer reader.Close()

	newName, err := koofr.FilesPut(primaryMountId, "/", "YouTube to Koofr/"+name, reader)
	if err != nil {
		return "", err
	}

	remotePath := "/YouTube to Koofr/" + newName

	shortUrl, err = KoofrCreateShortUrl(koofr, primaryMountId, remotePath)
	if err != nil {
		return "", err
	}

	return shortUrl, nil
}
