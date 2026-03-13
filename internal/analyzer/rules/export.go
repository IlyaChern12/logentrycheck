package rules

var (
	CheckEnglish  = checkEnglish
	IsEnglishOnly = isEnglishOnly

	CheckLowercase = checkLowercase

	ContainsSensitiveKeyword = containsSensitiveKeyword
	CheckSensitive           = checkSensitive
	CustomKeywordsFlag       = &customKeywordsFlag

	HasSpecialChars   = hasSpecialChars
	CheckSpecialChars = checkSpecialChars

	IsContextMethod    = isContextMethod
	ExtractLogMessage  = extractLogMessage
	ExtractMessageArg  = extractMessageArg
	ExtractPackageName = extractPackageName
	IsLoggerCall       = isLoggerCall
)
