package mediautil

import (
	"fmt"
	"io"
	"strings"

	"github.com/yapingcat/gomedia/go-mp4"
)

func split(mime string) (primary string, sub string, ok bool) {
	types := strings.Split(mime, "/")

	if len(types) != 2 {
		return "", "", false
	}

	return types[0], types[1], true
}

func IsVideo(mime string) bool {
	primary, _, ok := split(mime)

	return primary == "video" && ok
}

func IsAudio(mime string) bool {
	primary, _, ok := split(mime)

	return primary == "audio" && ok
}

func IsImage(mime string) bool {
	primary, _, ok := split(mime)

	return primary == "image" && ok
}

// GetMP4Info returns duration, width, height, error
// Test patch: original tdl accepted only H.264 MP4 tracks. This fallback also
// accepts MP4 tracks with width/height metadata, including VP9-in-MP4.
func GetMP4Info(r io.ReadSeeker) (int, int, int, error) {
	d := mp4.CreateMp4Demuxer(r)

	tracks, err := d.ReadHead()
	if err != nil {
		return 0, 0, 0, err
	}

	info := d.GetMp4Info()
	duration := int(info.Duration / info.Timescale)

	// Preserve original behavior when H.264 exists.
	for _, track := range tracks {
		if track.Cid == mp4.MP4_CODEC_H264 {
			return duration, int(track.Width), int(track.Height), nil
		}
	}

	// Fallback for VP9/other MP4 video tracks.
	for _, track := range tracks {
		if track.Width > 0 && track.Height > 0 {
			return duration, int(track.Width), int(track.Height), nil
		}
	}

	return 0, 0, 0, fmt.Errorf("no video track with width/height found")
}
