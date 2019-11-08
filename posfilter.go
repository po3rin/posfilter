package posfilter

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/msnoigrs/gosudachi"
	"github.com/msnoigrs/gosudachi/data"
)

// ModeA : split mode A is equivalent to UniDic short unit.
const ModeA = "A"

// ModeB : split mode B is equivalent to proper expression.
const ModeB = "B"

// ModeC : split mode C is intermediate unit between A and C.
const ModeC = "C"

var defaultTargetPos = map[string]struct{}{
	"名詞,普通名詞,一般":      struct{}{},
	"名詞,普通名詞,サ変可能":    struct{}{},
	"名詞,普通名詞,形状詞可能":   struct{}{},
	"名詞,普通名詞,サ変形状詞可能": struct{}{},
	"名詞,普通名詞,副詞可能":    struct{}{},
	"名詞,固有名詞,一般":      struct{}{},
	"名詞,固有名詞,人名":      struct{}{},
	"名詞,固有名詞,地名":      struct{}{},
	"名詞,固有名詞,組織名":     struct{}{},
}

func parseSettings(basePath string, settingfile string) (gosudachi.Settings, gosudachi.PluginMaker, error) {
	settings := gosudachi.NewSettingsJSON()

	var settingsreader io.Reader

	if settingfile != "" {
		var err error
		if !filepath.IsAbs(settingfile) {
			settingfile, err = filepath.Abs(settingfile)
			if err != nil {
				return nil, nil, err
			}
		}
		settingsfd, err := os.OpenFile(settingfile, os.O_RDONLY, 0644)
		if err != nil {
			return nil, nil, err
		}
		defer settingsfd.Close()
		settingsreader = settingsfd
	} else {
		settingsf, err := data.Assets.Open("sudachi.json")
		if err != nil {
			return nil, nil, err
		}
		defer settingsf.Close()
		settingsreader = settingsf
	}

	err := settings.ParseSettingsJSON(basePath, settingsreader)
	if err != nil {
		return nil, nil, err
	}

	return settings, settings, nil
}

func initDict(settingfile, resourcesdir string) (*gosudachi.JapaneseDictionary, error) {
	if resourcesdir == "" {
		ex, err := os.Executable()
		if err != nil {
			return nil, err
		}
		resourcesdir = filepath.Dir(ex)
	}

	settings, pluginmaker, err := parseSettings(resourcesdir, settingfile)
	if err != nil {
		return nil, err
	}

	settings.GetBaseConfig().Utf16String = true

	inputTextPlugins, err := pluginmaker.GetInputTextPluginArray(makeInputTextPlugin)
	if err != nil {
		return nil, err
	}
	oovProviderPlugins, err := pluginmaker.GetOovProviderPluginArray(makeOovProviderPlugin)
	if err != nil {
		return nil, err
	}
	pathRewritePlugins, err := pluginmaker.GetPathRewritePluginArray(makePathRewritePlugin)
	if err != nil {
		return nil, err
	}
	editConnectionCostPlugins, err := pluginmaker.GetEditConnectionCostPluginArray(makeEditConnectionCostPlugin)
	if err != nil {
		return nil, err
	}

	return gosudachi.NewJapaneseDictionary(
		settings.GetBaseConfig(),
		inputTextPlugins,
		oovProviderPlugins,
		pathRewritePlugins,
		editConnectionCostPlugins,
	)
}

func isTargetPos(pos []string, targetPos map[string]struct{}) bool {
	if len(pos) < 3 {
		newFeatures := []string{"*", "*", "*"}
		copy(newFeatures, pos)
		pos = newFeatures
	}
	key := fmt.Sprintf("%s,%s,%s", pos[0], pos[1], pos[2])

	_, ok := targetPos[key]
	return ok
}

func posFilter(tokenizer *gosudachi.JapaneseTokenizer, mode string, text string, targetPos map[string]struct{}) ([]string, error) {
	ms, err := tokenizer.Tokenize(mode, text)
	if err != nil {
		return nil, err
	}
	results := make([]string, 0)
	for i := 0; i < ms.Length(); i++ {
		m := ms.Get(i)

		fmt.Printf("%s\t%s\n",
			m.Surface(),
			strings.Join(m.PartOfSpeech(), ","),
		)
		if isTargetPos(m.PartOfSpeech(), targetPos) {
			results = append(results, m.Surface())
		}
	}
	return results, nil
}

// PosFilter filter part of speech with target pos map.
type PosFilter struct {
	tokenizer        *gosudachi.JapaneseTokenizer
	dict             *gosudachi.JapaneseDictionary
	targetPos        map[string]struct{}
	settingFilePath  string
	resourcesDirPath string
	mode             string
}

// OptionFunc sets settings.
type OptionFunc func(*PosFilter)

// NewPosFilter inits PosFilter.
func NewPosFilter(options ...OptionFunc) (*PosFilter, error) {

	p := &PosFilter{
		targetPos: defaultTargetPos,
		mode:      ModeC,
	}

	for _, option := range options {
		option(p)
	}

	dict, err := initDict(p.settingFilePath, p.resourcesDirPath)
	if err != nil {
		return nil, err
	}

	p.dict = dict
	p.tokenizer = dict.Create()

	return p, nil
}

// ModeOption sets split mode(A or B or C). default mode is C.
func ModeOption(mode string) OptionFunc {
	return func(p *PosFilter) {
		p.mode = mode
	}
}

// TargetPosOption sets customs target of part of speech.
func TargetPosOption(targets []string) OptionFunc {
	return func(p *PosFilter) {
		targetPos := make(map[string]struct{}, len(targets))
		for _, t := range targets {
			targetPos[t] = struct{}{}
		}
		p.targetPos = targetPos
	}
}

// SetSettingFilePath sets file path of json settings.
func SetSettingFilePath(settingFilePath string) OptionFunc {
	return func(p *PosFilter) {
		p.settingFilePath = settingFilePath
	}
}

// SetResourcesDirPath sets path of root directory of resources.
func SetResourcesDirPath(resourcesDirPath string) OptionFunc {
	return func(p *PosFilter) {
		p.resourcesDirPath = resourcesDirPath
	}
}

// Do exec tokenize & filter part of speech.
func (p *PosFilter) Do(text string) ([]string, error) {
	return posFilter(p.tokenizer, p.mode, text, p.targetPos)
}

// Close closes dictionary.
func (p *PosFilter) Close() {
	p.dict.Close()
}
