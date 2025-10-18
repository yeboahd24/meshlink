# 🔧 Fix Video Corruption Issue

## Problem Identified:
**Pixel format mismatch** - FFmpeg encodes in one format, GUI decodes in another, causing the green striped corruption.

## Root Cause:
```go
// Current FFmpeg encoding (in pkg/media/ffmpeg.go):
"format": "h264"           // Raw H.264 stream
"pix_fmt": "yuv420p"       // YUV color space

// GUI decoder expects:
"format": "rawvideo"       // Raw RGB frames  
"pix_fmt": "rgb24"         // RGB color space
```

## The Fix:

### Option 1: Fix Encoder Output (Recommended)
Change FFmpeg to output RGB frames instead of H.264:

```go
// In pkg/media/ffmpeg.go - Change output format:
Output("pipe:",
    ffmpeg.KwArgs{
        "format":  "rawvideo",    // Raw video instead of h264
        "pix_fmt": "rgb24",       // RGB instead of YUV
        "s":       "640x480",     // Fixed size
        "r":       "30",          // Frame rate
    })
```

### Option 2: Fix Decoder Input
Change decoder to properly handle H.264:

```go
// In pkg/media/decoder.go - Add H.264 decoding:
Input("pipe:",
    ffmpeg.KwArgs{
        "format": "h264",         // Expect H.264 input
    }).
Output("pipe:",
    ffmpeg.KwArgs{
        "format":  "rawvideo",    // Output raw frames
        "pix_fmt": "rgb24",       // Convert to RGB
        "s":       "640x480",     // Fixed size
    })
```

## Quick Test Fix:

### Update pkg/media/ffmpeg.go:
```go
err := stream.
    Output("pipe:",
        ffmpeg.KwArgs{
            "format":  "rawvideo",  // Change from "h264"
            "pix_fmt": "rgb24",     // Change from yuv420p
            "s":       "640x480",   // Add fixed size
            "r":       "30",        // Frame rate
        }).
    WithOutput(f.output).
    Run()
```

## Why This Happens:
1. **FFmpeg encodes** H.264 compressed video
2. **GUI expects** raw RGB pixel data
3. **Mismatch causes** corruption (green stripes)
4. **Fix aligns** both sides to same format

## Expected Result After Fix:
- ✅ **Clear video** instead of corruption
- ✅ **Proper colors** (no green stripes)
- ✅ **Real camera** feed (when available)
- ✅ **Test pattern** displays correctly

The P2P streaming works perfectly - we just need to fix the video format mismatch!