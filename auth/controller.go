package auth

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"job_board/models"
)

func Login(auth *Authenticator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userType := ctx.Query("type")
		if userType == "" {
			ctx.String(http.StatusBadRequest, "user type is required")
			return
		}

		var role models.RoleAllowed
		switch userType {
		case string(models.PosterRole):
			role = models.PosterRole
		case string(models.UserRole):
			role = models.UserRole
		default:
			ctx.String(http.StatusBadRequest, fmt.Sprintf("user type can be either %v or %v", string(models.PosterRole), string(models.UserRole)))
			return
		}

		state, err := generateRandomState()
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		// Save the state inside the session.
		session := sessions.Default(ctx)
		session.Set("state", state)
		session.Set("type", role)
		if err := session.Save(); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.Redirect(http.StatusTemporaryRedirect, auth.AuthCodeURL(state))
	}
}

func Callback(auth *Authenticator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		if ctx.Query("state") != session.Get("state") {
			ctx.String(http.StatusBadRequest, "Invalid state parameter.")
			return
		}

		// Exchange an authorization code for a token.
		token, err := auth.Exchange(ctx.Request.Context(), ctx.Query("code"))
		if err != nil {
			ctx.String(http.StatusUnauthorized, "Failed to exchange an authorization code for a token.")
			return
		}

		idToken, err := auth.VerifyIDToken(ctx.Request.Context(), token)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Failed to verify ID Token.")
			return
		}
		sub := strings.Split(idToken.Subject, "|")[0]
		// var profile map[string]interface{}
		// if err := idToken.Claims(&profile); err != nil {
		// 	ctx.String(http.StatusInternalServerError, err.Error())
		// 	return
		// }
		profile, err := decoder(idToken, sub)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		session.Set("access_token", token.AccessToken)
		session.Set("profile", profile)
		session.Set("subject", sub)
		if err := session.Save(); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		// Redirect to logged in page.
		ctx.Redirect(http.StatusTemporaryRedirect, "/authorize")
	}
}

func IsAuthenticated(ctx *gin.Context) {
	if sessions.Default(ctx).Get("profile") == nil {
		ctx.Redirect(http.StatusSeeOther, "/")
	} else {
		ctx.Next()
	}
}

func Authorize(ctx *gin.Context) {
	//decode the profile a normal variable
	//populate the user with the type and the profile data
	//i want to do a firstorCreate method to check if the user is already in the database and return it instead of creating a new one every
	//generate a token with only read access and encode the email address of the user in it
	session := sessions.Default(ctx)
	subject, ok := session.Get("subject").(string)
	if !ok {
		ctx.String(http.StatusInternalServerError, "invalid auth0 subject")
		return
	}
	fmt.Println(subject)
	profile, err := handleUser(subject, session)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	token, err := generateToken()
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	//save to db and if user don't exist redirect to token
	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"profile":      profile,
			"subject":      subject,
			"access_token": token.AccessToken,
		},
	})
}

func User(ctx *gin.Context) {
	//i just want to get the user data
	ctx.JSON(http.StatusOK, gin.H{
		"data": "user data",
	})
}

func Protect(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"data": "hello",
	})
}

func Logout(ctx *gin.Context) {
	logoutUrl, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/v2/logout")
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}

	returnTo, err := url.Parse(scheme + "://" + ctx.Request.Host)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	parameters := url.Values{}
	parameters.Add("returnTo", returnTo.String())
	parameters.Add("client_id", os.Getenv("AUTH0_CLIENT_ID"))
	logoutUrl.RawQuery = parameters.Encode()

	ctx.Redirect(http.StatusTemporaryRedirect, logoutUrl.String())
}
