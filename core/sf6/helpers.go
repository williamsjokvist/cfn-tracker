package sf6

import (
	"context"
	"fmt"

	"github.com/williamsjokvist/cfn-tracker/core/common"
)

// Make a SF6Tracker and expose it as a GameTracker
func MakeSF6Tracker(ctx context.Context, browser *common.Browser, username string, password string) (common.GameTracker, error) {
	fmt.Println(`making sf6 tracker`)
	sf6Tracker := NewSF6Tracker(ctx, browser)
	err := sf6Tracker.Authenticate(username, password)
	if err != nil {
		fmt.Printf(`auth err: %v`, err)
		return nil, err
	}

	var gt common.GameTracker = sf6Tracker
	return gt, nil
}
