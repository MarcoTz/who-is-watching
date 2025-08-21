use std::fmt;

#[derive(Debug)]
pub enum Error {
    Askama(askama::Error),
}

impl fmt::Display for Error {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match self {
            Error::Askama(err) => write!(f, "Error in Askama:\n{err}"),
        }
    }
}

impl std::error::Error for Error {}

impl From<askama::Error> for Error {
    fn from(err: askama::Error) -> Error {
        Error::Askama(err)
    }
}
