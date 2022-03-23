# Scanberry

[![builds.sr.ht status](https://builds.sr.ht/~bascht/scanberry.svg)](https://builds.sr.ht/~bascht/scanberry?)

First iteration of `Scanberry`, a little web wrapper around `scanimage` which can run on a Raspberry Pi.
Scanberry just provides a little web UI, all the heavy lifting is done by `scanimage`, `ocrmypdf` and `imagemagick`. 

![Scanberry Screencapture](https://img.bascht.com/2022-03-scanberry/2022-03-24-000251-scanberry.gif)

## Functionality of the web UI
- Trigger the scanning process
- Set a document name and whether it's duplex or simplex
- Prepend a timestamp prefix
- Offer the file as a download
- Copy it to a pre-determined location (e.g. my Nextcloud share)

## Forges

In the spirit of a distributed VCS, this project lives on

- https://git.bascht.space/bascht/scanberry
- https://git.sr.ht/~bascht/scanberry
- https://github.com/bascht/scanberry
