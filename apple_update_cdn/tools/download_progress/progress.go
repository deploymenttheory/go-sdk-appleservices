// Package download_progress provides a terminal progress bar for streaming
// file downloads.
package download_progress

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

const barWidth = 30

// Bar writes a download progress bar to w (defaults to os.Stderr when nil).
// Call New once per download, then pass Bar.Callback to DownloadFileV1.
//
// Example output (overwrites the same line):
//
//	Downloading UniversalMac_26.4.1_25E253_Restore.ipsw  [████████████░░░░░░░░░░░░░░░░]  4.5 GB / 19.7 GB  22.8%
type Bar struct {
	w    io.Writer
	mu   sync.Mutex
	done bool
}

// New creates a new Bar that writes to w. Pass nil to write to os.Stderr.
func New(w io.Writer) *Bar {
	if w == nil {
		w = os.Stderr
	}
	return &Bar{w: w}
}

// Callback satisfies the cdn.ProgressFunc signature.
// Pass this to DownloadFileV1 as the progress callback.
func (b *Bar) Callback(filename string) func(written, total int64) {
	return func(written, total int64) {
		b.mu.Lock()
		defer b.mu.Unlock()

		if total <= 0 {
			writtenGB := float64(written) / 1e9
			fmt.Fprintf(b.w, "\rDownloading %-40s  %.2f GB downloaded", filename, writtenGB)
			return
		}

		pct := float64(written) / float64(total)
		filled := int(pct * barWidth)
		if filled > barWidth {
			filled = barWidth
		}
		empty := barWidth - filled

		bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
		writtenGB := float64(written) / 1e9
		totalGB := float64(total) / 1e9

		fmt.Fprintf(b.w, "\rDownloading %-40s [%s]  %5.2f GB / %.2f GB  %5.1f%%",
			filename, bar, writtenGB, totalGB, pct*100)

		if written >= total && !b.done {
			b.done = true
			fmt.Fprintln(b.w)
		}
	}
}
