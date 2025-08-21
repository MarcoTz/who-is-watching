use crate::schema::Table;
use std::fmt;

#[derive(Debug, Clone, Copy)]
pub enum ColumnName {
    Id,
    Name,
    ShowId,
    SeasonNum,
    NumEpisodes,
    WatcherId,
}

impl ColumnName {
    pub fn qualified(&self, table: &Table) -> String {
        format!("{}.{}", table, self)
    }
}

pub struct ColumnDescription {
    name: ColumnName,
    ty: String,
    nullable: bool,
}

impl ColumnDescription {
    pub fn new(name: ColumnName, ty: &str, nullable: bool) -> ColumnDescription {
        ColumnDescription {
            name,
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

impl fmt::Display for ColumnName {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match self {
            ColumnName::Id => f.write_str("id"),
            ColumnName::Name => f.write_str("name"),
            ColumnName::SeasonNum => f.write_str("season_num"),
            ColumnName::NumEpisodes => f.write_str("num_episodes"),
            ColumnName::ShowId => f.write_str("show_id"),
            ColumnName::WatcherId => f.write_str("watcher_id"),
        }
    }
}
