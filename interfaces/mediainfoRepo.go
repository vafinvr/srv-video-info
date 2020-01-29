package interfaces

type MediaInfoRepo interface {
	Open(filename string) (m MediaRepo, err error)
}

type MediaRepo interface {
	GetStreamCount() int
	GetStream(index int) (s StreamRepo, err error)
	Close()
}

type StreamRepo interface {
	GetAudioInfo() (bitrate int, duration float64, codecName string, err error)
	GetVideoInfo() (width, height, bitrate int, duration float64, codecName string, err error)
	IsAudioStream() bool
	IsVideoStream() bool
	GetStreamType() int
}
