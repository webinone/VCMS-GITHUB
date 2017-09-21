package libs

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"bytes"
	"errors"
	appConfig "VCMS/apps/config"
	apiModel "VCMS/apps/models/api"
	"encoding/json"
	"path/filepath"
	"os"
)

type WowzaHttpClient struct {
	APIUrl	string
}

// 어플리케이션 생성
func (wowza_client WowzaHttpClient) CreateApplication(app_name string, app_type string, stream_type string) (string, error) {

	//app_name += "_"+  strings.ToLower(app_type)

	url := appConfig.Config.WOWZA.ApiUrl + "/v2/servers/_defaultServer_/vhosts/_defaultVHost_/applications/" + app_name

	storage_dir := appConfig.Config.APP.ContentsRoot + "/" + app_name

	reqBody := []byte(`
		{
		   "restURI": "`+url+`",
		   "name": "`+app_name+`",
		   "appType": "`+app_type+`",
		   "clientStreamReadAccess": "*",
		   "clientStreamWriteAccess": "*",
		   "description": "A basic live application",
		   "httpCORSHeadersEnabled":true,
		   "mediaReaderRandomAccessReaderClass": "",
		   "mediaReaderBufferSeekIO": false,
		   "streamConfig": {
		      "restURI": "`+url+`/streamconfiguration",
		      "streamType": "`+stream_type+`",
		      "storageDirExists": true,
		      "createStorageDir":true,
		      "storageDir": "`+storage_dir+`"
		   }
		}
	`)

	//"storageDir": "${com.wowza.wms.context.VHostConfigHome}/content/`+app_name+`"

	var jsonStr = []byte(reqBody)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type","application/json")
	req.Header.Set("Accept","application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(respBody))

	data := &apiModel.APIWowzaResult{}

	json.Unmarshal([]byte(respBody), data)

	fmt.Println("success : ", data.Success)

	if !data.Success {
		return "", errors.New("Wowza Internal Error : " + string(respBody))
	}

	return string(respBody), nil
}

// 어플리케이션 삭제
func (wowza_client WowzaHttpClient) DeleteApplication(app_name string, app_type string) (string, error) {

	//app_name += "_"+ strings.ToLower(app_type)

	url := appConfig.Config.WOWZA.ApiUrl + "/v2/servers/_defaultServer_/vhosts/_defaultVHost_/applications/" + app_name

	req, err := http.NewRequest("DELETE", url, nil)
	req.Header.Set("Content-Type","application/json")
	req.Header.Set("Accept","application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(respBody))

	data := &apiModel.APIWowzaResult{}

	json.Unmarshal([]byte(respBody), data)

	fmt.Println("success : ", data.Success)

	if !data.Success {
		return "", errors.New("Wowza Internal Error : " + string(respBody))
	}
	// VOD 삭제
	err = os.RemoveAll(appConfig.Config.APP.ContentsRoot + string(filepath.Separator) + "vod" + string(filepath.Separator) + app_name)

	if err != nil {
		return "", err
	}

	return string(respBody), nil
}

// Schedule Smil 생성
func (wowza_client WowzaHttpClient) CreateScheduleSmil(app_name string) (string, error) {

	//app_name += app_name + "_live"

	url := wowza_client.APIUrl + "/v2/servers/_defaultServer_/vhosts/_defaultVHost_/applications/" + app_name + "/smilfiles/streamschedule"

	reqBody := []byte(`
		{
		  "restURI": "`+url+`",
		  "smilStreams": [
		    {
		    "systemLanguage": "en",
		    "src": "sample.mp4",
		    "systemBitrate": "50000",
		    "type": "video",
		    "audioBitrate": "44100",
		    "videoBitrate": "750000",
		    "restURI": "`+url+`",
		    "width": "640",
		    "height": "360"
		    }
		  ]
		}
	`)

	var jsonStr = []byte(reqBody)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type","application/json")
	req.Header.Set("Accept","application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {

		fmt.Println(err.Error())
		return "", err
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(respBody))

	data := &apiModel.APIWowzaResult{}

	json.Unmarshal([]byte(respBody), data)

	fmt.Println("success : ", data.Success)

	if !data.Success {
		return "", errors.New("Wowza Internal Error : " + string(respBody))
	}

	return string(respBody), nil
}

type SwitchSmil struct {
	RestURI 	string 		`json:"restURI"`
	SmilStreams 	[]SmilStream 	`json:"smilStreams"`
}

type SmilStream struct {
	Src 		string		`json:"src"`
	SystemBirate 	string		`json:"systemBitrate"`
	Type		string 		`json:"type"`
	RestURI		string 		`json:"restURI"`
}


// Schedule Switch Smil 파일 생성
func (wowza_client WowzaHttpClient) CreateScheduleSwitchSmil(app_name string, stream_name string, bitrate string) (string, error) {

	//app_name += "_"+ "live"

	url := wowza_client.APIUrl + "/v2/servers/_defaultServer_/vhosts/_defaultVHost_/applications/" + app_name + "/smilfiles/switch"

	reqBody := []byte(`
		{
		  "restURI": "`+url+`",
		  "smilStreams": [
		    {
		    "src": "` + stream_name + `",
		    "systemBitrate": "` + bitrate + `",
		    "type" : "video",
		    "restURI": "`+url+`"
		    }
		  ]
		}
	`)

	var jsonStr = []byte(reqBody)

	fmt.Println(string(jsonStr))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type","application/json")
	req.Header.Set("Accept","application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {

		fmt.Println(err.Error())
		return "", err
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(respBody))

	data := &apiModel.APIWowzaResult{}

	json.Unmarshal([]byte(respBody), data)

	fmt.Println("success : ", data.Success)
	fmt.Println("respBody : ", string(respBody))

	if !data.Success {
		return "", errors.New("Wowza Internal Error : " + string(respBody))
	}

	return string(respBody), nil
}

// Smil File Modify
func (wowza_client WowzaHttpClient)  UpdateScheduleSwitchSmil(switchSmil SwitchSmil) error {


	reqBody, _ := json.Marshal(switchSmil)

	var jsonStr = []byte(reqBody)

	fmt.Println(string(jsonStr))

	url := switchSmil.RestURI

	fmt.Println(string(jsonStr))
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type","application/json")
	req.Header.Set("Accept","application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {

		fmt.Println(err.Error())
		return err
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(respBody))

	data := &apiModel.APIWowzaResult{}

	json.Unmarshal([]byte(respBody), data)

	fmt.Println("success : ", data.Success)
	fmt.Println("respBody : ", string(respBody))

	if !data.Success {
		return errors.New("Wowza Internal Error : " + string(respBody))
	}

	return nil
}



// Schedule Smil Reload
func (wowza_client WowzaHttpClient) ReloadScheduleSmil(app_name string) (string, error) {

	//app_name += "_"+ "live"

	url := wowza_client.APIUrl + "/scheduleloader?action=load&app="+app_name
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(respBody))

	data := &apiModel.APIWowzaResult{}

	json.Unmarshal([]byte(respBody), data)

	fmt.Println("success : ", data.Success)

	if !data.Success {
		return "", errors.New("Wowza Internal Error : " + string(respBody))
	}

	return string(respBody), nil
}

// Schedule Smil Get
func (wowza_client WowzaHttpClient) GetScheduleSmil(app_name string) (string, error) {



	return "", nil
}


//{
//"serverName": "_defaultServer_",
//"uptime": 166656,
//"bytesIn": 0,
//"bytesOut": 0,
//"bytesInRate": 0,
//"bytesOutRate": 0,
//"totalConnections": 0,
//"connectionCount": {
//"WEBM": 0,
//"DVRCHUNKS": 0,
//"RTMP": 0,
//"MPEGDASH": 0,
//"CUPERTINO": 0,
//"SANJOSE": 0,
//"SMOOTH": 0,
//"RTP": 0
//}
//}


// 어플리케이션 접속자
func (wowza_client WowzaHttpClient) GetApplicationConnections(app_name string) (*apiModel.APIWowzaStatisticResult, error) {

	//url := appConfig.Config.WOWZA.ApiUrl + "/v2/servers/_defaultServer_/vhosts/_defaultVHost_/applications/" + app_name
	url := wowza_client.APIUrl + "/v2/servers/_defaultServer_/vhosts/_defaultVHost_/applications/"+app_name+"/monitoring/current"

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type","application/json")
	req.Header.Set("Accept","application/json")


	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(respBody))

	data := &apiModel.APIWowzaStatisticResult{}

	json.Unmarshal([]byte(respBody), data)

	return data, nil
}


// 채널 접속자

func (wowza_client WowzaHttpClient) GetChannelConnections(app_name string, channel_name string) (*apiModel.APIWowzaStatisticResult, error) {

	//url := appConfig.Config.WOWZA.ApiUrl + "/v2/servers/_defaultServer_/vhosts/_defaultVHost_/applications/" + app_name
	url := wowza_client.APIUrl + "/v2/servers/_defaultServer_/vhosts/_defaultVHost_/applications/"+app_name+"/instances/_definst_/incomingstreams/"+channel_name+"/monitoring/current"

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type","application/json")
	req.Header.Set("Accept","application/json")


	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(respBody))

	data := &apiModel.APIWowzaStatisticResult{}

	json.Unmarshal([]byte(respBody), data)

	return data, nil
}