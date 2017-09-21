package libs

import (
	"testing"
	"fmt"
)

//func TestFFMpegExec_RunCreateThumbNail(t *testing.T) {
//
//	FFMpegExec{}.RunCreateThumbNail("/home/foresight/비디오/images/sample%d.png", "00:00:23", "hd720")
//
//}
//
//func TestWowzaHttpClient_CreateApplication(t *testing.T) {
//	fmt.Println("11111111111111111")
//
//
//	//wowzaHttpClient.CreateApplication("test01", "Live")
//	//WowzaHttpClient{}.CreateTenant()
//}

//func TestWowzaHttpClient_CreateScheduleSmil(t *testing.T) {
//	wowzaHttpClient := &WowzaHttpClient{APIUrl:"http://localhost:8087"}
//	wowzaHttpClient.CreateScheduleSmil("livetest")
//}

//func TestWowzaHttpClient_ReloadScheduleSmil(t *testing.T) {
//	//wowzaHttpClient := &WowzaHttpClient{APIUrl:"http://localhost:8086"}
//	//wowzaHttpClient.ReloadScheduleSmil("livetest")
//}

//func TestWowzaHttpClient_CreateScheduleSwitchSmil(t *testing.T) {
//	//wowzaHttpClient := &WowzaHttpClient{APIUrl:"http://localhost:8087"}
//	////wowzaHttpClient.ReloadScheduleSmil("livetest")
//	//wowzaHttpClient.CreateScheduleSwitchSmil("livetest", "myStream1", "450000")
//}

func TestWowzaHttpClient_GetApplicationConnections(t *testing.T) {
	wowzaHttpClient := &WowzaHttpClient{APIUrl:"http://222.231.29.47:8087"}
	fmt.Println(wowzaHttpClient.GetApplicationConnections("test"))
}

func TestWowzaHttpClient_GetChannelConnections(t *testing.T) {
	wowzaHttpClient := &WowzaHttpClient{APIUrl:"http://222.231.29.47:8087"}
	fmt.Println(wowzaHttpClient.GetChannelConnections("test", "a059479c-d0e6-4758-9ee2-200411f1eb24"))
}
