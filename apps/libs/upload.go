package libs

import (
	"github.com/labstack/echo"
	"io"
	"net/http"
	"strings"
	"path/filepath"
	"fmt"
	"github.com/Sirupsen/logrus"
	"os"
	"github.com/satori/go.uuid"
	appConfig "VCMS/apps/config"
	"errors"
	"math"
)

type UploadInfo struct {
	OriginName 	string 	`json:"origin_name"`
	GeneratedName   string  `json:"generated_name"`
	FileSize   	int64  	`json:"file_size"`
	FilePath 	string 	`json:"file_path"`
	FileType        string  `json:"file_type"`
	FileExt        	string  `json:"file_ext"`
	FileDuration    string 	`json:"file_duration"`
	Stream		StreamUrl  `json:"stream_url"`
}

type StreamUrl struct {
	MpegDash	string 	`json:"mpeg_dash"`
	RTMP		string  `json:"rtmp"`
	HLS		string  `json:"hls"`
	HDS             string	  `json:"hds"`
	Mobile          MobileUrl `json:"mobile"`
}

type MobileUrl struct {
	IOS		string 	`json:"ios"`
	Android		string 	`json:"android"`
}


func (upload *UploadInfo) UploadContent (c echo.Context) error {

	// Source
	tenant_id := c.FormValue("tenant_id")

	if tenant_id == "" {
		return errors.New("TENANT NOT EXIST !!")
	}

	originFile, err := c.FormFile("file")

	if err != nil {
		return err
	}
	src, err := originFile.Open()

	if err != nil {
		return err
	}
	defer src.Close()

	// TODO : VOD에만 올리고 나중에 스케쥴 잡히면 복사한다.
	dstDir := appConfig.Config.APP.ContentsRoot + string(filepath.Separator) + "vod" + string(filepath.Separator) + tenant_id + string(filepath.Separator)

	// Destination
	dst, err := os.Create(dstDir + uuid.NewV4().String() +".mp4")

	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	generatedFile, _ := dst.Stat()


	// Create a buffer to store the header of the file in
	fileHeader := make([]byte, 512)

	n, err := src.Read(fileHeader)
	if err != nil && err != io.EOF {
		return err
	}

	originFileName 	:= originFile.Filename
	fileName 	:= generatedFile.Name()
	fileSize 	:= generatedFile.Size()
	filePath	:= dstDir
	fileType        := http.DetectContentType(fileHeader[:n])
	fileExt         := originFileName[strings.Index(originFileName, ".")+1:]	// 확장자

	// TODO : 썸네일 생성
	// 썸네일 생성 디렉토리
	//thumbPath := string(filepath.Separator) + "thumbnail" + string(filepath.Separator) + fileName[:strings.Index(fileName, ".")] + ".png"
	videoPath := filePath + string(filepath.Separator) + fileName

	// 5초를 표준으로 가져간다.
	//err = appLibs.FFMpegExec{}.RunCreateThumbNail(videoPath, thumbPath, "00:00:05", "hd720")

	// Duration 가져온다.
	duration, err := FFMpegExec{}.RunGetDuration(videoPath)

	if err != nil {
		fmt.Errorf(">>>>> Error : ", err.Error())
	}

	fmt.Println(">>>>>>>>> Final Duration : ", duration)

	logrus.Debug("originFileName : ", originFileName)
	logrus.Debug("fileName : ", fileName)
	logrus.Debug("fileSize : ", fileSize)
	logrus.Debug("filePath : ", filePath)
	logrus.Debug("fileType : ", fileType)
	logrus.Debug("fileExt : ",  fileExt)
	logrus.Debug("duration : ", duration)

	upload.OriginName = originFileName
	upload.GeneratedName = fileName
	upload.FileSize = fileSize
	upload.FilePath = filePath
	upload.FileType = fileType
	upload.FileExt = fileExt
	upload.FileDuration = duration[:strings.Index(duration, ".")]

	mpegDashUrl := "http://" + appConfig.Config.WOWZA.StreamUrl + "/" + appConfig.Config.WOWZA.VcmsVodName + "/_definst_/mp4:"+tenant_id+"/"+fileName+"/manifest.mpd"
	rtmpUrl := "rtmp://" + appConfig.Config.WOWZA.StreamUrl + "/" + appConfig.Config.WOWZA.VcmsVodName +"/_definst_/mp4:"+tenant_id+"/"+fileName
	hlsUrl := "http://" + appConfig.Config.WOWZA.StreamUrl + "/" + appConfig.Config.WOWZA.VcmsVodName + "/_definst_/mp4:"+tenant_id+"/"+fileName+"/manifest.m3u8"
	hdsUrl := "http://" + appConfig.Config.WOWZA.StreamUrl + "/" + appConfig.Config.WOWZA.VcmsVodName + "/_definst_/mp4:"+tenant_id+"/"+fileName+"/manifest.f4m"
	iosUrl := "http://" + appConfig.Config.WOWZA.StreamUrl + "/" + appConfig.Config.WOWZA.VcmsVodName + "/_definst_/mp4:"+tenant_id+"/"+fileName+"/playlist.m3u8"
	androidUrl := "rtsp://" + appConfig.Config.WOWZA.StreamUrl + "/" + appConfig.Config.WOWZA.VcmsVodName +"/_definst_/"+tenant_id+"/"+fileName
	/*
	mpegdash-> http://118.128.155.121:1935/vod/_definst_/mp4:chomdan/3297f9b5-cff3-4dc3-afb0-f5b43a755e72.mp4/manifest.mpd
	rtmp-> rtmp://118.128.155.121:1935/vod/_definst_/mp4:chomdan/3297f9b5-cff3-4dc3-afb0-f5b43a755e72.mp4
	hds : http://118.128.155.121:1935/vod/_definst_/mp4:chomdan/3297f9b5-cff3-4dc3-afb0-f5b43a755e72.mp4/manifest.f4m
	iOS
	http://118.128.155.121:1935/vod/_definst_/mp4:chomdan/3297f9b5-cff3-4dc3-afb0-f5b43a755e72.mp4/playlist.m3u8
	Android/Other
	rtsp://118.128.155.121:1935/vod/_definst_/chomdan/3297f9b5-cff3-4dc3-afb0-f5b43a755e72.mp4
	*/

	upload.Stream.MpegDash = mpegDashUrl
	upload.Stream.RTMP = rtmpUrl
	upload.Stream.HLS = hlsUrl
	upload.Stream.HDS = hdsUrl
	upload.Stream.Mobile.IOS = iosUrl
	upload.Stream.Mobile.Android = androidUrl


	return nil
}

