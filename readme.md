# Ray Tracing in One Weekend - Golang

This is a chapter-by-chapter progression through the excellent,
free, [Raytracing in One Weekend](https://drive.google.com/drive/folders/14yayBb9XiL16lmuhbYhhvea8mKUUK77W)
book by Peter Shirley.

Final render:

![cover image](https://user-images.githubusercontent.com/364501/51394607-bf056180-1b08-11e9-8968-d319697d40ae.png)

```bash
$ git clone https://github.com/hunterloftis/oneweekend.git
$ cd oneweekend
$ go build ./cmd/trace
$ ./trace > cover.ppm
$ open cover.ppm
```