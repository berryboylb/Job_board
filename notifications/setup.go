package notifications

import (
	novu "github.com/novuhq/go-novu/lib"

	"github.com/google/uuid"
	"github.com/joho/godotenv"

	"context"
	"fmt"
	"log"
	"os"
)

var novuClient *novu.APIClient

var ctx = context.Background()

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	novuApiKey := os.Getenv("NOVU_API_KEY")

	//check if key isn't empty
	if novuApiKey == "" {
		log.Fatal("Error loading db port env variables")
	}

	novuClient = novu.NewAPIClient(novuApiKey, &novu.Config{})
}

func CreateSubscriber(userDetails Subscriber) error {
	subscriberID := uuid.New().String()
	subscriber := novu.SubscriberPayload{
		FirstName: userDetails.Name,
		Email:     userDetails.Email,
		Avatar:    userDetails.Avatar,
		Data:      userDetails.Data,
	}
	resp, err := novuClient.SubscriberApi.Identify(ctx, subscriberID, subscriber)
	if err != nil {
		return err
	}
	fmt.Println(resp)
	return nil
}

func UpdateSubscriber(name string) error {
	subscriberID := uuid.New().String()
	updateSubscriber := novu.SubscriberPayload{FirstName: name}
	resp, err := novuClient.SubscriberApi.Update(ctx, subscriberID, updateSubscriber)
	if err != nil {
		return err
	}
	fmt.Println(resp)
	return nil
}

func TriggerNotification(payload Trigger) error {
	to := map[string]interface{}{
		"lastName":     "",
		"firstName":    payload.Name,
		"subscriberId": payload.SubscriberID,
		"email":        payload.Email,
	}
	data := map[string]interface{}{
		"name": payload.Title,
		"organization": map[string]interface{}{
			"logo": payload.Logo,
		},
	}
	triggerResp, err := novuClient.EventApi.Trigger(ctx, payload.EventID, novu.ITriggerPayloadOptions{
		To:      to,
		Payload: data,
	})
	if err != nil {
		return err
	}

	fmt.Println(triggerResp)
	return nil
}

func CreateTopic()  {
	
}