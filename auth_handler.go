package main

import "fmt"
import "encoding/json"
import "golang.org/x/oauth2"
import "golang.org/x/oauth2/google"
import "net/http"
import "github.com/satori/go.uuid"
import "github.com/gorilla/sessions"

type GoogleProfile struct {
	Name string `json:"name"`
}

var store = sessions.NewCookieStore(uuid.NewV4().Bytes())

func IsAuthorized(response http.ResponseWriter, request *http.Request) bool {
	session, err := store.Get(request, "minimoozie")
	check(err)
	log.Info(session)

	if session.IsNew {
		return false
	} else {
		token := oauth2.Token{AccessToken: session.Values["token"].(string)}
		user := getUser(&token)
		if user.Name == session.Values["name"] {
			return true
		} else {
			log.Error(fmt.Sprintf("cookie username (%s) doesnt match received username (%s)", session.Values["name"], user.Name))
		}
	}

	return false
}

func Authorize(response http.ResponseWriter, request *http.Request) {
	if !IsAuthorized(response, request) {
		oauthconf := &oauth2.Config{
			ClientID:     Conf.OauthClientId,
			ClientSecret: Conf.OauthClientSecret,
			RedirectURL:  Conf.RedirectURL,
			Scopes:       []string{"profile"},
			Endpoint:     google.Endpoint,
		}
		url := oauthconf.AuthCodeURL("state")
		http.Redirect(response, request, url, http.StatusFound)
	}
}

func OauthCallbackHandler(response http.ResponseWriter, request *http.Request) {
	oauthconf := &oauth2.Config{
		ClientID:     Conf.OauthClientId,
		ClientSecret: Conf.OauthClientSecret,
		RedirectURL:  Conf.RedirectURL,
		Scopes:       []string{"profile"},
		Endpoint:     google.Endpoint,
	}
	code := request.FormValue("code")

	tok, err := oauthconf.Exchange(oauth2.NoContext, code)
	check(err)

	user := getUser(tok)

	session, err := store.Get(request, "minimoozie")
	check(err)
	session.Values["token"] = tok.AccessToken
	session.Values["name"] = user.Name
	err = session.Save(request, response)
	check(err)

	http.Redirect(response, request, "/", http.StatusFound)
}

func getUser(token *oauth2.Token) GoogleProfile {
	oauthconf := &oauth2.Config{
		ClientID:     Conf.OauthClientId,
		ClientSecret: Conf.OauthClientSecret,
		RedirectURL:  Conf.RedirectURL,
		Scopes:       []string{"profile"},
		Endpoint:     google.Endpoint,
	}
	client := oauthconf.Client(oauth2.NoContext, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v1/userinfo?alt=json")
	check(err)

	result := new(GoogleProfile)
	err = json.NewDecoder(resp.Body).Decode(result)
	check(err)

	return *result
}
