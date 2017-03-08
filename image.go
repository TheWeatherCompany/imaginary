package main

import (
	"encoding/json"
	"errors"

	"gopkg.in/h2non/bimg.v1"
)

// Image stores an image binary buffer and its MIME type
type Image struct {
	Body []byte
	Mime string
}

// Operation implements an image transformation runnable interface
type Operation func([]byte, ImageOptions) (Image, error)

// Run performs the image transformation
func (o Operation) Run(buf []byte, opts ImageOptions) (Image, error) {
	return o(buf, opts)
}

// ImageInfo represents an image details and additional metadata
type ImageInfo struct {
	Width       int    `json:"Width (in pixels) of image area to extract/resize."`
	Height      int    `json:"height"`
	Type        string `json:"type"`
	Space       string `json:"space"`
	Alpha       bool   `json:"hasAlpha"`
	Profile     bool   `json:"hasProfile"`
	Channels    int    `json:"channels"`
	Orientation int    `json:"orientation"`
}

// @Title info
// @Description retrieves orders for given customer defined by customer ID
// @Accept  application/json
// @Produce application/json
// @Success 200 {array}  ImageInfo
// @Failure 400 {object} Error    "Cannot retrieve image medatata"
// @Router /info [get]
func Info(buf []byte, o ImageOptions) (Image, error) {
	// We're not handling an image here, but we reused the struct.
	// An interface will be definitively better here.
	image := Image{Mime: "application/json"}

	meta, err := bimg.Metadata(buf)
	if err != nil {
		return image, NewError("Cannot retrieve image medatata: %s" + err.Error(), BadRequest)
	}

	info := ImageInfo{
		Width:       meta.Size.Width,
		Height:      meta.Size.Height,
		Type:        meta.Type,
		Space:       meta.Space,
		Alpha:       meta.Alpha,
		Profile:     meta.Profile,
		Channels:    meta.Channels,
		Orientation: meta.Orientation,
	}

	body, _ := json.Marshal(info)
	image.Body = body

	return image, nil
}

// @Title resize
// @Description Resize an image by width or height. Image aspect ratio is maintained.
// @Accept  image/*
// @Produce  image/*
// @Param   width       query    int     true        "Width (in pixels) of image area to extract/resize."
// @Param   height      query    int     false        "Height (in pixels) of image area to extract/resize."
// @Param   quality     query    int     false        "JPEG image quality between 1-100. Defaults to `80` (type: 'jpeg' ONLY)"
// @Param   compression query    int     false        "PNG compression level. Default: `6` (type: 'png' ONLY)"
// @Param   type        query    string  false        "Specify the image format to output. Possible values are: `jpeg`, `png` and `webp`"
// @Param   file        query    string  false        "Use image from server local file path. In order to use this you must pass the -mount=<dir> flag (GET only)."
// @Param   url         query    string  false        "Fetch the image from a remove HTTP server. In order to use this you must pass the -enable-url-source flag (GET only)."
// @Param   force       query    bool    false        "Force image transformation size. Default: `false`"
// @Param   rotate      query    int     false        "Image rotation angle. Must be multiple of `90`. Example: `180`"
// @Param   embed       query    bool    false        "Embded"
// @Param   norotation  query    bool    false        "Disable auto rotation based on EXIF orientation. Defaults to `false`"
// @Param   noprofile   query    bool    false        "Disable adding ICC profile metadata. Defaults to `false`"
// @Param   flip        query    bool    false        "Transform the resultant image with flip operation. Default: `false`"
// @Param   flop        query    bool    false        "Transform the resultant image with flop operation. Default: `false`"
// @Param   extend      query    string  false        "Extend represents the image extend mode used when the edges of an image are extended. Allowed values are:`black`, `copy`, `mirror`, `white` and `background`. If background value is specified, you can define the desired extend RGB color via background param, such as ?extend=background&background=250,20,10. For more info, see libvips docs."
// @Param   background  query    string  false        "Background RGB decimal base color to use when flattening transparent PNGs. Example: `255,200,150`"
// @Param   colorspace  query    string  false        "Use a custom color space for the output image. Allowed values are: `srgb` or `bw` (black&white)"
// @Param   gravity     query    string  false        "Gravity *Need to confirm whether allowed?"
// @Param   field       query    string  false        "Form Field. Custom image form field name if using `multipart/form` (POST only). Defaults to: `file`"
// @Success 200 {array}  Image
// @Failure 400 {object} Error   "Some error"
// @Router /resize [get]
func Resize(buf []byte, o ImageOptions) (Image, error) {
	if o.Width == 0 && o.Height == 0 {
		return Image{}, NewError("Missing required param: height or width", BadRequest)
	}

	opts := BimgOptions(o)
	opts.Embed = true

	if o.NoCrop == false {
		opts.Crop = true
	}

	return Process(buf, opts)
}

