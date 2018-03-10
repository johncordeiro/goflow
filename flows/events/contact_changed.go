package events

import "github.com/nyaruka/goflow/flows"
import "encoding/json"

// TypeContactChanged is the type of our set contact event
const TypeContactChanged string = "contact_changed"

// ContactChangedEvent events are created to set a contact on a session
//
// ```
//   {
//     "type": "contact_changed",
//     "created_on": "2006-01-02T15:04:05Z",
//     "contact": {
//       "uuid": "0e06f977-cbb7-475f-9d0b-a0c4aaec7f6a",
//       "name": "Bob",
//       "urns": ["tel:+11231234567"]
//     }
//   }
// ```
//
// @event contact_changed
type ContactChangedEvent struct {
	BaseEvent
	Contact json.RawMessage `json:"contact"`
}

// Type returns the type of this event
func (e *ContactChangedEvent) Type() string { return TypeContactChanged }

// AllowedOrigin determines where this event type can originate
func (e *ContactChangedEvent) AllowedOrigin() flows.EventOrigin { return flows.EventOriginEither }

// Validate validates our event is valid and has all the assets it needs
func (e *ContactChangedEvent) Validate(assets flows.SessionAssets) error {
	return nil
}

// Apply applies this event to the given run
func (e *ContactChangedEvent) Apply(run flows.FlowRun) error {
	contact, err := flows.ReadContact(run.Session(), e.Contact)
	if err != nil {
		return err
	}

	run.SetContact(contact)
	run.Session().SetContact(contact)
	return nil
}
