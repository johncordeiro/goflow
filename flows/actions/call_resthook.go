package actions

import (
	"net/http"
	"strings"

	"github.com/nyaruka/goflow/assets/rest/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

func init() {
	RegisterType(TypeCallResthook, func() flows.Action { return &CallResthookAction{} })
}

// TypeCallResthook is the type for the call resthook action
const TypeCallResthook string = "call_resthook"

// CallResthookAction can be used to call a resthook.
//
// A [event:webhook_called] event will be created for each subscriber of the resthook with the results
// of the HTTP call. If the action has `result_name` set, a result will
// be created with that name, and if the resthook returns valid JSON, that will be accessible
// through `extra` on the result.
//
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "call_resthook",
//     "resthook": "new-registration"
//   }
//
// @action call_resthook
type CallResthookAction struct {
	BaseAction
	onlineAction

	Resthook   string `json:"resthook" validate:"required"`
	ResultName string `json:"result_name,omitempty"`
}

// Type returns the type of this action
func (a *CallResthookAction) Type() string { return TypeCallResthook }

// Validate validates our action is valid and has all the assets it needs
func (a *CallResthookAction) Validate(assets flows.SessionAssets) error {
	return nil
}

// Execute runs this action
func (a *CallResthookAction) Execute(run flows.FlowRun, step flows.Step) error {
	// if resthook doesn't exist, treat it like an existing one with no subscribers
	resthook := run.Session().Assets().Resthooks().FindBySlug(a.Resthook)
	if resthook == nil {
		resthook = flows.NewResthook(types.NewResthook(a.Resthook, nil))
	}

	// build our payload
	payload, err := run.EvaluateTemplateAsString(flows.DefaultWebhookPayload, false)
	if err != nil {
		a.logError(run, step, err)
	}

	// make a call to each subscriber URL
	webhooks := make([]*flows.WebhookCall, 0, len(resthook.Subscribers()))

	for _, url := range resthook.Subscribers() {
		req, err := http.NewRequest("POST", url, strings.NewReader(payload))
		if err != nil {
			a.logError(run, step, err)
			return nil
		}

		req.Header.Add("Content-Type", "application/json")

		webhook, err := flows.MakeWebhookCall(run.Session(), req)
		if err != nil {
			a.logError(run, step, err)
		} else {
			webhooks = append(webhooks, webhook)
			a.log(run, step, events.NewWebhookCalledEvent(webhook, a.Resthook))
		}
	}

	asResult := a.pickResultWebhook(webhooks)
	if asResult != nil && a.ResultName != "" {
		a.saveWebhookResult(run, step, a.ResultName, asResult)
	}

	return nil
}

// picks one of the resthook calls to become the result generated by this action
func (a *CallResthookAction) pickResultWebhook(calls []*flows.WebhookCall) *flows.WebhookCall {
	var lastFailure, asResult *flows.WebhookCall
	for _, call := range calls {
		if call.Status() == flows.WebhookStatusSuccess {
			asResult = call
		} else {
			lastFailure = call
		}
	}
	if lastFailure != nil {
		asResult = lastFailure
	}
	return asResult
}
