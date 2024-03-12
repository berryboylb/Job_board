package auth

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"gorm.io/gorm"

	"job_board/helpers"
	"job_board/jwt"
	"job_board/models"
	"job_board/notifications"
	"job_board/user"
)

func Login(auth *Authenticator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userType := ctx.Query("type")
		if userType == "" {
			helpers.CreateResponse(ctx, helpers.Response{
				Message:    "user type is required",
				StatusCode: http.StatusBadRequest,
				Data:       nil,
			})
			return
		}

		var role models.RoleAllowed
		switch userType {
		case string(models.PosterRole):
			role = models.PosterRole
		case string(models.UserRole):
			role = models.UserRole
		default:
			helpers.CreateResponse(ctx, helpers.Response{
				Message:    fmt.Sprintf("user type can be either %v or %v", string(models.PosterRole), string(models.UserRole)),
				StatusCode: http.StatusBadRequest,
				Data:       nil,
			})
			return
		}

		state, err := generateRandomState()
		if err != nil {
			helpers.CreateResponse(ctx, helpers.Response{
				Message:    err.Error(),
				StatusCode: http.StatusInternalServerError,
				Data:       nil,
			})
			return
		}

		// Save the state inside the session.
		session := sessions.Default(ctx)
		session.Set("state", state)
		session.Set("type", role)
		if err := session.Save(); err != nil {
			helpers.CreateResponse(ctx, helpers.Response{
				Message:    err.Error(),
				StatusCode: http.StatusInternalServerError,
				Data:       nil,
			})
			return
		}

		ctx.Redirect(http.StatusTemporaryRedirect, auth.AuthCodeURL(state))
	}
}

func Callback(auth *Authenticator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		if ctx.Query("state") != session.Get("state") {
			helpers.CreateResponse(ctx, helpers.Response{
				Message:    "Invalid state parameter.",
				StatusCode: http.StatusBadRequest,
				Data:       nil,
			})
			return
		}

		audience := "https://" + os.Getenv("AUTH0_DOMAIN") + "/"
		// Exchange an authorization code for a token.
		token, err := auth.Exchange(ctx.Request.Context(), ctx.Query("code"), oauth2.SetAuthURLParam("audience", audience))
		if err != nil {
			helpers.CreateResponse(ctx, helpers.Response{
				Message:    "Failed to exchange an authorization code for a token.",
				StatusCode: http.StatusUnauthorized,
				Data:       nil,
			})
			return
		}

		idToken, err := auth.VerifyIDToken(ctx.Request.Context(), token)
		if err != nil {
			helpers.CreateResponse(ctx, helpers.Response{
				Message:    "Failed to verify ID Token.",
				StatusCode: http.StatusInternalServerError,
				Data:       nil,
			})
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
			helpers.CreateResponse(ctx, helpers.Response{
				Message:    err.Error(),
				StatusCode: http.StatusInternalServerError,
				Data:       nil,
			})
			return
		}

		session.Set("access_token", token)
		session.Set("profile", profile)
		session.Set("subject", sub)
		if err := session.Save(); err != nil {
			helpers.CreateResponse(ctx, helpers.Response{
				Message:    err.Error(),
				StatusCode: http.StatusInternalServerError,
				Data:       nil,
			})
			return
		}

		// Redirect to logged in page.
		ctx.Redirect(http.StatusTemporaryRedirect, "/api/v1/auth/authorize")
	}
}

func IsAuthenticated(ctx *gin.Context) {
	if sessions.Default(ctx).Get("profile") == nil {
		ctx.Redirect(http.StatusSeeOther, "/")
		return

	} else {
		ctx.Next()
	}
}

