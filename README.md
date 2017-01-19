#urlChecker

This is a simple tool for checking URLs. It takes a CSV file with URLs in first parameter and it ouputs a CSV file with KO or OK regarding http requests responses.

```
Usage: urlchecker [OPTIONS] SRC DEST

Checks urls from a 'csv' FILE and write to a DEST

Arguments:
  SRC="urls.csv"            File with urls
  DEST="urls-results.csv"   Destination file with results

Options:
  -v, --version     Show the version and exit
  -l, --logs=true   Print result of requests
```