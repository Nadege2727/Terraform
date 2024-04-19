package configs

import "fmt"

type WorkspaceDeprecations struct {
	ModuleDeprecationInfos []*ModuleDeprecationInfo
}

type ModuleDeprecationInfo struct {
	SourceName           string
	RegistryDeprecation  *RegistryModuleDeprecation
	ExternalDependencies []*ModuleDeprecationInfo
}

type RegistryModuleDeprecation struct {
	Version string
	Link    string
}

func (i *WorkspaceDeprecations) BuildDeprecationWarningString() string {
	modDeprecationStrings := []string{}
	for _, modDeprecationInfo := range i.ModuleDeprecationInfos {
		if modDeprecationInfo != nil && modDeprecationInfo.RegistryDeprecation != nil {
			// Link is an optional field, if unset it is an empty string by default
			if modDeprecationInfo.RegistryDeprecation.Link != "" {
				modDeprecationStrings = append(modDeprecationStrings, fmt.Sprintf("Version %s of \"%s\" \nTo learn more visit: %s\n", modDeprecationInfo.RegistryDeprecation.Version, modDeprecationInfo.SourceName, modDeprecationInfo.RegistryDeprecation.Link))
			} else {
				modDeprecationStrings = append(modDeprecationStrings, fmt.Sprintf("Version %s of \"%s\" \n", modDeprecationInfo.RegistryDeprecation.Version, modDeprecationInfo.SourceName))
			}
		}
		modDeprecationStrings = append(modDeprecationStrings, buildChildDeprecationWarnings(modDeprecationInfo.ExternalDependencies, []string{modDeprecationInfo.SourceName})...)
	}
	deprecationsMessage := ""
	for _, deprecationString := range modDeprecationStrings {
		deprecationsMessage += deprecationString + "\n"
	}

	return deprecationsMessage
}

func buildChildDeprecationWarnings(modDeprecations []*ModuleDeprecationInfo, parentMods []string) []string {
	modDeprecationStrings := []string{}
	for _, deprecation := range modDeprecations {
		if deprecation.RegistryDeprecation != nil {
			modDeprecationStrings = append(modDeprecationStrings, fmt.Sprintf("Version %s of \"%s\" %s \nTo learn more visit: %s\n", deprecation.RegistryDeprecation.Version, deprecation.SourceName, buildModHierarchy(parentMods, deprecation.SourceName), deprecation.RegistryDeprecation.Link))
		}
		newParentMods := append(parentMods, deprecation.SourceName)
		modDeprecationStrings = append(modDeprecationStrings, buildChildDeprecationWarnings(deprecation.ExternalDependencies, newParentMods)...)
	}
	return modDeprecationStrings
}

func buildModHierarchy(parentMods []string, modName string) string {
	heirarchy := ""
	for _, parent := range parentMods {
		heirarchy += fmt.Sprintf("%s -> ", parent)
	}
	heirarchy += modName
	return fmt.Sprintf("(Root: %s)", heirarchy)
}