// 썸네일
func (upload *UploadInfo) UploadThumbnail (c echo.Context) error {

	content_id := c.FormValue("content_id")

	originFile, err := c.FormFile("file")

	if err != nil {
		return err
	}
	src, err := originFile.Open()

	if err != nil {
		return err
	}
	defer src.Close()

	dstDir := appConfig.Config.APP.ThumbNailRoot
	originFileName := originFile.Filename
	originFileExt := originFileName[strings.Index(originFileName, ".")+1:]

	generateName := "thumb_" + content_id + "." + originFileExt

	// Destination
	dst, err := os.Create(dstDir + string(filepath.Separator) + generateName)

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	// 기존 파일들 삭제

	generatedFile, _ := dst.Stat()

	upload.OriginName = originFile.Filename
	upload.GeneratedName = generatedFile.Name()
	upload.FileSize = generatedFile.Size()
	upload.FilePath = dstDir
	upload.FileExt = originFileExt

	return nil
}


// UPLOAD 디렉토리 사이즈
func DirSizeMB(path string) float64 {
	var dirSize int64 = 0

	readSize := func(path string, file os.FileInfo, err error) error {
		if !file.IsDir() {
			dirSize += file.Size()
		}

		return nil
	}

	filepath.Walk(path, readSize)

	sizeMB := toFixed(float64(dirSize) / 1024.0 / 1024.0, 2)

	return sizeMB
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num * output)) / output
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}