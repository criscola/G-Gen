package consts

/**
 * Session (session name, keys...)
 */
const (
	SessionName          = "UserSession"
	SessionGeneratorJob  = "SessionGeneratorJob"
)

/**
 * HTTP REST
 */
const (
	HttpContentType   					= "Content-Type"
	HttpContentLength 					= "Content-Length"
	HttpContentDisposition				= "Content-Disposition"
	HttpMimeTextPlain 					= "text/plain"
	HttpMimeApplicationOctetStream  	= "application/octet-stream"
)

/**
 * Form values keys
 */
const (
	FormImage       	= "image"
	FormScaleFactor 	= "scaleFactor"
	FormModelThickness 	= "modelThickness"
	FormTravelSpeed 	= "travelSpeed"
	FormGCodeDialect	= "gcodeDialect"
	FormRepRap			= "RepRap"
	FormUltimaker		= "Ultimaker"
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