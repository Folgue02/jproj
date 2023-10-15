package version

import (
	"fmt"
	"strconv"
	"strings"
)


// Represents a version number.
type Version struct {
    Major int 
    Minor int
    Patch int 
}

// Returns the current version of jproj.
func GetJprojVersion() Version {
    return Version { 0, 1, 5 }
}

// Parses the string passed and creates a Version object from it.
// The format of the string must be 'x.x.x' (x being integers), 
// otherwise, this method will return an error.
func VersionFromString(rawString string) (*Version, error) {
    nums := strings.Split(rawString, ".")

    if len(nums) != 3 {
        return nil, fmt.Errorf("Wrong format for version")
    }

    var version Version

    for i, num := range nums {
        parsedNum, err := strconv.Atoi(num)

        if err != nil {
            return nil, err
        }
        switch i {
        case 0:
            version.Major = parsedNum
        case 1:
            version.Minor = parsedNum
        case 2:
            version.Patch = parsedNum
        }
    }

    return &version, nil
}

// Returns a version object.
func NewVersion(major, minor, patch int) Version {
    return Version { 
        Major: major,
        Minor: minor,
        Patch: patch,
    }
}

func (v Version) String() string {
    return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// Represents the value of the version, this is used
// for comparing version numbers.
func (v Version) getSum() int {
    return (v.Major * 100) + (v.Minor * 10) + v.Patch
}

// Checks if the current version is greater than the 
// one passed by arguments.
func (v Version) IsGreaterThan(other Version) bool {
    return v.getSum() > other.getSum()
}
