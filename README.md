<p align="center">
  <img src="icon.svg" width="300" />
</p>

# Anon — A UNIX Command To Anonymise Data
[![Build Status](https://travis-ci.org/intenthq/anon.svg?branch=master)](https://travis-ci.org/intenthq/anon) <a href="https://codecov.io/gh/intenthq/anon">
  <img src="https://codecov.io/gh/intenthq/anon/branch/master/graph/badge.svg" />
</a> [![Go Report Card](https://goreportcard.com/badge/github.com/intenthq/anon)](https://goreportcard.com/report/github.com/intenthq/anon) [![License](https://img.shields.io/npm/l/express.svg)](https://github.com/intenthq/anon/LICENSE)
![GitHub release](https://img.shields.io/github/release/intenthq/anon.svg)

Anon is a tool for taking delimited files and anonymising or transforming columns until the output is useful for applications where sensitive information cannot be exposed.

## Installation

Releases of Anon are available as pre-compiled static binaries [on the corresponding GitHub release](https://github.com/intenthq/anon/releases). Simply download the appropriate build for your machine and make sure it's in your `PATH` (or use it directly).

## Usage

```sh
anon [--config <path to config file, default is ./config.json>]
     [--output <path to output to, default is STDOUT>]
```

Anon is designed to take input from `STDIN` and by default will output the anonymised file to `STDOUT`:

```sh
anon < some_file.csv > some_file_anonymised.csv
```

### Configuration

In order to be useful, Anon needs to be told what you want to do to each column of the CSV. The config is defined as a JSON file (defaults to a file called `config.json` in the current directory):

```json5
{
  "csv": {
    "delimiter": ","
  },
  // Optionally define a number of rows to randomly sample down to.
  // To do it, it will hash (using FNV-1 32 bits) the column with the ID
  // in it and will mod the result by the value specified to decide if the
  // row is included or not -> include = hash(idColumn) % mod == 0
  "sampling": {
    // Number used to mod the hash of the id and determine if the row
    // has to be included in the sample or not
    "mod": 30000
    // Specify in which a column a unique ID exists on which the sampling can
    // be performed. Indices are 0 based, so this would sample on the first
    // column.
    "idColumn": 0
  },
  // An array of actions to take on each column - indices are 0 based, so index
  // 0 in this array corresponds to column 1, and so on.
  //
  // There must be an action for every column in the CSV.
  "actions": [
    {
      // The no-op, leaves the input unchanged.
      "name": "nothing"
    },
    {
      // Takes a UK format postcode (eg. W1W 8BE) and just keeps the outcode
      // (eg. W1W).
      "name": "outcode"
    },
    {
      // Hash (SHA1) the input.
      "name": "hash"
    },
    {
      // Given a date, just keep the year.
      "name": "year",
      "dateConfig": {
        // Define the format of the input date here.
        "format": "YYYYmmmdd"
      }
    },
    {
      // Summarise a range of values.
      "name": "range",
      "rangeConfig": {
        "ranges": [
          // For example, this will take values between 0 and 100, and convert
          // them to the string "0-100".
          // You can use one of (gt, gte) and (lt, lte) but not both at the
          // same time.
          // You also need to define at least one of (gt, gte, lt, lte).
          {
            "gte": 0,
            "lt": 100,
            "output": "0-100"
          }
        ]
      }
    }
  ]
}
```

## How to contribute

Any contribution is welcome, raise a bug (and fix it! :-)) request or add a new feature... Don't be shy and raise a pull request, anything on the following topics will be very welcome:
- New actions to anonymise data
- New input formats (JSON?)
- Bug fixes

You can also take a look at the [issues](https://github.com/intenthq/anon/issues) and pick the one you like better.

If you are going to contribute, we ask you to do the following:
- Use `gofmt` to format your code
- Check your code with `go vet`, `gocyclo`, `golint`
- Cover the logic with enough tests

# License

This project is [licensed under the MIT license](LICENSE).

The icon is by [Pixel Perfect](https://www.flaticon.com/authors/pixel-perfect) from [Flaticon](https://www.flaticon.com/), and is licensed under a [Creative Commons 3.0 BY](http://creativecommons.org/licenses/by/3.0/) license.
