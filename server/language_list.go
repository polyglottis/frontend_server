package server

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"sort"

	goi18n "github.com/nicksnyder/go-i18n/i18n"

	"github.com/polyglottis/platform/i18n"
	"github.com/polyglottis/platform/language"
)

const languageListFile = "language_list.json"

var languageList []language.Code

func init() {
	err := loadLanguageList()
	if err != nil {
		log.Println("Failed to load language list:", err)
	}
}

func saveLanguageList() error {
	toSave, err := json.Marshal(languageList)
	if err != nil {
		return err
	}

	f, err := os.Create(languageListFile)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(toSave)
	return err
}

func loadLanguageList() error {
	f, err := ioutil.ReadFile(languageListFile)
	if err != nil {
		return err
	}
	var ll []language.Code
	err = json.Unmarshal(f, &ll)
	if err != nil {
		return err
	}
	languageList = ll
	return nil
}

func SetLanguageList(list []language.Code) error {
	languageList = list
	return saveLanguageList()
}

func getLanguageList() ([]language.Code, error) {
	if languageList == nil {
		return nil, errors.New("Language list not set")
	}
	return languageList, nil
}

func (args *TmplArgs) GetLanguageOptions() ([]*FormOption, error) {
	return getSortedLanguageOptions(args.Localizer.Locale())
}

func (args *TmplArgs) GetLanguage(code language.Code) string {
	return args.GetText(i18n.Key("lang_" + code))
}

func getSortedLanguageOptions(locale string) ([]*FormOption, error) {
	options, err := getLanguageOptions(locale)
	if err != nil {
		return nil, err
	}
	sort.Sort(optionsByText(options))
	return options, nil
}

func getLanguageOptions(locale string) ([]*FormOption, error) {
	codes, err := getLanguageList()
	if err != nil {
		return nil, err
	}

	T, err := goi18n.Tfunc(locale)
	if err != nil {
		return nil, err
	}
	options := make([]*FormOption, len(codes)+1)
	options[0] = PleaseSelect
	for i, code := range codes {
		options[i+1] = &FormOption{
			Value: string(code),
			Text:  T("lang_" + string(code)),
		}
	}
	return options, nil
}

type optionsByText []*FormOption

func (o optionsByText) Len() int           { return len(o) }
func (o optionsByText) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o optionsByText) Less(i, j int) bool { return o[i].Text < o[j].Text }
