package identities

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"math/rand"
	"regexp"
	"time"

	"github.com/empijei/identigen/identities/lists"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Person represents a person object. It must be initialized by the generator.
type Person struct {
	firstName, lastName           string
	genderIsFemale                bool
	birthDate                     time.Time
	town, townCode, birthDistrict string
	residence                     string
	drv                           *DrivingLicense
	fiscalCode                    string
	partitaIva                    string
	locationCode                  int
	partitaIvaCounty              string
	cc                            *CartaCredito
	mobilePhone                   string
	id                            string
	iban                          *Iban
	up                            *Credentials
	nationality                   string
}

// NewPerson generates a new identity
func NewPerson(minage, maxage int, country string) *Person {
	person := &Person{}
	person.genderIsFemale = rand.Int()%2 == 0
	var names, surnames []string

	names, surnames, person.nationality = namesAndNation(country, person.genderIsFemale)

	person.firstName = names[rand.Int()%len(names)]
	person.lastName = surnames[rand.Int()%len(surnames)]

	var age int
	if minage == maxage {
		age = minage
	} else {
		age = rand.Int()%(maxage-minage) + minage
	}
	person.birthDate = time.Date(time.Now().Year()-age, time.Month(rand.Int()%12+1), rand.Int()%28+1, 12, 0, 0, 0, time.UTC)

	// TODO: adapt to multi-country
	birthInfo := lists.BirthInfo[rand.Int()%len(lists.BirthInfo)]
	person.town = birthInfo.Paese
	person.townCode = birthInfo.CodiceCatasto
	person.birthDistrict = birthInfo.Provincia
	person.mobilePhone = "3" + randString([]rune("1234567890"), 9)
	return person
}

// FirstName is a person first name
func (p *Person) FirstName() string {
	return p.firstName
}

// LastName is a person last name
func (p *Person) LastName() string {
	return p.lastName
}

// Gender is the gender of a person in Italian
func (p *Person) Gender() string {
	if p.genderIsFemale {
		return "Donna"
	}
	return "Uomo"
}

// BirthDate is a person birth date formatted using the globally specified format
func (p *Person) BirthDate() string {
	return p.birthDate.Format(LocalizDate.Format())
}

// BirthTown is a person birth town
func (p *Person) BirthTown() string {
	return p.town
}

// BirthDistrict is the name and label of the city the person was birth in
func (p *Person) BirthDistrict() string {
	return p.birthDistrict
}

// Phone is the phone number (Without the +39 italian prefix)
func (p *Person) Phone() string {
	return p.mobilePhone
}

// ID is the identity card number
func (p *Person) ID() string {
	if p.id != "" {
		return p.id
	}
	p.id = fmt.Sprintf("A%s%d", string("QWERTYUIOPASDFGHJKLZXCVBNM"[rand.Int()%26]), rand.Int()%10000000)
	return p.id
}

//String representation, the human readable serialization of a Person object
func (p Person) String() string {
	m := p.toMap()
	re := regexp.MustCompile("([a-z])([A-Z]+)")
	var buf bytes.Buffer
	for _, field := range fields {
		_, _ = buf.WriteString(re.ReplaceAllString(field, "$1 $2"))
		_, _ = buf.WriteString(": ")
		_, _ = buf.WriteString(m[field])
		_, _ = buf.WriteString(",\n")
	}
	return buf.String()
}

// MarshalJSON is the implementation of encoding/json.Marshaler
func (p *Person) MarshalJSON() (b []byte, err error) {
	return json.Marshal(p.toMap())
}

// MarshalXML is the implementation of encoding/xml.Marshaler
func (p *Person) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	_panic := func(e error) {
		if err != nil {
			panic(err)
		}
	}

	_panic(e.EncodeToken(start))

	for key, value := range p.toMap() {
		_panic(e.EncodeToken(xml.StartElement{Name: xml.Name{Local: key}}))
		_panic(e.EncodeToken(xml.CharData(value)))
		_panic(e.EncodeToken(xml.EndElement{Name: xml.Name{Local: key}}))
	}

	_panic(e.EncodeToken(xml.EndElement{Name: start.Name}))

	// flush to ensure tokens are written
	return e.Flush()
}

// MarshalCSV returns a []string that can be passed to an encoding/csv.Writer.Write() call
func (p Person) MarshalCSV() []string {
	m := p.toMap()
	var out []string
	for _, f := range fields {
		out = append(out, m[f])
	}
	return out
}

func (p *Person) toMap() map[string]string {
	toret := make(map[string]string)
	for _, f := range fields {
		toret[f] = printerMap[f](p)
	}
	return toret
}
