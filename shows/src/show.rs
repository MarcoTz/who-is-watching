use crate::SeasonEpisode;
use rand::Rng;
use std::fmt;

pub struct Show {
    pub id: u32,
    pub name: String,
    pub latest_episode: SeasonEpisode,
}

impl Show {
    pub fn new(id: u32, name: &str, season_num: u32, num_episodes: u32) -> Show {
        Show {
            id,
            name: name.to_owned(),
            latest_episode: SeasonEpisode::new(season_num, num_episodes),
        }
    }

    pub fn create(name: &str, season_num: u32, num_episodes: u32, used_ids: &[u32]) -> Show {
        let mut rng = rand::rng();
        let mut new_id = rng.random::<u32>();
        while used_ids.contains(&new_id) {
            new_id = rng.random::<u32>();
        }
        Show {
            id: new_id,
            name: name.to_owned(),
            latest_episode: SeasonEpisode::new(season_num, num_episodes),
        }
    }
}
impl fmt::Display for Show {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "{} ({}) {}", self.name, self.id, self.latest_episode)
    }
}
