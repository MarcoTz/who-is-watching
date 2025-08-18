use crate::schema::{ColumnName, Table};
use std::fmt;

#[derive(Debug)]
pub enum StatementType {
    Select,
    Insert,
}

#[derive(Debug)]
pub enum SqlAction {
    BuildDatabase,
    ConnectToDatabase,
    CreateTable(String),
    PrepareStatement(StatementType),
    ExecuteStatement(StatementType),
    GetNextRow,
    GetValue(String),
}

#[derive(Debug)]
pub enum Error {
    Turso {
        during: SqlAction,
        location: String,
        err: turso::Error,
    },
    CouldNotCast {
        table: Table,
        col: ColumnName,
        expected: String,
    },
    NegativeId {
        table: Table,
        col: ColumnName,
    },
    ShowNotFound {
        id: u32,
    },
    NoSeasons {
        show: String,
    },
    NoRows {
        query: String,
    },
    ShowExists {
        show_name: String,
    },
    SeasonExists {
        show_name: String,
        season_nr: u32,
    },
}

impl Error {
    pub fn turso(err: turso::Error, act: SqlAction, location: &str) -> Error {
        Error::Turso {
            err,
            during: act,
            location: location.to_owned(),
        }
    }

    pub fn cast(table: &Table, col: &ColumnName, exp: &str) -> Error {
        Error::CouldNotCast {
            table: *table,
            col: *col,
            expected: exp.to_owned(),
        }
    }

    pub fn neg_id(table: &Table, col: &ColumnName) -> Error {
        Error::NegativeId {
            table: *table,
            col: *col,
        }
    }

    pub fn no_seasons(show: &str) -> Error {
        Error::NoSeasons {
            show: show.to_owned(),
        }
    }

    pub fn show_not_found(id: u32) -> Error {
        Error::ShowNotFound { id }
    }

    pub fn no_rows(query: &str) -> Error {
        Error::NoRows {
            query: query.to_owned(),
        }
    }

    pub fn show_exists(show: &str) -> Error {
        Error::ShowExists {
            show_name: show.to_owned(),
        }
    }

    pub fn season_exists(show: &str, season: u32) -> Error {
        Error::SeasonExists {
            show_name: show.to_owned(),
            season_nr: season,
        }
    }
}

impl fmt::Display for Error {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match self {
            Error::Turso {
                err,
                during,
                location,
            } => {
                write!(f, "Error in turso during {during} at {location}:\n\t{err}")
            }
            Error::CouldNotCast {
                table,
                col,
                expected,
            } => {
                write!(f, "Could not cast column {table}.{col} to type {expected}")
            }
            Error::NegativeId { table, col } => write!(f, "Found negative {col} in table {table}"),
            Error::NoSeasons { show } => write!(f, "Could not find any seasons for show {show}"),
            Error::ShowNotFound { id } => write!(f, "No show with id {id}"),
            Error::NoRows { query } => write!(f, "Query returned no rows {query}"),
            Error::ShowExists { show_name } => {
                write!(f, "Cannot create show {show_name}, show already exists")
            }
            Error::SeasonExists {
                show_name,
                season_nr,
            } => write!(
                f,
                "Cannot create season {season_nr} for show {show_name}, season already exists"
            ),
        }
    }
}

impl fmt::Display for SqlAction {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match self {
            SqlAction::BuildDatabase => f.write_str("Building Database"),
            SqlAction::ConnectToDatabase => f.write_str("Connecting to Database"),
            SqlAction::CreateTable(tbl) => write!(f, "Creating Table {tbl}"),
            SqlAction::PrepareStatement(stmt) => write!(f, "Preparing {stmt} statement"),
            SqlAction::ExecuteStatement(stmt) => write!(f, "Executing {stmt} statement"),
            SqlAction::GetNextRow => write!(f, "Getting Next Row"),
            SqlAction::GetValue(val) => write!(f, "Getting value {val}"),
        }
    }
}

impl fmt::Display for StatementType {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match self {
            StatementType::Insert => f.write_str("Insert"),
            StatementType::Select => f.write_str("Select"),
        }
    }
}

impl std::error::Error for Error {}
