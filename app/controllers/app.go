package controllers

import (
	"fmt"
	"net/http"

	"github.com/bancek/youtube-to-koofr/app/models"
	koofrclient "github.com/koofr/go-koofrclient"
	"github.com/revel/revel"
	"golang.org/x/oauth2"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	me := map[string]interface{}{}

	if koofr, ok := c.koofr(); ok {
		info, err := koofr.UserInfo()

		if err != nil {
			revel.ERROR.Println(err)
		} else {
			me["name"] = info.FirstName + " " + info.LastName
		}
	}

	authUrl := models.KoofrOAuthConfig.AuthCodeURL("")

	return c.Render(me, authUrl)
}

func (c App) Auth(code string) revel.Result {
	token, err := models.KoofrOAuthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		revel.ERROR.Println(err)
		return c.Redirect(App.Index)
	}

	user := c.user()
	user.KoofrAccessToken = token.AccessToken
	return c.Redirect(App.Index)
}

type ConvertResult struct {
	url   string
	koofr *koofrclient.KoofrClient
}

func (r *ConvertResult) Apply(req *revel.Request, resp *revel.Response) {
	resp.WriteHeader(http.StatusOK, "text/html; charset=utf-8")

	writer := resp.GetWriter()
	flusher := writer.(http.Flusher)

	for i := 0; i < 4096; i++ {
		writer.Write([]byte(" "))
	}
	flusher.Flush()

	writer.Write([]byte("<pre><code>\n"))

	logger := func(line string) {
		writer.Write([]byte(line))
		for i := 0; i < 1024; i++ {
			writer.Write([]byte(" "))
		}
		writer.Write([]byte("\n"))
		flusher.Flush()

		if revel.DevMode {
			fmt.Println(line)
		}

		flusher.Flush()
	}

	shortUrl, err := models.Convert(r.url, r.koofr, logger)

	if err != nil {
		revel.ERROR.Println(err)
		return
	}

	resp.Out.Write([]byte("</pre></code>\n"))

	resp.Out.Write([]byte("<br /><a href=\"" + shortUrl + "\">" + shortUrl + "</a>"))

	return
}

func (c App) Convert(url string) revel.Result {
	if koofr, ok := c.koofr(); ok {
		return &ConvertResult{
			url:   url,
			koofr: koofr,
		}
	} else {
		return c.Redirect(App.Index)
	}
}

func (c App) user() *models.User {
	return c.ViewArgs["user"].(*models.User)
}

func (c App) koofr() (*koofrclient.KoofrClient, bool) {
	user := c.user()

	if user.KoofrAccessToken == "" {
		return nil, false
	}

	koofr := models.GetKoofrClient(user.KoofrAccessToken)

	return koofr, true
}

func setuser(c *revel.Controller) revel.Result {
	var user *models.User
	if _, ok := c.Session["uid"]; ok {
		user = models.GetUser(c.Session["uid"])
	}
	if user == nil {
		user = models.NewUser()
		c.Session["uid"] = user.Id
	}
	c.ViewArgs["user"] = user
	return nil
}

func init() {
	revel.InterceptFunc(setuser, revel.BEFORE, &App{})
}
