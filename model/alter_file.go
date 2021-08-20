package model

import "os"

func ChangePath(old,new string) error {
	return os.Rename(old, new)
}

