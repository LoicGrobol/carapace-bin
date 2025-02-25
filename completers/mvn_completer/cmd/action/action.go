package action

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace-bin/pkg/util"
)

type Plugin struct {
	XMLName    xml.Name `xml:"plugin"`
	GoalPrefix string   `xml:"goalPrefix"`
	Mojos      []struct {
		Goal        string `xml:"goal"`
		Description string `xml:"description"`
	} `xml:"mojos>mojo"`
}

func (p Plugin) FormattedGoals() map[string]string {
	goals := make(map[string]string)
	for _, mojo := range p.Mojos {
		goal := fmt.Sprintf("%v:%v", p.GoalPrefix, mojo.Goal)
		description := strings.SplitAfter(mojo.Description, ".")[0]
		if len(description) > 60 {
			description = description[:57] + "..."
		}
		goals[goal] = description
	}
	return goals
}

type Artifact struct {
	GroupId    string `xml:"groupId"`
	ArtifactId string `xml:"artifactId"`
	Version    string `xml:"version"`
}

func (a Artifact) Location(repository string) string {
	return fmt.Sprintf("%v/%v/%v/%v/%v-%v.jar", repository, strings.Replace(a.GroupId, ".", "/", -1), a.ArtifactId, a.Version, a.ArtifactId, a.Version)
}

type Project struct {
	// TODO parent pom plugins
	// TODO plugins locatad in pluginmanagement and profiles
	XMLName  xml.Name   `xml:"project"`
	Plugins  []Artifact `xml:"build>plugins>plugin"`
	Profiles []string   `xml:"profiles>profile>id"`
}

func repositoryLocation() string {
	// TODO environment variable / settings override
	if repoLocation, err := homedir.Expand("~/.m2/repository"); err == nil {
		return repoLocation
	}
	return "" // TODO handle error
}

func ActionGoalsAndPhases(file string) carapace.Action {
	// TODO caching
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		if project, err := loadProject(file); err != nil {
			return carapace.ActionMessage(err.Error())
		} else {
			goals := make(map[string]string)
			for _, plugin := range project.Plugins {
				if plugin := loadPlugin(plugin.Location(repositoryLocation())); plugin != nil {
					for key, value := range plugin.FormattedGoals() {
						goals[key] = value
					}
				}
			}

			for key, value := range defaultGoalsAndPhases() {
				goals[key] = value
			}

			vals := make([]string, 0)
			for key, value := range goals {
				vals = append(vals, key, value)
			}

			return carapace.ActionValuesDescribed(vals...)
		}
	})
}

func defaultGoalsAndPhases() (goals map[string]string) {
	goals = map[string]string{
		// clean lifecycle
		"pre-clean":  "execute processes needed prior to the actual project cleaning",
		"clean":      "remove all files generated by the previous build",
		"post-clean": "execute processes needed to finalize the project cleanin",

		// default lifecycle
		"validate":                "validate the project is correct and all necessary information is available.",
		"initialize":              "initialize build state, e.g. set properties or create directories.",
		"generate-sources":        "generate any source code for inclusion in compilation.",
		"process-sources":         "process the source code, for example to filter any values.",
		"generate-resources":      "generate resources for inclusion in the package.",
		"process-resources":       "copy and process the resources into the destination directory, ready for packaging.",
		"compile":                 "compile the source code of the project.",
		"process-classes":         "post-process the generated files from compilation.",
		"generate-test-sources":   "generate any test source code for inclusion in compilation.",
		"process-test-sources":    "process the test source code, for example to filter any values.",
		"generate-test-resources": "create resources for testing.",
		"process-test-resources":  "copy and process the resources into the test destination directory.",
		"test-compile":            "compile the test source code into the test destination directory",
		"process-test-classes":    "post-process the generated files from test compilation.",
		"test":                    "run tests using a suitable unit testing framework.",
		"prepare-package":         "perform any operations necessary to prepare a package before the actual packaging.",
		"package":                 "take the compiled code and package it in its distributable format, such as a JAR.",
		"pre-integration-test":    "perform actions required before integration tests are executed.",
		"integration-test":        "process and deploy the package into an environment where integration tests can be run.",
		"post-integration-test":   "perform actions required after integration tests have been executed.",
		"verify":                  "run any checks to verify the package is valid and meets quality criteria.",
		"install":                 "install the package into the local repository.",
		"deploy":                  "copies the final package to the remote repository.",

		// site lifecycle
		"pre-site":    "execute processes needed prior to the actual project site generation",
		"site":        "generate the project's site documentation",
		"post-site":   "execute processes needed to finalize the site generation",
		"site-deploy": "deploy the generated site documentation to the specified web server",
	}

	// TODO this traverses over all versions of a plugin while only the latest should be used
	_ = filepath.Walk(repositoryLocation()+"/org/apache/maven/plugins/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".jar") {
			if plugin := loadPlugin(path); plugin != nil {
				for key, value := range plugin.FormattedGoals() {
					goals[key] = value
				}
			}
		}
		return nil
	})
	return
}

func ActionProfiles(file string) carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		if project, err := loadProject(file); err != nil {
			return carapace.ActionMessage(err.Error())
		} else {
			return carapace.ActionValues(project.Profiles...)
		}
	})
}

func locatePom(file string) (pom string) {
	if file != "" {
		return file
	}
	pom, _ = util.FindReverse("", "pom.xml")
	return
}

func loadProject(file string) (project *Project, err error) {
	var content []byte
	if content, err = os.ReadFile(locatePom(file)); err == nil {
		err = xml.Unmarshal(content, &project)
	}
	return
}

func loadPlugin(file string) (plugin *Plugin) {
	if reader, err := zip.OpenReader(file); err == nil {
		defer reader.Close()
		for _, f := range reader.File {
			if f.Name == "META-INF/maven/plugin.xml" {
				if pluginFile, err := f.Open(); err == nil {
					defer pluginFile.Close()
					if content, err := io.ReadAll(pluginFile); err == nil {
						_ = xml.Unmarshal(content, &plugin)
					}
				}
			}
		}
	}
	return
}
