package main

import (
	"fmt"
	"image"

	"github.com/mineatar-io/skin-render"
)

var (
	RenderTypeFullBody  = "fullbody"
	RenderTypeFrontBody = "frontbody"
	RenderTypeBackBody  = "backbody"
	RenderTypeLeftBody  = "leftbody"
	RenderTypeRightBody = "rightbody"
	RenderTypeFace      = "face"
	RenderTypeHead      = "head"
)

// Render will render the image using the specified details and return the result.
func Render(renderType, uuid string, rawSkin *image.NRGBA, isSlim bool, opts *QueryParams) ([]byte, bool, error) {
	if config.Cache.EnableLocks {
		mutex := r.NewMutex(fmt.Sprintf("render-lock:%s-%d-%t-%s", renderType, opts.Scale, opts.Overlay, uuid))
		mutex.Lock()

		defer mutex.Unlock()
	}

	// Fetch the existing render from cache if it exists
	{
		cache, err := GetCachedRenderResult(renderType, uuid, opts)

		if err != nil {
			return nil, false, err
		}

		if cache != nil {
			return cache, true, nil
		}
	}

	var (
		result     *image.NRGBA
		renderOpts skin.Options = skin.Options{
			Overlay: opts.Overlay,
			Slim:    isSlim,
			Scale:   opts.Scale,
			Square:  opts.Square,
		}
	)

	// Render the image based on the type provided
	{
		switch renderType {
		case RenderTypeFullBody:
			{
				result = skin.RenderBody(rawSkin, renderOpts)

				break
			}
		case RenderTypeFrontBody:
			{
				result = skin.RenderFrontBody(rawSkin, renderOpts)

				break
			}
		case RenderTypeBackBody:
			{
				result = skin.RenderBackBody(rawSkin, renderOpts)

				break
			}
		case RenderTypeLeftBody:
			{
				result = skin.RenderLeftBody(rawSkin, renderOpts)

				break
			}
		case RenderTypeRightBody:
			{
				result = skin.RenderRightBody(rawSkin, renderOpts)

				break
			}
		case RenderTypeHead:
			{
				result = skin.RenderHead(rawSkin, renderOpts)

				break
			}
		case RenderTypeFace:
			{
				result = skin.RenderFace(rawSkin, renderOpts)

				break
			}
		default:
			panic(fmt.Errorf("unknown render type: %s", renderType))
		}
	}

	var (
		data []byte
		err  error
	)

	// Encode the image into a PNG in byte-array format
	{
		data, err = EncodeImage(result, opts)

		if err != nil {
			return nil, false, err
		}
	}

	// Put the result into the cache for later use
	{
		if err = SetCachedRenderResult(renderType, uuid, opts, data); err != nil {
			return nil, false, err
		}
	}

	return data, false, nil
}
