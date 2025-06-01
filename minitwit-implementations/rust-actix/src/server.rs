use actix_files as fs;
use actix_session::{storage::CookieSessionStore, SessionMiddleware};
use actix_web::cookie::Key;
// use actix_web::middleware::Logger;
use actix_web::{web, App, HttpServer};
use std::num::NonZeroUsize;
use std::sync::Arc;
use waect_rust::api::middleware::AuthMiddleware;
use waect_rust::api::services::api_services;
use waect_rust::database;
use waect_rust::frontend::services::page_services;

pub async fn start(port: u16) -> std::io::Result<()> {
    let default_actix_threads = std::thread::available_parallelism().map_or(2, NonZeroUsize::get);
    println!(
        "Starting Rust-Actix server ! \nProbable amount of threads: {}",
        default_actix_threads
    );
    let pool = Arc::new(database::establish_pool().await.unwrap());
    HttpServer::new(move || {
        App::new()
            //.wrap(Logger::default())
            .app_data(web::Data::new(pool.clone()))
            .wrap(AuthMiddleware)
            .service(api_services())
            .wrap(
                SessionMiddleware::builder(CookieSessionStore::default(), Key::from(&[0; 64]))
                    .cookie_secure(false)
                    .cookie_http_only(false)
                    .build(),
            )
            .service(fs::Files::new("/static", "./src/frontend/static/").index_file("index.html"))
            .service(page_services())
    }).workers(9)
    .bind(("0.0.0.0", port))?
    .run()
    .await
}
