package infrastructure

import (
	"fmt"
	"github.com/giorgisio/goav/avcodec"
	"github.com/giorgisio/goav/avformat"
	"srv-video-info/interfaces"
	"unsafe"
)

const (
	StreamTypeUnknown = iota
	StreamTypeVideo
	StreamTypeAudio
)

type codec struct {
	codec    *avcodec.Context
	codecCtx *avformat.CodecContext
}

type stream struct {
	s     *avformat.Stream
	index int
}

type media struct {
	ctx *avformat.Context
}

type mediaInfo struct{}

func NewMediaInfo() *mediaInfo {
	return &mediaInfo{}
}

func (mi *mediaInfo) Open(filename string) (m interfaces.MediaRepo, err error) {
	pFormatContext := avformat.AvformatAllocContext()
	if avformat.AvformatOpenInput(&pFormatContext, filename, nil, nil) != 0 {
		return nil, fmt.Errorf("unable to open file %s\n", filename)
	}

	if pFormatContext.AvformatFindStreamInfo(nil) < 0 {
		return nil, fmt.Errorf("couldn't find stream information")
	}

	m = &media{
		ctx: pFormatContext,
	}
	return
}

func (m *media) GetStreamCount() int {
	return int(m.ctx.NbStreams())
}

func (m *media) GetStream(index int) (s interfaces.StreamRepo, err error) {
	if index >= m.GetStreamCount() {
		err = fmt.Errorf("index out of range streams")
		return
	}

	s = &stream{
		index: index,
		s:     m.ctx.Streams()[index],
	}

	return
}

func (m *media) Close() {
	m.ctx.AvformatCloseInput()
}

func (s *stream) GetStreamType() int {
	switch s.s.CodecParameters().AvCodecGetType() {
	case avformat.AVMEDIA_TYPE_VIDEO:
		return StreamTypeVideo
	case avformat.AVMEDIA_TYPE_AUDIO:
		return StreamTypeAudio
	}
	return StreamTypeUnknown
}

func (s *stream) IsVideoStream() bool {
	if s.GetStreamType() == StreamTypeVideo {
		return true
	}
	return false
}

func (s *stream) IsAudioStream() bool {
	if s.GetStreamType() == StreamTypeAudio {
		return true
	}
	return false
}

func (s *stream) getDuration() float64 {
	return float64(s.s.Duration()) * (float64(s.s.TimeBase().Num()) / float64(s.s.TimeBase().Den()))
}

func (s *stream) GetVideoInfo() (width, height, bitrate int, duration float64, codecName string, err error) {
	c, err := s.getCodec()
	if err != nil {
		return
	}
	defer c.close()

	height = c.codecCtx.GetHeight()
	width = c.codecCtx.GetWidth()
	duration = s.getDuration()
	bitrate = c.codec.BitRate()
	codecName = avcodec.AvcodecGetName(avcodec.CodecId(c.codecCtx.GetCodecId()))

	return
}

func (s *stream) GetAudioInfo() (bitrate int, duration float64, codecName string, err error) {
	c, err := s.getCodec()
	if err != nil {
		return
	}
	defer c.close()

	duration = s.getDuration()
	bitrate = c.codec.BitRate()
	codecName = avcodec.AvcodecGetName(avcodec.CodecId(c.codecCtx.GetCodecId()))

	return
}

func (s *stream) getCodec() (c *codec, err error) {
	c = new(codec)

	c.codecCtx = s.s.Codec()

	pCodec := avcodec.AvcodecFindDecoder(avcodec.CodecId(c.codecCtx.GetCodecId()))
	if pCodec == nil {
		return nil, fmt.Errorf("unsupported codec")

	}

	c.codec = pCodec.AvcodecAllocContext3()
	if c.codec.AvcodecCopyContext((*avcodec.Context)(unsafe.Pointer(c.codecCtx))) != 0 {
		return nil, fmt.Errorf("couldn't copy codec context")
	}

	if c.codec.AvcodecOpen2(pCodec, nil) < 0 {
		return nil, fmt.Errorf("could not open codec")
	}

	return
}

func (c *codec) close() {
	c.codec.AvcodecClose()
	(*avcodec.Context)(unsafe.Pointer(c.codecCtx)).AvcodecClose()
}