func Authorize(ctx *gin.Context) {
	session := sessions.Default(ctx)
	subject, ok := session.Get("subject").(string)
	if !ok {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    "invalid auth0 subject",
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	profile, isNew, err := handleUser(subject, session)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	fmt.Println("got here")
	userAgent := ctx.Request.Header.Get("User-Agent")
	isMobile := strings.Contains(userAgent, "Android") || strings.Contains(userAgent, "iPhone") || strings.Contains(userAgent, "iPad")
	token, err := jwt.GenerateJWT(profile.ProviderID, isMobile)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	//save to db and if user don't exist redirect to token
	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"profile":      profile,
			"subject":      subject,
			"access_token": token,
		},
	})

	go func() {
		if isNew {
			log.Print("Try creating subscriber ")
			subscriber := notifications.Subscriber{
				SubscriberID: profile.SubscriberID,
				Name:         profile.Name,
				Email:        profile.Email,
				Avatar:       profile.Picture,
				Data:         map[string]interface{}{},
			}

			if _, err := notifications.CreateSubscriber(subscriber); err != nil {
				log.Printf("Failed to create subscriber: %v", err)
				return
			}
			log.Printf("finished creating subscriber ID: %v", subscriber.SubscriberID)
			if profile.Email != "" {
				log.Print("Send notification ")
				notification := notifications.Trigger{
					Name:         profile.Name,
					Email:        profile.Email,
					Title:        "Welcome to Jobby",
					SubscriberID: profile.SubscriberID,
					EventID:      "welcome",
					Logo:         "https://via.placeholder.com/200x200",
				}

				if _, err := notifications.SendNotification(notification); err != nil {
					log.Printf("Failed to send notification: %v", err)
					return
				}
				log.Print("Successfully sent welcome notification")
			} else {
				log.Print("No email found will be skipping")
			}
		}
	}()
}

func LoginAdmin(ctx *gin.Context) {
	var req Admin
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}

	newuser, err := user.GetAltSingleUser(models.User{
		Email: req.Email,
	})

	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	if newuser.RoleName != models.AdminRole && newuser.RoleName != models.SuperAdminRole {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    "Unauthorized access",
			StatusCode: http.StatusUnauthorized,
			Data:       nil,
		})
		return
	}

	isMatch := helpers.CheckPasswordHash(req.Password, newuser.Password)
	if !isMatch {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    "Invalid credentials",
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "Successfully sent OTP",
		StatusCode: http.StatusOK,
		Data:       nil,
	})
	go func() {
		otp := GenerateOtp(4)
		expiryTime := time.Now().Add(24 * time.Hour) //  24 hours

		_, err = user.UpdateSingleUser(newuser.ID, models.User{
			VerificationToken: otp,
			ExpiresAt:         expiryTime,
		})
		if err != nil {
			log.Println("Failed to update user:", err)
			return
		}

		log.Print("Send otp notification ")
		notification := notifications.Trigger{
			Name:         newuser.Name,
			Email:        newuser.Email,
			Title:        "You have the power",
			SubscriberID: newuser.SubscriberID,
			EventID:      "otp",
			Logo:         "https://via.placeholder.com/200x200",
			To: map[string]interface{}{
				"subscriberId": newuser.SubscriberID,
				"email":        newuser.Email,
			},
			Data: map[string]interface{}{
				"companyName": "Jobby",
				"otp":         otp,
			},
		}

		if _, err := notifications.SendNotification(notification); err != nil {
			log.Printf("Failed to send  otp notification: %v", err)
			return
		}
		log.Printf("Send otp notification: %v ", otp)
	}()

}

func ConfirmLoginAdmin(ctx *gin.Context) {
	var req OtpDto
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	dbUser, err := user.GetVerificationToken(req.Otp)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helpers.CreateResponse(ctx, helpers.Response{
				Message:    "Invalid otp try login again.",
				StatusCode: http.StatusBadRequest,
				Data:       nil,
			})
			return
		}
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	fmt.Println("GOT HERE", dbUser.ProviderID)
	userAgent := ctx.Request.Header.Get("User-Agent")
	isMobile := strings.Contains(userAgent, "Android") || strings.Contains(userAgent, "iPhone") || strings.Contains(userAgent, "iPad")
	token, err := jwt.GenerateJWT(dbUser.ProviderID, isMobile)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "Successfully loggedIn",
		StatusCode: http.StatusOK,
		Data:       token,
	})

	go func() {
		_, err = user.UpdateSingleUser(dbUser.ID, models.User{
			VerificationToken: "nil",
			ExpiresAt:         time.Time{},
		})
		if err != nil {
			log.Println("Failed to update user:", err)
			return
		}
	}()
}

func Protect(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"data": "hello",
	})
}

func Logout(ctx *gin.Context) {
	logoutUrl, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/v2/logout")
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}

	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}

	returnTo, err := url.Parse(scheme + "://" + ctx.Request.Host)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}

	parameters := url.Values{}
	parameters.Add("returnTo", returnTo.String())
	parameters.Add("client_id", os.Getenv("AUTH0_CLIENT_ID"))
	logoutUrl.RawQuery = parameters.Encode()

	ctx.Redirect(http.StatusTemporaryRedirect, logoutUrl.String())
}
