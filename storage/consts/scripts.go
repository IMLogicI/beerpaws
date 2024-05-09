package consts

import "fmt"

var (
	GetRules            = fmt.Sprintf("SELECT * FROM %s", PointRulesTable)
	MakeRequest         = fmt.Sprintf("INSERT INTO \"%s\" (\"%s\",\"%s\",\"%s\",\"%s\",\"%s\") VALUES", RequestPointTable, RuleIDField, UserIDField, ScreenshotLinkField, ApprovedField, ClosedField)
	GetUserByDiscordID  = fmt.Sprintf("SELECT %s, %s, %s FROM \"%s\" WHERE %s = ", IDField, NicknameField, RoleNameField, UserTable, DiscordIDField)
	SaveUserFromDiscord = fmt.Sprintf("INSERT INTO \"%s\" (\"%s\",\"%s\",\"%s\",\"%s\") VALUES", UserTable, NicknameField, DiscordIDField, PasswordField, RoleNameField)
	AddNewRule          = fmt.Sprintf("INSERT INTO \"%s\" (\"%s\", \"%s\", \"%s\", \"%s\") VALUES", PointRulesTable, CountField, NameField, DescriptionField, IsEarnedField)
	GetOpenedRequest    = "select rp.\"id\" , pr.\"name\" , u.nickname  , rp.screenshot_link , rp.approved  from request_point rp \njoin \"user\" u on u.id  = rp.user_id \njoin point_rules pr on pr.id = rp.id \nwhere rp.closed = false "
	ApproveRequest      = "update request_point set approved = true \nwhere id = "
	AddPoints           = "insert into point_history (request_id, \"time\") values (%d, '%s')"
	GetRequestByID      = "select * from request_point rp where id = %d"
	CloseRequest        = "update request_point set closed = true \nwhere id = "
	GetPoints           = "select sum(pr.count) from point_history ph \njoin request_point rp2 on rp2.id = ph.request_id \njoin point_rules pr on pr.id = rp2.rule_id \nwhere rp2.user_id  = %d"
)
