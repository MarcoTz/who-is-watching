use crate::SeasonEpisode;

pub struct ShowProgress {
    pub show_id: u32,
    pub last_watched: SeasonEpisode,
}

impl ShowProgress {
    pub fn new(show_id: u32, season_nr: u32, episode_nr: u32) -> ShowProgress {
        ShowProgress {
            show_id,
            last_watched: SeasonEpisode::new(season_nr, episode_nr),
        }
    }
}
