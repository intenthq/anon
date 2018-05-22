<p align="center">
  <img src="icon.svg" width="300" />
</p>

# [![Build Status](https://travis-ci.org/intenthq/anon.svg?branch=master)](https://travis-ci.org/intenthq/anon) Anon — A UNIX Command To Anonymise Data

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
    "delimiter": ",",
    "quote": "\""
  },
  // Optionally define a number of rows to randomly sample down to.
  "sampling": {
    "num": 30000
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

# License

This project is [licensed under the MIT license](LICENSE).

The icon is by [Pixel Perfect](https://www.flaticon.com/authors/pixel-perfect) from [Flaticon](https://www.flaticon.com/), and is licensed under a [Creative Commons 3.0 BY](http://creativecommons.org/licenses/by/3.0/) license.