// @Title enlarge
// @Description Enlarges the image by a given width and height.
// @Accept  image/*
// @Produce  image/*
// @Param   width       query    int     true         "Width (in pixels) of image area to extract/resize."
// @Param   height      query    int     true         "Height (in pixels) of image area to extract/resize."
// @Param   quality     query    int     false        "JPEG image quality between 1-100. Defaults to `80` (type: 'jpeg' ONLY)"
// @Param   compression query    int     false        "PNG compression level. Default: `6` (type: 'png' ONLY)"
// @Param   type        query    string  false        "Specify the image format to output. Possible values are: `jpeg`, `png` and `webp`"
// @Param   file        query    string  false        "Use image from server local file path. In order to use this you must pass the -mount=<dir> flag (GET only)."
// @Param   url         query    string  false        "Fetch the image from a remove HTTP server. In order to use this you must pass the -enable-url-source flag (GET only)."
// @Param   embed       query    bool    false        "Embded"
// @Param   force       query    bool    false        "Force image transformation size. Default: `false`"
// @Param   rotate      query    int     false        "Image rotation angle. Must be multiple of `90`. Example: `180`"
// @Param   norotation  query    bool    false        "Disable auto rotation based on EXIF orientation. Defaults to `false`"
// @Param   noprofile   query    bool    false        "Disable adding ICC profile metadata. Defaults to `false`"
// @Param   flip        query    bool    false        "Transform the resultant image with flip operation. Default: `false`"
// @Param   flop        query    bool    false        "Transform the resultant image with flop operation. Default: `false`"
// @Param   extend      query    string  false        "Extend represents the image extend mode used when the edges of an image are extended. Allowed values are:`black`, `copy`, `mirror`, `white` and `background`. If background value is specified, you can define the desired extend RGB color via background param, such as ?extend=background&background=250,20,10. For more info, see libvips docs."
// @Param   background  query    string  false        "Background RGB decimal base color to use when flattening transparent PNGs. Example: `255,200,150`"
// @Param   colorspace  query    string  false        "Use a custom color space for the output image. Allowed values are: `srgb` or `bw` (black&white)"
// @Param   gravity     query    string  false        "Gravity *Need to confirm whether allowed?"
// @Param   field       query    string  false        "Form Field. Custom image form field name if using `multipart/form` (POST only). Defaults to: `file`"
// @Success 200 {array}  Image
// @Failure 400 {object} Error   "Some error"
// @Router /enlarge [get]
func Enlarge(buf []byte, o ImageOptions) (Image, error) {
	if o.Width == 0 || o.Height == 0 {
		return Image{}, NewError("Missing required params: height, width", BadRequest)
	}

	opts := BimgOptions(o)
	opts.Enlarge = true

	if o.NoCrop == false {
		opts.Crop = true
	}

	return Process(buf, opts)
}

