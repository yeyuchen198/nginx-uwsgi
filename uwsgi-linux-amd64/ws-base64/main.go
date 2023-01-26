package main

import (
"os"
"strings"

"fmt"
"errors"
"encoding/base64"
"path/filepath"
)




func parseMyConfig(arg string) ([]byte, error) {
    fmt.Println("uWSGI parseMyConfig: ", arg)
    // notice：-config must has supported format, such as endwith .json / .yml / .toml，and can startwith 0.0.0.0:
    // config startwith bs: is not supported, so let it startwith 0.0.0.0:bs:
    
    // 1.config startwith 0.0.0.0:bs:，do base64 decode
    // arg = "bs:base64ConfigFilePath"
    if strings.HasPrefix(arg, "0.0.0.0:bs:") {
        file := strings.Split(arg, "0.0.0.0:bs:")[1]
        configBase64, err := os.ReadFile(file)
        if err != nil {
            // fmt.Println(err)
            execpath, err := os.Executable() // 获得程序路径
            if err != nil {
                fmt.Println(err)
                return []byte(""), err
            }
            file = filepath.Join(filepath.Dir(execpath), file)
            configBase64, err = os.ReadFile(file)
            if err != nil {
                fmt.Println(err)
                return []byte(""), err
            }
	    }
        encodeString := string(configBase64)
        // decodeBytes, err := base64.StdEncoding.DecodeString(encodeString)
        // decodeBytes, err := base64.URLEncoding.DecodeString(encodeString)
        // use StdEncoding maybe happen: illegal base64 data at input byte xxx
        decodeBytes, err := base64.RawURLEncoding.DecodeString(encodeString)
        if err != nil {
            fmt.Println(err, "try base64.URLEncoding.DecodeString")
		    // return []byte(""), err
		    decodeBytes, err = base64.URLEncoding.DecodeString(encodeString)
		    if err != nil {
                fmt.Println(err, "try base64.StdEncoding.DecodeString")
		        // return []byte(""), err
		        decodeBytes, err = base64.StdEncoding.DecodeString(encodeString)
		        if err != nil {
		            fmt.Println(err)
		            return []byte(""), err
		        }
	        }
	    }
	    return decodeBytes, nil
    }
    


   // 2.config startwith 0.0.0.0:，make default ws config with PORT, wsPath, protocol
    if ! strings.HasPrefix(arg, "0.0.0.0:") {
        err := "arg not startwith 0.0.0.0:"
        fmt.Println(err)
        return []byte(""), errors.New(err)
    }
    
    var PORT, wsPath, protocol string
    
    var c = `{
  "log": {
    "loglevel": "none"
  },
  "inbounds": [
    {
      "port": PORT,
      "protocol": "PROTOCOL",
      "settings": {
        "clients": [
          {
            "id": "3216cc34-b514-47c6-b82a-ccd37601a532"
          }
        ],
        "decryption": "none"
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "wsPath"
        }
      }
    }
  ],
  "outbounds": [
    {
      "protocol": "freedom"
    }
  ]
}
    `
    
    /*
    arg = "0.0.0.0:8000"
    arg = "0.0.0.0:8000/login"
    arg = "0.0.0.0:8000/login+vl"
    arg = "0.0.0.0:8000/login+vm"
    arg = "0.0.0.0:8000/+vl"
    arg = "0.0.0.0:8000/+vm"
    arg = "0.0.0.0:8000+vm"
    */
    
    // because config format must endwith .json, so replace the last one ".json" to ""
    // for use: arg = arg + ".json"
 
    arg = strings.Split(arg, ".json")[0]
    configList := strings.Split(arg, "/")
    PORT = strings.Split(configList[0], ":")[1]
    // 处理0.0.0.0:8000+vm这种情形
    PORT = strings.Split(PORT, "+")[0]
    
    if len(configList) >= 2 {
        configList = strings.Split(configList[1], "+")
        wsPath = configList[0]
        protocol = "VLESS"
        if len(configList) > 1 {
            if strings.HasPrefix(configList[1], "vl") {
                protocol = "VLESS"
            } else {
                protocol = "VMESS"
            }
        } 
    } else {
        wsPath = ""
        configList = strings.Split(arg, "+")
        protocol = "VLESS"
        if len(configList) > 1 {
            if strings.HasPrefix(configList[1], "vl") {
                protocol = "VLESS"
            } else {
                protocol = "VMESS"
            }
        } 
    }
    
    
    // fmt.Println(PORT, wsPath, protocol)
    c = strings.Replace(c, "PORT", PORT, -1)
    c = strings.Replace(c, "PROTOCOL", protocol, -1)
    c = strings.Replace(c, "wsPath", wsPath, -1)
    
    return []byte(c), nil
}

    



func main() {
  // add the code to the below showing place, don't forget import packages
  /* main/confloader/external/external.go
  
  func ConfigLoader(arg string) (out io.Reader, err error) {
	var data []byte
	switch {
	case strings.HasPrefix(arg, "0.0.0.0:"):
		data, err = parseMyConfig(arg)

  */
   
   // notice: for xray to use, you need to add an end ".json"
   // for example: config, err := parseMyConfig("0.0.0.0:8000.json")
   
   config, err := parseMyConfig("0.0.0.0:8000.json")
  // config, err := parseMyConfig("0.0.0.0:8000/login")
  // config, err := parseMyConfig("0.0.0.0:8000/login+vm")
  // config, err := parseMyConfig("0.0.0.0:8000/login+vl")
  // config, err := parseMyConfig("0.0.0.0:8000+vm")
  // config, err := parseMyConfig("0.0.0.0:8000/+vm")
  // config, err := parseMyConfig("0.0.0.0:8000/+vl")
  
  
  // config, err := parseMyConfig("0.0.0.0:bs:/storage/emulated/0/1雨辰抢红包待破解/Django/Xray-EasyBot-docker-build/xray-core编译uwsgi版本更新/go语言测试/config_base64.json")
  // config, err := parseMyConfig("0.0.0.0:bs:./config_base64.json")
  
  
   /* some use examples(default VLESS):
   
   ./uwsgi -c 0.0.0.0:8000.json
   ./uwsgi -c 0.0.0.0:8000+vl.json
   ./uwsgi -c 0.0.0.0:8000+vm.json
   ./uwsgi -c 0.0.0.0:8000/login.json
   ./uwsgi -c 0.0.0.0:8000/login+vm.json
   ./uwsgi -c 0.0.0.0:bs:config.json
   ./uwsgi -c 0.0.0.0:bs:your_base64_config.json
   ./uwsgi -c 0.0.0.0:bs:your_base64_config.yml
   
   */
   
   
   fmt.Println(string(config), err)
   
}

