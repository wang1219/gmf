package gmf

/*
#cgo pkg-config: libavcodec libavutil

#include <libavutil/avutil.h>
#include <libavutil/channel_layout.h>
#include <libavutil/dict.h>
#include <libavutil/pixdesc.h>
#include <libavutil/opt.h>
#include <libavutil/frame.h>
#include <libavutil/parseutils.h>
#include <libavutil/common.h>
#include <libavutil/eval.h>
#include <libavutil/audio_fifo.h>
#include "libavutil/error.h"
#include "libavutil/mathematics.h"
#include "libavutil/rational.h"
#include "libavutil/samplefmt.h"
#include "libavcodec/avcodec.h"
#include "libavutil/imgutils.h"

#ifdef AV_LOG_TRACE
#define GO_AV_LOG_TRACE AV_LOG_TRACE
#else
#define GO_AV_LOG_TRACE AV_LOG_DEBUG
#endif

#ifdef AV_PIX_FMT_XVMC_MPEG2_IDCT
#define GO_AV_PIX_FMT_XVMC_MPEG2_IDCT AV_PIX_FMT_XVMC_MPEG2_MC
#else
#define GO_AV_PIX_FMT_XVMC_MPEG2_IDCT 0
#endif

#ifdef AV_PIX_FMT_XVMC_MPEG2_MC
#define GO_AV_PIX_FMT_XVMC_MPEG2_MC AV_PIX_FMT_XVMC_MPEG2_MC
#else
#define GO_AV_PIX_FMT_XVMC_MPEG2_MC 0
#endif

static const AVDictionaryEntry *go_av_dict_next(const AVDictionary *m, const AVDictionaryEntry *prev)
{
 return av_dict_get(m, "", prev, AV_DICT_IGNORE_SUFFIX);
}

static const int go_av_dict_has(const AVDictionary *m, const char *key, int flags)
{
 if (av_dict_get(m, key, NULL, flags) != NULL)
 {
   return 1;
 }
 return 0;
}

static int go_av_expr_parse2(AVExpr **expr, const char *s, const char * const *const_names, int log_offset, void *log_ctx)
{
 return av_expr_parse(expr, s, const_names, NULL, NULL, NULL, NULL, log_offset, log_ctx);
}

static const int go_av_errno_to_error(int e)
{
 return AVERROR(e);
}
uint32_t return_int (int num) {
	return (uint32_t)(num);
}

uint8_t * gmf_alloc_buffer(int32_t fmt, int width, int height) {
	int numBytes = av_image_get_buffer_size(fmt, width, height, 0);
	return (uint8_t *) av_malloc(numBytes*sizeof(uint8_t));
}

#cgo pkg-config: libavutil
*/
import "C"

import (
	"bytes"
	"errors"
	"fmt"
	"syscall"
	"unsafe"
)

type ChannelLayout uint64

