## Synopsis

A simple web crawler implemented in Go.

## Details

I started this project mainly to learn more about Go. Currently this will connect to a website,
scrape and parse the HTML, then send the links it finds through a channel. Currently it will only
register links for the current domain (or a sub domain) of the initial base URL.

In the future I hope to recursively follow each of the links that are found, register a map of the
site. Another good project would be to expand parsing and maybe even index the text found on each
page, like a real world search index crawler.

## Author

Matthew Emerson

## License

Released under MIT License.
