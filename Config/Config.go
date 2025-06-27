package Config

import "fmt"

func Init() {
	webInit()
	databaseInit()
}

func GetPort() string {
	return webConfig.Port
}
func GetHost() string {
	return webConfig.Host
}
func GetHttps() bool {
	return webConfig.Https
}
func GetBuild() bool {
	return webConfig.Build
}
func GetStaticFolders() []string {
	return webConfig.StaticFolders
}
func GetCssFolders() []string {
	return webConfig.CssFolders
}
func GetJsFolders() []string {
	return webConfig.JsFolders
}
func GetViewFolder() string {
	return webConfig.ViewFolder
}
func GetMaxSessionCount() int {
	return webConfig.MaxSessionCount
}
func SetPort(p string) {
	webConfig.Port = p
}
func SetHost(h string) {
	webConfig.Host = h
}
func SetHttps(h bool) {
	webConfig.Https = h
}
func SetBuild(b bool) {
	webConfig.Build = b
}
func SetStaticFolders(folders []string) {
	webConfig.StaticFolders = folders
}
func SetCssFolders(folders []string) {
	webConfig.CssFolders = folders
}
func SetJsFolders(folders []string) {
	webConfig.JsFolders = folders
}
func SetViewFolder(folder string) {
	webConfig.ViewFolder = folder
}
func SetMaxSessionCount(count int) {
	webConfig.MaxSessionCount = count
}
func GetWebConfig() *webConfigStruct {
	return &webConfig
}

// function to get dsn of the database
func GetDSN() string {
	if databaseConfig.Driver == "mysql" {
		return fmt.Sprintf("%s:%s@%s(%s:%s)/%s",
			databaseConfig.User,
			databaseConfig.Password,
			databaseConfig.Protocol,
			databaseConfig.Host,
			databaseConfig.Port,
			databaseConfig.Database,
		)
	}
	return ""
}
func GetDatabaseDriver() string {
	return databaseConfig.Driver
}

func GetDatabaseConfig() *databaseConfigStruct {
	return &databaseConfig
}
