package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "os"
    "path/filepath"
    "strings"
)

type Config struct {
    // Format & Filter flags
    AsJson        bool
    Colored       bool
    Env           bool
    Inventory     bool
    Matched       bool
    OnlyCheats    bool
    OnlyLinks     bool
    OnlyGlossary  bool
    Path          bool
    
    // Storage
    Topic         string
    Terms         string
    EnvPaths      []string
}


func PrintResults(c *Config, res *DocResource) {
    if c.AsJson {
        if data, err := json.Marshal(res); err == nil {
            fmt.Println(string(data))
        }
    } else if !c.OnlyCheats && !c.OnlyLinks && !c.OnlyGlossary {
        res.PrintByName("cheats", c.Terms, c.Colored, c.Matched)
        res.PrintByName("links", c.Terms, c.Colored, c.Matched)
        res.PrintByName("glossary", c.Terms, c.Colored, c.Matched)
    } else {
        if c.OnlyCheats {
            res.PrintByName("cheats", c.Terms, c.Colored, c.Matched)
        } else if c.OnlyLinks {
            res.PrintByName("links", c.Terms, c.Colored, c.Matched)
        } else if c.OnlyGlossary {
            res.PrintByName("glossary", c.Terms, c.Colored, c.Matched)
        }
    }
}

func GetTopicFromPath(path string) string {
    topic := filepath.Base(path)
    return topic[:len(topic)-len(filepath.Ext(topic))]
}

func ListTopics(c *Config) {
    for _, root := range c.EnvPaths {
        fmt.Println(root)
        path := fmt.Sprintf("%s/*.yaml", root)
        topics, _ := filepath.Glob(path)
        for  _, topic := range topics {
            fmt.Println(GetTopicFromPath(topic))
        }
        fmt.Println("")
    }
}

func ShowEnv(c *Config) {
    fmt.Println("\n[*] DOCSEARCH_PATH")
    for _, p := range c.EnvPaths {
        fmt.Println(p)
    }
    fmt.Println("")
    if os.Getenv("DOCSEARCH_COLORED") != "" {
        fmt.Println("[*] DOCSEARCH_COLORED")
        mode := "disabled"
        if os.Getenv("DOCSEARCH_COLORED") == "1" {
            mode = "enabled"
        }
        fmt.Printf("color mode %s\n\n", mode)
    }
}

func ShowPath(c *Config) {
    for _, path := range c.EnvPaths {
        topic := fmt.Sprintf("%s/%s.yaml", path, c.Topic)
        if _, err := os.Stat(topic); err == nil {
            fmt.Println(topic)
        }
    }
}

func ShowTopicContent(c *Config) {
    var res DocResource
    for _, path := range c.EnvPaths {
        topic := fmt.Sprintf("%s/%s.yaml", path, c.Topic)
        if _, err := os.Stat(topic); err == nil {
            var tmp DocResource
            if err := tmp.LoadYAML(topic); err == nil {
                if c.Terms != "" {
                    res.UpdateOnMatch("cheats", c.Terms, &tmp)
                    res.UpdateOnMatch("links", c.Terms, &tmp)
                    res.UpdateOnMatch("glossary", c.Terms, &tmp)
                } else {
                    res.Set("cheats", tmp.Cheats)
                    res.Set("links", tmp.Links)
                    res.Set("glossary", tmp.Glossary)
                }
                res.Path = tmp.Path
            }
        }
    }
    PrintResults(c, &res)
}

func SearchAndShowTopicContent(c *Config) {
    var results []DocResource
    for _, root := range c.EnvPaths {
        path := fmt.Sprintf("%s/*.yaml", root)
        topics, _ := filepath.Glob(path)
        for  _, topic := range topics {
            var tmp DocResource
            var res DocResource
             if err := tmp.LoadYAML(topic); err == nil {
                 res.UpdateOnMatch("cheats", c.Terms, &tmp)
                 res.UpdateOnMatch("links", c.Terms, &tmp)
                 res.UpdateOnMatch("glossary", c.Terms, &tmp)
                 if res.HasEntries() {
                     res.Path = topic
                     results = append(results, res)
                 }
             }
        }
    }
    for _, result := range results {
        PrintResults(c, &result)
    }
}

func main() {
    // Guard 
    if os.Getenv("DOCSEARCH_PATH") == "" {
        fmt.Printf("You need to declare DOCSEARCH_PATH environment variable.")
        os.Exit(2)
    }

    var c = Config{}

    // Filters
    flag.BoolVar(&c.Env, "env", false, "Show useful DOCSEARCH_* environment variables")
    flag.BoolVar(&c.Env, "e", false, "-env (shorthand)")
    flag.BoolVar(&c.Inventory, "i", false, "List all availabled topics")
    flag.BoolVar(&c.Path, "p", false, "Show matched topics fullpath")
    
    //  Format
    flag.BoolVar(&c.OnlyCheats, "C", false, "Restrict search on cheatsheets terms")
    flag.BoolVar(&c.OnlyLinks, "L", false, "Restrict search on links terms")
    flag.BoolVar(&c.OnlyGlossary, "G", false, "Restrict search on glossary terms")
    flag.BoolVar(&c.Colored, "c", true, "Enable colored output")
    flag.BoolVar(&c.AsJson, "j", false, "JSON output")
    flag.BoolVar(&c.Matched, "m", false, "Enable colored match")

    // Data
    flag.StringVar(&c.Terms, "s", "", "Keyword or term to search")
    flag.StringVar(&c.Topic, "t", "", "Search on a specific topic")
    flag.Parse()

    c.EnvPaths = strings.Split(os.Getenv("DOCSEARCH_PATH"), ":")

    // Dispatcher
    if c.Inventory {
        ListTopics(&c)
    } else if c.Env {
        ShowEnv(&c)
    } else if c.Topic != "" {
        if c.Path {
            ShowPath(&c)
        } else {
            ShowTopicContent(&c)
        }     
    } else if c.Terms != "" {
        SearchAndShowTopicContent(&c)
    }
}
