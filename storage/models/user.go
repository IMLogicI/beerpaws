package models

type User struct {
	ID        int64  `db:"id"`
	Nickname  string `db:"nickname"`
	DiscordID string `db:"discord_id"`
	Password  string `db:"password"`
	Role      string `db:"role_name"`
}
