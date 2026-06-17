package config

type Config struct {
	ProjectName    string
	PythonVersion  string
	PackageManager string
	Database       string // empty means no db
	ORM            string
	Migrations     string
	UseDocker      bool
	UseCI          bool
	Components     []string
}
