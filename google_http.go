package main

import (
	"context"
	"errors"
	"fmt"
	getFileList "github.com/tanaikech/go-getfilelist"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/docs/v1"
	"google.golang.org/api/drive/v3"
	oauth22 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	// TODO: randomize it
	oauthStateString = "pseudo-random"
	googleOauthConfig *oauth2.Config
	client           *http.Client
	config           *oauth2.Config
)

func init() {
	fmt.Println("GOOGLE_CLIENT_ID:", os.Getenv("GOOGLE_CLIENT_ID"))
	fmt.Println("GOOGLE_CLIENT_SECRET:", os.Getenv("GOOGLE_CLIENT_SECRET"))

	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{
			docs.DriveScope,
			oauth22.UserinfoEmailScope,
			oauth22.UserinfoProfileScope,},
		Endpoint:     google.Endpoint,
	}
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	var htmlIndex = `<html>
<body>
	<a href="/login">Google Log In</a>
</body>
</html>`
	fmt.Fprintf(w, htmlIndex)
}


func GetOauthConfig(ctx context.Context, state, code string) error {
	if state != oauthStateString {
		return fmt.Errorf("invalid oauth state")
	}
	token, err := googleOauthConfig.Exchange(ctx, code)
	if err != nil {
		return fmt.Errorf("code exchange failed: %s", err.Error())
	}
	client = googleOauthConfig.Client(ctx, token)
	return nil
}

// login
func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	fmt.Println("AuthCodeURL: ", url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// callback
func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {


	state := r.FormValue("state")
	code := r.FormValue("code")
	content, err := getUserInfo(state, code, r.Context())
	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Fprintf(w, "Content: %s\n", content)
}


func getUserInfo(state string, code string,  ctx context.Context) ([]byte, error) {
	fmt.Printf("[getUserInfo] state=%s, code=%s\n", state, code)
	if state != oauthStateString {
		return nil, fmt.Errorf("invalid oauth state")
	}
	token, err := googleOauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}
	client = googleOauthConfig.Client(ctx, token)

	ListFiles()

	fmt.Printf("[getUserInfo] token=%s\n", token)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}
	fmt.Println("token.AccessToken:", token.AccessToken)
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	fmt.Printf("[getUserInfo] response.Body=%v\n", response.Body)
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}
	return contents, nil
}


func ListFiles() (string, error) {
	fmt.Println("listing files >>>> ")
	if client == nil {
		return "", fmt.Errorf("client expired")
	}
	api, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		fmt.Println("Unable to retrieve Drive Service:", err)
		return "", fmt.Errorf("unable to access drive: check your access")
	}
	r, err := api.Files.List().
		Q("mimeType = 'application/vnd.google-apps.folder' and name contains 'templates-gen'").
		PageSize(1).
		Fields("nextPageToken, files(id, name, properties)").Do()
	if err != nil {
		fmt.Printf("Unable to retrieve files: %v", err)
		return "", errors.New(fmt.Sprintf("Unable to retrieve files: %v", err))
	}

	//fileLength := len(r.Files)
	for _, i := range r.Files {
		fmt.Printf("%s (%s) %+v \n", i.Name, i.Id, i.Properties)
		res, err := getFileList.Folder(i.Id).Fields("files(id, name, properties)").Do(client)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for _, i := range res.FileList {
			for _, j := range i.Files {
				fmt.Printf("file name: %s, %s \n", j.Name, j.Id)
			}
		}
	}

	return "", errors.New(fmt.Sprintf("Multiple UDS Roots found."))

}


func main() {
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/login", handleGoogleLogin)
	http.HandleFunc("/callback", handleGoogleCallback)
	http.ListenAndServe(":8080", nil)
}
