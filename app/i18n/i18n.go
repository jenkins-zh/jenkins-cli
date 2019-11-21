package i18n

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/chai2010/gettext-go/gettext"
)

var knownTranslations = map[string][]string{
	"jcli": {
		"default",
		"en_US",
		"zh_CN",
	},
	// only used for unit tests.
	"test": {
		"default",
		"en_US",
	},
}

func loadSystemLanguage() string {
	// Implements the following locale priority order: LC_ALL, LC_MESSAGES, LANG
	// Similarly to: https://www.gnu.org/software/gettext/manual/html_node/Locale-Environment-Variables.html
	langStr := os.Getenv("LC_ALL")
	if langStr == "" {
		langStr = os.Getenv("LC_MESSAGES")
	}
	if langStr == "" {
		langStr = os.Getenv("LANG")
	}

	if langStr == "" {
		//klog.V(3).Infof("Couldn't find the LC_ALL, LC_MESSAGES or LANG environment variables, defaulting to en_US")
		return "default"
	}
	pieces := strings.Split(langStr, ".")
	if len(pieces) != 2 {
		//klog.V(3).Infof("Unexpected system language (%s), defaulting to en_US", langStr)
		return "default"
	}
	return pieces[0]
}

func findLanguage(root string, getLanguageFn func() string) string {
	langStr := getLanguageFn()

	translations := knownTranslations[root]
	for ix := range translations {
		if translations[ix] == langStr {
			return langStr
		}
	}
	//klog.V(3).Infof("Couldn't find translations for %s, using default", langStr)
	return "default"
}

// LoadTranslations loads translation files. getLanguageFn should return a language
// string (e.g. 'en-US'). If getLanguageFn is nil, then the loadSystemLanguage function
// is used, which uses the 'LANG' environment variable.
func LoadTranslations(root string, getLanguageFn func() string) error {
	if getLanguageFn == nil {
		getLanguageFn = loadSystemLanguage
	}

	langStr := findLanguage(root, getLanguageFn)
	translationFiles := []string{
		//"jcli/zh_CN/LC_MESSAGES/jcli.mo",
		"jcli/zh_CN/LC_MESSAGES/jcli.po",
	}

	//klog.V(3).Infof("Setting language to %s", langStr)
	// TODO: list the directory and load all files.
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	// Make sure to check the error on Close.
	for _, file := range translationFiles {
		filename := file
		f, err := w.Create(file)
		if err != nil {
			return err
		}
		data, err := Asset(filename)
		if err != nil {
			return err
		}
		if _, err := f.Write(data); err != nil {
			return nil
		}
	}
	if err := w.Close(); err != nil {
		return err
	}
	gettext.BindTextdomain("jcli", root+".zip", buf.Bytes())
	gettext.Textdomain("jcli")
	gettext.SetLocale(langStr)
	return nil
}

var i18nLoaded = false

// T translates a string, possibly substituting arguments into it along
// the way. If len(args) is > 0, args1 is assumed to be the plural value
// and plural translation is used.
func T(defaultValue string, args ...int) string {
	if !i18nLoaded {
		i18nLoaded = true
		if err := LoadTranslations("jcli", nil); err != nil {
			fmt.Println(err)
		}
	}

	if len(args) == 0 {
		return gettext.PGettext("", defaultValue)
	}
	return fmt.Sprintf(gettext.PNGettext("", defaultValue, defaultValue+".plural", args[0]),
		args[0])
}

// Errorf produces an error with a translated error string.
// Substitution is performed via the `T` function above, following
// the same rules.
func Errorf(defaultValue string, args ...int) error {
	return errors.New(T(defaultValue, args...))
}
