package internal

import (
	"fmt"
	"path/filepath"
	"plugin"
	"reflect"
	"strings"
	"time"
)

// Constants for reflection method names
const (
	methodName = "Name"
	methodRun  = "Run"
	symbolName = "PluginInstance"
)

// loadedPlugin wraps the plugin instance and its result container
type loadedPlugin struct {
	Instance reflect.Value
	Result   PluginResult
}

// UsePlugin loads and executes plugins from the specified directory
func UsePlugin(pluginPaths string, dataBytes []byte) []PluginResult {
	var pluginResults []PluginResult

	if pluginPaths == "" {
		return pluginResults
	}

	plugins := loadPlugins(pluginPaths)

	for _, p := range plugins {
		// If plugin loading failed, report it and skip execution
		if p.Result.LoadError != "" {
			pluginResults = append(pluginResults, p.Result)
			continue
		}

		start := time.Now()
		msgs, warns, errs := callRun(p.Instance, dataBytes)
		elapsed := time.Since(start)

		p.Result.Messages = msgs
		p.Result.Warnings = warns
		p.Result.Errors = errs
		p.Result.ExecutionTime = elapsed

		pluginResults = append(pluginResults, p.Result)
	}

	return pluginResults
}

// loadPlugins reads a directory and loads all `.so` plugins
func loadPlugins(pluginList string) []loadedPlugin {
	var plugins []loadedPlugin

	// Support both comma-separated and newline-separated input
	rawPaths := strings.FieldsFunc(pluginList, func(r rune) bool {
		return r == ',' || r == '\n'
	})

	for _, rawPath := range rawPaths {
		path := strings.TrimSpace(rawPath)
		if path == "" || filepath.Ext(path) != ".so" {
			continue
		}
		plugin := loadSinglePlugin(path, filepath.Base(path))
		plugins = append(plugins, plugin)
	}

	return plugins
}

// loadSinglePlugin loads and validates a single plugin file
func loadSinglePlugin(path, name string) loadedPlugin {
	result := PluginResult{Name: name}

	p, err := plugin.Open(path)
	if err != nil {
		return pluginErrorResult(name, fmt.Sprintf("Failed to open: %v", err))
	}

	sym, err := p.Lookup(symbolName)
	if err != nil {
		return pluginErrorResult(name, fmt.Sprintf("PluginInstance not found: %v", err))
	}

	val := reflect.ValueOf(sym).Elem()
	if !hasMethods(val, methodName, methodRun) {
		return pluginErrorResult(name, fmt.Sprintf("Missing %s() or %s()", methodName, methodRun))
	}

	result.Name = callName(val)

	return loadedPlugin{
		Instance: val,
		Result:   result,
	}
}

// callName uses reflection to call the Name() method on a plugin
func callName(v reflect.Value) string {
	if m := v.MethodByName(methodName); m.IsValid() {
		if result := m.Call(nil); len(result) > 0 {
			if name, ok := result[0].Interface().(string); ok {
				return name
			}
		}
	}
	return "UnknownPlugin"
}

// callRun safely executes the Run method and returns messages, warnings, and errors
func callRun(v reflect.Value, data []byte) (msgs, warns, errs []string) {
	defer func() {
		if r := recover(); r != nil {
			errs = []string{fmt.Sprintf("Plugin panicked: %v", r)}
		}
	}()

	m := v.MethodByName(methodRun)
	if !m.IsValid() {
		return nil, nil, []string{"Run method not found"}
	}

	results := m.Call([]reflect.Value{
		reflect.ValueOf(data),
	})

	if len(results) != 3 {
		return nil, nil, []string{fmt.Sprintf("Plugin returned %d values, expected 3", len(results))}
	}

	// Safely extract results
	if m := results[0].Interface(); m != nil {
		msgs, _ = m.([]string)
	}
	if w := results[1].Interface(); w != nil {
		warns, _ = w.([]string)
	}
	if e := results[2].Interface(); e != nil {
		errs, _ = e.([]string)
	}

	return
}

// pluginErrorResult creates a failed plugin result
func pluginErrorResult(name, errMsg string) loadedPlugin {
	return loadedPlugin{
		Result: PluginResult{
			Name:      name,
			LoadError: errMsg,
		},
	}
}

// hasMethods checks whether all specified methods exist in the plugin
func hasMethods(val reflect.Value, methods ...string) bool {
	for _, m := range methods {
		if !val.MethodByName(m).IsValid() {
			return false
		}
	}
	return true
}
