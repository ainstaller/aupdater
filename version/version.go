package version

import (
	"regexp"
	"strconv"
)

var re = regexp.MustCompile(`\"v(\d+)\.(\d+)\"`) // "v1234.5678"

// Version from ahud
type Version struct {
	Month int
	Year  int
	Day   int
}

// New returns a new version already parsing the string in the args
func New(str string) (*Version, error) {
	v := new(Version)
	err := v.Parse(str)
	return v, err
}

// Parse gets year, month and day from a string (mainmenuoverride.res)
func (v *Version) Parse(str string) error {
	m := re.FindAllStringSubmatch(str, -1)
	if len(m) == 0 {
		return nil
	}

	// loop through matches
	for i, match := range m[0] {
		var err error
		// 1 == year
		if i == 1 {
			// year
			v.Year, err = strconv.Atoi(match)

			// 2 == month and day
		} else if i == 2 {
			// month
			v.Month, err = strconv.Atoi(match[0:2])
			if err != nil {
				return err
			}

			// day
			v.Day, err = strconv.Atoi(match[2:4])
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// After checks if current version is before the given version
func (v *Version) After(v2 *Version) bool {
	return v.Year > v2.Year || v.Month > v2.Month || v.Day > v2.Day
}

// Before checks if current version is before the given version
func (v *Version) Before(v2 *Version) bool {
	return v.Year < v2.Year || v.Month < v2.Month || v.Day < v2.Day
}

// Equal checks if year, month and day are equal
func (v *Version) Equal(v2 *Version) bool {
	return v.Year == v2.Year && v.Month == v2.Month && v.Day == v2.Day
}
