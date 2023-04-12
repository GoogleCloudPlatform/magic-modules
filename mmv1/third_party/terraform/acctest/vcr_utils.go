package acctest

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var SourcesLock = sync.RWMutex{}

var Sources map[string]VcrSource

// VcrSource is a source for a given VCR test with the value that seeded it
type VcrSource struct {
	Seed   int64
	source rand.Source
}

func IsVcrEnabled() bool {
	envPath := os.Getenv("VCR_PATH")
	vcrMode := os.Getenv("VCR_MODE")
	return envPath != "" && vcrMode != ""
}

func readSeedFromFile(fileName string) (int64, error) {
	// Max number of digits for int64 is 19
	data := make([]byte, 19)
	f, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	_, err = f.Read(data)
	if err != nil {
		return 0, err
	}
	// Remove NULL characters from seed
	data = bytes.Trim(data, "\x00")
	seed := string(data)
	return StringToFixed64(seed)
}

// Retrieves a unique test name used for writing files
// replaces all `/` characters that would cause filepath issues
// This matters during tests that dispatch multiple tests, for example TestAccLoggingFolderExclusion
func VcrSeedFile(path, name string) string {
	return filepath.Join(path, fmt.Sprintf("%s.Seed", VcrFileName(name)))
}

func VcrFileName(name string) string {
	return strings.ReplaceAll(name, "/", "_")
}

// Produces a rand.Source for VCR testing based on the given mode.
// In RECORDING mode, generates a new seed and saves it to a file, using the seed for the source
// In REPLAYING mode, reads a seed from a file and creates a source from it
func vcrSource(t *testing.T, path, mode string) (*VcrSource, error) {
	SourcesLock.RLock()
	s, ok := Sources[t.Name()]
	SourcesLock.RUnlock()
	if ok {
		return &s, nil
	}
	tflog.Debug(context.Background(), fmt.Sprintf("VCR_MODE: %s", mode))
	switch mode {
	case "RECORDING":
		seed := rand.Int63()
		s := rand.NewSource(seed)
		vcrSource := VcrSource{Seed: seed, source: s}
		SourcesLock.Lock()
		Sources[t.Name()] = vcrSource
		SourcesLock.Unlock()
		return &vcrSource, nil
	case "REPLAYING":
		seed, err := readSeedFromFile(VcrSeedFile(path, t.Name()))
		if err != nil {
			return nil, fmt.Errorf("no cassette found on disk for %s, please replay this testcase in recording mode - %w", t.Name(), err)
		}
		s := rand.NewSource(seed)
		vcrSource := VcrSource{Seed: seed, source: s}
		SourcesLock.Lock()
		Sources[t.Name()] = vcrSource
		SourcesLock.Unlock()
		return &vcrSource, nil
	default:
		log.Printf("[DEBUG] No valid environment var set for VCR_MODE, expected RECORDING or REPLAYING, skipping VCR. VCR_MODE: %s", mode)
		return nil, errors.New("No valid VCR_MODE set")
	}
}

func StringToFixed64(v string) (int64, error) {
	return strconv.ParseInt(v, 10, 64)
}
