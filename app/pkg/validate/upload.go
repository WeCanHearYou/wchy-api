package validate

import (
	"context"

	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/pkg/i18n"
	"github.com/goenning/imagic"
)

// MaxDimensionSize is the max width/height of an image. If image is bigger than this, it'll be resized.
const MaxDimensionSize = 1500

// MultiImageUploadOpts arguments to validate mulitple image upload process
type MultiImageUploadOpts struct {
	MaxUploads   int
	IsRequired   bool
	MinWidth     int
	MinHeight    int
	ExactRatio   bool
	MaxKilobytes int
}

// ImageUploadOpts arguments to validate given upload
type ImageUploadOpts struct {
	IsRequired   bool
	MinWidth     int
	MinHeight    int
	ExactRatio   bool
	MaxKilobytes int
}

//MultiImageUpload validates multiple image uploads
func MultiImageUpload(ctx context.Context, currentAttachments []string, uploads []*dto.ImageUpload, opts MultiImageUploadOpts) ([]string, error) {
	if currentAttachments == nil {
		currentAttachments = []string{}
	}

	totalCount := len(currentAttachments)

	for _, upload := range uploads {
		if upload.Remove {
			for _, attachment := range currentAttachments {
				if attachment == upload.BlobKey {
					totalCount--
				}
			}
		} else if upload.Upload != nil {
			totalCount++
		}

		messages, err := ImageUpload(ctx, upload, ImageUploadOpts{
			IsRequired:   opts.IsRequired,
			MinWidth:     opts.MinWidth,
			MinHeight:    opts.MinHeight,
			ExactRatio:   opts.ExactRatio,
			MaxKilobytes: opts.MaxKilobytes,
		})
		if err != nil {
			return nil, err
		}
		if len(messages) > 0 {
			return messages, nil
		}
	}

	if totalCount > opts.MaxUploads {
		return []string{i18n.T(ctx, "validation.custom.maxattachments", i18n.Params{"number": opts.MaxUploads})}, nil
	}

	return []string{}, nil
}

//ImageUpload validates given image upload
func ImageUpload(ctx context.Context, upload *dto.ImageUpload, opts ImageUploadOpts) ([]string, error) {
	messages := []string{}

	if opts.IsRequired {
		if upload == nil || (upload.BlobKey == "" && upload.Upload == nil) || upload.Remove {
			messages = append(messages, i18n.T(ctx, "validation.required",
				i18n.Params{"name": i18n.T(ctx, "property.image")},
			))
		}
	}

	if upload != nil && upload.Upload != nil && len(upload.Upload.Content) > 0 {
		logo, err := imagic.Parse(upload.Upload.Content)
		if err != nil {
			if err == imagic.ErrNotSupported {
				messages = append(messages, i18n.T(ctx, "validation.custom.unsupportedfileformat"))
			} else {
				return nil, err
			}
		} else {

			if logo.Width < opts.MinWidth || logo.Height < opts.MinHeight {
				messages = append(messages, i18n.T(ctx, "validation.custom.minimagedimensions",
					i18n.Params{"width": opts.MinWidth},
					i18n.Params{"height": opts.MinHeight},
				))
			}

			if opts.ExactRatio && logo.Width != logo.Height {
				messages = append(messages, i18n.T(ctx, "validation.custom.imagesquareratio"))
			}

			if logo.Size > (opts.MaxKilobytes * 1024) {
				messages = append(messages, i18n.T(ctx, "validation.custom.maximagesize",
					i18n.Params{"kilobytes": opts.MaxKilobytes},
				))
			}

			if logo.Height > MaxDimensionSize && logo.Width > MaxDimensionSize {
				newImageBytes, err := imagic.Apply(upload.Upload.Content, imagic.Resize(MaxDimensionSize))
				if err != nil {
					return nil, err
				}
				upload.Upload.Content = newImageBytes
			}
		}
	}

	return messages, nil
}
