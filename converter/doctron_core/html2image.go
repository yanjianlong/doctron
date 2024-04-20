package doctron_core

import (
	"context"
	"errors"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

//FormatPng Image compression format (defaults to png).
const FormatPng = page.CaptureScreenshotFormatPng

//FormatJpeg Image compression format (defaults to png).
const FormatJpeg = page.CaptureScreenshotFormatJpeg

//DefaultQuality Compression quality from range [0..100] (jpeg only).
const DefaultQuality = 0

//DefaultFromSurface Capture the screenshot from the surface, rather than the view. Defaults to true.
const DefaultFromSurface = true

//DefaultViewportX Capture the screenshot of a given region only.
// X offset in device independent pixels (dip).
const DefaultViewportX = 0

//DefaultViewportY Y offset in device independent pixels (dip).
const DefaultViewportY = 0

//DefaultViewportWidth Rectangle width in device independent pixels (dip).
const DefaultViewportWidth = 996

//DefaultViewportHeight Rectangle height in device independent pixels (dip).
const DefaultViewportHeight = 996

//DefaultViewportScale Page scale factor.
const DefaultViewportScale = 1

//Html2ImageParams print page as PDF.
type Html2ImageParams struct {
	page.CaptureScreenshotParams
	CustomClip  bool
	WaitingTime int // Waiting time after the page loaded. Default 0 means not wait. unit:Millisecond
}

//NewDefaultHtml2ImageParams default html convert to image params
func NewDefaultHtml2ImageParams() Html2ImageParams {
	return Html2ImageParams{
		CustomClip: false,
		CaptureScreenshotParams: page.CaptureScreenshotParams{
			Format:  FormatPng,
			Quality: DefaultQuality,
			Clip: &page.Viewport{
				X:      DefaultViewportX,
				Y:      DefaultViewportY,
				Width:  DefaultViewportWidth,
				Height: DefaultViewportHeight,
				Scale:  DefaultViewportScale,
			},
			FromSurface: DefaultFromSurface,
		},
		WaitingTime: DefaultWaitingTime,
	}
}

type html2image struct {
	Doctron
}

func (ins *html2image) GetConvertElapsed() time.Duration {
	return ins.convertElapsed
}

func (ins *html2image) Convert() ([]byte, error) {
	start := time.Now()
	defer func() {
		ins.convertElapsed = time.Since(start)
	}()
	var params Html2ImageParams
	params, ok := ins.cc.Params.(Html2ImageParams)
	if !ok {
		return nil, errors.New("wrong html2image params given")
	}
	// options := []chromedp.ExecAllocatorOption{
	// 	chromedp.UserAgent(userAgent),
	// }
	// allocCtx, cancel := chromedp.NewExecAllocator(ins.ctx, options...)
	ctx, cancel := chromedp.NewContext(ins.ctx)
	defer cancel()
	headers := map[string]interface{}{
		"User-Agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 10_3 like Mac OS X) AppleWebKit/602.1.50 (KHTML, like Gecko) CriOS/56.0.2924.75 Mobile/14E5239e Safari/602.1",
	}
	// hbs, _ := json.Marshal(headers)
	if err := chromedp.Run(ctx,
		network.SetExtraHTTPHeaders(network.Headers(headers)),
		chromedp.Navigate(ins.cc.Url),
		chromedp.Sleep(time.Duration(params.WaitingTime)*time.Millisecond),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// if !params.CustomClip {
			// 	// get layout metrics
			// 	_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			// 	if err != nil {
			// 		return err
			// 	}
			// 	params.Clip.X = contentSize.X
			// 	params.Clip.Y = contentSize.Y
			// 	params.Clip.Width = contentSize.Width
			// 	params.Clip.Height = contentSize.Height
			// }

			//width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

			// force viewport emulation
			err := emulation.SetDeviceMetricsOverride(int64(params.Clip.Width), int64(params.Clip.Height), 1, true).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			if err != nil {
				return err
			}
			// capture screenshot
			ins.buf, err = params.Do(ctx)
			return err
		}),
	); err != nil {
		return nil, err
	}

	return ins.buf, nil
}
