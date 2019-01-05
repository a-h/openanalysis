package read

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/a-h/openanalysis/statistics"
)

func UserStatistics(dir string) (stats map[string]*statistics.Statistics, err error) {
	stats = make(map[string]*statistics.Statistics)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}
	for _, f := range files {
		if path.Ext(f.Name()) != ".json" {
			continue
		}
		b, rerr := ioutil.ReadFile(dir + "/" + f.Name())
		if rerr != nil {
			err = fmt.Errorf("failed to read: %v: %v", f.Name(), rerr)
			return
		}
		var s *statistics.Statistics
		uerr := json.Unmarshal(b, &s)
		if uerr != nil {
			err = fmt.Errorf("failed to unmarshal: %v: %v", f, uerr)
			return
		}
		_, n := filepath.Split(f.Name())
		n = strings.TrimSuffix(n, path.Ext(n))
		stats[n] = s
	}
	return
}
