package main

import (
    "fmt"
    "os"
    "testing"
)

var docr DocResource

func init() {
    root, _ := os.Getwd()
    path := fmt.Sprintf("%s/testdata/git.yaml", root)
    docr.LoadYAML(path)
}

func Test_LoadYAML(t *testing.T) {
    root, _ := os.Getwd()
    var tests = []struct {
        path     string
        expected error
    }{
        {fmt.Sprintf("%s/testdata/EmptyFile.yaml", root), nil},
        {fmt.Sprintf("%s/testdata/NotYAML.json", root), nil},
        {fmt.Sprintf("%s/testdata/vbox.yaml", root), nil},
    }

    for _, test := range tests {
        var res DocResource
        if err := res.LoadYAML(test.path); err != test.expected {
            t.Errorf("Test Failed: i=%s, e=%v, r=%v", test.path, test.expected, err)
        }
    }
}

func Test_Get(t *testing.T) {
    var tests = []struct {
        name     string
        length   int
    }{
        {"cheats", 17},
        {"links", 5},
        {"glossary", 1},
        {"bibop", 0},
    }

    for _, test := range tests {
        if entry := docr.Get(test.name); len(entry) != test.length {
            t.Errorf("Test Failed: i=%s e=%v r=%v", test.name, test.length, len(entry)) 
        }
    }
}

func Test_Append(t *testing.T) {
    var tests = []struct {
        name    string
        length  int
    }{
        {"cheats", 19},
        {"links", 5},
        {"glossary", 2},
    }

    data := []string{"test1", "test2"}
    entry := Entry{"this is a test", data}

    var tmp DocResource = docr
    tmp.Append("cheats", entry)
    tmp.Append("cheats", entry)
    tmp.Append("glossary", entry)

    for _, test := range tests {
        if entry := tmp.Get(test.name); len(entry) != test.length {
            t.Errorf("Test Failed: i=%s e=%v r=%v", test.name, test.length, len(entry)) 
        }
    }
}

func Test_Set(t *testing.T) {
    var tests = []struct {
        name     string
        length   int
    }{
        {"cheats", 17},
        {"links", 5},
        {"glossary", 0},
    }
    
    var tmp DocResource
    tmp.Set("cheats", docr.Cheats)
    tmp.Set("links", docr.Links)

    for _, test := range tests {
        if entry := tmp.Get(test.name); len(entry) != test.length {
            t.Errorf("Test Failed: i=%s e=%v r=%v", test.name, test.length, len(entry))
        }
    }
}

func Test_UpdateOnMatch(t *testing.T) {
    var tests = []struct {
        name  string
        terms string
        length int 
    }{
       { "cheats", "delete", 2},
       { "links", "emoji",  1},
       { "glossary", "rebase", 1},
    }

    for _, test := range tests {
        var tmp DocResource
        tmp.UpdateOnMatch(test.name, test.terms, &docr)

        if entry := tmp.Get(test.name); len(entry) != test.length {
            t.Errorf("Test Failed: i=%s e=%v r=%v", test.name, test.length, len(entry))
        }
    }
}

func Test_HasEntries(t *testing.T) {
    var tests = []struct {
        res DocResource
        expected bool
    }{
        { docr, true },
        { DocResource{}, false },
    }

    for _, test := range tests {
        if test.res.HasEntries() != test.expected {
            t.Errorf("Test Failed: i=%v e=%v", test.res, test.expected)
        }
    }
}
