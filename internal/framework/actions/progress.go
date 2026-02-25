// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package actions

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/action"
)

type SendProgressFunc func(string, ...any)

func NewSendProgressFunc(response *action.InvokeResponse) SendProgressFunc {
	return func(format string, a ...any) {
		response.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf(format, a...),
		})
	}
}
