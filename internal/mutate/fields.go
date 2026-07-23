package mutate

import (
	"fmt"

	"github.com/colchuck-ai/intent/internal/model"
)

// setField writes value into the named scalar field of el. It rejects unknown
// fields, list fields, and edges (the last go through Link) with a message that
// names what is settable, so a wrong field is a clear error rather than a silent
// no-op.
func setField(el any, addr, field, value string) error {
	switch v := el.(type) {
	case *model.Product:
		switch field {
		case "name":
			v.Name = value
		case "summary":
			v.Summary = value
		case "detail":
			v.Detail = value
		default:
			return unsettable(addr, field, "name, summary, detail")
		}
	case *model.Engineering:
		switch field {
		case "name":
			v.Name = value
		case "summary":
			v.Summary = value
		case "detail":
			v.Detail = value
		default:
			return unsettable(addr, field, "name, summary, detail")
		}
	case *model.Job:
		switch field {
		case "name":
			v.Name = value
		case "story":
			v.Story = value
		case "detail":
			v.Detail = value
		default:
			return unsettable(addr, field, "name, story, detail")
		}
	case *model.Outcome:
		switch field {
		case "name":
			v.Name = value
		case "statement":
			v.Statement = value
		case "detail":
			v.Detail = value
		case "type":
			if err := checkType(value); err != nil {
				return err
			}
			v.Type = value
		default:
			return unsettable(addr, field, "name, statement, detail, type")
		}
	case *model.Risk:
		switch field {
		case "name":
			v.Name = value
		case "statement":
			v.Statement = value
		case "detail":
			v.Detail = value
		default:
			return unsettable(addr, field, "name, statement, detail")
		}
	case *model.Requirement:
		switch field {
		case "name":
			v.Name = value
		case "statement":
			v.Statement = value
		case "detail":
			v.Detail = value
		case "type":
			if err := checkType(value); err != nil {
				return err
			}
			v.Type = value
		case "mitigates", "dependsOn":
			return errEdgeField(addr, field)
		case "acceptance_criteria":
			return errListField(addr, field)
		default:
			return unsettable(addr, field, "name, statement, detail, type")
		}
	case *model.Component:
		switch field {
		case "name":
			v.Name = value
		case "responsibility":
			v.Responsibility = value
		case "detail":
			v.Detail = value
		case "type":
			if err := checkType(value); err != nil {
				return err
			}
			v.Type = value
		case "data_model":
			v.DataModel = value
		case "interfaces":
			v.Interfaces = value
		case "behavior":
			v.Behavior = value
		case "fulfills", "relationships":
			return errEdgeField(addr, field)
		case "edge_cases", "success_criteria":
			return errListField(addr, field)
		default:
			return unsettable(addr, field, "name, responsibility, detail, type, data_model, interfaces, behavior")
		}
	case *model.DecisionRecord:
		switch field {
		case "name":
			v.Name = value
		case "summary":
			v.Summary = value
		case "context":
			v.Context = value
		case "decision":
			v.Decision = value
		case "affects":
			return errEdgeField(addr, field)
		case "options", "consequences":
			return errListField(addr, field)
		default:
			return unsettable(addr, field, "name, summary, context, decision")
		}
	case *model.ChangeRecord:
		switch field {
		case "name":
			v.Name = value
		case "change":
			v.Change = value
		case "rationale":
			v.Rationale = value
		case "affects":
			return errEdgeField(addr, field)
		default:
			return unsettable(addr, field, "name, change, rationale")
		}
	case *string:
		// A principle or constraint: its value is its statement.
		if field != "statement" {
			return unsettable(addr, field, "statement")
		}
		*v = value
	default:
		return fmt.Errorf("cannot set fields on %s", addr)
	}
	return nil
}

func checkType(value string) error {
	if value != "inline" && value != "document" {
		return fmt.Errorf("type must be inline or document, got %q", value)
	}
	return nil
}

func unsettable(addr, field, settable string) error {
	return fmt.Errorf("%s has no settable field %q (settable: %s)", addr, field, settable)
}

// errEdgeField explains that an edge field is managed with link/unlink, not set.
func errEdgeField(addr, field string) error {
	return fmt.Errorf("%s.%s is an edge; use `intent link` / `intent unlink` instead", addr, field)
}

// errListField explains that a list field isn't editable via set in this phase.
func errListField(addr, field string) error {
	return fmt.Errorf("%s.%s is a list; not settable via `set` in this phase", addr, field)
}
