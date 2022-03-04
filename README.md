# Scanberry

[![builds.sr.ht status](https://builds.sr.ht/~bascht/scanberry.svg)](https://builds.sr.ht/~bascht/scanberry?)

First iteration of `Scanberry`, a little web wrapper around `scanimage` which can run on a Raspberry Pi.
Scanberry just provides a little web UI, all the heavy lifting is done by `scanimage`, `ocrmypdf` and `imagemagick`. 

## Functionality of the web UI
- Trigger the scanning process
- Set a document name and whether it's duplex or simplex
- Prepend a timestamp prefix
- Offer the file as a download
- Copy it to a pre-determined location (e.g. my Nextcloud share)