const (
	ChannelLayoutMono            ChannelLayout = C.AV_CH_LAYOUT_MONO
	ChannelLayoutStereo          ChannelLayout = C.AV_CH_LAYOUT_STEREO
	ChannelLayout2Point1         ChannelLayout = C.AV_CH_LAYOUT_2POINT1
	ChannelLayout21              ChannelLayout = C.AV_CH_LAYOUT_2_1
	ChannelLayoutSurround        ChannelLayout = C.AV_CH_LAYOUT_SURROUND
	ChannelLayout3Point1         ChannelLayout = C.AV_CH_LAYOUT_3POINT1
	ChannelLayout4Point0         ChannelLayout = C.AV_CH_LAYOUT_4POINT0
	ChannelLayout4Point1         ChannelLayout = C.AV_CH_LAYOUT_4POINT1
	ChannelLayout22              ChannelLayout = C.AV_CH_LAYOUT_2_2
	ChannelLayoutQuad            ChannelLayout = C.AV_CH_LAYOUT_QUAD
	ChannelLayout5Point0         ChannelLayout = C.AV_CH_LAYOUT_5POINT0
	ChannelLayout5Point1         ChannelLayout = C.AV_CH_LAYOUT_5POINT1
	ChannelLayout5Point0Back     ChannelLayout = C.AV_CH_LAYOUT_5POINT0_BACK
	ChannelLayout5Point1Back     ChannelLayout = C.AV_CH_LAYOUT_5POINT1_BACK
	ChannelLayout6Point0         ChannelLayout = C.AV_CH_LAYOUT_6POINT0
	ChannelLayout6Point0Front    ChannelLayout = C.AV_CH_LAYOUT_6POINT0_FRONT
	ChannelLayoutHexagonal       ChannelLayout = C.AV_CH_LAYOUT_HEXAGONAL
	ChannelLayout6Point1         ChannelLayout = C.AV_CH_LAYOUT_6POINT1
	ChannelLayout6Point1Back     ChannelLayout = C.AV_CH_LAYOUT_6POINT1_BACK
	ChannelLayout6Point1Front    ChannelLayout = C.AV_CH_LAYOUT_6POINT1_FRONT
	ChannelLayout7Point0         ChannelLayout = C.AV_CH_LAYOUT_7POINT0
	ChannelLayout7Point0Front    ChannelLayout = C.AV_CH_LAYOUT_7POINT0_FRONT
	ChannelLayout7Point1         ChannelLayout = C.AV_CH_LAYOUT_7POINT1
	ChannelLayout7Point1Wide     ChannelLayout = C.AV_CH_LAYOUT_7POINT1_WIDE
	ChannelLayout7Point1WideBack ChannelLayout = C.AV_CH_LAYOUT_7POINT1_WIDE_BACK
	ChannelLayoutOctagonal       ChannelLayout = C.AV_CH_LAYOUT_OCTAGONAL
	ChannelLayoutStereoDownmix   ChannelLayout = C.AV_CH_LAYOUT_STEREO_DOWNMIX
)

func FindChannelLayoutByName(name string) (ChannelLayout, bool) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	cChannelLayout := C.av_get_channel_layout(cName)
	return ChannelLayout(cChannelLayout), (cChannelLayout != 0)
}

func FindDefaultChannelLayout(numberOfChannels int) (ChannelLayout, bool) {
	cl := C.av_get_default_channel_layout(C.int(numberOfChannels))
	if cl <= 0 {
		return 0, false
	}
	return ChannelLayout(cl), true
}

func (cl ChannelLayout) NumberOfChannels() int {
	return int(C.av_get_channel_layout_nb_channels((C.uint64_t)(cl)))
}

func (cl ChannelLayout) Name() string {
	str, _ := cl.NameOk()
	return str
}

func (cl ChannelLayout) NameOk() (string, bool) {
	for index := C.unsigned(0); ; index++ {
		var cCL C.uint64_t
		var cName *C.char
		if C.av_get_standard_channel_layout(index, &cCL, &cName) != 0 {
			break
		}
		if ChannelLayout(cCL) == cl {
			return cStringToStringOk(cName)
		}
	}
	return "", false
}

func (cl ChannelLayout) DescriptionOk() (string, bool) {
	return cStringToStringOk(C.av_get_channel_description((C.uint64_t)(cl)))
}

func cStringToStringOk(cStr *C.char) (string, bool) {
	if cStr == nil {
		return "", false
	}
	return C.GoString(cStr), true
}

func ChannelLayouts() []ChannelLayout {
	var cls []ChannelLayout
	for index := C.unsigned(0); ; index++ {
		var cCL C.uint64_t
		if C.av_get_standard_channel_layout(index, &cCL, nil) != 0 {
			break
		}
		cls = append(cls, ChannelLayout(cCL))
	}
	return cls
}

type AVRational C.struct_AVRational

type AVR struct {
	Num int
	Den int
}

const (
	AVERROR_EOF = -541478725
	// AV_ROUND_PASS_MINMAX = 8192
)

var (
	AV_TIME_BASE   int        = C.AV_TIME_BASE
	AV_TIME_BASE_Q AVRational = AVRational{1, C.int(AV_TIME_BASE)}
)

func (a AVR) AVRational() AVRational {
	return AVRational{C.int(a.Num), C.int(a.Den)}
}

func (a AVR) String() string {
	return fmt.Sprintf("%d/%d", a.Num, a.Den)
}

func (a AVR) Av2qd() float64 {
	return float64(a.Num) / float64(a.Den)
}

func (a AVR) Invert() AVR {
	return AVR{Num: a.Den, Den: a.Num}
}

func (a AVRational) AVR() AVR {
	return AVR{Num: int(a.num), Den: int(a.den)}
}

