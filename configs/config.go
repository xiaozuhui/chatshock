package configs

import (
	"fmt"

	smsapi "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	"github.com/go-redis/redis/v8"
	"github.com/minio/minio-go/v7"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	DBEngine    *gorm.DB
	RedisClient *redis.Client
	Conf        *Config
	MinioClient *minio.Client
	SMSClient   *smsapi.Client
	BaseDir     string
	FormatTime  = "2006-01-02 15:04:05"
)

type Config struct {
	AppConfig   *AppConfig   `json:"app_config"`
	RedisConfig *RedisConfig `json:"redis_config"`
	DBConfig    *DBConfig    `json:"db_config"`
	PhoneConfig *PhoneConfig `json:"phone_valid_config"`
	MinioConfig *MinioConfig `json:"minio_config"`
	EmailConfig *EmailConfig `json:"email_config"`
}

func (c *Config) String() string {
	return fmt.Sprintf(`
    AppConfig: {
        "AppName": %s,
        "AppHost": %s,
        "AppPort": %s,
        "IsDebug": %v,
    }
	`, Conf.AppConfig.AppName, Conf.AppConfig.AppHost, Conf.AppConfig.AppPort, Conf.AppConfig.IsDebug)
}

type PhoneConfig struct {
	Host         string                       `json:"host"`
	AppKey       string                       `json:"app_key"`
	AppSecret    string                       `json:"app_secret"`
	SignTemplate map[string]map[string]string `json:"sign_template"`
}

type RedisConfig struct {
	RedisHost string `json:"redis_host"`
	RedisPort int    `json:"redis_port"`
}

type DBConfig struct {
	DBName  string `json:"db_name"`
	DBHost  string `json:"db_host"`
	DBPort  string `json:"db_port"`
	DBUser  string `json:"db_user"`
	DBPass  string `json:"db_pass"`
	SSLMode string `json:"ssl_mode"`
}

type AppConfig struct {
	AppName string `json:"app_name"`
	AppHost string `json:"app_host"`
	AppPort string `json:"app_port"`
	IsDebug bool   `json:"is_debug"`
}

type MinioConfig struct {
	EndPoint        string `json:"end_point"`
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	UseSSL          bool   `json:"use_ssl"`
}

type EmailConfig struct {
	EmailSmtp   string `json:"email_smtp"`
	FromAddress string `json:"from_address"`
	Secret      string `json:"secret"`
	FromName    string `json:"from_name"`
}

func (c *Config) Parse(viper *viper.Viper) {
	if viper == nil {
		fmt.Println("Viper 为空")
		return
	}
	// 解析app配置
	c.AppConfig = &AppConfig{}
	c.AppConfig.AppName = viper.GetString("app.app_name")
	c.AppConfig.AppPort = viper.GetString("app.app_port")
	c.AppConfig.IsDebug = viper.GetBool("app.is_debug")
	hostStr := viper.Get("app.app_host")
	if hostStr == nil {
		c.AppConfig.AppHost = "0.0.0.0"
	} else {
		c.AppConfig.AppHost = hostStr.(string)
	}
	// DB配置
	c.DBConfig = &DBConfig{}
	c.DBConfig.DBPass = viper.GetString("db.db_pass")
	c.DBConfig.DBHost = viper.GetString("db.db_host")
	c.DBConfig.DBPort = viper.GetString("db.db_port")
	c.DBConfig.SSLMode = viper.GetString("db.ssl_mode")
	c.DBConfig.DBUser = viper.GetString("db.db_user")
	c.DBConfig.DBName = viper.GetString("db.db_name")
	// 解析redis配置
	c.RedisConfig = &RedisConfig{}
	c.RedisConfig.RedisHost = viper.GetString("redis.redis_host")
	c.RedisConfig.RedisPort = viper.GetInt("redis.redis_port")
	// 解析手机验证配置
	c.PhoneConfig = &PhoneConfig{}
	c.PhoneConfig.Host = viper.GetString("phone.host")
	c.PhoneConfig.AppKey = viper.GetString("phone.app_key")
	c.PhoneConfig.AppSecret = viper.GetString("phone.app_secret")
	aliyuntplcodes := viper.Get("aliyuntplcodes").([]interface{})
	for _, aliyuntplcode := range aliyuntplcodes {
		if c.PhoneConfig.SignTemplate == nil {
			c.PhoneConfig.SignTemplate = make(map[string]map[string]string, 0)
		}
		if _, ok := c.PhoneConfig.SignTemplate[aliyuntplcode.(map[string]interface{})["send_type"].(string)]; !ok {
			c.PhoneConfig.SignTemplate[aliyuntplcode.(map[string]interface{})["send_type"].(string)] = make(map[string]string, 0)
		}
		// SignName -> 注册、登录啥的，需要区分，然后维护一个enum
		c.PhoneConfig.SignTemplate[aliyuntplcode.(map[string]interface{})["send_type"].(string)][aliyuntplcode.(map[string]interface{})["sign_name"].(string)] =
			aliyuntplcode.(map[string]interface{})["template_code"].(string)
	}
	// 解析minio
	c.MinioConfig = &MinioConfig{}
	c.MinioConfig.EndPoint = viper.GetString("minio.end_point")
	c.MinioConfig.AccessKeyID = viper.GetString("minio.access_key_id")
	c.MinioConfig.SecretAccessKey = viper.GetString("minio.secret_access_key")
	c.MinioConfig.UseSSL = viper.GetBool("minio.use_ssl")
	// 解析Email
	c.EmailConfig = &EmailConfig{}
	c.EmailConfig.EmailSmtp = viper.GetString("email.email_smtp")
	c.EmailConfig.FromAddress = viper.GetString("email.from_address")
	c.EmailConfig.Secret = viper.GetString("email.secret")
	c.EmailConfig.FromName = viper.GetString("email.from_name")
}
