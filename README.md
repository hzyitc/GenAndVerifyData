# GenAndVerifyData
[![GitHub release](https://img.shields.io/github/v/tag/hzyitc/GenAndVerifyData?label=release)](https://github.com/hzyitc/GenAndVerifyData/releases)

[README](README.md) | [中文文档](README_zh.md)

## Introduction

`GenAndVerifyData` is A tool to write and/or verify the file with some "random" data.

It could be used to test disk

## Usage

```
Usage: 
  GenAndVerifyData [-write|-verify] {path} [begin [end]]
    -write    Write only
    -verify   Verify only
    path      Path
    begin     From 0. Align to 4096. Include
    end       From 0. Align to 4096. Not include
```

Example:

Run a test for /dev/sdz:

```
./GenAndVerifyData /dev/sdz
```