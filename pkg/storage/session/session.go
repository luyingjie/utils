// Package gsession implements manager and storage features for sessions.
package session

import (
	"errors"

	"github.com/luyingjie/utils/pkg/generates/uid"
)

var (
	ErrorDisabled = errors.New("this feature is disabled in this storage")
)

// NewSessionId creates and returns a new and unique session id string,
// which is in 36 bytes.
func NewSessionId() string {
	return uid.S()
}
