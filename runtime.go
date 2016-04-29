package template

import (
	"github.com/apaxa-io/strconvhelper"
	"io/ioutil"
	"regexp"
)

const (
	OpeningTag = "<?"
	CloseTag   = "?>"
)

var re = regexp.MustCompile(regexp.QuoteMeta(OpeningTag) + "(?:[[:space:]]+([[:alnum:]]*))?[[:space:]]+" + regexp.QuoteMeta(CloseTag))

func ParseString(s string) (r map[string]string) {
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
	}

	return
}

func ParseFile(fileName string) (map[string]string, error) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return ParseString(string(b)), nil
}
