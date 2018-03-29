package consts

/**
 * Session (session name, keys...)
 */
const (
	SessionName          = "UserSession"
	//SessionImageFilename = "SessionImageFilename"
	SessionGeneratorJob  = "SessionGeneratorJob"
)

/**
 * HTTP REST
 */
const (
	HttpContentType   = "Content-Type"
	HttpContentLength = "Content-Length"
	HttpContentDisposition = "Content-Disposition"
	HttpMimeTextPlain = "text/plain"
	HttpMimeImageJpeg = "image/jpeg"
	HttpMimeImagePng  = "image/png"
	HttpMimeImageGif  				= "image/gif"
	HttpMimeApplicationOctetStream  = "application/octet-stream"
)

/**
 * Form values keys
 */
const (
	FormJobId			= "jobId"
	FormImage       	= "image"
	FormScaleFactor 	= "scaleFactor"
	FormModelThickness 	= "modelThickness"
	FormTravelSpeed 	= "travelSpeed"
)

/**
 * Http requests parameters
 */
const (
	RequestFilename = "filename"
	RequestJobId    = "jobId"
)

/**
 * General default values
 */
 const (
 	DefaultImageExtension = "png"
 	DefaultCookiesMaxAge  = 259200 // 3 days in seconds
 )