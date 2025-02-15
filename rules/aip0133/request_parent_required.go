package aip0133

import (
	"fmt"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
	"github.com/stoewer/go-strcase"
)

var requestParentRequired = &lint.MessageRule{
	Name:   lint.NewRuleName(133, "request-parent-required"),
	OnlyIf: utils.IsCreateRequestMessage,
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		if m.FindFieldByName("parent") == nil {
			// Sanity check: If the resource has a pattern, and that pattern
			// contains only one variable, then a parent field is not expected.
			//
			// In order to parse out the pattern, we get the resource message
			// from the request, then get the resource annotation from that,
			// and then inspect the pattern there (oy!).
			singular := getResourceMsgNameFromReq(m)
			if field := m.FindFieldByName(strcase.SnakeCase(singular)); field != nil {
				if !utils.HasParent(utils.GetResource(field.GetMessageType())) {
					return nil
				}
			}

			// Nope, this is not the unusual case, and a parent field is expected.
			return []lint.Problem{{
				Message:    fmt.Sprintf("Message %q has no `parent` field", m.GetName()),
				Descriptor: m,
			}}
		}

		return nil
	},
}
