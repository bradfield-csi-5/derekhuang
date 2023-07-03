# Completion Exercise

Please complete one of the following two exercises from *The Go Programming Language*:

* **Exercise 4.12**: The popular web comic *xkcd* has a JSON interface. For example, a request to https://xkcd.com/571/info.0.json produces a detailed description of comic 571, one of many favorites. Download each URL (once!) and build an offline index. Write a tool `xkcd` that, using this index, prints the URL and transcript of each comic that matches a search term provided on the command line.
* **Exercise 8.7**: Write a concurrent program that creates a local mirror of a web site, fetching each reachable page and writing it to a directory on the local disk. Only pages within the original domain (for instance, `golang.org`) should be fetched. URLs within the mirrored pages should be altered as needed so that they refer to the mirrored page, not the original.

In addition, please add some tests and benchmarks to familiarize yourself with the capabilities of the `testing` package.
