package constant

var (
	TimeOut5s  = 5
	TimeOut20s = 20
	TimeOut5m  = 300

	DateLayout         = "2006-01-02" // or use time.DateOnly while go version >= 1.20
	DefaultDate        = "1970-01-01"
	DateTimeLayout     = "2006-01-02 15:04:05" // or use time.DateTime while go version >= 1.20
	DateTimeSlimLayout = "20060102150405"
)

var (
	Red   = "\033[31m"
	Green = "\033[32m"
	Reset = "\033[0m"
)
