package modifiers

import (
	"encoding/json"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	RegisterType(TypeGroups, readGroupsModifier)
}

// TypeGroups is the type of our groups modifier
const TypeGroups string = "groups"

// GroupsModification is the type of modification to make
type GroupsModification string

// the supported types of modification
const (
	GroupsAdd    GroupsModification = "add"
	GroupsRemove GroupsModification = "remove"
)

// GroupsModifier modifies the group membership of the contact
type GroupsModifier struct {
	baseModifier

	groups       []*flows.Group
	modification GroupsModification
}

// NewGroupsModifier creates a new groups modifier
func NewGroupsModifier(groups []*flows.Group, modification GroupsModification) *GroupsModifier {
	return &GroupsModifier{
		baseModifier: newBaseModifier(TypeGroups),
		groups:       groups,
		modification: modification,
	}
}

// Apply applies this modification to the given contact
func (m *GroupsModifier) Apply(env utils.Environment, assets flows.SessionAssets, contact *flows.Contact, log func(flows.Event)) {
	diff := make([]*flows.Group, 0, len(m.groups))
	if m.modification == GroupsAdd {
		for _, group := range m.groups {

			// ignore group if contact is already in it
			if contact.Groups().FindByUUID(group.UUID()) != nil {
				continue
			}

			contact.Groups().Add(group)
			diff = append(diff, group)
		}

		// only generate event if contact's groups change
		if len(diff) > 0 {
			log(events.NewContactGroupsChangedEvent(diff, nil))
		}
	} else if m.modification == GroupsRemove {
		for _, group := range m.groups {
			// ignore group if contact isn't actually in it
			if contact.Groups().FindByUUID(group.UUID()) == nil {
				continue
			}

			contact.Groups().Remove(group)
			diff = append(diff, group)
		}

		// only generate event if contact's groups change
		if len(diff) > 0 {
			log(events.NewContactGroupsChangedEvent(nil, diff))
		}
	}
}

var _ flows.Modifier = (*GroupsModifier)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type groupsModifierEnvelope struct {
	utils.TypedEnvelope
	Groups       []*assets.GroupReference `json:"groups" validate:"required,dive"`
	Modification GroupsModification       `json:"modification" validate:"eq=add|eq=remove"`
}

func readGroupsModifier(assets flows.SessionAssets, data json.RawMessage) (flows.Modifier, error) {
	e := &groupsModifierEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	groups := make([]*flows.Group, len(e.Groups))
	var err error
	for g, groupRef := range e.Groups {
		groups[g], err = assets.Groups().Get(groupRef.UUID)
		if err != nil {
			return nil, err
		}
	}

	return NewGroupsModifier(groups, e.Modification), nil
}

func (m *GroupsModifier) MarshalJSON() ([]byte, error) {
	groupRefs := make([]*assets.GroupReference, len(m.groups))
	for g := range m.groups {
		groupRefs[g] = m.groups[g].Reference()
	}

	return json.Marshal(&groupsModifierEnvelope{
		TypedEnvelope: utils.TypedEnvelope{Type: m.Type()},
		Groups:        groupRefs,
		Modification:  m.modification,
	})
}
