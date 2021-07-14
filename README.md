<p align="center"><br><br><br><br>
:coffin:<br>
<b>GitHub Issues closer</b>
</p>

<p align="center">
This is a CLI tool that closes all issues that match your search. <br>
There is an option to close only duplicates and support for GitHub Enterprise.
</p><br><br><br><br>

Installation
--

To install, use `go get`:

```sh
$ go get github.com/linyows/github-issues-closer
```

Usage
--

Example:

```sh
$ export GITHUB_TOKEN=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
# github-issues-closer <repository> <search-word>
$ github-issues-closer linyows/github-issues-closer 17a7a5a250601fe9
```

If github enteprise:

```sh
$ github-issues-closer -e https://mygithub.com/api/v3 linyows/github-issues-closer 17a7a5a250601fe9
```

If duplicate only:

```sh
$ github-issues-closer -d linyows/github-issues-closer 17a7a5a250601fe9
```

Contribution
------------

1. Fork ([https://github.com/linyows/github-issues-closer/fork](https://github.com/linyows/github-issues-closer/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

Author
------

[linyows](https://github.com/linyows)

