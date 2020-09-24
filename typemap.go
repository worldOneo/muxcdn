package main

// The Http Content type of this
const (
	TypeAACaudio                           = "audio/aac"
	TypeAbiWorddocument                    = "application/x-abiword"
	TypeArchivedocument                    = "application/x-freearc"
	TypeAVI                                = "video/x-msvideo"
	TypeAmazonKindleeBookformat            = "application/vnd.amazon.ebook"
	TypeAnykindofbinarydata                = "application/octet-stream"
	Type2BitmapGraphics                    = "image/bmp"
	TypeBZiparchive                        = "application/x-bzip"
	TypeBZip2archive                       = "application/x-bzip2"
	TypeCShellscript                       = "application/x-csh"
	TypeCascadingStyleSheets               = "text/css"
	TypeCommaSeparatedvalues               = "text/csv"
	TypeMicrosoftWord                      = "application/msword"
	TypeMicrosoftWordOpenXML               = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	TypeMSEmbeddedOpenTypefonts            = "application/vnd.ms-fontobject"
	TypeElectronicpublication              = "application/epub+zip"
	TypeGZipCompressedArchive              = "application/gzip"
	TypeGraphicsInterchangeFormat          = "image/gif"
	TypeHyperTextMarkupLanguage            = "text/html"
	TypeIconformat                         = "image/vnd.microsoft.icon"
	TypeiCalendarformat                    = "text/calendar"
	TypeJavaArchive                        = "application/java-archive"
	TypeJPEGimages                         = "image/jpeg"
	TypeJavaScript                         = "text/javascript, per the following specifications:"
	TypeJSONformat                         = "application/json"
	TypeJSONLDformat                       = "application/ld+json"
	TypeMusicalInstrumentDigitalInterface  = "audio/midi"
	TypeJavaScriptmodule                   = "text/javascript"
	TypeMP3audio                           = "audio/mpeg"
	TypeMPEGVideo                          = "video/mpeg"
	TypeAppleInstallerPackage              = "application/vnd.apple.installer+xml"
	TypeOpenDocumentpresentationdocument   = "application/vnd.oasis.opendocument.presentation"
	TypeOpenDocumentspreadsheetdocument    = "application/vnd.oasis.opendocument.spreadsheet"
	TypeOpenDocumenttextdocument           = "application/vnd.oasis.opendocument.text"
	TypeOGGaudio                           = "audio/ogg"
	TypeOGGvideo                           = "video/ogg"
	TypeOGG                                = "application/ogg"
	TypeOpusaudio                          = "audio/opus"
	TypeOpenTypefont                       = "font/otf"
	TypePortableNetworkGraphics            = "image/png"
	TypeAdobePortableDocumentFormat        = "application/pdf"
	TypeHypertextPreprocessor              = "application/x-httpd-php"
	TypeMicrosoftPowerPoint                = "application/vnd.ms-powerpoint"
	TypeMicrosoftPowerPointOpenXML         = "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	TypeRARarchive                         = "application/vnd.rar"
	TypeRichTextFormat                     = "application/rtf"
	TypeBourneshellscript                  = "application/x-sh"
	TypeScalableVectorGraphics             = "image/svg+xml"
	TypeSmallwebformatorAdobeFlashdocument = "application/x-shockwave-flash"
	TypeTapeArchive                        = "application/x-tar"
	TypeTaggedImageFileFormat              = "image/tiff"
	TypeMPEGtransportstream                = "video/mp2t"
	TypeTrueTypeFont                       = "font/ttf"
	TypeText                               = "text/plain"
	TypeMicrosoftVisio                     = "application/vnd.visio"
	TypeWaveformAudioFormat                = "audio/wav"
	TypeWEBMaudio                          = "audio/webm"
	TypeWEBMvideo                          = "video/webm"
	TypeWEBPimage                          = "image/webp"
	TypeWebOpenFontFormat                  = "font/woff"
	TypeWebOpenFontFormat2                 = "font/woff2"
	TypeXHTML                              = "application/xhtml+xml"
	TypeMicrosoftExcel                     = "application/vnd.ms-excel"
	TypeMicrosoftExcelOpenXML              = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	TypeXML                                = "application/xml "
	TypeXUL                                = "application/vnd.mozilla.xul+xml"
	TypeZIParchive                         = "application/zip"
	Type3GPPaudio                          = "video/3gpp"
	Type3GPP2audio                         = "video/3gpp2"
	Type7ziparchive                        = "application/x-7z-compressed"
)

