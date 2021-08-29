# GODOCSEARCH

Like docsearch.rb but with Golang experimentations

## Install

```
$ go build && go install
```

## Configuration

```
export GODOCSEARCH_ROOT="$HOME/perso/git/godocsearch"

if [ -x "$GOBIN"/godocsearch ]; then
    alias gdcs="$GOBIN"/godocsearch
    # autocompletion
    if [ -f "$GODOCSEARCH_ROOT/gdcs-completion.bash" ]; then
        . "$GODOCSEARCH_ROOT/gdcs-completion.bash"
    fi
fi
```

## Usage

```
Usage of godocsearch:
  -C	Restrict search on cheatsheets terms
  -G	Restrict search on glossary terms
  -L	Restrict search on links terms
  -c	Enable colored output (default true)
  -e	-env (shorthand)
  -env
    	Show useful DOCSEARCH_* environment variables
  -i	List all availabled topics
  -j	JSON output
  -m	Enable colored match
  -p	Show matched topics fullpath
  -s string
    	Keyword or term to search
  -t string
    	Search on a specific topic
```
