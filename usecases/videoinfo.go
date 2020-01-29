package usecases

import (
	"fmt"
	"log"
	"srv-video-info/interfaces"
)

type StreamInfo struct {
	Name     string `json:"name"`
	Width    *int   `json:"width,omitempty"`
	Height   *int   `json:"height,omitempty"`
	BitRate  int    `json:"bitRate"`
	Duration string `json:"duration"`
}

type VideoInfo struct {
	Video StreamInfo `json:"video,omitempty"`
	Audio StreamInfo `json:"audio,omitempty"`
}

type VideoInfoInteractor struct {
	MediaInfo interfaces.MediaInfoRepo
}

func NewVideoInfoInteractor(mediaInfo interfaces.MediaInfoRepo) *VideoInfoInteractor {
	return &VideoInfoInteractor{MediaInfo: mediaInfo}
}

func (v *VideoInfoInteractor) Get(filename string) (videoInfo VideoInfo, err error) {
	media, err := v.MediaInfo.Open(filename)
	if err != nil {
		return
	}

	for i := 0; i < int(media.GetStreamCount()); i++ {
		stream, err := media.GetStream(i)
		if err != nil {
			log.Printf("failed get stream %s\n", err.Error())
			continue
		}

		if stream.IsVideoStream() {
			width, height, bitrate, duration, codec, err := stream.GetVideoInfo()
			if err != nil {
				log.Printf("failed get video stream info %s\n", err.Error())
				continue
			}

			videoInfo.Video.Duration = fmt.Sprintf("%fs", duration)
			videoInfo.Video.Width = &width
			videoInfo.Video.Height = &height
			videoInfo.Video.BitRate = bitrate
			videoInfo.Video.Name = codec
		}

		if stream.IsAudioStream() {
			bitrate, duration, codec, err := stream.GetAudioInfo()
			if err != nil {
				log.Printf("failed get audio stream info %s\n", err.Error())
				continue
			}

			videoInfo.Audio.Duration = fmt.Sprintf("%fs", duration)
			videoInfo.Audio.BitRate = bitrate
			videoInfo.Audio.Name = codec
		}

	}

	return
}
