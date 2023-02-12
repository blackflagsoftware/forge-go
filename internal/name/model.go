package name

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"unicode"

	"github.com/Masterminds/sprig"
)

type (
	Name struct {
		RawName    string `json:"-"` // name given by the user/script
		Lower      string `json:"-"`
		Camel      string `json:"-"`
		LowerCamel string `json:"-"`
		Abbr       string `json:"-"`
		AllLower   string `json:"-"`
		Upper      string `json:"-"`
		EnvVar     string `json:"-"`
	}

	Format struct {
		Name string
	}
)

func (n *Name) BuildName(name string, knownAliases []string) string {
	n.RawName = name
	runeName := []rune(n.RawName)
	idxFirst := 0
	idxLast := len(runeName)
	for i := range runeName {
		if !unicode.IsLetter(runeName[i]) {
			idxFirst = i + 1
			continue
		}
		break
	}
	for i := len(runeName) - 1; i > 0; i-- {
		if !unicode.IsLetter(runeName[i]) {
			idxLast = i
			continue
		}
		break
	}

	rawName := string(runeName[idxFirst:idxLast])

	n.Camel = BuildAltName(rawName, "pascalCase")
	n.LowerCamel = BuildAltName(rawName, "camelCase")
	n.Lower = BuildAltName(rawName, "lowerCase")
	n.AllLower = strings.ToLower(n.Camel)
	n.Upper = strings.ToUpper(n.Camel)
	n.EnvVar = strings.ToUpper(rawName)
	return n.DetermineAbbr(knownAliases)
}

func BuildAltName(name, mode string) string {
	f := Format{Name: name}
	var t *template.Template
	var err error
	switch mode {
	case "snakeCase":
		t, err = template.New("format").Funcs(sprig.GenericFuncMap()).Parse("{{.Name | snakecase}}")
	case "kebabCase":
		t, err = template.New("format").Funcs(sprig.GenericFuncMap()).Parse("{{.Name | kebabcase}}")
	case "camelCase":
		t, err = template.New("format").Funcs(sprig.GenericFuncMap()).Parse("{{.Name | camelcase}}")
	case "pascalCase":
		t, err = template.New("format").Funcs(sprig.GenericFuncMap()).Parse("{{.Name | camelcase}}")
	case "upperCase":
		t, err = template.New("format").Funcs(sprig.GenericFuncMap()).Parse("{{.Name | upper}}")
	default:
		// lowerCase
		t, err = template.New("format").Funcs(sprig.GenericFuncMap()).Parse("{{.Name | lower}}")
	}
	if err != nil {
		fmt.Println(err)
		return name
	}
	b := bytes.NewBufferString("")
	errE := t.Execute(b, f)
	if errE != nil {
		fmt.Println(errE)
		return name
	}
	if mode == "camelCase" {
		// finish off the camel case functionality
		n := []rune(b.String())
		return string(append([]rune{unicode.ToLower(n[0])}, n[1:]...))
	}
	return b.String()
}

func (n *Name) DetermineAbbr(knownAliases []string) string {
	// check all known aliases and create the best Name.Abbr
	if len(n.Lower) < 3 {
		n.Abbr = n.Lower
		return n.Abbr
	}
	try := 0
Loop:
	for {
		try++
		var retry bool
		n.Abbr, retry = nameVersion(try, n.Lower)
		if !retry {
			break
		}
		for _, ka := range knownAliases {
			if ka == n.Abbr {
				continue Loop
			}
		}
		break
	}
	return n.Abbr
}

func nameVersion(tried int, name string) (abbr string, retry bool) {
	switch tried {
	case 1:
		// first 3 character
		abbr = string(name[:3])
		retry = true
	case 2:
		// first 2 characters
		abbr = string(name[:2])
		retry = true
	case 3, 4:
		// first letters of each snake case
		charNumber := tried - 2
		split := strings.Split(name, "_")
		if len(split) == 1 {
			abbr = name
			return
		}
		for _, c := range split {
			abbr += string(c[:charNumber])
		}
		retry = true
	default:
		abbr = name
		retry = false
	}
	return
}
