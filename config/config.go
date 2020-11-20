package config

import (
	"encoding/json"
	"log"
	"os"

	"github.com/Treblex/simple-daily/tools/upload"
)

// Global 全局配置
var Global *configType = initConfig()

func initConfig() *configType {
	t := &configType{}
	//读取配置文件
	if err := t.ReadConfig(); err != nil {
		panic(err)
	}
	return t
}

type configType struct {
	BaseURL string            `json:"base_url"` // 网站根目录
	Port    int               `json:"port"`     //端口
	AliOss  upload.AliOssConf `json:"ali_oss"`  //阿里云oss
}

// ReadConfig 读取配置 初始化时运行 绑定为全局变量
// 在我使用 ReadConfig 命名函数的时候 编辑器提示了错误， 函数应该和结构体configType保存一直的大写或者小写 以保证其他包的调用者可以使用这个函数
func (c *configType) ReadConfig() (err error) {
	f, err := os.Open("./config/config.json")
	defer f.Close()
	if err != nil {
		log.Fatalln("打开配置文件错误，请创建 config/config.json 参考(config-defaultjson")
		return
	}
	conf := &configType{}
	if err = json.NewDecoder(f).Decode(c); err != nil {
		return err
	}

	c = conf
	return nil
}
