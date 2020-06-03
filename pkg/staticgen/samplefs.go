package staticgen

import "github.com/phogolabs/parcello"

// FS returns embedded filesystem
func FS(path string) parcello.FileSystemManager {
	return parcello.ManagerAt(path)
}
