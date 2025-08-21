use std::fmt;

#[derive(Debug)]
pub enum Error {
    Askama(askama::Error),
    Sql(sqlite::Error),
}

impl fmt::Display for Error {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match self {
            Error::Askama(err) => write!(f, "Error in Askama:\n{err}"),
            Error::Sql(err) => write!(f, "Error in Sqlite:\n{err}"),
        }
    }
}

impl std::error::Error for Error {}

impl From<askama::Error> for Error {
    fn from(err: askama::Error) -> Error {
        Error::Askama(err)
    }
}

impl From<sqlite::Error> for Error {
    fn from(err: sqlite::Error) -> Error {
        Error::Sql(err)
    }
}
