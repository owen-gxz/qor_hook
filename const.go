package hook

var (
	TypeList = []string{
		"string",
		"checkbox",
		"number",
		"datetime",
		"hidden",
		"password",
		"rich_editor",
		"select_one",
		"single_edit",
		"file",
	}

	SqlType = []string{
		"VARCHAR",
		"DATETIME",
		"INT",
	}
	Mysql = "mysql"

	disTables = map[string]bool{
		"auth_identities":       true,
		"qor_activities":        true,
		"qor_jobs":              true,
		"resource_models":       true,
		"resource_table_models": true,
		"scheduled_events":      true,
		"translations":          true,
	}
)
