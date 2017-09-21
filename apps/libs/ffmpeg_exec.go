package libs

import (
	"fmt"
	"os/exec"
	"strings"
	"os"
)

type FFMpegExec struct {

}

func printCommand(cmd *exec.Cmd) {
	fmt.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
}

func printError(err error) {
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
	}
}

func printOutput(outs []byte) {
	if len(outs) > 0 {
		fmt.Printf("==> Output: %s\n", string(outs))
	}
}

// 썸네일 생성
func (ffmpegExec FFMpegExec) RunCreateThumbNail(_video_path string, _thumb_path string, _time string, _size string) error {

	fmt.Println(">>>>>>>>>>>>> RunCreateThumbNail !!!")

	// 기존에 존재하면 삭제
	var _, err = os.Stat(_thumb_path)
	if os.IsExist(err) {
		os.Remove(_thumb_path)
	}

	//ffmpeg -i sample.mp4 -r 1 -ss 00:00:10 -s 320x240 -vframes 1 -f image2 /home/foresight/비디오/www_320x240.png

	// ffmpeg -ss 00:00:10 -i /MEDIA_DATA/VCMS/8719321.mp4 -r 1 -s 320x240 -vcodec png -vframes 3 /MEDIA_DATA/VCMS/8719321.png

	cmd := exec.Command("ffmpeg",
		"-ss",
		_time,
		"-i",
		_video_path,
		"-r",
		"1",
		"-s",
		//"640x320",
		//"320x240"
		//"hd720",
		_size,
		//"00:00:05",
		"-vcodec",
		"png",
		"-vframes",
		"3",
		//"-f",
		//"image2",
		_thumb_path)
		//"/home/foresight/비디오/wwwww.png")

	printCommand(cmd)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return err
	}

	if len(output) > 0 {
		//fmt.Printf("==> Output: %s\n", string(output))
	}
	//printError(err)
	//printOutput(output) // => go version go1.3 darwin/amd64
	return nil

}


// 재생시간 구하기
func (ffmpegExec FFMpegExec) RunGetDuration(_path string) (string, error) {

	fmt.Println(">>>>>>>>>>>>> RunGetDuration !!!")
	fmt.Println(">>>>>>>>>>>>> _path : " , _path)

	//ffmpeg -i "/home/foresight/비디오/sample.mp4" 2>&1 | grep "Duration"

	cmd := exec.Command(
		"ffmpeg",
		"-i",
		_path )
		//"2>&1",
		//"|",
		//"grep",
		//"Duration")

	printCommand(cmd)

	output, err := cmd.CombinedOutput()

	printError(err)
	printOutput(output) // => go version go1.3 darwin/amd64

	if len(output) > 0 {
		//fmt.Printf("==> Output: %s\n", string(output))
	}

	strOutput := string(output)

	i := strings.Index(strOutput, "Duration");
	z := strings.Index(strOutput, ", start:");

	fmt.Println(">>>>>> i : '", i, "'")
	fmt.Println(">>>>>> z : '", z, "'")

	substring := strOutput[i+10:z]

	substring = strings.TrimSpace(substring)

	fmt.Println(">>>>>> substring : '", substring, "'")

	return substring, err


}
