package notification

// import (
// 	"log"

// 	"github.com/appleboy/go-fcm"

// 	"context"

// 	firebase "firebase.google.com/go"
// 	"google.golang.org/api/option"
// )

// func init() {
// 	ctx := context.Background()

// 	opt := option.WithCredentialsJSON([]byte{})
// 	app, err := firebase.NewApp(ctx, nil, opt)
// 	if err != nil {
// 		panic(err)
// 	}

// 	client, err := app.Messaging(ctx)

// 	panic(client)
// 	// client.
// }

// type FirebaseNotification struct {
// 	client *fcm.Client
// }

// func NewFirebaseNotification(apikey string) (Notification, error) {
// 	client, err := fcm.NewClient(apikey)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &FirebaseNotification{
// 		client: client,
// 	}, nil
// }

// func (n *FirebaseNotification) Send() error {

// }

// func init() {
// 	// Create the message to be sent.
// 	msg := &fcm.Message{
// 		To: "sample_device_token",
// 		Data: map[string]interface{}{
// 			"foo": "bar",
// 		},
// 		Notification: &fcm.Notification{
// 			Title: "title",
// 			Body:  "body",
// 		},
// 	}

// 	// Create a FCM client to send the message.
// 	client, err := fcm.NewClient("sample_api_key")
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	// Send the message and receive the response without retries.
// 	response, err := client.Send(msg)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	log.Printf("%#v\n", response)
// }
