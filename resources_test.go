package main

import (
    "fmt"
    "os"
    "testing"
)

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

var docr DocResource

func init() {
    root, _ := os.Getwd()
    path := fmt.Sprintf("%s/testdata/git.yaml", root)
    docr.LoadYAML(path)
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
    
    docr.Append("cheats", entry)
    docr.Append("cheats", entry)
    docr.Append("glossary", entry)

    for _, test := range tests {
        if entry := docr.Get(test.name); len(entry) != test.length {
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
    
    var res DocResource
    res.Set("cheats", docr.Cheats)
    res.Set("links", docr.Links)

    for _, test := range tests {
        if entry := res.Get(test.name); len(entry) != test.length {
            t.Errorf("Test Failed: i=%s e=%v r=%v", test.name, test.length, len(entry))
        }
    }
}
