# Global Names Finder [![Build Status][travis-img]][travis] [![Doc Status][doc-img]][doc]


Finds scientific names using dictionary and nlp approaches.

## Development

There are dictionaries used by the project that are too big to handle with
"plain git". Dictionaries are stored using [lfs] extention to git. You need to
setup git-lfs to pull dictionaries to your local copy of the repository

You will also need some additional go packages

```
go get github.com/json-iterator/go
go get github.com/rakyll/statik
```

To update dictionaries `cd` to the projects' root directory and run

```
go generate
```

### Testing

Install [ginkgo], a [BDD] testing framefork for Go.

```bash
go get github.com/onsi/ginkgo/ginkgo
go get github.com/onsi/gomega
```

To run tests go to root directory of the project and run

```bash
ginkgo

#or

go test
```

[lfs]: https://git-lfs.github.com/
[travis-img]: https://travis-ci.org/gnames/gnfinder.svg?branch=master
[travis]: https://travis-ci.org/gnames/gnfinder
[doc-img]: https://godoc.org/github.com/gnames/gnfinder?status.png
[doc]: https://godoc.org/github.com/gnames/gnfinder
