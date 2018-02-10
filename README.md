# nrm

Upload the screenshot to github.
Then, embed tags in your clipboard.

## Usage

![video](https://github.com/pyama86/nrm/blob/master/misc/nrm.gif?raw=true)

please set github auth token
```
export GITHUB_AUTH_TOKEN=xxxxxxxxxxxxxxxxxx:
```

```
Usage of nrm:
  -o string
        owner(Short)
  -owner string
        owner
  -r string
        repo(Short)
  -repo string
        repo
  -version
        Print version information and quit.
```

### Github Enterprise

plese set endpoint url to GITHUB_API_ENDPOINT environment variables.
```
export GITHUB_API_ENDPOINT=https://git.xxx.com/api/v3
```

```bash
$ go get -d github.com/pyama86/nrm
```

## Contribution

1. Fork ([https://github.com/pyama86/nrm/fork](https://github.com/pyama86/nrm/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[pyama86](https://github.com/pyama86)
