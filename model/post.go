package model

import (
	"strings"
	"regexp"
	"unicode/utf8"
	"net/http"
	"errors"
	"io"
	"encoding/json"
	"sync"
)

const (
	POST_SYSTEM_MESSAGE_PREFIX  = "system_"
	POST_DEFAULT                = ""
	POST_SLACK_ATTACHMENT       = "slack_attachment"
	POST_SYSTEM_GENERIC         = "system_generic"
	POST_JOIN_LEAVE             = "system_join_leave" // Deprecated, use POST_JOIN_CHANNEL or POST_LEAVE_CHANNEL instead
	POST_JOIN_CHANNEL           = "system_join_channel"
	POST_GUEST_JOIN_CHANNEL     = "system_guest_join_channel"
	POST_LEAVE_CHANNEL          = "system_leave_channel"
	POST_JOIN_TEAM              = "system_join_team"
	POST_LEAVE_TEAM             = "system_leave_team"
	POST_AUTO_RESPONDER         = "system_auto_responder"
	POST_ADD_REMOVE             = "system_add_remove" // Deprecated, use POST_ADD_TO_CHANNEL or POST_REMOVE_FROM_CHANNEL instead
	POST_ADD_TO_CHANNEL         = "system_add_to_channel"
	POST_ADD_GUEST_TO_CHANNEL   = "system_add_guest_to_chan"
	POST_REMOVE_FROM_CHANNEL    = "system_remove_from_channel"
	POST_MOVE_CHANNEL           = "system_move_channel"
	POST_ADD_TO_TEAM            = "system_add_to_team"
	POST_REMOVE_FROM_TEAM       = "system_remove_from_team"
	POST_HEADER_CHANGE          = "system_header_change"
	POST_DISPLAYNAME_CHANGE     = "system_displayname_change"
	POST_CONVERT_CHANNEL        = "system_convert_channel"
	POST_PURPOSE_CHANGE         = "system_purpose_change"
	POST_CHANNEL_DELETED        = "system_channel_deleted"
	POST_CHANNEL_RESTORED       = "system_channel_restored"
	POST_EPHEMERAL              = "system_ephemeral"
	POST_CHANGE_CHANNEL_PRIVACY = "system_change_chan_privacy"
	POST_ADD_BOT_TEAMS_CHANNELS = "add_bot_teams_channels"
	POST_FILEIDS_MAX_RUNES      = 150
	POST_FILENAMES_MAX_RUNES    = 4000
	POST_HASHTAGS_MAX_RUNES     = 1000
	POST_MESSAGE_MAX_RUNES_V1   = 4000
	POST_MESSAGE_MAX_BYTES_V2   = 65535                         // Maximum size of a TEXT column in MySQL
	POST_MESSAGE_MAX_RUNES_V2   = POST_MESSAGE_MAX_BYTES_V2 / 4 // Assume a worst-case representation
	POST_PROPS_MAX_RUNES        = 8000
	POST_PROPS_MAX_USER_RUNES   = POST_PROPS_MAX_RUNES - 400 // Leave some room for system / pre-save modifications
	POST_CUSTOM_TYPE_PREFIX     = "custom_"
	POST_ME                     = "me"
	PROPS_ADD_CHANNEL_MEMBER    = "add_channel_member"

	POST_PROPS_ADDED_USER_ID       = "addedUserId"
	POST_PROPS_DELETE_BY           = "deleteBy"
	POST_PROPS_OVERRIDE_ICON_URL   = "override_icon_url"
	POST_PROPS_OVERRIDE_ICON_EMOJI = "override_icon_emoji"

	POST_PROPS_MENTION_HIGHLIGHT_DISABLED = "mentionHighlightDisabled"
	POST_PROPS_GROUP_HIGHLIGHT_DISABLED   = "disable_group_highlight"
	POST_SYSTEM_WARN_METRIC_STATUS        = "warn_metric_status"
)