// TypeMap ending html content type map
var TypeMap map[string]string = map[string]string{"aac": TypeAACaudio,
	"abw":    TypeAbiWorddocument,
	"arc":    TypeArchivedocument,
	"avi":    TypeAVI,
	"azw":    TypeAmazonKindleeBookformat,
	"bin":    TypeAnykindofbinarydata,
	"bmp":    Type2BitmapGraphics,
	"bz":     TypeBZiparchive,
	"bz2":    TypeBZip2archive,
	"csh":    TypeCShellscript,
	"css":    TypeCascadingStyleSheets,
	"csv":    TypeCommaSeparatedvalues,
	"doc":    TypeMicrosoftWord,
	"docx":   TypeMicrosoftWord,
	"eot":    TypeMSEmbeddedOpenTypefonts,
	"epub":   TypeElectronicpublication,
	"gz":     TypeGZipCompressedArchive,
	"gif":    TypeGraphicsInterchangeFormat,
	"htm":    TypeHyperTextMarkupLanguage,
	"html":   TypeHyperTextMarkupLanguage,
	"ico":    TypeIconformat,
	"ics":    TypeiCalendarformat,
	"jar":    TypeJavaArchive,
	"jpeg":   TypeJPEGimages,
	"js":     TypeJavaScript,
	"json":   TypeJSONformat,
	"jsonld": TypeJSONLDformat,
	"mid":    TypeMusicalInstrumentDigitalInterface,
	"mjs":    TypeJavaScriptmodule,
	"mp3":    TypeMP3audio,
	"mpeg":   TypeMPEGVideo,
	"mpkg":   TypeAppleInstallerPackage,
	"odp":    TypeOpenDocumentpresentationdocument,
	"ods":    TypeOpenDocumentspreadsheetdocument,
	"odt":    TypeOpenDocumenttextdocument,
	"oga":    TypeOGGaudio,
	"ogv":    TypeOGGvideo,
	"ogx":    TypeOGG,
	"opus":   TypeOpusaudio,
	"otf":    TypeOpenTypefont,
	"png":    TypePortableNetworkGraphics,
	"pdf":    TypeAdobePortableDocumentFormat,
	"php":    TypeHypertextPreprocessor,
	"ppt":    TypeMicrosoftPowerPoint,
	"pptx":   TypeMicrosoftPowerPoint,
	"rar":    TypeRARarchive,
	"rtf":    TypeRichTextFormat,
	"sh":     TypeBourneshellscript,
	"svg":    TypeScalableVectorGraphics,
	"swf":    TypeSmallwebformatorAdobeFlashdocument,
	"tar":    TypeTapeArchive,
	"tiff":   TypeTaggedImageFileFormat,
	"ts":     TypeMPEGtransportstream,
	"ttf":    TypeTrueTypeFont,
	"txt":    TypeText,
	"vsd":    TypeMicrosoftVisio,
	"wav":    TypeWaveformAudioFormat,
	"weba":   TypeWEBMaudio,
	"webm":   TypeWEBMvideo,
	"webp":   TypeWEBPimage,
	"woff":   TypeWebOpenFontFormat,
	"woff2":  TypeWebOpenFontFormat,
	"xhtml":  TypeXHTML,
	"xls":    TypeMicrosoftExcel,
	"xlsx":   TypeMicrosoftExcel,
	"xml":    TypeXML,
	"xul":    TypeXUL,
	"zip":    TypeZIParchive,
	"3gp":    Type3GPPaudio,
	"3g2":    Type3GPP2audio,
	"7z":     Type7ziparchive}
