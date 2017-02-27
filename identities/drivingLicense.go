package identities

import (
	"strings"
)

//Returns a Italian driving license number
func (p *Person) DrivingLicense() string {
	if p.drvLicense != "" {
		return p.drvLicense
	}
	p.drvLicense = strings.Split(p.birthDistrict, "(")[1][:2] +
		randString([]rune("0123456789"), 7) +
		randString([]rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ"), 1)
	return p.drvLicense
}