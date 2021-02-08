package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "os"
    "path/filepath"
    "strings"
   
    "io/ioutil"
    "gopkg.in/yaml.v2"
)

type Entry struct {
    Description string`yaml:"description",json:"description"`
    Data        []string `yaml:"data",json:"data"`
}

type DocResource struct {
    Path        string
    Cheats      []Entry`yaml:"cheats",json:"cheats"`
    Links       []Entry`yaml:"links",json:"links"`
    Glossary    []Entry`yaml:"glossary",json:glossary`
}

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

func (r *DocResource) LoadYAML(path string) error {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return err
    }
    r.Path = path
    return yaml.Unmarshal(data, r)
}

func (r *DocResource) PrintByName(name string, terms string, colored bool, matched bool) {
    var entries []Entry

    switch name {
        case "cheats":
            entries = r.Cheats
        case "links":
            entries = r.Links
        case "glossary":
            entries = r.Glossary
        default:
            entries = r.Cheats
    }

    for _, entry := range entries {
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

func (c *Config) PrintResults(res *DocResource) {
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
        } 
        if c.OnlyLinks {
            res.PrintByName("links", c.Terms, c.Colored, c.Matched)
        }
        if c.OnlyGlossary {
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
                    for _, cheat := range tmp.Cheats {
                        if strings.Contains(cheat.Description, c.Terms) {
                            res.Cheats = append(res.Cheats, cheat)
                        }
                    }
                    for _, link := range tmp.Links {
                        if strings.Contains(link.Description, c.Terms) {
                            res.Links = append(res.Links, link)
                        }
                    }
                    for _, glossary := range tmp.Glossary {
                        if strings.Contains(glossary.Description, c.Terms) {
                            res.Glossary = append(res.Glossary, glossary)
                        }
                    }
                } else {
                    res.Cheats = append(res.Cheats, tmp.Cheats...)
                    res.Links = append(res.Links, tmp.Links...)
                    res.Glossary = append(res.Glossary, tmp.Glossary...)
                }
                res.Path = tmp.Path
            }
        }
    }
    c.PrintResults(&res)
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
                 for _, cheat := range tmp.Cheats {
                     if strings.Contains(cheat.Description, c.Terms) {
                         res.Cheats = append(res.Cheats, cheat)
                     }
                 }
                 for _, link := range tmp.Links {
                     if strings.Contains(link.Description, c.Terms) {
                         res.Links = append(res.Links, link)
                     }
                 }
                 for _, glossary := range tmp.Glossary {
                     if strings.Contains(glossary.Description, c.Terms) {
                         res.Glossary = append(res.Glossary, glossary)
                     }
                 }
                 if len(res.Cheats) > 0 || len(res.Links) > 0 || len(res.Glossary) > 0 {
                     res.Path = topic
                     results = append(results, res)
                 }
             }
        }
    }
    for _, result := range results {
        c.PrintResults(&result)
    }
}

func Dispatcher(c *Config) {
    if c.Inventory {
        ListTopics(c)
    } else if c.Env {
        ShowEnv(c)
    } else if c.Topic != "" {
        if c.Path {
            ShowPath(c)
        } else {
            ShowTopicContent(c)
        }     
    } else if c.Terms != "" {
        SearchAndShowTopicContent(c)
    }
}

func main() {
    
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
    flag.BoolVar(&c.OnlyGlossary, "G", false, "Restrict search on glossary terms")
    flag.BoolVar(&c.OnlyLinks, "L", false, "Restrict search on links terms")
    flag.BoolVar(&c.Colored, "c", true, "Enable colored output")
    flag.BoolVar(&c.AsJson, "j", false, "JSON output")
    flag.BoolVar(&c.Matched, "m", false, "Enable colored match")

    // Data
    flag.StringVar(&c.Terms, "s", "", "Keyword or term to search")
    flag.StringVar(&c.Topic, "t", "", "Search on a specific topic")
    flag.Parse()

    c.EnvPaths = strings.Split(os.Getenv("DOCSEARCH_PATH"), ":")
    Dispatcher(&c)
}
