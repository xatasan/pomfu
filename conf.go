package pomfu

import (
	"bufio"
	"fmt"
	"io"
	"net/mail"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// looks to find a configuration file, as described in the Setup
// docstring
func findConf() (*os.File, error) {
	confDirName := os.Getenv("XDG_CONFIG_HOME")
	if confDirName == "" {
		confDirName = filepath.Join(os.Getenv("HOME"), ".config")
	}

	fi, err := os.Stat(confDirName)
	if os.IsNotExist(err) || !fi.Mode().IsDir() {
		confDirName = os.Getenv("HOME")
	}

	return os.Open(filepath.Join(confDirName, "pomfu"))
}

// parses a configuration file, solely based on it's input. Each record
// is separated via a line matching "[serverName]". Each following line
// until the next header has to have a key/value structure seperated
// with a "="-sign. The key has to have an alphanumeric value, value can
// contain anything, even an "="-sign.
//
// Currently, the following keys are parsed, and expect the following
// values:
//
// | Key                           | Pomf-struct attribute | Datatype |
// |-------------------------------+-----------------------+----------|
// | html, htmlallowed             | .HtmlAllowed          | Boolean  |
// | disabled, noupload, off       | .Disabled             | Boolean  |
// | max, maxsize, maximum         | .MaxSize              | Int      |
// | upload, uploadto              | .Upload               | URL      |
// | about, info, admin            | .About                | URL      |
// | webmaster, website, ownersite | .Webmaster            | URL      |
// | email, mail, contact          | .Email                | Email    |
// | owner, ownername              | .Owner                | String   |
// | name, specialname             | .Name                 | String   |
//
// If a value is malformed, and error is returned.
func parseFile(in io.Reader) (map[string]*Pomf, error) {
	var (
		current string
		server  *Pomf
		lineNr  = 1
		result  = make(map[string]*Pomf)

		newHeader = regexp.MustCompile("^\\[(\\w+)\\]$")
		newOption = regexp.MustCompile("^(\\w+)\\s*=\\s*(.*)$")

		saveRecord = func() error {
			if current != "" {
				if server.Upload == nil {
					return fmt.Errorf("Entry %s underspecified: missing 'upload' key", current)
				}
				if server.MaxSize <= 0 {
					return fmt.Errorf("Entry %s malspecified: 'maxsize' key less or equal 0", current)
				}

				result[current] = server
			}
			return nil
		}
	)

	scan := bufio.NewScanner(in)
	for scan.Scan() {
		line := strings.TrimSpace(scan.Text())

		switch {
		case len(line) == 0, strings.HasPrefix(line, "#"):
			continue
		case newHeader.MatchString(line):
			if err := saveRecord(); err != nil {
				return nil, err
			}

			current = newHeader.FindStringSubmatch(line)[1]
			server = &Pomf{Name: current} // default name to keyword
		case newOption.MatchString(line):
			option := strings.ToLower(newOption.FindStringSubmatch(line)[1])
			value := newOption.FindStringSubmatch(line)[2]

			var err error
			switch option {
			case "html", "htmlallowed":
				_, err = fmt.Sscan(value, &server.HtmlAllowed)
			case "disabled", "noupload", "off":
				_, err = fmt.Sscan(value, &server.Disabled)
			case "max", "maxsize", "maximum":
				_, err = fmt.Sscan(value, &server.MaxSize)
			case "upload", "uploadto":
				server.Upload, err = url.Parse(value)
			case "about", "info", "admin":
				server.About, err = url.Parse(value)
			case "webmaster", "website", "ownersite":
				server.Webmaster, err = url.Parse(value)
			case "email", "mail", "contact":
				server.Email, err = mail.ParseAddress(value)
			case "owner", "ownername":
				server.Owner = value
			case "name", "specialname":
				server.Name = value
			}
			if err != nil {
				e := fmt.Errorf("Error in configuration file at line %d: %s", lineNr, err)
				return nil, e
			}
		default:
			err := fmt.Errorf("Error in configuration file at line %d", lineNr)
			return nil, err
		}
		lineNr++
	}
	if err := saveRecord(); err != nil {
		return nil, err
	}

	return result, nil
}

// Setup locates, starts processing pomfu's configuration files, which
// are either to be found in $XDG_CONFIG_HOME/pomfu, ~/.config/pomfu or
// in the users home directory
//
// Has to be manually called, otherwise any configuration is ignored
func Setup() error {
	conf, err := findConf()
	if err != nil {
		return err
	}
	defer conf.Close()

	servers, err := parseFile(conf)
	if err == nil {
		for k, v := range servers {
			Servers[k] = v
		}
	}
	return err
}
