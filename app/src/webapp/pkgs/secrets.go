package pkgs

type WebAppSecrets struct {
}

var WebAppSecretStore *WebAppSecrets

func newWebAppSecrets() *WebAppSecrets {
	return &WebAppSecrets{}
}

func InitWebAppSecretStore() {
	if WebAppSecretStore == nil {
		WebAppSecretStore = newWebAppSecrets()
	}
}
