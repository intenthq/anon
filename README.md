<p align="center">
  <img src="icon.svg" width="300" />
</p>

# Anon — A UNIX Command To Anonymise Data
[![Build Status](https://travis-ci.org/intenthq/anon.svg?branch=master)](https://travis-ci.org/intenthq/anon) <a href="https://codecov.io/gh/intenthq/anon">
  <img src="https://codecov.io/gh/intenthq/anon/branch/master/graph/badge.svg" />
</a> [![Go Report Card](https://goreportcard.com/badge/github.com/intenthq/anon)](https://goreportcard.com/report/github.com/intenthq/anon) [![License](https://img.shields.io/npm/l/express.svg)](https://github.com/intenthq/anon/LICENSE)
![GitHub release](https://img.shields.io/github/release/intenthq/anon.svg)

Anon is a tool for taking delimited files and anonymising or transforming columns/fields until the output is useful for applications where sensitive information cannot be exposed. Currently this tools supports both CSV and JSON files (with one level of depth).

## Installation

Releases of Anon are available as pre-compiled static binaries [on the corresponding GitHub release](https://github.com/intenthq/anon/releases). Simply download the appropriate build for your machine and make sure it's in your `PATH` (or use it directly).

## Usage

```sh
anon [--config <path to config file, default is ./config.json>]
     [--output <path to output to, default is STDOUT>]
```

Anon is designed to take input from `STDIN` and by default will output the anonymised file to `STDOUT`:

```sh
anon < some_file > some_file_anonymised
```

### Configuration

In order to be useful, Anon needs to be told what you want to do to each column/field of the input. The config is defined as a JSON file (defaults to a file called `config.json` in the current directory):

```json5
{
  // Name of the format of the input file
  // Currently supports "csv" and "json"
  "formatName": {
    // Options for the format you have picked go here.
    // See the documentation for the format you choose below.
  },
  // Optionally define a number of rows to randomly sample down to.
  // To do it, it will hash (using FNV-1 32 bits) the column with the ID
  // in it and will mod the result by the value specified to decide if the
  // row is included or not -> include = hash(idColumn) % mod == 0
  "sampling": {
    // Number used to mod the hash of the id and determine if the row
    // has to be included in the sample or not
    "mod": 30000
  },
  // An array of actions to take on each column - indices are 0 based, so index
  // 0 in this array corresponds to column 1, and so on.
  //
  // If anonymising a CSV, there must be an action for every column in it.
  // If anonymising a JSON, there must be an action for each field that needs to
  // be anonymised. If there is no action defined for a specific field, this
  // field value will be left untouched.
  "actions": [
    {
      // The no-op, leaves the input unchanged.
      "name": "nothing"
    },
    {
      // Takes a UK format postcode (eg. W1W 8BE) and just keeps the outcode
      // (eg. W1W).
      "name": "outcode",
      // what field in the json this action needs to be applied. If a field in
      // the json doesn't have an action defined, then it will be left untouched.
      "jsonField": "postcode"
    },
    {
      // Hash (SHA1) the input.
      "name": "hash",
      // Optional salt that will be appened to the input.
      // If not defined, a random salt will be generated
      "salt": "salt"
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

## Formats

You can use CSV or JSON files as input.

### CSV

For a CSV file you will need a config like this:

```json5
"csv": {
  "delimiter": ",",
  // Specify in which column a unique ID exists on which the sampling can
  // be performed. Indices are 0 based, so this would sample on the first
  // column.
  "idColumn": "0"
}
```

### JSON

For a JSON file you will need to define config like this: 

```json5
"json": {
  // Specify in which field a unique ID exists on which the sampling can
  // be performed.
  "idField": "id"
}
```

## Contributing

Any contribution will be welcome, please refer to our [contributing guidelines](CONTRIBUTING.md) for more information.

## License

This project is [licensed under the MIT license](LICENSE).

The icon is by [Pixel Perfect](https://www.flaticon.com/authors/pixel-perfect) from [Flaticon](https://www.flaticon.com/), and is licensed under a [Creative Commons 3.0 BY](http://creativecommons.org/licenses/by/3.0/) license.