func AvError(averr int) error {
	errlen := 1024
	b := make([]byte, errlen)

	C.av_strerror(C.int(averr), (*C.char)(unsafe.Pointer(&b[0])), C.size_t(errlen))

	return errors.New(string(b[:bytes.Index(b, []byte{0})]))
}

func AvErrno(ret int) syscall.Errno {
	if ret < 0 {
		ret = -ret
	}

	return syscall.Errno(ret)
}

func RescaleQ(a int64, encBase AVRational, stBase AVRational) int64 {
	return int64(C.av_rescale_q(C.int64_t(a), C.struct_AVRational(encBase), C.struct_AVRational(stBase)))
}

func RescaleQRnd(a int64, encBase AVRational, stBase AVRational) int64 {
	return int64(C.av_rescale_q_rnd(C.int64_t(a), C.struct_AVRational(encBase), C.struct_AVRational(stBase), C.AV_ROUND_NEAR_INF|C.AV_ROUND_PASS_MINMAX))
}

func CompareTimeStamp(aTimestamp int, aTimebase AVRational, bTimestamp int, bTimebase AVRational) int {
	return int(C.av_compare_ts(C.int64_t(aTimestamp), C.struct_AVRational(aTimebase),
		C.int64_t(bTimestamp), C.struct_AVRational(bTimebase)))
}
func RescaleDelta(inTb AVRational, inTs int64, fsTb AVRational, duration int, last *int64, outTb AVRational) int64 {
	return int64(C.av_rescale_delta(C.struct_AVRational(inTb), C.int64_t(inTs), C.struct_AVRational(fsTb), C.int(duration), (*C.int64_t)(unsafe.Pointer(&last)), C.struct_AVRational(outTb)))
}

func Rescale(a, b, c int64) int64 {
	return int64(C.av_rescale(C.int64_t(a), C.int64_t(b), C.int64_t(c)))
}

func RescaleTs(pkt *Packet, encBase AVRational, stBase AVRational) {
	C.av_packet_rescale_ts(&pkt.avPacket, C.struct_AVRational(encBase), C.struct_AVRational(stBase))
}

func GetSampleFmtName(fmt int32) string {
	return C.GoString(C.av_get_sample_fmt_name(fmt))
}

func AvInvQ(q AVRational) AVRational {
	avr := q.AVR()
	return AVRational{C.int(avr.Den), C.int(avr.Num)}
}

// Synthetic video generator. It produces 25 iteratable frames.
// Used for tests.
func GenSyntVideoNewFrame(w, h int, fmt int32) chan *Frame {
	yield := make(chan *Frame)

	go func() {
		defer close(yield)
		for i := 0; i < 25; i++ {
			frame := NewFrame().SetWidth(w).SetHeight(h).SetFormat(fmt)

			if err := frame.ImgAlloc(); err != nil {
				return
			}

			for y := 0; y < h; y++ {
				for x := 0; x < w; x++ {
					frame.SetData(0, y*frame.LineSize(0)+x, x+y+i*3)
				}
			}

			// Cb and Cr
			for y := 0; y < h/2; y++ {
				for x := 0; x < w/2; x++ {
					frame.SetData(1, y*frame.LineSize(1)+x, 128+y+i*2)
					frame.SetData(2, y*frame.LineSize(2)+x, 64+x+i*5)
				}
			}

			yield <- frame
		}
	}()
	return yield
}

// tmp
func GenSyntVideoN(N, w, h int, fmt int32) chan *Frame {
	yield := make(chan *Frame)

	go func() {
		defer close(yield)
		for i := 0; i < N; i++ {
			frame := NewFrame().SetWidth(w).SetHeight(h).SetFormat(fmt)

			if err := frame.ImgAlloc(); err != nil {
				return
			}

			for y := 0; y < h; y++ {
				for x := 0; x < w; x++ {
					frame.SetData(0, y*frame.LineSize(0)+x, x+y+i*3)
				}
			}

			// Cb and Cr
			for y := 0; y < h/2; y++ {
				for x := 0; x < w/2; x++ {
					frame.SetData(1, y*frame.LineSize(1)+x, 128+y+i*2)
					frame.SetData(2, y*frame.LineSize(2)+x, 64+x+i*5)
				}
			}

			yield <- frame
		}
	}()
	return yield
}
