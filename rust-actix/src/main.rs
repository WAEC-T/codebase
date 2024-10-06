mod server;
use std::thread;
use env_logger::Env;

fn main() {
    env_logger::init_from_env(Env::default().default_filter_or("info"));

    let workers = thread::spawn(server::start);
    let _ = workers.join().unwrap();
}
