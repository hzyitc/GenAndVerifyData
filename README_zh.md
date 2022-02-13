
# GenAndVerifyData
[![GitHub release](https://img.shields.io/github/v/tag/hzyitc/GenAndVerifyData?label=release)](https://github.com/hzyitc/GenAndVerifyData/releases)

[README](README.md) | [中文文档](README_zh.md)

## 介绍

`GenAndVerifyData`是一个使用“随机”数据来写入并/或校验指定文件的工具。

可用来测试硬盘。

## 使用指南

```
Usage: 
  GenAndVerifyData [-write|-verify] {path} [begin [end]]
    -write    仅写
    -verify   仅校验
    path      路径
    begin     始于0. 对齐4096. 包含该值
    end       始于0. 对齐4096. 不包含该值
```

对/dev/sdz运行一次测试:

```
./GenAndVerifyData /dev/sdz
```