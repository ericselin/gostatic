# `gostatic` - Build static sites with docker

This is a golang package and accompanying CLI tool for pulling (or cloning) a git repository and building the contained static site using an arbitrary docker image.

Usage:

```golang
gostatic.Build(
  "https://github.com/ericselin/ericselin",
  "example",
  "denoland/deno",
  []string{"run", "-A", "bob.ts"},
  ctx,
)
```
