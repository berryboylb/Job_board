package user

import (
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"io"
	"log"
	"net/http"

	"job_board/helpers"
	"job_board/models"
	"job_board/notifications"
)

func User(ctx *gin.Context) {
	//i just want to get the user data

	value, exists := ctx.Get("user")
	if !exists {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    "User not found in session",
			StatusCode: http.StatusUnauthorized,
			Data:       nil,
		})
		return
	}

	user, ok := value.(models.User)
	if !ok {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    "Mismatching types",
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	// Now 'user' contains the user if it exists, and you can proceed with further processing
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully fetched user",
		StatusCode: http.StatusOK,
		Data:       user,
	})
}

func GetAllUsers(ctx *gin.Context) {
	query := FilterDetails{
		Name:         ctx.Query("name"),
		Email:        ctx.Query("email"),
		Picture:      ctx.Query("picture"),
		RoleName:     ctx.Query("role"),
		ProviderID:   ctx.Query("provider_id"),
		SubscriberID: ctx.Query("subscriber_id"),
		MobileNumber: ctx.Query("mobile_number"),
	}
	//i just want to get the user data
	users, total, page, perPage, err := GetUsers(query, ctx.Query("page"), ctx.Query("pageNumber"))

	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}

	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully fetched users",
		StatusCode: http.StatusOK,
		Data: map[string]interface{}{
			"data":     users,
			"total":    total,
			"page":     page,
			"per_page": perPage,
		},
	})
}

func CreateAdmin(ctx *gin.Context) {
	var req Admin
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	subscriberID := uuid.NewString()
	providerID := "admin|" + subscriberID
	newUser := models.User{
		Name:         req.Name,
		MobileNumber: &req.MobileNumber,
		Picture:      req.Picture,
		RoleName:     models.AdminRole,
		ProviderID:   providerID,
		SubscriberID: subscriberID,
		Password:     req.Password,
		Email:        req.Email,
	}

	user, err := CreateAdminUser(newUser)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}

	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully created admin",
		StatusCode: http.StatusOK,
		Data:       user,
	})

	go func() {
		log.Print("Try creating subscriber ")
		subscriber := notifications.Subscriber{
			SubscriberID: user.SubscriberID,
			Name:         user.Name,
			Email:        user.Email,
			Avatar:       user.Picture,
			Data:         map[string]interface{}{},
		}

		if _, err := notifications.CreateSubscriber(subscriber); err != nil {
			log.Printf("Failed to create subscriber: %v", err)
			return
		}
		log.Printf("finished creating subscriber ID: %v", subscriber.SubscriberID)
		log.Print("Send notification ")
		notification := notifications.Trigger{
			// Name:         user.Name,
			// Email:        user.Email,
			// Title:        "You have the power",
			// SubscriberID: user.SubscriberID,
			EventID:      "wel",
			To: map[string]interface{}{
				"subscriberId": user.SubscriberID,
				"phone":        user.MobileNumber,
				"email":        user.Email,
			},
			Data: map[string]interface{}{
				"companyName": "Jobby",
				"password":    req.Password,
			},
		}

		if _, err := notifications.SendNotification(notification); err != nil {
			log.Printf("Failed to send notification: %v", err)
			return
		}
		log.Print("Successfully sent welcome notification to admin")

	}()
}

func UpdateUser(ctx *gin.Context) {
	value, exists := ctx.Get("user")
	if !exists {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    "User not found in session",
			StatusCode: http.StatusUnauthorized,
			Data:       nil,
		})
		return
	}

	user, ok := value.(models.User)
	if !ok {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    "Mismatching types",
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}

	fmt.Println("GOT HERE")

	var req UserDetails
	if err := ctx.ShouldBindJSON(&req); err != nil {
		if err == io.EOF {
			helpers.CreateResponse(ctx, helpers.Response{
				Message:    fmt.Sprintf("please add least one value: %v", err.Error()),
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
	newUser, err := UpdateSingleUser(user.ID, req)
	session := sessions.Default(ctx)
	session.Set(newUser.ProviderID, newUser)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}

	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully updated user",
		StatusCode: http.StatusOK,
		Data:       newUser,
	})
}

func DeleteUser(ctx *gin.Context) {
	value, exists := ctx.Get("user")
	if !exists {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    "User not found in session",
			StatusCode: http.StatusUnauthorized,
			Data:       nil,
		})
		return
	}

	user, ok := value.(models.User)
	if !ok {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    "Mismatching types",
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	err := DeleteSingleUser(user.ID)

	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusUnauthorized,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully deleted user",
		StatusCode: http.StatusOK,
		Data:       nil,
	})
}

func ReinStateAccount(ctx *gin.Context) {
	userID := ctx.Param("id")
	user, err := Reinstate(userID)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully reinstated user",
		StatusCode: http.StatusOK,
		Data:       user,
	})
}
