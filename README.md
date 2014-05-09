#cairo [![GoDoc](https://godoc.org/github.com/jimmyfrasche/cairo?status.png)](https://godoc.org/github.com/jimmyfrasche/cairo)
Package cairo wraps libcairo, a 2D graphics library with support for multiple output devices.

Download:
```shell
go get github.com/jimmyfrasche/cairo
```

* * *
Package cairo wraps libcairo, a 2D graphics library with support for multiple
output devices.
Libcairo is is designed to produce consistent output on all output media
while taking advantage of display hardware acceleration when available.

The cairo API provides operations similar to the drawing operators
of PostScript and PDF.
Operations in cairo including stroking and filling cubic Bézier splines,
transforming and compositing translucent images, and antialiased text
rendering.
All drawing operations can be transformed by any affine transformation
(scale, rotation, shear, etc.)

##Reference semantics
Many types in this and related packages are references, that behave much
like file handles, to underlying libcairo objects.
These handles are values that can be copied at will, but the libcairo
object referenced must be freed explicitly with Close (or Unmap,
in the case of MappedImageSurface).

##Naming Conventions
Cairo refers to this package and its related packages.
Libcairo refers to the C library that this package is a binding to.

##Libcairo version
This package requires libcairo version 1.12 or greater.
Libcairo must be compiled with:

```
CAIRO_HAS_IMAGE_SURFACE
```

Related packages, such as cairo/ps, may require further options compiled
in to libcairo, but they will be documented.

Libcairo can be found at http://cairographics.org .

##Xtensions
Many types, functions, and methods are prefixed by Xtension.
You may ignore these unless you are writing an extension.
An extension is a package that integrates another portion of libcairo or
binds with a library that supports libcairo integration.



* * *
Automatically generated by [autoreadme](https://github.com/jimmyfrasche/autoreadme) on 2014.05.08
