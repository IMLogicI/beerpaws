package consts

import "fmt"

var (
	GetRules            = fmt.Sprintf("SELECT * FROM %s", PointRulesTable)
	MakeRequest         = fmt.Sprintf("INSERT INTO \"%s\" (\"%s\",\"%s\",\"%s\") VALUES", RequestPointTable, RuleIDField, UserIDField, ScreenshotLinkField)
	GetUserByDiscordID  = fmt.Sprintf("SELECT %s, %s FROM \"%s\" WHERE %s = ", IDField, NicknameField, UserTable, DiscordIDField)
	SaveUserFromDiscord = fmt.Sprintf("INSERT INTO \"%s\" (\"%s\",\"%s\",\"%s\") VALUES", UserTable, NicknameField, DiscordIDField, PasswordField)
)
