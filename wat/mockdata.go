package wat

func MockModelPayload(inputSource string, plugin Plugin) ModelPayload {
	mconfig := ModelConfiguration{}
	switch plugin.Name {
	case "fragilitycurveplugin":
		paths := make([]string, 1)
		paths[0] = inputSource + "fc.json"
		mconfig.Name = "levee_failures"
		mconfig.ModelConfigurationPaths = paths
	case "hydrographscaler":
		paths := make([]string, 1)
		paths[0] = inputSource + "hsm.json"
		mconfig.Name = "hydrographs"
		mconfig.ModelConfigurationPaths = paths
	}
	payload := ModelPayload{
		TargetPlugin:       plugin.Name,
		PluginImageAndTag:  plugin.ImageAndTag,
		ModelConfiguration: mconfig,
	}
	return payload
}
