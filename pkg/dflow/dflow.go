package dflow

import (
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

var projectID = "projects/homeautomation-e931b"

func DoSignIn(session string) *dialogflowpb.DetectIntentRequest {
	req := &dialogflowpb.DetectIntentRequest{
		Session: session,
		QueryInput: &dialogflowpb.QueryInput{
			Input: &dialogflowpb.QueryInput_Event{
				Event: &dialogflowpb.EventInput{
					Name:         "actions_intent_SIGN_IN",
					LanguageCode: "de",
				},
			},
		},
	}
	return req
}
