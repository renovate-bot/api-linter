// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aip0214

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestResourceExpiry(t *testing.T) {
	for _, test := range []struct {
		name                  string
		FieldBehavior         string
		AddtlField            string
		problems              testutils.Problems
		AddResourceDefinition bool
	}{
		{"ValidTtl", "", "google.protobuf.Duration ttl = 3;", testutils.Problems{}, true},
		{"ValidOutputOnly", "[(google.api.field_behavior) = OUTPUT_ONLY]", "", testutils.Problems{}, true},
		{"Invalid", "", "", testutils.Problems{{Message: "ttl"}}, true},
		{"SkipANonResource", "", "", testutils.Problems{}, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/field_behavior.proto";
				import "google/api/resource.proto";
				import "google/protobuf/timestamp.proto";
				import "google/protobuf/duration.proto";

				message Book {
					{{if .AddResourceDefinition}}
					option (google.api.resource) = { type: "library.googleapis.com/Book" };
					{{end}}
					string name = 1;
					google.protobuf.Timestamp expire_time = 2 {{.FieldBehavior}};
					{{.AddtlField}}
				}
			`, test)
			field := f.GetMessageTypes()[0].GetFields()[1]
			if diff := test.problems.SetDescriptor(field).Diff(resourceExpiry.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