// @Title extract
// @Description Enlarges the image by a given width and height.
// @Accept  image/*, multipart/form-data
// @Produce  image/*
// @Param   top         query    int     true         "Top edge of area to extract (in pixels). Example: `100`"
// @Param   left        query    int     false        "Left edge of area to extract (in pixels). Example: `100`"
// @Param   areawidth   query    int     true         "Height area to extract (in pixels). Example: `100`"
// @Param   areaheight  query    int     false        "Width area to extract (in pixels). Example: `100`"
// @Param   width       query    int     false        "Width (in pixels) of image area to extract/resize."
// @Param   height      query    int     false        "Height (in pixels) of image area to extract/resize."
// @Param   quality     query    int     false        "JPEG image quality between 1-100. Defaults to `80` (type: 'jpeg' ONLY)"
// @Param   compression query    int     false        "PNG compression level. Default: `6` (type: 'png' ONLY)"
// @Param   type        query    string  false        "Specify the image format to output. Possible values are: `jpeg`, `png` and `webp`"
// @Param   file        query    string  false        "Use image from server local file path. In order to use this you must pass the -mount=<dir> flag (GET only)."
// @Param   url         query    string  false        "Fetch the image from a remove HTTP server. In order to use this you must pass the -enable-url-source flag (GET only)."
// @Param   embed       query    bool    false        "Embded"
// @Param   force       query    bool    false        "Force image transformation size. Default: `false`"
// @Param   rotate      query    int     false        "Image rotation angle. Must be multiple of `90`. Example: `180`"
// @Param   norotation  query    bool    false        "Disable auto rotation based on EXIF orientation. Defaults to `false`"
// @Param   noprofile   query    bool    false        "Disable adding ICC profile metadata. Defaults to `false`"
// @Param   flip        query    bool    false        "Transform the resultant image with flip operation. Default: `false`"
// @Param   flop        query    bool    false        "Transform the resultant image with flop operation. Default: `false`"
// @Param   extend      query    string  false        "Extend represents the image extend mode used when the edges of an image are extended. Allowed values are:`black`, `copy`, `mirror`, `white` and `background`. If background value is specified, you can define the desired extend RGB color via background param, such as ?extend=background&background=250,20,10. For more info, see libvips docs."
// @Param   background  query    string  false        "Background RGB decimal base color to use when flattening transparent PNGs. Example: `255,200,150`"
// @Param   colorspace  query    string  false        "Use a custom color space for the output image. Allowed values are: `srgb` or `bw` (black&white)"
// @Param   gravity     query    string  false        "Gravity *Need to confirm whether allowed?"
// @Param   field       query    string  false        "Form Field. Custom image form field name if using `multipart/form` (POST only). Defaults to: `file`"
// @Success 200 {array}  Image
// @Failure 400 {object} Error   "Some error"
// @Router /extract [get]
func Extract(buf []byte, o ImageOptions) (Image, error) {
	if o.AreaWidth == 0 || o.AreaHeight == 0 {
		return Image{}, NewError("Missing required params: areawidth or areaheight", BadRequest)
	}

	opts := BimgOptions(o)
	opts.Top = o.Top
	opts.Left = o.Left
	opts.AreaWidth = o.AreaWidth
	opts.AreaHeight = o.AreaHeight

	return Process(buf, opts)
}

// @Title crop
// @Description Crop the image by a given width or height. Image ratio is maintained.
// @Accept  image/*, multipart/form-data
// @Produce  image/*
// @Param   width       query    int     false        "Width (in pixels) of image area to extract/resize."
// @Param   height      query    int     false        "Height (in pixels) of image area to extract/resize."
// @Param   quality     query    int     false        "JPEG image quality between 1-100. Defaults to `80` (type: 'jpeg' ONLY)"
// @Param   compression query    int     false        "PNG compression level. Default: `6` (type: 'png' ONLY)"
// @Param   type        query    string  false        "Specify the image format to output. Possible values are: `jpeg`, `png` and `webp`"
// @Param   file        query    string  false        "Use image from server local file path. In order to use this you must pass the -mount=<dir> flag (GET only)."
// @Param   url         query    string  false        "Fetch the image from a remove HTTP server. In order to use this you must pass the -enable-url-source flag (GET only)."
// @Param   force       query    bool    false        "Force image transformation size. Default: `false`"
// @Param   rotate      query    int     false        "Image rotation angle. Must be multiple of `90`. Example: `180`"
// @Param   embed       query    bool    false        "Embded"
// @Param   norotation  query    bool    false        "Disable auto rotation based on EXIF orientation. Defaults to `false`"
// @Param   noprofile   query    bool    false        "Disable adding ICC profile metadata. Defaults to `false`"
// @Param   flip        query    bool    false        "Transform the resultant image with flip operation. Default: `false`"
// @Param   flop        query    bool    false        "Transform the resultant image with flop operation. Default: `false`"
// @Param   extend      query    string  false        "Extend represents the image extend mode used when the edges of an image are extended. Allowed values are:`black`, `copy`, `mirror`, `white` and `background`. If background value is specified, you can define the desired extend RGB color via background param, such as ?extend=background&background=250,20,10. For more info, see libvips docs."
// @Param   background  query    string  false        "Background RGB decimal base color to use when flattening transparent PNGs. Example: `255,200,150`"
// @Param   colorspace  query    string  false        "Use a custom color space for the output image. Allowed values are: `srgb` or `bw` (black&white)"
// @Param   gravity     query    string  false        "Define the crop operation gravity. Supported values are: `north`, `south`, `centre`, `west` and `east`. Defaults to `centre`"
// @Param   field       query    string  false        "Form Field. Custom image form field name if using `multipart/form` (POST only). Defaults to: `file`"
// @Success 200 {array}  Image
// @Failure 400 {object} Error   "Customer ID must be specified"
// @Router /crop [get]
func Crop(buf []byte, o ImageOptions) (Image, error) {
	if o.Width == 0 && o.Height == 0 {
		return Image{}, NewError("Missing required param: height or width", BadRequest)
	}

	opts := BimgOptions(o)
	opts.Crop = true
	return Process(buf, opts)
}

