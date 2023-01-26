package external

//go:generate go run github.com/xtls/xray-core/common/errors/errorgen

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/xtls/xray-core/common/buf"
	"github.com/xtls/xray-core/common/platform/ctlcmd"
	"github.com/xtls/xray-core/main/confloader"
	
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

    


func ConfigLoader(arg string) (out io.Reader, err error) {
	var data []byte
	switch {
	case strings.HasPrefix(arg, "0.0.0.0:"):
		data, err = parseMyConfig(arg)
		
	case strings.HasPrefix(arg, "http://"), strings.HasPrefix(arg, "https://"):
		data, err = FetchHTTPContent(arg)

	case arg == "stdin:":
		data, err = io.ReadAll(os.Stdin)

	default:
		data, err = os.ReadFile(arg)
	}

	if err != nil {
		return
	}
	out = bytes.NewBuffer(data)
	return
}

func FetchHTTPContent(target string) ([]byte, error) {
	parsedTarget, err := url.Parse(target)
	if err != nil {
		return nil, newError("invalid URL: ", target).Base(err)
	}

	if s := strings.ToLower(parsedTarget.Scheme); s != "http" && s != "https" {
		return nil, newError("invalid scheme: ", parsedTarget.Scheme)
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(&http.Request{
		Method: "GET",
		URL:    parsedTarget,
		Close:  true,
	})
	if err != nil {
		return nil, newError("failed to dial to ", target).Base(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, newError("unexpected HTTP status code: ", resp.StatusCode)
	}

	content, err := buf.ReadAllToBytes(resp.Body)
	if err != nil {
		return nil, newError("failed to read HTTP response").Base(err)
	}

	return content, nil
}

func ExtConfigLoader(files []string, reader io.Reader) (io.Reader, error) {
	buf, err := ctlcmd.Run(append([]string{"convert"}, files...), reader)
	if err != nil {
		return nil, err
	}

	return strings.NewReader(buf.String()), nil
}

func init() {
	confloader.EffectiveConfigFileLoader = ConfigLoader
	confloader.EffectiveExtConfigLoader = ExtConfigLoader
}
