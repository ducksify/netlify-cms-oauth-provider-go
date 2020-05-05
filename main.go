package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/pat"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/bitbucket"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gitlab"

	"github.com/joho/godotenv"
)

var (
	port         = "3000"
	githubHost   = "https://github.com"
	callbackHost = "https://localhost:3000"
)

const (
	script = `<!DOCTYPE html><html><head><script>
	if (!window.opener) {
  	window.opener = {
  	  postMessage: function(action, origin) {
  	    console.log(action, origin);
  	  }
  	}
	}	
	(function(status, provider, result) {
	  function receiveMessage(e) {
			console.log("Receive message:", e);
			
			msg = "authorization:" + provider + ":" + status + ":" + JSON.stringify(result)
			console.log("Sending message:", msg);
	    // send message to main window with da app
	    window.opener.postMessage(msg, e.origin);
	  }
	  window.addEventListener("message", receiveMessage, false);
	  // Start handshake with parent
	  console.log("Sending message:", provider);
	  window.opener.postMessage("authorizing:" + provider, "*");
	})("%s", "%s", %s)
	</script></head><body></body></html>`
)

// GET /
func handleMain(res http.ResponseWriter, req *http.Request) {
	log.Printf("handling root route '%s'\n", req)
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(``))
}

// GET /auth Page  redirecting after provider get param
func handleAuth(res http.ResponseWriter, req *http.Request) {
	url := fmt.Sprintf("auth/%s", req.FormValue("provider"))
	log.Printf("redirect to %s\n", url)
	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}

// GET /auth/provider  Initial page redirecting by provider
func handleAuthProvider(res http.ResponseWriter, req *http.Request) {
	log.Printf("handling /auth/provider\n")
	gothic.BeginAuthHandler(res, req)
}

// GET /callback/{provider}  Called by provider after authorization is granted
func handleCallbackProvider(res http.ResponseWriter, req *http.Request) {
	var (
		status string
		result string
	)
	provider, errProvider := gothic.GetProviderName(req)
	user, errAuth := gothic.CompleteUserAuth(res, req)
	status = "error"
	if errProvider != nil {
		log.Printf("provider failed with '%s'\n", errProvider)
		result = fmt.Sprintf("%s", errProvider)
	} else if errAuth != nil {
		log.Printf("auth failed with '%s'\n", errAuth)
		result = fmt.Sprintf("%s", errAuth)
	} else {
		log.Printf("logged in as %s user: %s (%s)\n", user.Provider, user.Email, user.AccessToken)
		status = "success"
		result = fmt.Sprintf(`{"token":"%s", "provider":"%s"}`, user.AccessToken, user.Provider)
	}
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(fmt.Sprintf(script, status, provider, result)))
}

// GET /refresh
func handleRefresh(res http.ResponseWriter, req *http.Request) {
	log.Printf("refresh with '%s'\n", req)
	res.Write([]byte(""))
}

// GET /success
func handleSuccess(res http.ResponseWriter, req *http.Request) {
	log.Printf("success with '%s'\n", req)
	res.Write([]byte(""))
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not present. Loading directly from environment")
	}
	if portEnv, ok := os.LookupEnv("PORT"); ok {
		port = portEnv
	}
	if callbackEnv, ok := os.LookupEnv("CALLBACK_HOST"); ok {
		callbackHost = callbackEnv
	}

	if gitlabServer, ok := os.LookupEnv("GITLAB_SERVER"); ok {
		goth.UseProviders(
			gitlab.NewCustomisedURL(
				os.Getenv("GITLAB_KEY"), os.Getenv("GITLAB_SECRET"),
				fmt.Sprintf("%s/callback/gitlab", callbackHost),
				fmt.Sprintf("https://%s/oauth/authorize", gitlabServer),
				fmt.Sprintf("https://%s/oauth/token", gitlabServer),
				fmt.Sprintf("https://%s/api/v3/user", gitlabServer),
			),
		)
	}
	if githubHost, ok := os.LookupEnv("GITHUB_HOST"); ok {
		goth.UseProviders(
			github.NewCustomisedURL(
				os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"),
				fmt.Sprintf("%s/callback/github", callbackHost),
				fmt.Sprintf("%s/login/oauth/authorize", githubHost),
				fmt.Sprintf("%s/login/oauth/access_token", githubHost),
				fmt.Sprintf("%s/api/v3/user", githubHost),
				fmt.Sprintf("%s/api/v3/user/emails", githubHost),
				"repo",
			),
		)
	} else {
		goth.UseProviders(
			github.New(
				os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"),
				fmt.Sprintf("%s/callback/github", callbackHost, ),
				"repo",
			),
		)
	}
	goth.UseProviders(
		bitbucket.New(
			os.Getenv("BITBUCKET_KEY"), os.Getenv("BITBUCKET_SECRET"),
			fmt.Sprintf("%s/callback/bitbucket", callbackHost),
		),
	)
}

func main() {
	router := pat.New()
	router.Get("/callback/{provider}", handleCallbackProvider)
	router.Get("/auth/{provider}", handleAuthProvider)
	router.Get("/auth", handleAuth)
	router.Get("/refresh", handleRefresh)
	router.Get("/success", handleSuccess)
	router.Get("/", handleMain)

	http.Handle("/", router)

	log.Printf("started running on %s\n", ":"+port)
	log.Println(http.ListenAndServe(":"+port, nil))
}
