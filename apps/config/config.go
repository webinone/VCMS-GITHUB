package config

import (
	"io/ioutil"
	"encoding/json"
	"flag"
	"fmt"
	"log"
)

type Configuration struct {
	APP 	appConfig
	RDB 	[]rdbConfig
	REDIS   redisConfig
	WOWZA   wowzaConfig
	AUTH 	authConfig
	EMAIL   emailConfig
	HELLOT  helloTConfig
}

type appConfig struct {
	Port 			string	`json:"port"`
	ContentsRoot 	string 	`json:"contents_root"`
	ThumbNailRoot	string	`json:"thumbnail_root"`
	BannerRoot		string	`json:"banner_root"`
	DownloadRoot	string	`json:"download_root"`
	ExcelRoot       string	`json:"excel_root"`
}

type rdbConfig struct {
	Type            string	`json:"type"`
	Product 	string		`json:"product"`
	ConnectString 	string	`json:"connect_string"`
	MaxIdleConns	int		`json:"max_idle_conns"`
	MaxOpenConns	int		`json:"max_open_conns"`
	Debug 	bool			`json:"debug"`
	Migrate bool			`json:"migrate"`
}

type redisConfig struct {
	Url			string 	`json:"url"`
	Password	string 	`json:"password"`
	Db			int		`json:"db"`
	ExpiredTime	int64		`json:"expired_time"`
}

type wowzaConfig struct {
	ApiUrl 		string		`json:"api_url"`
	CustomApiUrl    string 		`json:"custom_api_url"`
	StreamUrl       string 		`json:"stream_url"`
	VcmsVodName     string		`json:"vcms_vod_name"`
}

type authConfig struct {
	SSOUrl		string 		`json:"sso_url"`
	AuthKey		string 		`json:"auth_key"`
	JwtKey		string 		`json:"jwt_key"`
}

type emailConfig struct {
	SmtpUrl			string		`json:"smtp_url"`
	SmtpPort		int			`json:"smtp_port"`
	User			string		`json:"user"`
	Password		string		`json:"password"`
	Template		string		`json:"template"`
}

type helloTConfig struct {
	ApiUrl 		string		`json:"api_url"`
}

var Config  = &Configuration{}

func LoadPathConfig(path string) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Config File Missing. ", err)

	}

	err = json.Unmarshal(file, &Config)
	if err != nil {
		log.Fatal("Config Parse Error: ", err)
	}
}

func LoadAutoConfig () {

	configRoot := flag.String("mode", "foo", "development mode")
	flag.Parse()

	var path string
	// 개발 환경 셋팅
	if *configRoot == "foo" {
		path = "app.json"
	} else if *configRoot == "dev" {
		//ConfigRoot = *configRoot
		// 개발 환경
		path = "app.json"
	} else if *configRoot == "prod" {
		//ConfigRoot = *configRoot
		// 운영환경 셋팅
		path = "app-product.json"
	} else {
		panic("Development mode error !!")
	}

	fmt.Println("config file path : ", path)

	LoadPathConfig(path)
}