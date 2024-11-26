// Description: This package is used to stop the program gracefully.

package stopper

import (
	"os"
	"os/signal"
	"sync/atomic"

	"github.com/rs/zerolog/log"
)

// flag if the program should stop. 1 - stop, 0 - continue
var stop int64

// Stop returns true if the program should stop
// used atomic operation to read the flag
func Stop() bool {
	return atomic.LoadInt64(&stop) == 1
}

// WaitForCancel - blocks the execution until the the user presses Ctrl+C.
// Then sets the stop flag to 1
func WaitForCancel() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Info().Msg("Received an interrupt signal. Stopping the application.")
	atomic.StoreInt64(&stop, 1)
}
