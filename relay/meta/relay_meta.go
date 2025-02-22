package meta

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/songquanpeng/one-api/common/ctxkey"
	"github.com/songquanpeng/one-api/model"
	"github.com/songquanpeng/one-api/relay/channeltype"
	"github.com/songquanpeng/one-api/relay/relaymode"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Meta struct {
	Mode         int
	ChannelType  int
	ChannelId    int
	TokenId      int
	TokenName    string
	UserId       int
	Group        string
	ModelMapping map[string]string
	// BaseURL is the proxy url set in the channel config
	BaseURL  string
	APIKey   string
	APIType  int
	Config   model.ChannelConfig
	IsStream bool
	// OriginModelName is the model name from the raw user request
	OriginModelName string
	// ActualModelName is the model name after mapping
	ActualModelName    string
	RequestURLPath     string
	PromptTokens       int // only for DoResponse
	ForcedSystemPrompt string
	StartTime          time.Time
}

func GetByContext(c *gin.Context) *Meta {
	meta := Meta{
		Mode:               relaymode.GetByPath(c.Request.URL.Path),
		ChannelType:        c.GetInt(ctxkey.Channel),
		ChannelId:          c.GetInt(ctxkey.ChannelId),
		TokenId:            c.GetInt(ctxkey.TokenId),
		TokenName:          c.GetString(ctxkey.TokenName),
		UserId:             c.GetInt(ctxkey.Id),
		Group:              c.GetString(ctxkey.Group),
		ModelMapping:       c.GetStringMapString(ctxkey.ModelMapping),
		OriginModelName:    c.GetString(ctxkey.RequestModel),
		BaseURL:            c.GetString(ctxkey.BaseURL),
		APIKey:             strings.TrimPrefix(c.Request.Header.Get("Authorization"), "Bearer "),
		RequestURLPath:     c.Request.URL.String(),
		ForcedSystemPrompt: c.GetString(ctxkey.SystemPrompt),
		StartTime:          time.Now(),
	}
	cfg, ok := c.Get(ctxkey.Config)
	if ok {
		meta.Config = cfg.(model.ChannelConfig)
	}
	if meta.BaseURL == "" {
		meta.BaseURL = channeltype.ChannelBaseURLs[meta.ChannelType]
	}
	meta.APIType = channeltype.ToAPIType(meta.ChannelType)
	return &meta
}

func (m *Meta) CopyToSpan(span trace.Span) {
	span.SetAttributes(
		attribute.Int("meta.mode", m.Mode),
		attribute.Int("meta.channel_type", m.ChannelType),
		attribute.Int("meta.channel_id", m.ChannelId),
		attribute.Int("meta.token_id", m.TokenId),
		attribute.String("meta.token_name", m.TokenName),
		attribute.Int("meta.user_id", m.UserId),
		attribute.String("meta.group", m.Group),
		attribute.String("meta.base_url", m.BaseURL),
		attribute.String("meta.api_key", m.APIKey),
		attribute.Int("meta.api_type", m.APIType),
		attribute.Bool("meta.is_stream", m.IsStream),
		attribute.String("meta.origin_model_name", m.OriginModelName),
		attribute.String("meta.actual_model_name", m.ActualModelName),
		attribute.String("meta.request_url_path", m.RequestURLPath),
		attribute.String("meta.forced_system_prompt", m.ForcedSystemPrompt),
		attribute.Int("meta.prompt_tokens", m.PromptTokens),
	)
}
