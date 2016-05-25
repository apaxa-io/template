package template

import (
	"errors"
	"github.com/apaxa-io/strconvhelper"
	"io/ioutil"
	"log"
	"regexp"
)

const (
	OpeningTag = "<?"
	CloseTag   = "?>"
)

type Options struct {
	Prefix   bool
	Optional bool
	Suffix   bool
}

func defaultOptions() (o Options) {
	o.Prefix = true
	o.Optional, o.Suffix = false, false
	return
}

var re = regexp.MustCompile(regexp.QuoteMeta(OpeningTag) + "(?:[[:space:]]+([[:alnum:]]*))?(?:[[:space:]]+(\\+noprefix))?(?:[[:space:]]+(\\+optional))?(?:[[:space:]]+(\\+suffix))?[[:space:]]+" + regexp.QuoteMeta(CloseTag))

func ParseString(s string) (strings map[string]string, order []string, options map[string]Options, err error) {
	strings = make(map[string]string)
	options = make(map[string]Options)
	seps := re.FindAllStringSubmatchIndex(s, -1)

	// Check for omitted first splitter
	if len(seps) == 0 || seps[0][0] != 0 {
		seps = append([][]int{[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}, seps...)
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

		if _, ok := strings[name]; ok {
			err = errors.New("Duplicate template name detected.")
		}

		// options
		//options := defaultOptions()
		var opt Options
		opt.Prefix = !(sep[4] >= 0 && sep[5] >= 0)
		opt.Optional = sep[6] >= 0 && sep[7] >= 0
		opt.Suffix = sep[8] >= 0 && sep[9] >= 0

		// Value

		var end int
		if i+1 < len(seps) {
			end = seps[i+1][0]
		} else {
			end = len(s)
		}

		value := s[sep[1]:end]

		//

		strings[name] = value
		order = append(order, name)
		options[name] = opt
	}

	return
}

func ParseFile(filename string) (strings map[string]string, order []string, options map[string]Options, err error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, nil, nil, err
	}
	return ParseString(string(b))
}

func MustParseFile(filename string) (strings map[string]string, order []string, options map[string]Options) {
	strings, order, options, err := ParseFile(filename)
	if err != nil {
		log.Panic("Unable to parse template: ", err) // TODO change log.Panic to smth other?
	}
	return
}

/*
TODO implement this function
func CompileSimple(w http.ResponseWriter, filename string) error {
	t, o, err := template.ParseFile(filename)
	if err != nil {
		log.Print(err)
		return err
	}

	for _, name := range o {
		if _, err := w.Write([]byte(t[name])); err != nil {
			log.Print(err)
			return err
		}
	}
	return nil
}
*/
