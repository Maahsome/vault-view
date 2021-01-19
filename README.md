# vault-view

CLI tool with TUI interface used to browse Hashicorp Vault secret stores

The `derailed/k9s` project has inspired me.  A great number of CLI utilities
ran through my head.  A TUI CLI for browsing Hashicorp vault may not be something
that is needed, though it should be fun to try to create an experience similar to
the k9s interface.  So much coolness.

## Credits

|Github Repo|Notes|
|----------|---|
|derailed/k9s|inspiration mostly, hopefully I'll understand the code better when I'm done|
|skanehira/docui|I used this project as the framework for my work.  It was easier to follow than the k9s code|
|tjgq/clipboard|I most certainly wouldn't have wanted to re-invent this wheel!|
|rivo/tview & gdamore/tcell|Clearly cannot do TUI without these|

## Build

Check out the repository, and `go build` should do the trick.  Once things get
to a good stable point, there will be a Homebrew install.

## Usage

Currently looks for VAULT_ADDR and VAULT_TOKEN environment variables to be set
in order to connect to the Hashicorp Vault instance.

## Roadmap

- [ ] Implement _test.go tests
- [ ] Add Secret Version Navigation
- [ ] Open a view of RAW returned information from Vault for a "Data" object
- [ ] Generate BASH script lines from the selected/marked item(s)
- [ ] Generate Jenkinsfile lines from the selected/marked item(s)
- [ ] Diffing Versions of a "Data" object
- [ ] Indicator of "marked" items in sub folders
- [ ] TUI improvements
  - [ ] Suggestions for : commands
  - [ ] Improved searching: add || && logic

### Potentially View Other Vault Object Types

- [ ] - Policies
- [ ] - Roles
- [ ] - Tokens?
