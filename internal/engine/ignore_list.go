package engine

import (
	"path/filepath"

	"github.com/ThisaruGuruge/bestow/internal/config"
)

type IgnoreList struct {
	src   string
	items []string
}

func newIgnoreList(src string) (*IgnoreList, error) {
	list := &IgnoreList{src: src}

	//TODO: Should return custom error?
	// Load global ignore list
	configHome := config.AppConfigHome()
	if err := readIgnoreFile(configHome, &list.items); err != nil {
		return nil, err
	}

	//load package ignore list
	if err := readIgnoreFile(src, &list.items); err != nil {
		return nil, err
	}
	return list, nil
}

func (i *IgnoreList) forPackage(pkg string) ([]string, error) {
	result := append([]string(nil), i.items...)
	if err := readIgnoreFile(filepath.Join(i.src, pkg), &result); err != nil {
		return nil, err
	}
	return result, nil
}
