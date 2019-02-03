# Ray Tracing in Go

This is a chapter-by-chapter progression through the excellent
free [ray-tracing books](https://drive.google.com/drive/folders/14yayBb9XiL16lmuhbYhhvea8mKUUK77W)
by Peter Shirley:

- [Ray Tracing in One Weekend](https://github.com/hunterloftis/oneweekend/tree/oneweekend#ray-tracing-in-one-weekend)
- [Ray Tracing: the Next Week](https://github.com/hunterloftis/oneweekend/tree/master#ray-tracing-the-next-week)

[![GoDoc](https://godoc.org/github.com/hunterloftis/oneweekend/oneweekend?status.svg)](https://godoc.org/github.com/hunterloftis/oneweekend)

## Who is this for?

If you're interested in graphics and ray tracing,
this is a working example of a simple, easy-to-read ray tracer written in Go.
It is built up piece-by-piece in concert with the chapters of the
[original C++ books](https://drive.google.com/drive/folders/14yayBb9XiL16lmuhbYhhvea8mKUUK77W)
by Peter Shirley.

If you're interested in Go,
this is a fun, visual way to explore the language.
It's [fully documented](https://godoc.org/github.com/hunterloftis/oneweekend)
and easy to change in order to create your own ray traced images.

## Ray Tracing in One Weekend

```bash
$ git clone https://github.com/hunterloftis/oneweekend.git
$ cd oneweekend
$ git checkout oneweekend
$ go build ./cmd/trace
$ ./trace > cover.ppm && open cover.ppm
```

![cover image](https://user-images.githubusercontent.com/364501/51394607-bf056180-1b08-11e9-8968-d319697d40ae.png)

## Ray Tracing: the Next Week

```bash
$ git clone https://github.com/hunterloftis/oneweekend.git
$ cd oneweekend
$ go build ./cmd/trace
$ ./trace > cover.ppm && open cover.ppm
```

![cover image](https://user-images.githubusercontent.com/364501/52127550-5afe9500-2600-11e9-8c12-70b1aaae2e1d.png)
