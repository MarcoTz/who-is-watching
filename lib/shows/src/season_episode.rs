use std::fmt;

#[derive(Debug)]
pub struct SeasonEpisode {
    pub season_nr: u32,
    pub episode_nr: u32,
}

impl SeasonEpisode {
    pub fn new(season: u32, episode: u32) -> SeasonEpisode {
        SeasonEpisode {
            season_nr: season,
            episode_nr: episode,
        }
    }
}
impl fmt::Display for SeasonEpisode {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "S{:0<2}E{}", self.season_nr, self.episode_nr)
    }
}
