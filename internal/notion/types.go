package notion

type User struct {
	Object    string  `json:"object"`
	ID        string  `json:"id"`
	Type      *string `json:"type,omitempty"`
	Name      *string `json:"name,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	Person    *Person `json:"person,omitempty"`
	Bot       *Bot    `json:"bot,omitempty"`
}

type Person struct {
	Email string `json:"email"`
}

type Bot struct {
	Owner         Owner  `json:"owner"`
	WorkspaceName string `json:"workspace_name"`
}

type Owner struct {
	Type      string `json:"type"`
	Workspace bool   `json:"workspace"`
}

type Text struct {
	Content string `json:"content"`
	Link    struct {
		URL string `json:"url"`
	} `json:"link,omitempty"`
}
type Equation struct {
	Expression string `json:"expression"`
}

type Mention struct {
	Type               string              `json:"type"`
	DatabaseMention    *DatabaseMention    `json:"database,omitempty"`
	DateMention        *DateMention        `json:"date,omitempty"`
	LinkPreviewMention *LinkPreviewMention `json:"link_preview,omitempty"`
	PageMention        *PageMention        `json:"page,omitempty"`
	TemplateMention    *TemplateMention    `json:"template_mention,omitempty"`
	User               *UserMention        `json:"user,omitempty"`
}

type DatabaseMention struct {
	ID string `json:"id"`
}

type DateMention struct {
	Start string  `json:"start"`
	End   *string `json:"end"`
}

type LinkPreviewMention struct {
	URL string `json:"url"`
}

type PageMention struct {
	ID string `json:"id"`
}

type TemplateMention struct {
	Type                string  `json:"type"`
	TemplateMentionDate *string `json:"template_mention_date,omitempty"`
	TemplateMentionUer  *string `json:"template_mention_user,omitempty"`
}

type UserMention struct {
	ID     string `json:"id"`
	Object string `json:"object"`
}

type Annotations struct {
	Bold          bool   `json:"bold"`
	Italic        bool   `json:"italic"`
	Strikethrough bool   `json:"strikethrough"`
	Underline     bool   `json:"underline"`
	Code          bool   `json:"code"`
	Color         string `json:"color"`
}

type RichText struct {
	Type        string      `json:"type"`
	Text        *Text       `json:"text,omitempty"`
	Mention     *Mention    `json:"mention,omitempty"`
	Equation    *Equation   `json:"equation,omitempty"`
	Annotations Annotations `json:"annotations"`
	PlainText   string      `json:"plain_text"`
	Href        *string     `json:"href,omitempty"`
}

type ExternalFile struct {
	URL string `json:"url"`
}

type File struct {
	URL        string `json:"url"`
	ExpiryTime string `json:"expiry_time"`
}

type Emoji struct {
	Emoji string `json:"emoji"`
}

type Icon struct {
	Type     string        `json:"type"`
	External *ExternalFile `json:"external,omitempty"`
	File     *File         `json:"file,omitempty"`
	Emoji    *string       `json:"emoji,omitempty"`
}
type SelectOption struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Color       string  `json:"color"`
	Description *string `json:"description"`
}

type CheckboxProperty struct{}
type CreatedTimeProperty struct{}
type DateProperty struct{}
type EmailProperty struct{}
type FileProperty struct{}
type FormulaProperty struct {
	Expression string `json:"expression"`
}
type LastEditedByProperty struct{}
type LastEditedTimeProperty struct{}
type MultiSelectPropertyOption struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}
type MultiSelectProperty struct {
	Options []MultiSelectPropertyOption `json:"options"`
}
type NumberProperty struct {
	Format string `json:"format"`
}
type PeopleProperty struct{}
type PhoneNumberProperty struct{}
type RelationProperty struct {
	DatabaseID         string `json:"database_id"`
	SyncedPropertyName string `json:"synced_property_name"`
	SyncedPropertyID   string `json:"synced_property_id"`
}
type RichTextProperty struct{}
type RollupProperty struct {
	Function             string `json:"function"`
	RelationPropertyName string `json:"relation_property_name"`
	RelationPropertyID   string `json:"relation_property_id"`
	RollupPropertyName   string `json:"rollup_property_name"`
	RollupPropertyID     string `json:"rollup_property_id"`
}
type SelectPropertyOption struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}
type SelectProperty struct {
	Options []SelectPropertyOption `json:"options"`
}
type StatusPropertyOption struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}
type StatusPropertyGroup struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Color     string   `json:"color"`
	OptionIDs []string `json:"option_ids"`
}
type StatusProperty struct {
	Options []StatusPropertyOption `json:"options"`
	Groups  []StatusPropertyGroup  `json:"groups"`
}
type TitleProperty struct{}
type URLProperty struct{}

type DatabaseProperty struct {
	ID             string                  `json:"id"`
	Name           string                  `json:"name"`
	Type           string                  `json:"type"`
	Checkbox       *CheckboxProperty       `json:"checkbox,omitempty"`
	CreatedTime    *CreatedTimeProperty    `json:"created_time,omitempty"`
	Date           *DateProperty           `json:"date,omitempty"`
	Email          *EmailProperty          `json:"email,omitempty"`
	Files          *FileProperty           `json:"files,omitempty"`
	Formula        *FormulaProperty        `json:"formula,omitempty"`
	LastEditedBy   *LastEditedByProperty   `json:"last_edited_by,omitempty"`
	LastEditedTime *LastEditedTimeProperty `json:"last_edited_time,omitempty"`
	MultiSelect    *MultiSelectProperty    `json:"multi_select,omitempty"`
	Number         *NumberProperty         `json:"number,omitempty"`
	People         *PeopleProperty         `json:"people,omitempty"`
	PhoneNumber    *PhoneNumberProperty    `json:"phone_number,omitempty"`
	Relation       *RelationProperty       `json:"relation,omitempty"`
	RichText       *RichTextProperty       `json:"rich_text,omitempty"`
	Rollup         *RollupProperty         `json:"rollup,omitempty"`
	Select         *SelectProperty         `json:"select,omitempty"`
	Status         *StatusProperty         `json:"status,omitempty"`
	Title          *TitleProperty          `json:"title,omitempty"`
	URL            *URLProperty            `json:"url,omitempty"`
}

type PartialUser struct {
	Object string `json:"object"`
	ID     string `json:"id"`
}

type Parent struct {
	Type       string  `json:"type"`
	DatabaseID *string `json:"database_id,omitempty"`
	PageID     *string `json:"page_id,omitempty"`
	Workspace  *bool   `json:"workspace,omitempty"`
	BlockID    *string `json:"block_id,omitempty"`
}

type Database struct {
	Object         string      `json:"object"`
	ID             string      `json:"id"`
	CreatedTime    string      `json:"created_time"`
	CreatedBy      PartialUser `json:"created_by"`
	LastEditedTime string      `json:"last_edited_time"`
	LastEditedBy   PartialUser `json:"last_edited_by"`
	Title          []RichText  `json:"title"`
	Description    []RichText  `json:"description"`
	Icon           Icon        `json:"icon"`
	Cover          struct {
		Type     string       `json:"type"`
		External ExternalFile `json:"external"`
	} `json:"cover"`
	Properties map[string]DatabaseProperty `json:"properties"`
	Parent     Parent                      `json:"parent"`
	URL        string                      `json:"url"`
	Archived   bool                        `json:"archived"`
	IsInline   bool                        `json:"is_inline"`
	PublicURL  string                      `json:"public_url"`
}

type DatabaseQueryResponse struct {
	Object         string   `json:"object"`
	Results        []Page   `json:"results"`
	HasMore        bool     `json:"has_more"`
	NextCursor     string   `json:"next_cursor"`
	Type           string   `json:"type"`
	PageOrDatabase struct{} `json:"page_or_database"`
}

// Page Structs

type PageDateProperty struct {
	Start    string  `json:"start"`
	End      *string `json:"end"`
	Timezone *string `json:"timezone,omitempty"`
}

type PageFormulaProperty struct {
	Type    string            `json:"type"`
	Boolean *bool             `json:"boolean,omitempty"`
	Date    *PageDateProperty `json:"date,omitempty"`
	String  *string           `json:"string,omitempty"`
	Number  *float64          `json:"number,omitempty"`
}

type PageRelationProperty struct {
	ID string `json:"id"`
}

type PageRollupProperty struct {
	Type        string            `json:"type"`
	Number      *float64          `json:"number,omitempty"`
	Date        *PageDateProperty `json:"date,omitempty"`
	Array       *[]string         `json:"array,omitempty"`
	Incomplete  *bool             `json:"incomplete,omitempty"`
	Unsupported *interface{}      `json:"unsupported,omitempty"`
	Function    string            `json:"function"`
}

type PageSelectProperty struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type UniqueID struct {
	Number int     `json:"number"`
	Prefix *string `json:"prefix,omitempty"`
}

type Version struct {
	State      string            `json:"state"`
	VerfiyedBy *User             `json:"verified_by,omitempty"`
	Date       *PageDateProperty `json:"date,omitempty"`
}

type PageProperty struct {
	ID             string                  `json:"id"`
	Type           *string                 `json:"type,omitempty"`
	Checkbox       *bool                   `json:"checkbox,omitempty"`
	CreatedBy      *PartialUser            `json:"created_by,omitempty"`
	CreatedTime    *string                 `json:"created_time,omitempty"`
	Date           *PageDateProperty       `json:"date,omitempty"`
	Email          *string                 `json:"email,omitempty"`
	Files          *[]File                 `json:"files,omitempty"`
	Formula        *PageFormulaProperty    `json:"formula,omitempty"`
	LastEditedBy   *User                   `json:"last_edited_by,omitempty"`
	LastEditedTime *string                 `json:"last_edited_time,omitempty"`
	MultiSelect    *[]MultiSelectProperty  `json:"multi_select,omitempty"`
	Number         *float64                `json:"number,omitempty"`
	People         *[]PartialUser          `json:"people,omitempty"`
	PhoneNumber    *string                 `json:"phone_number,omitempty"`
	Relation       *[]PageRelationProperty `json:"relation,omitempty"`
	HasMore        *bool                   `json:"has_more,omitempty"`
	Rollup         *PageRollupProperty     `json:"rollup,omitempty"`
	RichText       *[]RichText             `json:"rich_text,omitempty"`
	Select         *PageSelectProperty     `json:"select,omitempty"`
	Title          *[]RichText             `json:"title,omitempty"`
	URL            *string                 `json:"url,omitempty"`
	UniqueID       *UniqueID               `json:"unique_id,omitempty"`
}

type Page struct {
	Object         string      `json:"object"`
	ID             string      `json:"id"`
	CreatedTime    string      `json:"created_time"`
	CreatedBy      PartialUser `json:"created_by"`
	LastEditedTime string      `json:"last_edited_time"`
	LastEditedBy   PartialUser `json:"last_edited_by"`
	Archived       bool        `json:"archived"`
	Icon           Icon        `json:"icon"`
	Cover          struct {
		Type     string       `json:"type"`
		External ExternalFile `json:"external"`
	} `json:"cover"`
	Properties map[string]PageProperty `json:"properties"`
	Parent     Parent                  `json:"parent"`
	URL        string                  `json:"url"`
	PublicURL  *string                 `json:"public_url"`
}
