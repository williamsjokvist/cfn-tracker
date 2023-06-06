package sfv

import (
	"context"
	"fmt"

	"github.com/williamsjokvist/cfn-tracker/core/common"
)

// Make a SFVTracker and expose it as a GameTracker
func MakeSFVTracker(ctx context.Context, browser *common.Browser, username string, password string) (common.GameTracker, error) {
	sfvTracker := NewSFVTracker(ctx, browser)
	err := sfvTracker.Authenticate(username, password)
	if err != nil {
		fmt.Printf(`auth err: %v`, err)
		return nil, err
	}

	var gt common.GameTracker = sfvTracker
	return gt, nil
}
