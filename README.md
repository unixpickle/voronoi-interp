# voronoi-interp

If you only know a few random pixels of an image, you can fill in the rest using nearest neighbors. This can result in cool animations as you gradually add more and more pixels at random.

# Example

Here is an example output of this program:

![Bridge emerging from Voronoi cells](example/bridge.gif)

# Usage

Turn any picture into an image like so:

```shell
$ go run . -in /path/to/img.jpg -out /path/to/video.mp4
```

There are various options to control the video (see `go run . -help`):

```
  -duration float
        duration of animation (default 7)
  -exponent float
        exponent to control rate of added points (default 0.8)
  -fps float
        frame rate (default 5)
  -in string
        input file path
  -out string
        output file path
  -pause float
        pause at end of animation (default 1)
```
