struct Show {
    id: u64,
    name: String,
    latest_episode: SeasonEpisode,
}

struct SeasonEpisode {
    season_nr: u64,
    episode_nr: u64,
}

pub struct ShowProgress {
    show_id: u64,
    last_watched: SeasonEpisode,
}