// @Title rotate
// @Description Rotates the image (with auto-rotate based on EXIF orientation).
// @Accept  image/*, multipart/form-data
// @Produce  image/*
// @Param   rotate      query    int     true         "Rotation degrees"
// @Param   width       query    int     false        "Width (in pixels) of image area to extract/resize."
// @Param   height      query    int     false        "Height (in pixels) of image area to extract/resize."
// @Param   quality     query    int     false        "JPEG image quality between 1-100. Defaults to `80` (type: 'jpeg' ONLY)"
// @Param   compression query    int     false        "PNG compression level. Default: `6` (type: 'png' ONLY)"
// @Param   type        query    string  false        "Specify the image format to output. Possible values are: `jpeg`, `png` and `webp`"
// @Param   file        query    string  false        "Use image from server local file path. In order to use this you must pass the -mount=<dir> flag (GET only)."
// @Param   url         query    string  false        "Fetch the image from a remove HTTP server. In order to use this you must pass the -enable-url-source flag (GET only)."
// @Param   embed       query    bool    false        "Embded"
// @Param   force       query    bool    false        "Force image transformation size. Default: `false`"
// @Param   rotate      query    int     false        "Image rotation angle. Must be multiple of `90`. Example: `180`"
// @Param   norotation  query    bool    false        "Disable auto rotation based on EXIF orientation. Defaults to `false`"
// @Param   noprofile   query    bool    false        "Disable adding ICC profile metadata. Defaults to `false`"
// @Param   flip        query    bool    false        "Transform the resultant image with flip operation. Default: `false`"
// @Param   flop        query    bool    false        "Transform the resultant image with flop operation. Default: `false`"
// @Param   extend      query    string  false        "Extend represents the image extend mode used when the edges of an image are extended. Allowed values are:`black`, `copy`, `mirror`, `white` and `background`. If background value is specified, you can define the desired extend RGB color via background param, such as ?extend=background&background=250,20,10. For more info, see libvips docs."
// @Param   background  query    string  false        "Background RGB decimal base color to use when flattening transparent PNGs. Example: `255,200,150`"
// @Param   colorspace  query    string  false        "Use a custom color space for the output image. Allowed values are: `srgb` or `bw` (black&white)"
// @Param   gravity     query    string  false        "Gravity *Need to confirm whether allowed?"
// @Param   field       query    string  false        "Form Field. Custom image form field name if using `multipart/form` (POST only). Defaults to: `file`"
// @Success 200 {array}  Image
// @Failure 400 {object} Error   "Some error"
// @Router /rotate [get]
func Rotate(buf []byte, o ImageOptions) (Image, error) {
	if o.Rotate == 0 {
		return Image{}, NewError("Missing required param: rotate", BadRequest)
	}

	opts := BimgOptions(o)
	return Process(buf, opts)
}

