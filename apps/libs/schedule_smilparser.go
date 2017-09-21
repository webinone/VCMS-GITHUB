package libs

import (
	"encoding/xml"
	"fmt"
	"os"
	"io/ioutil"
	"time"
	"github.com/metakeule/fmtdate"
	"github.com/go-xmlfmt/xmlfmt"
)

type Smil struct {
	XMLName		xml.Name	`xml:"smil"`
	Head		head
	Body		body
}

type head struct {
	XMLName		xml.Name	`xml:"head"`
}

type body struct {
	XMLName		xml.Name	`xml:"body"`
	Streams        	[]stream 	`xml:"stream"`
	Playlists	[]playlist	`xml:"playlist"`
	Switch 		switchEl	`xml:"switch"`
}

type switchEl struct {
	XMLName		xml.Name	`xml:"switch"`
	Videos		[]Video		`xml:"video"`
}

type stream struct {
	XMLName    	xml.Name 	`xml:"stream"`
	Name	 	string   	`xml:"name,attr"`
}

type playlist struct {
	XMLName    	xml.Name 	`xml:"playlist"`
	Name	 	string   	`xml:"name,attr"`
	PlayOnStream	string   	`xml:"playOnStream,attr"`
	Repeat	 	string   	`xml:"repeat,attr"`
	Scheduled	string   	`xml:"scheduled,attr"`
	Videos		[]Video		`xml:"video"`
}

type Video struct {
	XMLName    	xml.Name 	`xml:"video"`
	Src	 	string   	`xml:"src,attr"`
	Start		string   	`xml:"start,attr"`
	Length	 	string   	`xml:"length,attr"`
	SystemBitrate   string          `xml:"systemBitrate,attr"`
}

type ScheduleSmilParser struct {
	FilePath 	string
}

func (smilParser ScheduleSmilParser) UpdateScheduleSmil(_stream_id string, _repeat string, _schedule_id string, _schedule_datetime string, _videos []Video) error {


	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!! schedule_id : ", _schedule_id)


	// 스트림이 존재하는지를 따진다.
	var streamExists bool
	streamExists = false

	file, err := os.Open(smilParser.FilePath) // For read access.

	if err != nil {
		fmt.Printf("error: %v", err)
		return err
	}
	defer file.Close()

	doc, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("error: %v", err)
		return err
	}

	smilXmlObj := Smil{}
	err = xml.Unmarshal(doc, &smilXmlObj)
	if err != nil {
		fmt.Printf("error: %v", err)
		return  err
	}

	if len(smilXmlObj.Body.Streams) > 0 {
		for _, stream := range smilXmlObj.Body.Streams {

			if stream.Name == _stream_id {
				streamExists = true
			}

		}
	}


	// Stream 추가
	if !streamExists {
		smilXmlObj.Body.Streams = append(smilXmlObj.Body.Streams, stream{ Name:_stream_id })
	}


	if len(smilXmlObj.Body.Playlists) > 0 {
		// 존재한다.

		// 존재하면 가장 먼저 시간이 오래된 것은 지워버린다. 오래된 기준은 1일 이상 지난것으로 한다. 나중에 수정 가능
		offSet, _ := time.ParseDuration("+09.00h")

		now := time.Now().UTC().Add(offSet)

		fmt.Println("Today : ", fmtdate.Format("YYYY-MM-DD hh:mm:ss", now))

		//var longTimeAgo string

		for _, playlist := range smilXmlObj.Body.Playlists {

			fmt.Printf("Name: %s ", playlist.Name)

			longTimeAgo, _ := fmtdate.Parse("YYYY-MM-DD hh:mm:ss", playlist.Scheduled)

			longTimeAgo.UTC().Add(offSet)

			diff := now.Sub(longTimeAgo)

			fmt.Println("DIFF HOURS : ", int(diff.Hours()))


			if int(diff.Hours()) > 24 {
				// 삭제 해야 한다.
				removePlayListByScheduled(&smilXmlObj, &playlist)
			}

			// 기존에 존재하면 삭제 한다.

			if playlist.Name == _schedule_id {
				removePlayListBySchduleId(&smilXmlObj, &playlist)
			}
		}

		// play list 추가
		smilXmlObj.Body.Playlists = append(smilXmlObj.Body.Playlists, playlist{
			Name		: _schedule_id,
			PlayOnStream	: _stream_id,
			Repeat		: _repeat,
			Scheduled	: _schedule_datetime,
			Videos		: _videos,

		})


	} else {
		// Playlist 존재하지 않으므로 추가

		smilXmlObj.Body.Playlists = append(smilXmlObj.Body.Playlists, playlist{
			Name		: _schedule_id,
			PlayOnStream	: _stream_id,
			Repeat		: _repeat,
			Scheduled	: _schedule_datetime,
			Videos		: _videos,

		})

	}

	//fmt.Println(json.Marshal(smilXmlObj))

	res1B, _ := xml.Marshal(smilXmlObj)
	fmt.Println(string(res1B))

	prettyXml := xmlfmt.FormatXML(string(res1B), "\t", "  ")

	ioutil.WriteFile(smilParser.FilePath, []byte(prettyXml), 0777)

	//ioutil.WriteFile(smilParser.FilePath, res1B, 0777)

	return nil
}

