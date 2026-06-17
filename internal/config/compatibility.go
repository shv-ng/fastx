package config

var (
	compatibility map[string][]string = map[string][]string{
		"SQLAlchemy": {"Alembic", "None"},
		"SQLModel":   {"Alembic", "None"},
		"Tortoise":   {"Aerich", "None"},
		"RawAsyncpg": {"Yoyo", "None"},
	}
)

func ValidMigrations(orm string) []string {
	if mig, ok := compatibility[orm]; ok {
		return mig
	}
	return []string{"None"}
}