// @Title flip
// @Description Flips the image horizontally (with auto-flip based on EXIF metadata).
// @Accept  image/*, multipart/form-data
// @Produce  image/*
// @Param   width       query    int     false        "Width (in pixels) of image area to extract/resize."
// @Param   height      query    int     false        "Height (in pixels) of image area to extract/resize."
// @Param   quality     query    int     false        "JPEG image quality between 1-100. Defaults to `80` (type: 'jpeg' ONLY)"
// @Param   compression query    int     false        "PNG compression level. Default: `6` (type: 'png' ONLY)"
// @Param   type        query    string  false        "Specify the image format to output. Possible values are: `jpeg`, `png` and `webp`"
// @Param   file        query    string  false        "Use image from server local file path. In order to use this you must pass the -mount=<dir> flag (GET only)."
// @Param   url         query    string  false        "Fetch the image from a remove HTTP server. In order to use this you must pass the -enable-url-source flag (GET only)."
// @Param   force       query    bool    false        "Force image transformation size. Default: `false`"
// @Param   rotate      query    int     false        "Image rotation angle. Must be multiple of `90`. Example: `180`"
// @Param   embed       query    bool    false        "Embded"
// @Param   norotation  query    bool    false        "Disable auto rotation based on EXIF orientation. Defaults to `false`"
// @Param   noprofile   query    bool    false        "Disable adding ICC profile metadata. Defaults to `false`"
// @Param   flip        query    bool    false        "Transform the resultant image with flip operation. Default: `false`"
// @Param   flop        query    bool    false        "Transform the resultant image with flop operation. Default: `false`"
// @Param   extend      query    string  false        "Extend represents the image extend mode used when the edges of an image are extended. Allowed values are:`black`, `copy`, `mirror`, `white` and `background`. If background value is specified, you can define the desired extend RGB color via background param, such as ?extend=background&background=250,20,10. For more info, see libvips docs."
// @Param   background  query    string  false        "Background RGB decimal base color to use when flattening transparent PNGs. Example: `255,200,150`"
// @Param   colorspace  query    string  false        "Use a custom color space for the output image. Allowed values are: `srgb` or `bw` (black&white)"
// @Param   gravity     query    string  false        "Define the crop operation gravity. Supported values are: `north`, `south`, `centre`, `west` and `east`. Defaults to `centre`"
// @Param   field       query    string  false        "Form Field. Custom image form field name if using `multipart/form` (POST only). Defaults to: `file`"
// @Success 200 {array}  Image
// @Failure 400 {object} Error   "Customer ID must be specified"
// @Router /flip [get]
func Flip(buf []byte, o ImageOptions) (Image, error) {
	opts := BimgOptions(o)
	opts.Flip = true
	return Process(buf, opts)
}

// @Title flop
// @Description Flips/Flops the image vertically.
// @Accept  image/*, multipart/form-data
// @Produce  image/*
// @Param   width       query    int     false        "Width (in pixels) of image area to extract/resize."
// @Param   height      query    int     false        "Height (in pixels) of image area to extract/resize."
// @Param   quality     query    int     false        "JPEG image quality between 1-100. Defaults to `80` (type: 'jpeg' ONLY)"
// @Param   compression query    int     false        "PNG compression level. Default: `6` (type: 'png' ONLY)"
// @Param   type        query    string  false        "Specify the image format to output. Possible values are: `jpeg`, `png` and `webp`"
// @Param   file        query    string  false        "Use image from server local file path. In order to use this you must pass the -mount=<dir> flag (GET only)."
// @Param   url         query    string  false        "Fetch the image from a remove HTTP server. In order to use this you must pass the -enable-url-source flag (GET only)."
// @Param   force       query    bool    false        "Force image transformation size. Default: `false`"
// @Param   rotate      query    int     false        "Image rotation angle. Must be multiple of `90`. Example: `180`"
// @Param   embed       query    bool    false        "Embded"
// @Param   norotation  query    bool    false        "Disable auto rotation based on EXIF orientation. Defaults to `false`"
// @Param   noprofile   query    bool    false        "Disable adding ICC profile metadata. Defaults to `false`"
// @Param   flip        query    bool    false        "Transform the resultant image with flip operation. Default: `false`"
// @Param   flop        query    bool    false        "Transform the resultant image with flop operation. Default: `false`"
// @Param   extend      query    string  false        "Extend represents the image extend mode used when the edges of an image are extended. Allowed values are:`black`, `copy`, `mirror`, `white` and `background`. If background value is specified, you can define the desired extend RGB color via background param, such as ?extend=background&background=250,20,10. For more info, see libvips docs."
// @Param   background  query    string  false        "Background RGB decimal base color to use when flattening transparent PNGs. Example: `255,200,150`"
// @Param   colorspace  query    string  false        "Use a custom color space for the output image. Allowed values are: `srgb` or `bw` (black&white)"
// @Param   gravity     query    string  false        "Define the crop operation gravity. Supported values are: `north`, `south`, `centre`, `west` and `east`. Defaults to `centre`"
// @Param   field       query    string  false        "Form Field. Custom image form field name if using `multipart/form` (POST only). Defaults to: `file`"
// @Success 200 {array}  Image
// @Failure 400 {object} Error   "Customer ID must be specified"
// @Router /flop [get]
func Flop(buf []byte, o ImageOptions) (Image, error) {
	opts := BimgOptions(o)
	opts.Flop = true
	return Process(buf, opts)
}

