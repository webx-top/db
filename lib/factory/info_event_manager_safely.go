//go:build db_safelyeventmanager

package factory

import "github.com/webx-top/com"

var SafelyEventManager = com.GetenvBool(`DB_SAFELY_EVENT_MANAGER`, true)
