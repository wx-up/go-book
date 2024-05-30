package viperx

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"testing"
)

func TestAddOnChangeFunc(t *testing.T) {
	AddOnChangeFunc(func(in fsnotify.Event) {
		fmt.Println(11)
	})
	AddOnChangeFunc(func(in fsnotify.Event) {
		fmt.Println(22)
	})
	fmt.Println(GetOnChangeFs())
}
