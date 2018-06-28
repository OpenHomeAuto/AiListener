package dflow

type followupEvent struct {
	Name         string      `json:"name"`
	LanguageCode string      `json:"languageCode"`
	Parameters   interface{} `json:"parameters,omitempty"`
}

func DoSignIn() *followupEvent {
	req := &followupEvent{
		Name:         "actions_intent_SIGN_IN",
		LanguageCode: "de",
	}
	return req
}
