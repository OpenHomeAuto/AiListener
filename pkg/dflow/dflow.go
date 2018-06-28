package dflow

import (
	dflow "cloud.google.com/go/dialogflow/apiv2"
	"context"
	"github.com/OpenHomeAuto/AiListener/pkg/util"
	"google.golang.org/api/option"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

var projectID = "projects/homeautomation-e931b"

func DoSignIn(session string) (*dialogflowpb.DetectIntentResponse, error) {
	ctx := context.Background()
	c, err := dflow.NewSessionsClient(ctx, option.WithCredentialsFile(*util.ServiceAccountFilePath))
	if err != nil {
		return nil, err
	}
	req := &dialogflowpb.DetectIntentRequest{
		Session: session,
		QueryInput: &dialogflowpb.QueryInput{
			Input: &dialogflowpb.QueryInput_Event{
				Event: &dialogflowpb.EventInput{
					Name:         "actions.intent.SIGN_IN",
					LanguageCode: "de",
				},
			},
		},
	}
	return c.DetectIntent(ctx, req)
}
