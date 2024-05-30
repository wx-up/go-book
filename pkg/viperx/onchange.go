package viperx

import "github.com/fsnotify/fsnotify"

type OnChangeFs []func(in fsnotify.Event)

func (o *OnChangeFs) Add(f func(in fsnotify.Event)) {
	*o = append(*o, f)
}

var DefaultOnChangeFs OnChangeFs

func AddOnChangeFunc(f func(in fsnotify.Event)) {
	DefaultOnChangeFs.Add(f)
}

func GetOnChangeFs() []func(in fsnotify.Event) {
	return DefaultOnChangeFs
}
