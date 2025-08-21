use rand::Rng;
use shows::ShowProgress;

pub struct Watcher {
    pub id: u32,
    pub name: String,
    pub watching: Vec<ShowProgress>,
}

impl Watcher {
    pub fn new(id: u32, name: &str, watching: Vec<ShowProgress>) -> Watcher {
        Watcher {
            id,
            name: name.to_owned(),
            watching,
        }
    }

    pub fn create(name: &str, used_ids: &[u32]) -> Watcher {
        let mut rng = rand::rng();
        let mut new_id = rng.random::<u32>();
        while used_ids.contains(&new_id) {
            new_id = rng.random::<u32>();
        }
        Watcher {
            id: new_id,
            name: name.to_owned(),
            watching: vec![],
        }
    }
}
