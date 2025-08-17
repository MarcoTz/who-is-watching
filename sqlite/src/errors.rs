use std::fmt;

#[derive(Debug)]
pub enum Error {
    Turso(turso::Error),
}

impl fmt::Display for Error {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match self {
            Error::Turso(err) => write!(f, "Error in turso:\n\t{err}"),
        }
    }
}

impl From<turso::Error> for Error {
    fn from(err: turso::Error) -> Error {
        Error::Turso(err)
    }
}
