use crate::SeasonEpisode;

pub struct ShowProgress {
    pub show_id: u32,
    pub show_name: String,
    pub last_watched: SeasonEpisode,
}

impl ShowProgress {
    pub fn new(show_id: u32, show_name: &str, season_nr: u32, episode_nr: u32) -> ShowProgress {
        ShowProgress {
            show_id,
            show_name: show_name.to_owned(),
            last_watched: SeasonEpisode::new(season_nr, episode_nr),
        }
    }
}
