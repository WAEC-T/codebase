mod server;
use std::{num::NonZeroUsize, thread};
// use env_logger::Env;

fn main() {
    let default_actix_threads = std::thread::available_parallelism().map_or(2, NonZeroUsize::get);
    println!(
        "Starting Rust-Actix server ! \nProbable amount of threads: {}",
        default_actix_threads
    );
    // env_logger::init_from_env(Env::default().default_filter_or("info"));
    let workers = thread::spawn(server::start);
    let _ = workers.join().unwrap();
}
