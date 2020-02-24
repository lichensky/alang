# alang

Tool for cleaning docker images from remote repositories.

At the moment only GCR is supported.

## Config

Config requires repository list in YAML format.

Example config:

```
repositories:
  - name: gcr.io/my-repository
    gracePeriod: 10
    numberToKeep: 5
    cleanTags: true
    keepTags:
      - master
      - latest
      - ^v(0|\d+).(0|\d+).(0|\d+)$
```

Each repository has its own set of parameters:
* `name` - Full repository name. This field is required.
* `gracePeriod` - Relative grace period in days. All newer images won't be
deleted. Defaults to 0.
* `numberToKeep` - Number of latest images to keep. Defaults to 0.
* `cleanTags` - If set to true, tagged images will be deleted. Defaults to
false.
* `keepTags` - List of tag patterns to keep. Could be either exact string or
regular expression.

## Usage

Run the `alang` command with config path as `-config` flag.
The `-config` defaults to `./config.yaml`

```
alang -config myconfig.yaml
```

### Dry run

If you want to only check the result, but not actually delete images use
`-dry` flag.

```
alang -config myconfig.yaml -dry
```

## Why alang?

[Alang](https://en.wikipedia.org/wiki/Alang) is a town which have become a major
worldwide centre for ship breaking.
