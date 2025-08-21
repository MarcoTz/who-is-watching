use crate::errors::Error;
use askama::Template;

#[derive(Template)]
#[template(path = "shows.html")]
pub struct ShowsTemplate {}

impl ShowsTemplate {
    pub fn render() -> Result<String, Error> {
        Ok(ShowsTemplate {}.render()?)
    }
}