var AT_MENTION_PATTEN = regexp.MustCompile(`\B@`)

type Post struct {
	Id         string `json:"id"`
	CreateAt   int64  `json:"create_at"`
	UpdateAt   int64  `json:"update_at"`
	EditAt     int64  `json:"edit_at"`
	DeleteAt   int64  `json:"delete_at"`
	IsPinned   bool   `json:"is_pinned"`
	UserId     string `json:"user_id"`
	ChannelId  string `json:"channel_id"`
	RootId     string `json:"root_id"`
	ParentId   string `json:"parent_id"`
	OriginalId string `json:"original_id"`

	Message string `json:"message"`
	// MessageSource will contain the message as submitted by the user if Message has been modified
	// by Mattermost for presentation (e.g if an image proxy is being used). It should be used to
	// populate edit boxes if present.
	MessageSource string `json:"message_source,omitempty" db:"-"`

	Type          string          `json:"type"`
	propsMu       sync.RWMutex    `db:"-"`       // Unexported mutex used to guard Post.Props.
	Props         StringInterface `json:"props"` // Deprecated: use GetProps()
	Hashtags      string          `json:"hashtags"`
	Filenames     StringArray     `json:"filenames,omitempty"` // Deprecated, do not use this field any more
	FileIds       StringArray     `json:"file_ids,omitempty"`
	PendingPostId string          `json:"pending_post_id" db:"-"`
	HasReactions  bool            `json:"has_reactions,omitempty"`

	// Transient data populated before sending a post to the client
	ReplyCount int64         `json:"reply_count" db:"-"`
	Metadata   *PostMetadata `json:"metadata,omitempty" db:"-"`
}

// ShallowCopy is an utility function to shallow copy a Post to the given
// destination without touching the internal RWMutex.
func (o *Post) ShallowCopy(dst *Post) error {
	if dst == nil {
		return errors.New("dst cannot be nil")
	}
	o.propsMu.RLock()
	defer o.propsMu.RUnlock()
	dst.propsMu.Lock()
	defer dst.propsMu.Unlock()
	dst.Id = o.Id
	dst.CreateAt = o.CreateAt
	dst.UpdateAt = o.UpdateAt
	dst.EditAt = o.EditAt
	dst.DeleteAt = o.DeleteAt
	dst.IsPinned = o.IsPinned
	dst.UserId = o.UserId
	dst.ChannelId = o.ChannelId
	dst.RootId = o.RootId
	dst.ParentId = o.ParentId
	dst.OriginalId = o.OriginalId
	dst.Message = o.Message
	dst.MessageSource = o.MessageSource
	dst.Type = o.Type
	dst.Props = o.Props
	dst.Hashtags = o.Hashtags
	dst.Filenames = o.Filenames
	dst.FileIds = o.FileIds
	dst.PendingPostId = o.PendingPostId
	dst.HasReactions = o.HasReactions
	dst.ReplyCount = o.ReplyCount
	dst.Metadata = o.Metadata
	return nil
}

// Clone shallowly copies the post and returns the copy.
func (o *Post) Clone() *Post {
	copy := &Post{}
	o.ShallowCopy(copy)
	return copy
}

func (o *Post) ToJson() string {
	copy := o.Clone()
	copy.StripActionIntegrations()
	b, _ := json.Marshal(copy)
	return string(b)
}

func (o *Post) ToUnsanitizedJson() string {
	b, _ := json.Marshal(o)
	return string(b)
}

type GetPostsSinceOptions struct {
	ChannelId        string
	Time             int64
	SkipFetchThreads bool
}

type GetPostsOptions struct {
	ChannelId        string
	PostId           string
	Page             int
	PerPage          int
	SkipFetchThreads bool
}

func PostFromJson(data io.Reader) *Post {
	var o *Post
	json.NewDecoder(data).Decode(&o)
	return o
}

