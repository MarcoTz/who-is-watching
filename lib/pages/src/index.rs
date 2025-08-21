use crate::errors::Error;
use askama::Template;

#[derive(Template)]
#[template(path = "index.html")]
pub struct IndexTemplate<'a> {
    page_name: &'a str,
}

impl<'a> IndexTemplate<'a> {
    pub fn render() -> Result<String, Error> {
        Ok(IndexTemplate {
            page_name: "Watching Tracker",
        }
        .render()?)
    }
}
