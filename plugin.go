package posfilter

import (
	"reflect"

	"github.com/msnoigrs/gosudachi"
)

var (
	plugins = make(map[string]reflect.Type)

	pkgstr                            = "github.com/msnoigrs/gosudachi"
	defaultInputTextPlugin            = pkgstr + ".DefaultInputTextPlugin"
	prolongedSoundMarkInputTextPlugin = pkgstr + ".ProlongedSoundMarkInputTextPlugin"
	inhibitConnectionPlugin           = pkgstr + ".InhibitConnectionPlugin"
	meCabOovProviderPlugin            = pkgstr + ".MeCabOovProviderPlugin"
	simpleOovProviderPlugin           = pkgstr + ".SimpleOovProviderPlugin"
	joinNumericPlugin                 = pkgstr + ".JoinNumericPlugin"
	joinKatakanaOovPlugin             = pkgstr + ".JoinKatakanaOovPlugin"
)

func init() {
	register(gosudachi.DefaultInputTextPlugin{})
	register(gosudachi.ProlongedSoundMarkInputTextPlugin{})
	register(gosudachi.MeCabOovProviderPlugin{})
	register(gosudachi.SimpleOovProviderPlugin{})
	register(gosudachi.JoinNumericPlugin{})
	register(gosudachi.JoinKatakanaOovPlugin{})
	register(gosudachi.InhibitConnectionPlugin{})
}

func register(x interface{}) {
	t := reflect.TypeOf(x)
	n := t.PkgPath() + "." + t.Name()
	plugins[n] = t
}

func newPlugin(name string) (interface{}, bool) {
	t, ok := plugins[name]
	if !ok {
		return nil, false
	}
	v := reflect.New(t)
	return v.Interface(), true
}

func makeInputTextPlugin(k string) gosudachi.InputTextPlugin {
	switch k {
	case "DefaultInputTextPlugin", "com.worksap.nlp.sudachi.DefaultInputTextPlugin", defaultInputTextPlugin:
		plugin, ok := newPlugin(defaultInputTextPlugin)
		if !ok {
			return nil
		}
		rplugin, ok := plugin.(gosudachi.InputTextPlugin)
		if !ok {
			return nil
		}
		return rplugin
	case "ProlongedSoundMarkInputTextPlugin", "com.worksap.nlp.sudachi.ProlongedSoundMarkInputTextPlugin", prolongedSoundMarkInputTextPlugin:
		plugin, ok := newPlugin(prolongedSoundMarkInputTextPlugin)
		if !ok {
			return nil
		}
		rplugin, ok := plugin.(gosudachi.InputTextPlugin)
		if !ok {
			return nil
		}
		return rplugin
	}
	return nil
}

func makeOovProviderPlugin(k string) gosudachi.OovProviderPlugin {
	switch k {
	case "MeCabOovProviderPlugin", "com.worksap.nlp.sudachi.MeCabOovProviderPlugin", meCabOovProviderPlugin:
		plugin, ok := newPlugin(meCabOovProviderPlugin)
		if !ok {
			return nil
		}
		rplugin, ok := plugin.(gosudachi.OovProviderPlugin)
		if !ok {
			return nil
		}
		return rplugin
	case "SimpleOovProviderPlugin", "com.worksap.nlp.sudachi.SimpleOovProviderPlugin", simpleOovProviderPlugin:
		plugin, ok := newPlugin(simpleOovProviderPlugin)
		if !ok {
			return nil
		}
		rplugin, ok := plugin.(gosudachi.OovProviderPlugin)
		if !ok {
			return nil
		}
		return rplugin
	}
	return nil
}

func makePathRewritePlugin(k string) gosudachi.PathRewritePlugin {
	switch k {
	case "JoinNumericPlugin", "com.worksap.nlp.sudachi.JoinNumericPlugin", joinNumericPlugin:
		plugin, ok := newPlugin(joinNumericPlugin)
		if !ok {
			return nil
		}
		rplugin, ok := plugin.(gosudachi.PathRewritePlugin)
		if !ok {
			return nil
		}
		return rplugin
	case "JoinKatakanaOovPlugin", "com.worksap.nlp.sudachi.JoinKatakanaOovPlugin", joinKatakanaOovPlugin:
		plugin, ok := newPlugin(joinKatakanaOovPlugin)
		if !ok {
			return nil
		}
		rplugin, ok := plugin.(gosudachi.PathRewritePlugin)
		if !ok {
			return nil
		}
		return rplugin
	}
	return nil
}

func makeEditConnectionCostPlugin(k string) gosudachi.EditConnectionCostPlugin {
	switch k {
	case "InhibitConnectionPlugin", "com.worksap.nlp.sudachi.InhibitConnectionPlugin", inhibitConnectionPlugin:
		plugin, ok := newPlugin(inhibitConnectionPlugin)
		if !ok {
			return nil
		}
		rplugin, ok := plugin.(gosudachi.EditConnectionCostPlugin)
		if !ok {
			return nil
		}
		return rplugin
	}
	return nil
}
