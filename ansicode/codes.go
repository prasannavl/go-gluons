package ansicode

var (
	Escape      string = "\x1b"
	EscapeBlock string = Escape + "["
	EndBlock    string = "m"

	Reset            string = EscapeBlock + "0m"
	Bold             string = EscapeBlock + "1m"
	Italic           string = EscapeBlock + "3m"
	Underline        string = EscapeBlock + "4m"
	Blink            string = EscapeBlock + "5m"
	Inverse          string = EscapeBlock + "7m"
	Strikethrough    string = EscapeBlock + "9m"
	UnderlineOff     string = EscapeBlock + "24m"
	InverseOff       string = EscapeBlock + "27m"
	StrikethroughOff string = EscapeBlock + "29m"

	Black     string = EscapeBlock + "30m"
	Red       string = EscapeBlock + "31m"
	Green     string = EscapeBlock + "32m"
	Blue      string = EscapeBlock + "34m"
	Yellow    string = EscapeBlock + "33m"
	Magenta   string = EscapeBlock + "35m"
	Cyan      string = EscapeBlock + "36m"
	White     string = EscapeBlock + "37m"
	DefaultFg string = EscapeBlock + "39m"
	Grey      string = EscapeBlock + "90m"

	BlackBright   string = EscapeBlock + "30;1m"
	RedBright     string = EscapeBlock + "31;1m"
	GreenBright   string = EscapeBlock + "32;1m"
	BlueBright    string = EscapeBlock + "34;1m"
	YellowBright  string = EscapeBlock + "33;1m"
	MagentaBright string = EscapeBlock + "35;1m"
	CyanBright    string = EscapeBlock + "36;1m"
	WhiteBright   string = EscapeBlock + "37;1m"
	GreyBright    string = EscapeBlock + "90;1m"

	BlackBg   string = EscapeBlock + "40m"
	RedBg     string = EscapeBlock + "41m"
	GreenBg   string = EscapeBlock + "42m"
	BlueBg    string = EscapeBlock + "44m"
	YellowBg  string = EscapeBlock + "43m"
	MagentaBg string = EscapeBlock + "45m"
	CyanBg    string = EscapeBlock + "46m"
	WhiteBg   string = EscapeBlock + "47m"
	DefaultBg string = EscapeBlock + "49m"

	BlackBrightBg   string = EscapeBlock + "40;1m"
	RedBrightBg     string = EscapeBlock + "41;1m"
	GreenBrightBg   string = EscapeBlock + "42;1m"
	BlueBrightBg    string = EscapeBlock + "44;1m"
	YellowBrightBg  string = EscapeBlock + "43;1m"
	MagentaBrightBg string = EscapeBlock + "45;1m"
	CyanBrightBg    string = EscapeBlock + "46;1m"
	WhiteBrightBg   string = EscapeBlock + "47;1m"
)

func Make256Color(rawColorStr string) string {
	return EscapeBlock + "38;5;" + rawColorStr + "m"
}
