package main

import (
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// watch は memoid no nozuki の見張り。フォルダならその直下、ファイルならそのファイルが変わるたびに、
// 変わったものの名前の最後の部分を send に渡す。done が閉じたら止める。
func watch(path string, send func(name string), done <-chan struct{}) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	// ファイルは、エディタが保存のときに置き換えることがあるので、入っているフォルダを見張って名前で選ぶ (fsnotify の勧め)
	dir, only := path, ""
	if !info.IsDir() {
		dir, only = filepath.Dir(path), filepath.Clean(path)
	}
	if err := w.Add(dir); err != nil {
		w.Close()
		return err
	}
	go func() {
		defer w.Close()
		for {
			select {
			case e, ok := <-w.Events:
				if !ok {
					return
				}
				// 権限だけの変更は、Node.js の fs.watch も中身の変更として知らせないので飛ばす
				if e.Op == fsnotify.Chmod || (only != "" && filepath.Clean(e.Name) != only) {
					continue
				}
				send(filepath.Base(e.Name))
			case _, ok := <-w.Errors:
				if !ok {
					return
				}
			case <-done:
				return
			}
		}
	}()
	return nil
}
