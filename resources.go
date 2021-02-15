package main

import (
    "fmt"
    "io/ioutil"
    "strings"

    "gopkg.in/yaml.v2"
)

type Entry struct {
    Description string`yaml:"description"`
    Data        []string `yaml:"data"`
}

type DocResource struct {
    Path        string
    Cheats      []Entry`yaml:"cheats"`
    Links       []Entry`yaml:"links"`
    Glossary    []Entry`yaml:"glossary"`
}

func (r *DocResource) LoadYAML(path string) error {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return err
    }
    r.Path = path
    return yaml.Unmarshal(data, r)
}

func (r *DocResource) PrintByName(name string, terms string, colored bool, matched bool) {
    for _, entry := range r.Get(name) {
        description := entry.Description
        if matched {
            cterms := fmt.Sprintf("\033[31;1m%s\033[0m", terms)
            description = strings.Replace(entry.Description, terms, cterms, -1)
        }
        if colored {
            fmt.Printf("\033[33;1m[%s]\033[0m %s\n", GetTopicFromPath(r.Path), description)
        } else {
            fmt.Printf("[%s] %s\n", GetTopicFromPath(r.Path), description)
        }
        for _, data := range entry.Data {
            fmt.Println(data)
        }
        fmt.Println("")
    }
}

func (r *DocResource) Get(name string) []Entry {
    if name == "cheats" {
        return r.Cheats
    } else if name == "links" {
        return r.Links
    } else if name == "glossary" {
        return r.Glossary
    }
    return nil
}

func (r *DocResource) Append(name string, entry Entry) {
   if name == "cheats" {
        r.Cheats = append(r.Cheats, entry)
   } else if name == "links" {
        r.Links = append(r.Links, entry)
   } else if name == "glossary" {
        r.Glossary = append(r.Glossary, entry)
   }
}

func (r *DocResource) Set(name string, entries []Entry) {
   if name == "cheats" {
        r.Cheats = append(r.Cheats, entries...)
   } else if name == "links" {
        r.Links = append(r.Links, entries...)
   } else if name == "glossary" {
        r.Glossary = append(r.Glossary, entries...)
   }
}

func (r *DocResource) UpdateOnMatch(name string, terms string, from *DocResource) {
    for _, entry := range from.Get(name) {
        if strings.Contains(entry.Description, terms) {
            r.Append(name, entry)
        }
    }
}

func (r *DocResource) HasEntries() bool {
    return len(r.Cheats) > 0 || len(r.Links) > 0 || len(r.Glossary) > 0
}