func (smilParser ScheduleSmilParser) DeleteScheduleSmil(_schedule_id string) error {

	// 스트림이 존재하는지를 따진다.
	file, err := os.Open(smilParser.FilePath) // For read access.

	if err != nil {
		fmt.Printf("error: %v", err)
		return err
	}
	defer file.Close()

	doc, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("error: %v", err)
		return err
	}

	smilXmlObj := Smil{}
	err = xml.Unmarshal(doc, &smilXmlObj)
	if err != nil {
		fmt.Printf("error: %v", err)
		return  err
	}

	if len(smilXmlObj.Body.Playlists) > 0 {

		for _, playlist := range smilXmlObj.Body.Playlists {
			if playlist.Name == _schedule_id {
				// 삭제 해야 한다.
				removePlayListBySchduleId(&smilXmlObj, &playlist)
			}
		}

	}

	res1B, _ := xml.Marshal(smilXmlObj)
	fmt.Println(string(res1B))

	prettyXml := xmlfmt.FormatXML(string(res1B), "\t", "  ")

	ioutil.WriteFile(smilParser.FilePath, []byte(prettyXml), 0777)

	return nil
}


func removePlayListByScheduled (f *Smil, p *playlist) {

	for i, b := range f.Body.Playlists {

		fmt.Println(">>>>> b.Scheduled ", b.Scheduled)
		fmt.Println(">>>>> p.Scheduled ", p.Scheduled)
		fmt.Println(">>>>> len(f.Body.Playlists) ", len(f.Body.Playlists))

		if b.Scheduled == p.Scheduled {
			if len(f.Body.Playlists) > 1 {
				copy(f.Body.Playlists[i:], f.Body.Playlists[i+1:])  // shift
			}
			//f.Body.Playlists[len(f.Body.Playlists)-1] = nil     // remove reference
			f.Body.Playlists = f.Body.Playlists[:len(f.Body.Playlists)-1]
		}
	}
}

func removePlayListBySchduleId (f *Smil, p *playlist) {

	for i, b := range f.Body.Playlists {

		if b.Name == p.Name {

			copy(f.Body.Playlists[i:], f.Body.Playlists[i+1:])  // shift
			//f.Body.Playlists[len(f.Body.Playlists)-1] = nil     // remove reference
			f.Body.Playlists = f.Body.Playlists[:len(f.Body.Playlists)-1]
		}
	}
}


func (smilParser ScheduleSmilParser) UpdateScheduleSwitchSmil(_channel_id string, _bitrate string) error {

	// 1. smil 파일이 존재하는지를 따진다.
	file, err := os.Open(smilParser.FilePath) // For read access.

	if err != nil {
		fmt.Printf("error: %v", err)
		return err
	}
	defer file.Close()

	doc, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("error: %v", err)
		return err
	}

	smilXmlObj := Smil{}
	err = xml.Unmarshal(doc, &smilXmlObj)
	if err != nil {
		fmt.Printf("error: %v", err)
		return  err
	}

	isExistChannel := false;

	if len(smilXmlObj.Body.Switch.Videos) > 0 {
		for i, video := range smilXmlObj.Body.Switch.Videos {

			fmt.Println(video.Src)

			if _channel_id == video.Src {
				isExistChannel = true
				video.SystemBitrate = _bitrate
				smilXmlObj.Body.Switch.Videos[i] = video
			}
		}
	}

	if !isExistChannel {

		smilXmlObj.Body.Switch.Videos = append(smilXmlObj.Body.Switch.Videos, Video{
			Src: _channel_id,
			SystemBitrate: _bitrate,
		})

	}

	res1B, _ := xml.Marshal(smilXmlObj)
	fmt.Println(string(res1B))

	prettyXml := xmlfmt.FormatXML(string(res1B), "\t", "  ")

	ioutil.WriteFile(smilParser.FilePath, []byte(prettyXml), 0777)

	return nil
}


func (smilParser ScheduleSmilParser) DeleteScheduleSwitchSmil(_channel_id string) error {

	// 1. smil 파일이 존재하는지를 따진다.
	file, err := os.Open(smilParser.FilePath) // For read access.

	if err != nil {
		fmt.Printf("error: %v", err)
		return err
	}
	defer file.Close()

	doc, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("error: %v", err)
		return err
	}

	smilXmlObj := Smil{}
	err = xml.Unmarshal(doc, &smilXmlObj)
	if err != nil {
		fmt.Printf("error: %v", err)
		return  err
	}

	if len(smilXmlObj.Body.Switch.Videos) > 0 {
		for _, video := range smilXmlObj.Body.Switch.Videos {

			if _channel_id == video.Src {
				//smilXmlObj.Body.Switch.Videos[i] = nil
				removeVideo(&smilXmlObj, &video)
			}
		}
	}

	res1B, _ := xml.Marshal(smilXmlObj)
	fmt.Println(string(res1B))

	prettyXml := xmlfmt.FormatXML(string(res1B), "\t", "  ")

	ioutil.WriteFile(smilParser.FilePath, []byte(prettyXml), 0777)


	return nil
}

func removeVideo (f *Smil, p *Video) {

	for i, b := range f.Body.Switch.Videos {

		if b.Src == p.Src {
			copy(f.Body.Switch.Videos[i:], f.Body.Switch.Videos[i+1:])  // shift
			f.Body.Switch.Videos = f.Body.Switch.Videos[:len(f.Body.Switch.Videos)-1]
		}
	}
}