func (o *Post) Etag() string {
	return Etag(o.Id, o.UpdateAt)
}

func (o *Post) IsValid(maxPostSize int) *AppError {
	if !IsValidId(o.Id) {
		return NewAppError("Post.IsValid", "model.post.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}

	if o.CreateAt == 0 {
		return NewAppError("Post.IsValid", "model.post.is_valid.create_at.app_error", nil, "id="+o.Id, http.StatusBadRequest)
	}

	if o.UpdateAt == 0 {
		return NewAppError("Post.IsValid", "model.post.is_valid.update_at.app_error", nil, "id="+o.Id, http.StatusBadRequest)
	}

	if !IsValidId(o.UserId) {
		return NewAppError("Post.IsValid", "model.post.is_valid.user_id.app_error", nil, "", http.StatusBadRequest)
	}

	if !IsValidId(o.ChannelId) {
		return NewAppError("Post.IsValid", "model.post.is_valid.channel_id.app_error", nil, "", http.StatusBadRequest)
	}

	if !(IsValidId(o.RootId) || len(o.RootId) == 0) {
		return NewAppError("Post.IsValid", "model.post.is_valid.root_id.app_error", nil, "", http.StatusBadRequest)
	}

	if !(IsValidId(o.ParentId) || len(o.ParentId) == 0) {
		return NewAppError("Post.IsValid", "model.post.is_valid.parent_id.app_error", nil, "", http.StatusBadRequest)
	}

	if len(o.ParentId) == 26 && len(o.RootId) == 0 {
		return NewAppError("Post.IsValid", "model.post.is_valid.root_parent.app_error", nil, "", http.StatusBadRequest)
	}

	if !(len(o.OriginalId) == 26 || len(o.OriginalId) == 0) {
		return NewAppError("Post.IsValid", "model.post.is_valid.original_id.app_error", nil, "", http.StatusBadRequest)
	}

	if utf8.RuneCountInString(o.Message) > maxPostSize {
		return NewAppError("Post.IsValid", "model.post.is_valid.msg.app_error", nil, "id="+o.Id, http.StatusBadRequest)
	}

	if utf8.RuneCountInString(o.Hashtags) > POST_HASHTAGS_MAX_RUNES {
		return NewAppError("Post.IsValid", "model.post.is_valid.hashtags.app_error", nil, "id="+o.Id, http.StatusBadRequest)
	}

	switch o.Type {
	case
		POST_DEFAULT,
		POST_SYSTEM_GENERIC,
		POST_JOIN_LEAVE,
		POST_AUTO_RESPONDER,
		POST_ADD_REMOVE,
		POST_JOIN_CHANNEL,
		POST_GUEST_JOIN_CHANNEL,
		POST_LEAVE_CHANNEL,
		POST_JOIN_TEAM,
		POST_LEAVE_TEAM,
		POST_ADD_TO_CHANNEL,
		POST_ADD_GUEST_TO_CHANNEL,
		POST_REMOVE_FROM_CHANNEL,
		POST_MOVE_CHANNEL,
		POST_ADD_TO_TEAM,
		POST_REMOVE_FROM_TEAM,
		POST_SLACK_ATTACHMENT,
		POST_HEADER_CHANGE,
		POST_PURPOSE_CHANGE,
		POST_DISPLAYNAME_CHANGE,
		POST_CONVERT_CHANNEL,
		POST_CHANNEL_DELETED,
		POST_CHANNEL_RESTORED,
		POST_CHANGE_CHANNEL_PRIVACY,
		POST_ME,
		POST_ADD_BOT_TEAMS_CHANNELS,
		POST_SYSTEM_WARN_METRIC_STATUS:
	default:
		if !strings.HasPrefix(o.Type, POST_CUSTOM_TYPE_PREFIX) {
			return NewAppError("Post.IsValid", "model.post.is_valid.type.app_error", nil, "id="+o.Type, http.StatusBadRequest)
		}
	}

	if utf8.RuneCountInString(ArrayToJson(o.Filenames)) > POST_FILENAMES_MAX_RUNES {
		return NewAppError("Post.IsValid", "model.post.is_valid.filenames.app_error", nil, "id="+o.Id, http.StatusBadRequest)
	}

	if utf8.RuneCountInString(ArrayToJson(o.FileIds)) > POST_FILEIDS_MAX_RUNES {
		return NewAppError("Post.IsValid", "model.post.is_valid.file_ids.app_error", nil, "id="+o.Id, http.StatusBadRequest)
	}

	if utf8.RuneCountInString(StringInterfaceToJson(o.GetProps())) > POST_PROPS_MAX_RUNES {
		return NewAppError("Post.IsValid", "model.post.is_valid.props.app_error", nil, "id="+o.Id, http.StatusBadRequest)
	}

	return nil
}

func (o *Post) DelProp(key string) {
	o.propsMu.Lock()
	defer o.propsMu.Unlock()
	propsCopy := make(map[string]interface{}, len(o.Props)-1)
	for k, v := range o.Props {
		propsCopy[k] = v
	}
	delete(propsCopy, key)
	o.Props = propsCopy
}

func (o *Post) AddProp(key string, value interface{}) {
	o.propsMu.Lock()
	defer o.propsMu.Unlock()
	propsCopy := make(map[string]interface{}, len(o.Props)+1)
	for k, v := range o.Props {
		propsCopy[k] = v
	}
	propsCopy[key] = value
	o.Props = propsCopy
}

func (o *Post) GetProps() StringInterface {
	o.propsMu.RLock()
	defer o.propsMu.RUnlock()
	return o.Props
}

func (o *Post) SetProps(props StringInterface) {
	o.propsMu.Lock()
	defer o.propsMu.Unlock()
	o.Props = props
}

func (o *Post) GetProp(key string) interface{} {
	o.propsMu.RLock()
	defer o.propsMu.RUnlock()
	return o.Props[key]
}

func (o *Post) IsSystemMessage() bool {
	return len(o.Type) >= len(POST_SYSTEM_MESSAGE_PREFIX) && o.Type[:len(POST_SYSTEM_MESSAGE_PREFIX)] == POST_SYSTEM_MESSAGE_PREFIX
}

func (o *Post) Attachments() []*SlackAttachment {
	if attachments, ok := o.GetProp("attachments").([]*SlackAttachment); ok {
		return attachments
	}
	var ret []*SlackAttachment
	if attachments, ok := o.GetProp("attachments").([]interface{}); ok {
		for _, attachment := range attachments {
			if enc, err := json.Marshal(attachment); err == nil {
				var decoded SlackAttachment
				if json.Unmarshal(enc, &decoded) == nil {
					ret = append(ret, &decoded)
				}
			}
		}
	}
	return ret
}

func (o *Post) SanitizeProps() {
	membersToSanitize := []string{
		PROPS_ADD_CHANNEL_MEMBER,
	}

	for _, member := range membersToSanitize {
		if _, ok := o.GetProps()[member]; ok {
			o.DelProp(member)
		}
	}
}

func (o *Post) PreSave() {
	if o.Id == "" {
		o.Id = NewId()
	}

	o.OriginalId = ""

	if o.CreateAt == 0 {
		o.CreateAt = GetMillis()
	}

	o.UpdateAt = o.CreateAt
	o.PreCommit()
}

func (o *Post) PreCommit() {
	if o.GetProps() == nil {
		o.SetProps(make(map[string]interface{}))
	}

	if o.Filenames == nil {
		o.Filenames = []string{}
	}

	if o.FileIds == nil {
		o.FileIds = []string{}
	}

	o.GenerateActionIds()

	// There's a rare bug where the client sends up duplicate FileIds so protect against that
	o.FileIds = RemoveDuplicateStrings(o.FileIds)
}
