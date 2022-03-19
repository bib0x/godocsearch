# GODOCSEARCH

Like docsearch.rb but with Golang experimentations

## Build

### Manual

```bash
$ go build && go install
```

### Nixos

```
$ cat default.nix
{ lib, buildGoModule, fetchFromGitHub, installShellFiles }:

buildGoModule rec {
  pname = "godocsearch";
  version = "0.1.0";

  src = fetchFromGitHub {
    owner = "bib0x";
    repo = "godocsearch";
    rev = "v${version}";
    sha256 = "0kmvxc5miab2p0xyfmpgg8cv37gy82l54xhfgzm2c9fg92xcgwqh";
  };

  vendorSha256 = "11r1l5lcdfm3wymrkbddl5khpjmr30jln31l40mfyyy9msnqayf3";

  nativeBuildInputs = [ installShellFiles ];

  postInstall = ''
    installShellCompletion --bash --name godocsearch.bash gdcs-completion.bash
  '';

  meta = with lib; {
    description = "Simple cheatsheets CLI written in Go";
    homepage = "https://github.com/bib0x/godocsearch";
    license = licenses.mit;
    maintainers = with maintainers; [ bib0x ];
    platforms = platforms.linux;
  };

}
```

## Configuration

```bash
export GODOCSEARCH_ROOT="$HOME/perso/git/godocsearch"

if [ -x "$GOBIN"/godocsearch ]; then
    alias gdcs="$GOBIN"/godocsearch

    export DOCSEARCH_PATH="$HOME/perso/git/resources:$HOME/work/git/resources"
    export DOCSEARCH_COLORED=1

    # autocompletion if not using Nixos
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

## Bonus

You can validate that your resource files as well formatted, using the script `check.cue` such as:

```
$ nix-shell -p cue
$ cue vet check.cue $HOME/perso/git/resources/git.yaml
# if no error, your are good
```
