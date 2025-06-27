package Config

func Init() {
	webInit()
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
