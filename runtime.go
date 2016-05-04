package template

import (
	"github.com/apaxa-io/strconvhelper"
	"io/ioutil"
	"log"
	"regexp"
)

const (
	OpeningTag = "<?"
	CloseTag   = "?>"
)

var re = regexp.MustCompile(regexp.QuoteMeta(OpeningTag) + "(?:[[:space:]]+([[:alnum:]]*))?[[:space:]]+" + regexp.QuoteMeta(CloseTag))

func ParseString(s string) (r map[string]string, o []string) {
	r = make(map[string]string)
	seps := re.FindAllStringSubmatchIndex(s, -1)

	// Check for omitted first splitter
	if len(seps) == 0 || seps[0][0] != 0 {
		seps = append([][]int{[]int{0, 0, 0, 0}}, seps...)
	}

	var nextAutoName = 0

	for i, sep := range seps {
		// Name

		name := ""
		if sep[2] >= 0 && sep[3] >= 0 { // name may be omitted in template
			name = s[sep[2]:sep[3]]
		}

		if name == "" {
			name = strconvhelper.FormatInt(nextAutoName)
			nextAutoName++
		} else if v, err := strconvhelper.ParseInt(name); err == nil {
			nextAutoName = v + 1
		}

		// Value

		var end int
		if i+1 < len(seps) {
			end = seps[i+1][0]
		} else {
			end = len(s)
		}

		value := s[sep[1]:end]

		//

		r[name] = value
		o = append(o, name)
	}

	return
}

func ParseFile(filename string) (r map[string]string, o []string, err error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, nil, err
	}
	r, o = ParseString(string(b))
	return
}

func MustParseFile(filename string) (r map[string]string, o []string) {
	r, o, err := ParseFile(filename)
	if err != nil {
		log.Panic("Unable to parse template: ", err)
	}
	return
}