// @Title thumbnail
// @Description Create a thumbnail.
// @Accept  image/*, multipart/form-data
// @Produce  image/*
// @Param   width       query    int     false        "Width (in pixels) of image area to extract/resize."
// @Param   height      query    int     false        "Height (in pixels) of image area to extract/resize."
// @Param   quality     query    int     false        "JPEG image quality between 1-100. Defaults to `80` (type: 'jpeg' ONLY)"
// @Param   compression query    int     false        "PNG compression level. Default: `6` (type: 'png' ONLY)"
// @Param   type        query    string  false        "Specify the image format to output. Possible values are: `jpeg`, `png` and `webp`"
// @Param   file        query    string  false        "Use image from server local file path. In order to use this you must pass the -mount=<dir> flag (GET only)."
// @Param   url         query    string  false        "Fetch the image from a remove HTTP server. In order to use this you must pass the -enable-url-source flag (GET only)."
// @Param   embed       query    bool    false        "Embded"
// @Param   force       query    bool    false        "Force image transformation size. Default: `false`"
// @Param   rotate      query    int     false        "Image rotation angle. Must be multiple of `90`. Example: `180`"
// @Param   norotation  query    bool    false        "Disable auto rotation based on EXIF orientation. Defaults to `false`"
// @Param   noprofile   query    bool    false        "Disable adding ICC profile metadata. Defaults to `false`"
// @Param   flip        query    bool    false        "Transform the resultant image with flip operation. Default: `false`"
// @Param   flop        query    bool    false        "Transform the resultant image with flop operation. Default: `false`"
// @Param   extend      query    string  false        "Extend represents the image extend mode used when the edges of an image are extended. Allowed values are:`black`, `copy`, `mirror`, `white` and `background`. If background value is specified, you can define the desired extend RGB color via background param, such as ?extend=background&background=250,20,10. For more info, see libvips docs."
// @Param   background  query    string  false        "Background RGB decimal base color to use when flattening transparent PNGs. Example: `255,200,150`"
// @Param   colorspace  query    string  false        "Use a custom color space for the output image. Allowed values are: `srgb` or `bw` (black&white)"
// @Param   gravity     query    string  false        "Gravity *Need to confirm whether allowed?"
// @Param   field       query    string  false        "Form Field. Custom image form field name if using `multipart/form` (POST only). Defaults to: `file`"
// @Success 200 {array}  Image
// @Failure 400 {object} Error   "Some error"
// @Router /thumbnail [get]
func Thumbnail(buf []byte, o ImageOptions) (Image, error) {
	if o.Width == 0 && o.Height == 0 {
		return Image{}, NewError("Missing required params: width or height", BadRequest)
	}

	return Process(buf, BimgOptions(o))
}

// @Title zoom
// @Description Zooms into the image.
// @Accept  image/*, multipart/form-data
// @Produce  image/*
// @Param   factor      query    float32 true         "Zoom factor level. Example: `2`"
// @Param   width       query    int     false        "Width (in pixels) of image area to extract/resize."
// @Param   height      query    int     false        "Height (in pixels) of image area to extract/resize."
// @Param   quality     query    int     false        "JPEG image quality between 1-100. Defaults to `80` (type: 'jpeg' ONLY)"
// @Param   compression query    int     false        "PNG compression level. Default: `6` (type: 'png' ONLY)"
// @Param   type        query    string  false        "Specify the image format to output. Possible values are: `jpeg`, `png` and `webp`"
// @Param   file        query    string  false        "Use image from server local file path. In order to use this you must pass the -mount=<dir> flag (GET only)."
// @Param   url         query    string  false        "Fetch the image from a remove HTTP server. In order to use this you must pass the -enable-url-source flag (GET only)."
// @Param   embed       query    bool    false        "Embded"
// @Param   force       query    bool    false        "Force image transformation size. Default: `false`"
// @Param   rotate      query    int     false        "Image rotation angle. Must be multiple of `90`. Example: `180`"
// @Param   norotation  query    bool    false        "Disable auto rotation based on EXIF orientation. Defaults to `false`"
// @Param   noprofile   query    bool    false        "Disable adding ICC profile metadata. Defaults to `false`"
// @Param   flip        query    bool    false        "Transform the resultant image with flip operation. Default: `false`"
// @Param   flop        query    bool    false        "Transform the resultant image with flop operation. Default: `false`"
// @Param   extend      query    string  false        "Extend represents the image extend mode used when the edges of an image are extended. Allowed values are:`black`, `copy`, `mirror`, `white` and `background`. If background value is specified, you can define the desired extend RGB color via background param, such as ?extend=background&background=250,20,10. For more info, see libvips docs."
// @Param   background  query    string  false        "Background RGB decimal base color to use when flattening transparent PNGs. Example: `255,200,150`"
// @Param   colorspace  query    string  false        "Use a custom color space for the output image. Allowed values are: `srgb` or `bw` (black&white)"
// @Param   gravity     query    string  false        "Gravity *Need to confirm whether allowed?"
// @Param   field       query    string  false        "Form Field. Custom image form field name if using `multipart/form` (POST only). Defaults to: `file`"
// @Success 200 {array}  Image
// @Failure 400 {object} Error   "Some error"
// @Router /zoom [get]
func Zoom(buf []byte, o ImageOptions) (Image, error) {
	if o.Factor == 0 {
		return Image{}, NewError("Missing required param: factor", BadRequest)
	}

	opts := BimgOptions(o)

	if o.Top > 0 || o.Left > 0 {
		if o.AreaWidth == 0 && o.AreaHeight == 0 {
			return Image{}, NewError("Missing required params: areawidth, areaheight", BadRequest)
		}

		opts.Top = o.Top
		opts.Left = o.Left
		opts.AreaWidth = o.AreaWidth
		opts.AreaHeight = o.AreaHeight

		if o.NoCrop == false {
			opts.Crop = true
		}
	}

	opts.Zoom = o.Factor
	return Process(buf, opts)
}

