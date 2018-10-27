package serviceAutorization

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig *oauth2.Config
)

func returnEmail(code string) (string, error) {
	content, err := getUserInfo(code)
	if err != nil {
		return "", err
	}

	var infoUser contentParse

	err = json.Unmarshal(content, &infoUser)
	if err != nil {
		return "", err
	}

	return infoUser.Email, nil
}

func getUserInfo(code string) ([]byte, error) {

	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:5000/dummy/oauth2callback",
		ClientID:     "93334427286-9kkusop9sjl32iml5qasuc58dhht25q7.apps.googleusercontent.com", //os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: "UdoSX27W3kMgERgTY9lqH6Qv",                                                //os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, err
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

type contentParse struct {
	Id             string `json:"id"`
	Email          string `json:"email"`
	Verified_email bool   `json:"verified_email"`
	Name           string `json:"name"`
	Given_name     string `json:"given_name"`
	Link           string `json:"link"`
	Picture        string `json:"picture"`
}
