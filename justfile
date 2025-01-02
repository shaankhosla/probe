#!/usr/bin/env -S just --justfile

tag version:
    git tag -a {{version}} -m "Release {{version}}"
    git push origin {{version}}
