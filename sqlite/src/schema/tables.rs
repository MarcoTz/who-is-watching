use super::columns::{ColumnDescription, ColumnName};
use std::fmt;

#[derive(Debug, Copy, Clone)]
pub enum Table {
    Watchers,
    Shows,
    ShowsSeasons,
    ShowsWatchers,
}

impl Table {
    pub fn all() -> [Table; 4] {
        [
            Table::Watchers,
            Table::Shows,
            Table::ShowsSeasons,
            Table::ShowsWatchers,
        ]
    }

    pub fn name(&self) -> &str {
        match self {
            Table::Watchers => "watchers",
            Table::Shows => "shows",
            Table::ShowsSeasons => "shows_seasons",
            Table::ShowsWatchers => "shows_watchers",
        }
    }

    pub fn primary_key(&self) -> Option<ColumnName> {
        match self {
            Table::Watchers => Some(ColumnName::Id),
            Table::Shows => Some(ColumnName::Id),
            Table::ShowsSeasons => None,
            Table::ShowsWatchers => None,
        }
    }

    pub fn cols(&self) -> Vec<ColumnDescription> {
        match self {
            Table::Watchers => vec![
                ColumnDescription::new(ColumnName::Id, "integer", false),
                ColumnDescription::new(ColumnName::Name, "text", false),
            ],
            Table::Shows => vec![
                ColumnDescription::new(ColumnName::Id, "integer", false),
                ColumnDescription::new(ColumnName::Name, "text", false),
            ],
            Table::ShowsSeasons => vec![
                ColumnDescription::new(ColumnName::ShowId, "integer", false),
                ColumnDescription::new(ColumnName::SeasonNum, "integer", false),
                ColumnDescription::new(ColumnName::NumEpisodes, "integer", false),
            ],
            Table::ShowsWatchers => vec![
                ColumnDescription::new(ColumnName::WatcherId, "integer", false),
                ColumnDescription::new(ColumnName::ShowId, "integer", false),
                ColumnDescription::new(ColumnName::SeasonNum, "integer", false),
                ColumnDescription::new(ColumnName::NumEpisodes, "integer", false),
            ],
        }
    }

    pub fn create_sql(&self) -> String {
        let key_str = if let Some(col) = self.primary_key() {
            format!(", PRIMARY KEY ({col})")
        } else {
            "".to_owned()
        };
        format!(
            "CREATE TABLE {} ({}{});",
            self.name(),
            self.cols()
                .iter()
                .map(|col| col.create_sql())
                .collect::<Vec<String>>()
                .join(", "),
            key_str
        )
    }
}
impl fmt::Display for Table {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "{}", self.name())
    }
}
