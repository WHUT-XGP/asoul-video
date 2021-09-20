// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package route

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	log "unknwon.dev/clog/v2"

	"github.com/asoul-video/asoul-video/internal/context"
	"github.com/asoul-video/asoul-video/internal/db"
	"github.com/asoul-video/asoul-video/pkg/model"
)

type Source struct{}

// NewSourceHandler creates a new Source router.
func NewSourceHandler() *Source {
	return &Source{}
}

// VerifyKey is the middleware which verifies the crawler key.
// Only the real acao we can trust.
func (*Source) VerifyKey(key string) func(ctx context.Context) {
	return func(ctx context.Context) {
		if ctx.Request().Header.Get("Authorization") != key {
			ctx.Error(http.StatusForbidden, errors.New("error key"))
		}
	}
}

func (*Source) Report(ctx context.Context) {
	var req struct {
		Type model.ReportType    `json:"type"`
		Data jsoniter.RawMessage `json:"data"`
	}
	if err := jsoniter.NewDecoder(ctx.Request().Request.Body).Decode(&req); err != nil {
		ctx.Error(http.StatusBadRequest, err)
		return
	}

	switch req.Type {
	case model.ReportTypeUpdateMember:
		var updateMember model.UpdateMember
		if err := jsoniter.Unmarshal(req.Data, &updateMember); err != nil {
			ctx.Error(http.StatusBadRequest, err)
			return
		}

		if err := db.Members.Create(ctx.Request().Context(), db.CreateMemberOptions{
			SecUID:    updateMember.SecUID,
			UID:       updateMember.UID,
			UniqueID:  updateMember.UniqueID,
			ShortUID:  updateMember.ShortUID,
			Name:      updateMember.Name,
			AvatarURL: updateMember.AvatarURL,
			Signature: updateMember.Signature,
		}); err != nil {
			log.Error("Failed to create new member: %v", err)
			ctx.ServerError()
			return
		}

	case model.ReportTypeCreateVideo:
		var createVideo model.CreateVideo
		if err := jsoniter.Unmarshal(req.Data, &createVideo); err != nil {
			ctx.Error(http.StatusBadRequest, err)
			return
		}

		if err := db.Videos.Create(ctx.Request().Context(), createVideo.ID, db.CreateVideoOptions{
			AuthorSecUID:     createVideo.AuthorSecUID,
			Description:      createVideo.Description,
			TextExtra:        createVideo.TextExtra,
			OriginCoverURLs:  createVideo.OriginCoverURLs,
			DynamicCoverURLs: createVideo.DynamicCoverURLs,
			VideoHeight:      createVideo.VideoHeight,
			VideoWidth:       createVideo.VideoWidth,
			VideoDuration:    createVideo.VideoDuration,
			VideoRatio:       createVideo.VideoRatio,
			VideoURLs:        createVideo.VideoURLs,
			VideoCDNURL:      createVideo.VideoCDNURL,
			CreatedAt:        createVideo.CreatedAt,
		}); err != nil {
			log.Error("Failed to create new video: %v", err)
			ctx.ServerError()
			return
		}

	default:
		ctx.Error(http.StatusBadRequest, errors.Errorf("unexpected report type %q", req.Type))
		return
	}

	ctx.ResponseWriter().WriteHeader(http.StatusNoContent)
}