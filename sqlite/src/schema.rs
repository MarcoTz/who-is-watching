pub struct ColumnDescription {
    name: String,
    ty: String,
    nullable: bool,
}

impl ColumnDescription {
    pub fn new(name: &str, ty: &str, nullable: bool) -> ColumnDescription {
        ColumnDescription {
            name: name.to_owned(),
            ty: ty.to_owned(),
            nullable,
        }
    }

    pub fn create_sql(&self) -> String {
        format!(
            "{} {} {}",
            self.name,
            self.ty,
            if self.nullable { "" } else { "not null" }
        )
    }
}

pub enum Tables {
    Watchers,
    Shows,
    ShowsSeasons,
    ShowsWatchers,
}

impl Tables {
    pub fn all() -> [Tables; 4] {
        [
            Tables::Watchers,
            Tables::Shows,
            Tables::ShowsSeasons,
            Tables::ShowsWatchers,
        ]
    }
    pub fn name(&self) -> &str {
        match self {
            Tables::Watchers => "watchers",
            Tables::Shows => "shows",
            Tables::ShowsSeasons => "shows_seasons",
            Tables::ShowsWatchers => "shows_watchers",
        }
    }
    pub fn cols(&self) -> Vec<ColumnDescription> {
        match self {
            Tables::Watchers => vec![
                ColumnDescription::new("id", "integer", false),
                ColumnDescription::new("name", "text", false),
            ],
            Tables::Shows => vec![
                ColumnDescription::new("id", "integer", false),
                ColumnDescription::new("name", "text", false),
            ],
            Tables::ShowsSeasons => vec![
                ColumnDescription::new("show_id", "integer", false),
                ColumnDescription::new("season_num", "integer", false),
                ColumnDescription::new("num_episodes", "integer", false),
            ],
            Tables::ShowsWatchers => vec![
                ColumnDescription::new("watcher_id", "integer", false),
                ColumnDescription::new("show_id", "integer", false),
                ColumnDescription::new("season", "integer", false),
                ColumnDescription::new("episode", "integer", false),
            ],
        }
    }

    pub fn create_sql(&self) -> String {
        format!(
            "CREATE TABLE {} ({});",
            self.name(),
            self.cols()
                .iter()
                .map(|col| col.create_sql())
                .collect::<Vec<String>>()
                .join(", ")
        )
    }
}
