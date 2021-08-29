# GODOCSEARCH

Like docsearch.rb but with Golang experimentations

## Install

```
$ go build -o gdcs
$ go install
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
