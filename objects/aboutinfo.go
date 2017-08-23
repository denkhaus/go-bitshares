package objects

//go:generate ffjson $GOFILE

type AboutInfo struct {
	ClientVersion       string `json:"client_version"`
	GrapheneRevision    string `json:"graphene_revision"`
	FCRevision          string `json:"fc_revision"`
	CompileDate         string `json:"compile_date"`
	OpenSSLVersion      string `json:"openssl_version"`
	GrapheneRevisionAge string `json:"graphene_revision_age"`
	FCRevisionAge       string `json:"fc_revision_age"`
	BoostVersion        string `json:"boost_version"`
	Build               string `json:"build"`
}
