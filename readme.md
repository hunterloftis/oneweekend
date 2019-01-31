# Ray Tracing in Golang

This is a chapter-by-chapter progression through the excellent,
free, [raytracing books](https://drive.google.com/drive/folders/14yayBb9XiL16lmuhbYhhvea8mKUUK77W)
by Peter Shirley:

- Raytracing in One Weekend
- Raytracing: the Next Week

The `master` branch goes through the first book and
the `nextweek` branch continues through the second.

## Raytracing in One Weekend

```bash
$ git clone https://github.com/hunterloftis/oneweekend.git
$ cd oneweekend
$ go build ./cmd/trace
$ ./trace > cover.ppm && open cover.ppm
```

![cover image](https://user-images.githubusercontent.com/364501/51394607-bf056180-1b08-11e9-8968-d319697d40ae.png)

## Raytracing: the Next Week

```bash
$ git clone https://github.com/hunterloftis/oneweekend.git
$ cd oneweekend
$ git checkout nextweek
$ go build ./cmd/trace
$ ./trace > final.ppm && open final.ppm
```
