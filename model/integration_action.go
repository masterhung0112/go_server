package model

type PostAction struct {
	// A unique Action ID. If not set, generated automatically.
	Id string `json:"id,omitempty"`

	// The type of the interactive element. Currently supported are
	// "select" and "button".
	Type string `json:"type,omitempty"`

	// The text on the button, or in the select placeholder.
	Name string `json:"name,omitempty"`

	// If the action is disabled.
	Disabled bool `json:"disabled,omitempty"`

	// Style defines a text and border style.
	// Supported values are "default", "primary", "success", "good", "warning", "danger"
	// and any hex color.
	Style string `json:"style,omitempty"`

	// DataSource indicates the data source for the select action. If left
	// empty, the select is populated from Options. Other supported values
	// are "users" and "channels".
	DataSource string `json:"data_source,omitempty"`

	// Options contains the values listed in a select dropdown on the post.
	Options []*PostActionOptions `json:"options,omitempty"`

	// DefaultOption contains the option, if any, that will appear as the
	// default selection in a select box. It has no effect when used with
	// other types of actions.
	DefaultOption string `json:"default_option,omitempty"`

	// Defines the interaction with the backend upon a user action.
	// Integration contains Context, which is private plugin data;
	// Integrations are stripped from Posts when they are sent to the
	// client, or are encrypted in a Cookie.
	Integration *PostActionIntegration `json:"integration,omitempty"`
	Cookie      string                 `json:"cookie,omitempty" db:"-"`
}

// PostActionCookie is set by the server, serialized and encrypted into
// PostAction.Cookie. The clients should hold on to it, and include it with
// subsequent DoPostAction requests.  This allows the server to access the
// action metadata even when it's not available in the database, for ephemeral
// posts.
type PostActionCookie struct {
	Type        string                 `json:"type,omitempty"`
	PostId      string                 `json:"post_id,omitempty"`
	RootPostId  string                 `json:"root_post_id,omitempty"`
	ChannelId   string                 `json:"channel_id,omitempty"`
	DataSource  string                 `json:"data_source,omitempty"`
	Integration *PostActionIntegration `json:"integration,omitempty"`
	RetainProps map[string]interface{} `json:"retain_props,omitempty"`
	RemoveProps []string               `json:"remove_props,omitempty"`
}

type PostActionOptions struct {
	Text  string `json:"text"`
	Value string `json:"value"`
}

type PostActionIntegration struct {
	URL     string                 `json:"url,omitempty"`
	Context map[string]interface{} `json:"context,omitempty"`
}

type PostActionIntegrationRequest struct {
	UserId      string                 `json:"user_id"`
	UserName    string                 `json:"user_name"`
	ChannelId   string                 `json:"channel_id"`
	ChannelName string                 `json:"channel_name"`
	TeamId      string                 `json:"team_id"`
	TeamName    string                 `json:"team_domain"`
	PostId      string                 `json:"post_id"`
	TriggerId   string                 `json:"trigger_id"`
	Type        string                 `json:"type"`
	DataSource  string                 `json:"data_source"`
	Context     map[string]interface{} `json:"context,omitempty"`
}

type PostActionIntegrationResponse struct {
	Update           *Post  `json:"update"`
	EphemeralText    string `json:"ephemeral_text"`
	SkipSlackParsing bool   `json:"skip_slack_parsing"` // Set to `true` to skip the Slack-compatibility handling of Text.
}

type PostActionAPIResponse struct {
	Status    string `json:"status"` // needed to maintain backwards compatibility
	TriggerId string `json:"trigger_id"`
}

type Dialog struct {
	CallbackId       string          `json:"callback_id"`
	Title            string          `json:"title"`
	IntroductionText string          `json:"introduction_text"`
	IconURL          string          `json:"icon_url"`
	Elements         []DialogElement `json:"elements"`
	SubmitLabel      string          `json:"submit_label"`
	NotifyOnCancel   bool            `json:"notify_on_cancel"`
	State            string          `json:"state"`
}

type DialogElement struct {
	DisplayName string               `json:"display_name"`
	Name        string               `json:"name"`
	Type        string               `json:"type"`
	SubType     string               `json:"subtype"`
	Default     string               `json:"default"`
	Placeholder string               `json:"placeholder"`
	HelpText    string               `json:"help_text"`
	Optional    bool                 `json:"optional"`
	MinLength   int                  `json:"min_length"`
	MaxLength   int                  `json:"max_length"`
	DataSource  string               `json:"data_source"`
	Options     []*PostActionOptions `json:"options"`
}

type OpenDialogRequest struct {
	TriggerId string `json:"trigger_id"`
	URL       string `json:"url"`
	Dialog    Dialog `json:"dialog"`
}

type SubmitDialogRequest struct {
	Type       string                 `json:"type"`
	URL        string                 `json:"url,omitempty"`
	CallbackId string                 `json:"callback_id"`
	State      string                 `json:"state"`
	UserId     string                 `json:"user_id"`
	ChannelId  string                 `json:"channel_id"`
	TeamId     string                 `json:"team_id"`
	Submission map[string]interface{} `json:"submission"`
	Cancelled  bool                   `json:"cancelled"`
}

type SubmitDialogResponse struct {
	Error  string            `json:"error,omitempty"`
	Errors map[string]string `json:"errors,omitempty"`
}

func (o *Post) StripActionIntegrations() {
	attachments := o.Attachments()
	if o.GetProp("attachments") != nil {
		o.AddProp("attachments", attachments)
	}
	for _, attachment := range attachments {
		for _, action := range attachment.Actions {
			action.Integration = nil
		}
	}
}
