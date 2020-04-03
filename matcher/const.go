package matcher

// DefaultPort TODO
var DefaultPort = map[string]string{
	"http":  "80",
	"https": "443",
	"ssh":   "22",
	"ftp":   "21",
}

// Const TODO
const (
	Empty = ""
	Dot   = "."
	Colon = ":"
	Star  = "*"
)

// Const TODO
const (
	ByteDot   byte = '.'
	ByteColon byte = ':'
	ByteStar  byte = '*'
)
