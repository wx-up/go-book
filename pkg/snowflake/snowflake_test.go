package snowflake

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	tm, err := time.ParseInLocation(time.DateTime, "2024-04-01 00:00:00", time.Local)
	require.NoError(t, err)
	fmt.Println(tm.UnixMilli())

	fmt.Println(int64(-1) ^ (-1 << timestampBits))
	fmt.Println(int64(math.Pow(2, 41)))
}

func Test_Generate(t *testing.T) {
	sf := NewSnowFlake(1)
	fmt.Println(fmt.Sprintf("%b", sf.Generate()))
}
