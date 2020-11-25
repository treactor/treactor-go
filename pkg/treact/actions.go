package treact

import (
	"context"
	"fmt"
	"github.com/treactor/treactor-go/pkg/resource"
	"strconv"
	"time"
)

func mem(ctx context.Context, factorValue string) []byte {

	factor, _ := strconv.ParseFloat(factorValue,32)

	size := factor * 1024 * 1024

	return make([]byte, int(size))
}

func cpu(ctx context.Context, durationValue string) {
	duration, _ := strconv.ParseInt(durationValue, 10, 64)
	start := time.Now()
	for ; true; {
		elapsed := time.Now().Sub(start)
		select {
		case <- ctx.Done():
			resource.Logger.WarningF(ctx, "CPU Action cancelled after %dms (%dms)", elapsed, duration)
			return
		default:
			// nothing
		}
		if elapsed.Milliseconds() > duration {
			return
		}
		fmt.Println(ctx.Err())
	}
}
