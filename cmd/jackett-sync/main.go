package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
	"regexp"
	"text/template"

	"github.com/google/go-github/v45/github"
	"github.com/iancoleman/strcase"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v3"
)

// This file is used to generate trackers from https://github.com/Jackett/Jackett/tree/master/src/Jackett.Common/Definitions.
// https://stackoverflow.com/a/61656698

const (
	DefinitionsFolderName string = "src/Jackett.Common/Definitions"
	DefinitionsBranchName string = "master"
)

// Trackers describes available trackers on Jackett.
var Trackers []*TrackerDefinition

var (
	packageName string
	output      string
	githubToken string
	debug       bool

	re = regexp.MustCompile(`\d`)
)

func init() {
	flag.StringVar(&packageName, "package", "jackett", "package for generated file")
	flag.StringVar(&output, "output", "trackers.go", "generated output file path")
	flag.StringVar(&githubToken, "token", os.Getenv("GH_TOKEN"), "GitHub token to use")
	flag.BoolVar(&debug, "debug", false, "enable this to stop search trackers after 10")
}

func main() {
	flag.Parse()

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)

	// get the entire tree of the repo
	tree, _, err := client.Git.GetTree(context.Background(), "Jackett", "Jackett", DefinitionsBranchName, true)
	if err != nil {
		log.Fatalln(err)
	}

	definitionsTree := getDefinitionsEntry(tree.Entries)
	if definitionsTree == nil {
		log.Fatalf("cannot find folder [%s]. Aborting.", DefinitionsFolderName)
	}

	tree, _, err = client.Git.GetTree(context.Background(), "Jackett", "Jackett", definitionsTree.GetSHA(), false)
	if err != nil {
		log.Fatalln(err)
	}

	for i, entry := range tree.Entries {
		if err := processTracker(context.Background(), client, entry); err != nil {
			log.Printf("[ERROR] cannot process entry at path [%s]: %s", entry.GetPath(), err.Error())
		}

		if i >= 5 && debug && len(Trackers) > 3 {
			break
		}
	}

	if err := generateTemplate(); err != nil {
		log.Fatalln(err)
	}
}

// getDefinitionsEntry returns a subtree (directory) with our definitions files.
func getDefinitionsEntry(entries []*github.TreeEntry) *github.TreeEntry {
	for _, entry := range entries {
		if entry.GetPath() == DefinitionsFolderName {
			return entry
		}
	}

	return nil
}

func processTracker(ctx context.Context, client *github.Client, entry *github.TreeEntry) error {
	// get entry detail
	definition, _, err := client.Git.GetBlob(ctx, "Jackett", "Jackett", entry.GetSHA())
	if err != nil {
		return fmt.Errorf("cannot get tracker file content from GitHub Git service: %w", err)
	}

	// read file content as base64
	raw, err := base64.StdEncoding.DecodeString(definition.GetContent())
	if err != nil {
		return fmt.Errorf("cannot base64.DecodeString [%s] content: %w", entry.GetPath(), err)
	}

	// parse content as YAML
	var trackerDefinition TrackerDefinition
	if err := yaml.NewDecoder(bytes.NewReader(raw)).Decode(&trackerDefinition); err != nil {
		return fmt.Errorf("cannot yaml.Decode [%s] content: %w", entry.GetPath(), err)
	}

	// exclude trackers which contain number
	if re.MatchString(trackerDefinition.Name) {
		return nil
	}

	// write data to our map
	Trackers = append(Trackers, &trackerDefinition)
	return nil
}

func generateTemplate() error {
	tmpl, err := template.New("template").Funcs(template.FuncMap{
		"camelCase": func(v string) string {
			return strcase.ToCamel(v)
		},
	}).Parse(tmpl)
	if err != nil {
		return fmt.Errorf("cannot parse template: %w", err)
	}

	// write generated code
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, map[string]interface{}{
		"Package":  packageName,
		"Trackers": Trackers,
	}); err != nil {
		return fmt.Errorf("cannot execute template: %w", err)
	}

	// format generated code
	p, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("cannot format generated file: %w", err)
	}

	file, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("cannot open output [%s]: %w", output, err)
	}

	// write result to file
	if _, err := file.Write(p); err != nil {
		return fmt.Errorf("cannot write to file: %w", err)
	}

	if err := file.Close(); err != nil {
		return fmt.Errorf("cannot close file: %w", err)
	}

	return nil
}
