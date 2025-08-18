use shows::Show;
use sqlite::DBDriver;
use std::io::Write;
use turso::Builder;

#[tokio::main]
async fn main() -> Result<(), sqlite::errors::Error> {
    let old = Builder::new_local("old/watchers.db").build().await.unwrap();
    let old_conn = old.connect().unwrap();
    let mut stmt = old_conn.prepare("SELECT name FROM shows;").await.unwrap();
    let mut rows = stmt.query(()).await.unwrap();
    let mut old_shows = vec![];
    while let Some(row) = rows.next().await.unwrap() {
        let name_val = row.get_value(0).unwrap();
        let name = name_val.as_text().unwrap();
        old_shows.push(name.clone());
    }

    let drv = DBDriver::connect().await?;
    let shows = drv.load_shows().await?;
    let show_ids = shows.iter().map(|show| show.id).collect::<Vec<u32>>();
    let stdin = std::io::stdin();
    let mut stdout = std::io::stdout();
    for show_name in old_shows {
        if drv.show_exists(&show_name).await.unwrap() {
            continue;
        }
        let mut num_seasons = String::new();
        let mut episodes = String::new();
        println!("Creating {show_name}");
        print!("Enter number of seasons: ");
        stdout.flush().unwrap();
        stdin.read_line(&mut num_seasons).unwrap();
        let seasons = num_seasons.trim().parse::<u32>().unwrap();

        let mut seasons_episodes = vec![];
        for season in 1..=seasons {
            episodes.clear();
            print!("Enter number of episodes (season {season}): ");
            stdout.flush().unwrap();
            stdin.read_line(&mut episodes).unwrap();
            let episode = episodes.trim().parse::<u32>().unwrap();
            seasons_episodes.push((season, episode));
        }

        seasons_episodes.reverse();
        let (max_season, num_episodes) = seasons_episodes.remove(0);
        let new_show = Show::create(&show_name, max_season, num_episodes, &show_ids);
        drv.create_show(&new_show).await?;

        for (season, episodes) in seasons_episodes {
            drv.add_season(new_show.id, season, episodes).await?;
        }
    }

    Ok(())
}
