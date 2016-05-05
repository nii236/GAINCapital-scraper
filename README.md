# Project Name

This is a simple script written in Go that will download historical ticker data from GAINCapital for you.

## Installation

### Easy

Download the binary from the releases page.

### Harder

Make sure you have:

- Go 1.6
- Glide

```
$ git clone git@github.com:nii236/GAINCapital-scraper.git
$ cd GAINCapital-scraper
$ glide up
$ go build
```

## Usage

Create a `config.json` file. Here's an example:

```
{
  "from": 2013,
  "to": 2013,
  "pairs": ["AUD_USD", "EUR_USD", "GBP_USD", "USD_JPY"]
}

```

Run the program:

```
./GAINCapital-scraper
```

Enjoy!

## Contributing

1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request :D