// @Title convert
// @Description Converts an image from one type/format to another with additional quality/compression settings.
// @Accept  image/*, multipart/form-data
// @Produce  image/*
// @Param   type        query    float32 true         "Specify the image format to output. Possible values are: `jpeg`, `png` and `webp`"
// @Param   width       query    int     false        "Width (in pixels) of image area to extract/resize."
// @Param   height      query    int     false        "Height (in pixels) of image area to extract/resize."
// @Param   quality     query    int     false        "JPEG image quality between 1-100. Defaults to `80` (type: 'jpeg' ONLY)"
// @Param   compression query    int     false        "PNG compression level. Default: `6` (type: 'png' ONLY)"
// @Param   file        query    string  false        "Use image from server local file path. In order to use this you must pass the -mount=<dir> flag (GET only)."
// @Param   url         query    string  false        "Fetch the image from a remove HTTP server. In order to use this you must pass the -enable-url-source flag (GET only)."
// @Param   embed       query    bool    false        "Embded"
// @Param   force       query    bool    false        "Force image transformation size. Default: `false`"
// @Param   rotate      query    int     false        "Image rotation angle. Must be multiple of `90`. Example: `180`"
// @Param   norotation  query    bool    false        "Disable auto rotation based on EXIF orientation. Defaults to `false`"
// @Param   noprofile   query    bool    false        "Disable adding ICC profile metadata. Defaults to `false`"
// @Param   flip        query    bool    false        "Transform the resultant image with flip operation. Default: `false`"
// @Param   flop        query    bool    false        "Transform the resultant image with flop operation. Default: `false`"
// @Param   extend      query    string  false        "Extend represents the image extend mode used when the edges of an image are extended. Allowed values are:`black`, `copy`, `mirror`, `white` and `background`. If background value is specified, you can define the desired extend RGB color via background param, such as ?extend=background&background=250,20,10. For more info, see libvips docs."
// @Param   background  query    string  false        "Background RGB decimal base color to use when flattening transparent PNGs. Example: `255,200,150`"
// @Param   colorspace  query    string  false        "Use a custom color space for the output image. Allowed values are: `srgb` or `bw` (black&white)"
// @Param   gravity     query    string  false        "Gravity *Need to confirm whether allowed?"
// @Param   field       query    string  false        "Form Field. Custom image form field name if using `multipart/form` (POST only). Defaults to: `file`"
// @Success 200 {array}  Image
// @Failure 400 {object} Error   "Some error"
// @Router /convert [get]
func Convert(buf []byte, o ImageOptions) (Image, error) {
	if o.Type == "" {
		return Image{}, NewError("Missing required param: type", BadRequest)
	}
	if ImageType(o.Type) == bimg.UNKNOWN {
		return Image{}, NewError("Invalid image type: " + o.Type, BadRequest)
	}
	opts := BimgOptions(o)

	return Process(buf, opts)
}

