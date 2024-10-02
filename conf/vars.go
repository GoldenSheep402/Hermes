package conf

type GlobalConfig struct {
	MODE string `yaml:"Mode"`
	Port string `yaml:"Port"` // grpc和http服务监听端口
	Log  struct {
		LogPath string `yaml:"LogPath"`
		CLS     struct {
			Endpoint    string `yaml:"Endpoint"`
			AccessKey   string `yaml:"AccessKey"`
			AccessToken string `yaml:"AccessToken"`
			TopicID     string `yaml:"TopicID"`
		} `yaml:"CLS"`
	} `yaml:"Log"`
	B2 struct {
		BucketKeyID string `yaml:"BucketKeyID"`
		BucketKey   string `yaml:"BucketKey"`
		BucketName  string `yaml:"BucketName"`
	} `yaml:"B2"`
	Pyroscope struct {
		ApplicationName string `yaml:"ApplicationName"`
		ServerAddress   string `yaml:"ServerAddress"`
		BasicAuthUser   string `yaml:"BasicAuthUser"`
		BasicAuthPass   string `yaml:"BasicAuthPass"`
		TenantID        string `yaml:"TenantID"`
	} `yaml:"Pyroscope"`
	Uptrace struct {
		ServiceName    string `yaml:"ServiceName"`
		ServiceVersion string `yaml:"ServiceVersion"`
		DSN            string `yaml:"DSN"`
	} `yaml:"Uptrace"`
	SentryDsn string `yaml:"SentryDsn"`
	TrackerV1 struct {
		Endpoint       string   `yaml:"Endpoint"`
		AllowedSubnets []string `yaml:"AllowedSubnets"`
	} `yaml:"TrackerV1"`
}
