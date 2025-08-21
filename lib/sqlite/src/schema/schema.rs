use std::fmt;

pub enum ColumName {
    Id,
    Name,
    ShowId,
    SeasonNum,
    NumEpisodes,
    WatcherId,
    Season,
    Episode,
}

pub struct ColumnDescription {
    name: String,
    ty: String,
    nullable: bool,
    unique: bool,
}

impl ColumnDescription {
    pub fn new(name: &str, ty: &str, nullable: bool, unique: bool) -> ColumnDescription {
        ColumnDescription {
            name: name.to_owned(),
            ty: ty.to_owned(),
            nullable,
            unique,
        }
    }

    pub fn create_sql(&self) -> String {
        format!(
            "{} {} {} {}",
            self.name,
            self.ty,
            if self.unique { "unique" } else { "" },
            if self.nullable { "" } else { "not null" }
        )
    }
}

#[derive(Debug, Copy, Clone)]
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
                ColumnDescription::new("id", "integer", false, true),
                ColumnDescription::new("name", "text", false, false),
            ],
            Tables::Shows => vec![
                ColumnDescription::new("id", "integer", false, true),
                ColumnDescription::new("name", "text", false, false),
            ],
            Tables::ShowsSeasons => vec![
                ColumnDescription::new("show_id", "integer", false, false),
                ColumnDescription::new("season_num", "integer", false, false),
                ColumnDescription::new("num_episodes", "integer", false, false),
            ],
            Tables::ShowsWatchers => vec![
                ColumnDescription::new("watcher_id", "integer", false, false),
                ColumnDescription::new("show_id", "integer", false, false),
                ColumnDescription::new("season", "integer", false, false),
                ColumnDescription::new("episode", "integer", false, false),
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

impl fmt::Display for Tables {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "{}", self.name())
    }
}