// @Title watermark
// @Description Adds a custom watermark text to an image.
// @Accept  image/*, multipart/form-data
// @Produce  image/*
// @Param   text        query    string  true         "Watermark text content. Example: `copyright (c) 2189`"
// @Param   margin      query    int     false        "Text area margin for watermark. Example: `50`"
// @Param   dpi         query    int     false        "DPI value for watermark. Example: `150`"
// @Param   textwidth   query    int     false        "Text area width for watermark. Example: `200`"
// @Param   opacity     query    float32 false        "Opacity level for watermark text. Default: `0.2`"
// @Param   noreplicate query    float32 false        "Disable text replication in watermark. Defaults to `false`"
// @Param   font        query    string  false        "Watermark text font type and format. Example: `sans bold 12`"
// @Param   color       query    string  false        "Watermark text RGB decimal base color. Example: `255,200,150`"
// @Param   quality     query    int     false        "JPEG image quality between 1-100. Defaults to `80` (type: 'jpeg' ONLY)"
// @Param   compression query    int     false        "PNG compression level. Default: `6` (type: 'png' ONLY)"
// @Param   type        query    string  false        "Specify the image format to output. Possible values are: `jpeg`, `png` and `webp`"
// @Param   file        query    string  false        "Use image from server local file path. In order to use this you must pass the -mount=<dir> flag (GET only)."
// @Param   url         query    string  false        "Fetch the image from a remove HTTP server. In order to use this you must pass the -enable-url-source flag (GET only)."
// @Param   embed       query    bool    false        "Embded"
// @Param   force       query    bool    false        "Force image transformation size. Default: `false`"
// @Param   rotate      query    int     false        "Image rotation angle. Must be multiple of `90`. Example: `180`"
// @Param   norotation  query    bool    false        "Disable auto rotation based on EXIF orientation. Defaults to `false`"
// @Param   noprofile   query    bool    false        "Disable adding ICC profile metadata. Defaults to `false`"
// @Param   flip        query    bool    false        "Transform the resultant image with flip operation. Default: `false`"
// @Param   flop        query    bool    false        "Transform the resultant image with flop operation. Default: `false`"
// @Param   extend      query    string  false        "Extend represents the image extend mode used when the edges of an image are extended. Allowed values are:`black`, `copy`, `mirror`, `white` and `background`. If background value is specified, you can define the desired extend RGB color via background param, such as ?extend=background&background=250,20,10. For more info, see libvips docs."
// @Param   background  query    string  false        "Background RGB decimal base color to use when flattening transparent PNGs. Example: `255,200,150`"
// @Param   colorspace  query    string  false        "Use a custom color space for the output image. Allowed values are: `srgb` or `bw` (black&white)"
// @Param   field       query    string  false        "Form Field. Custom image form field name if using `multipart/form` (POST only). Defaults to: `file`"
// @Param   width       query    int     false        "Width *Need to confirm whether allowed?"
// @Param   height      query    int     false        "Height *Need to confirm whether allowed?"
// @Param   gravity     query    string  false        "Gravity *Need to confirm whether allowed?"
// @Success 200 {array}  Image
// @Failure 400 {object} Error   "Some error"
// @Router /watermark [get]
func Watermark(buf []byte, o ImageOptions) (Image, error) {
	if o.Text == "" {
		return Image{}, NewError("Missing required param: text", BadRequest)
	}

	opts := BimgOptions(o)
	opts.Watermark.DPI = o.DPI
	opts.Watermark.Text = o.Text
	opts.Watermark.Font = o.Font
	opts.Watermark.Margin = o.Margin
	opts.Watermark.Width = o.TextWidth
	opts.Watermark.Opacity = o.Opacity
	opts.Watermark.NoReplicate = o.NoReplicate

	if len(o.Color) > 2 {
		opts.Watermark.Background = bimg.Color{o.Color[0], o.Color[1], o.Color[2]}
	}

	return Process(buf, opts)
}

func Process(buf []byte, opts bimg.Options) (out Image, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch value := r.(type) {
			case error:
				err = value
			case string:
				err = errors.New(value)
			default:
				err = errors.New("libvips internal error")
			}
			out = Image{}
		}
	}()

	buf, err = bimg.Resize(buf, opts)
	if err != nil {
		return Image{}, err
	}

	mime := GetImageMimeType(bimg.DetermineImageType(buf))
	return Image{Body: buf, Mime: mime}, nil
}
