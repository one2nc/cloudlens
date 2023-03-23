package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

const (
	// DefaultDirMod default unix perms for k9s directory.
	DefaultDirMod os.FileMode = 0755
	// DefaultFileMod default unix perms for k9s files.
	DefaultFileMod os.FileMode = 0600
)

type (
	IsSwapHappen bool
)

// EnsurePath ensures a directory exist from the given path.
func EnsurePath(path string, mod os.FileMode) {
	dir := filepath.Dir(path)
	EnsureFullPath(dir, mod)
}

// EnsureFullPath ensures a directory exist from the given path.
func EnsureFullPath(path string, mod os.FileMode) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(path, mod); err != nil {
			log.Fatal().Msgf("Unable to create dir %q %v", path, err)
		}
	}
}

func LookupForValue(profiles []string, value string) bool {
	for _, got := range profiles {
		if strings.EqualFold(got, value) {
			return true
		}
	}
	return false
}

// SwapFirstIndexWithValue return swapped array if match found. If match not found returns same array and says match not found.
func SwapFirstIndexWithValue(array []string, value string) ([]string, IsSwapHappen) {
	if len(array) > 0 {
		var isSwapped IsSwapHappen
		for i, got := range array {
			if strings.EqualFold(got, value) {
				array[0], array[i] = array[i], array[0]
				isSwapped = true
				break
			}
		}
		return array, isSwapped
	} else {
		return nil, false
	}
}
