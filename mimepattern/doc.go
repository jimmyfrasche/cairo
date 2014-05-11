//Package mimepattern creates patterns from images that, when possible,
//directly embed the uncompressed image data into the target surface.
//
//Direct embedding is only possible when:
//	- The surface supports embedding (PDF, PS, SVG, and Win32 Printing)
//	- The particular variant of the uncompressed image data is not supported
//	  by that surface. (That is, a surface may support a mime type but not
//	  a particular compression method or extension used in that file).
//	- The pattern is being applied in a way that the surface
//	  cannot directly use the data, such as a compositing operation not natively
//	  supported by that surface.
//
//If direct embedding is not possible, the pattern still works as expected,
//but the image data will be re-encoded.
//If the image is being stored in a lossy format such as JPEG this could
//result in degradation of image quality.
//
//For SVG surfaces, there is the special case.
//Instead of embedding the uncompressed image contents, a url to the image
//is embedded in the document.
package mimepattern
