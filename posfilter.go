package posfilter

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/msnoigrs/gosudachi"
	"github.com/msnoigrs/gosudachi/data"
)

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

func initDict() (*gosudachi.JapaneseDictionary, error) {
	var (
		settingfile  string
		resourcesdir string
	)

	flag.StringVar(&settingfile, "r", "", "read settings from file (overrides -s)")
	flag.StringVar(&resourcesdir, "p", "", "root directory of resources")
	flag.Parse()

	if resourcesdir == "" {
		ex, err := os.Executable()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		resourcesdir = filepath.Dir(ex)
	}

	settings, pluginmaker, err := parseSettings(resourcesdir, settingfile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fail to parse settings: %s\n", err)
		os.Exit(1)
	}

	settings.GetBaseConfig().Utf16String = true

	inputTextPlugins, err := pluginmaker.GetInputTextPluginArray(makeInputTextPlugin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fail to cleate any InputTextPlugin: %s\n", err)
		os.Exit(1)
	}
	oovProviderPlugins, err := pluginmaker.GetOovProviderPluginArray(makeOovProviderPlugin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fail to cleate any OovProviderPlugin: %s\n", err)
		os.Exit(1)
	}
	pathRewritePlugins, err := pluginmaker.GetPathRewritePluginArray(makePathRewritePlugin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fail to cleate any PathRewritePlugin: %s\n", err)
		os.Exit(1)
	}
	editConnectionCostPlugins, err := pluginmaker.GetEditConnectionCostPluginArray(makeEditConnectionCostPlugin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fail to cleate any ConnectionCostPlugin: %s\n", err)
		os.Exit(1)
	}

	return gosudachi.NewJapaneseDictionary(
		settings.GetBaseConfig(),
		inputTextPlugins,
		oovProviderPlugins,
		pathRewritePlugins,
		editConnectionCostPlugins,
	)
}

func posFilter(tokenizer *gosudachi.JapaneseTokenizer, mode string, text string) ([]string, error) {
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

		results = append(results, m.Surface())
	}
	return results, nil
}

// PosFilter filter part of speech with target pos map.
type PosFilter struct {
	tokenizer *gosudachi.JapaneseTokenizer
	targetPos map[string]struct{}
	mode      string
}

// ModeA : split mode A is equivalent to UniDic short unit.
const ModeA = "A"

// ModeB : split mode B is equivalent to proper expression.
const ModeB = "B"

// ModeC : split mode C is intermediate unit between A and C.
const ModeC = "C"

// SetMode sets split mode(A or B or C). default mode is C.
func (p *PosFilter) SetMode(mode string) *PosFilter {
	p.mode = mode
	return p
}

// SetTargetPos sets customs target of part of speech.
func (p *PosFilter) SetTargetPos(targets []string) *PosFilter {
	targetPos := make(map[string]struct{}, len(targets))
	for _, t := range targets {
		targetPos[t] = struct{}{}
	}
	p.targetPos = targetPos
	return p
}

// Do exec tokenize & filter part of speech.
func (p *PosFilter) Do(text string) ([]string, error) {
	if p.tokenizer == nil {
		dict, err := initDict()
		if err != nil {
			return nil, err
		}
		defer dict.Close()
		p.tokenizer = dict.Create()
	}
	if p.mode == "" {
		p.mode = ModeC
	}
	if p.targetPos == nil {
		p.targetPos = defaultTargetPos
	}
	return posFilter(p.tokenizer, p.mode, text)
}
