package pomfu

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// an internal flag memorizing whether the
// config has already been parsed or not
var hasReadConf = false

// Setup locates, starts processing pomfu's configuration files, which
// are either to be found in $XDG_CONFIG_HOME/pomfu, ~/.config/pomfu/ or
// in the users home directory
//
// Has to be manually called
func Setup() error {
	if hasReadConf {
		return nil
	}

	confDirName := os.Getenv("XDG_CONFIG_HOME")
	if confDirName == "" {
		confDirName = filepath.Join(os.Getenv("HOME"), ".config")
	}

	fi, err := os.Stat(confDirName)
	if os.IsNotExist(err) || !fi.Mode().IsDir() {
		confDirName = os.Getenv("HOME")
	}

	conf, err := os.Open(filepath.Join(confDirName, "pomfu.json"))
	if os.IsNotExist(err) {
		return nil
	}
	defer conf.Close()

	dec := json.NewDecoder(conf)
	var servers []*Pomf
	err = dec.Decode(&servers)
	if err != nil {
		return err
	}
	for _, p := range servers {
		Servers[strings.ToLower(strings.Replace(p.Name, ".", "", -1))] = p
	}

	hasReadConf = true
	return nil
}
