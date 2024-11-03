use actix_identity::config::LogoutBehaviour;
use actix_identity::IdentityMiddleware;
use actix_web::middleware::Logger;
use actix_web::{App, HttpServer};
use waect_rust::api::services::api_services;
use waect_rust::frontend::services::page_services;
use actix_files as fs;
use actix_session::{storage::CookieSessionStore, SessionMiddleware};
use actix_web::cookie::Key;

#[actix_web::main]
pub async fn start() -> std::io::Result<()> {
    HttpServer::new(move || {
        App::new()
            .wrap(Logger::default())
            .service(api_services())
            .wrap(
                IdentityMiddleware::builder()
                    .logout_behaviour(LogoutBehaviour::DeleteIdentityKeys)
                    .build(),
            )
            .service(fs::Files::new("/static", "./src/frontend/static/").index_file("index.html"))
            .wrap(
                SessionMiddleware::builder(CookieSessionStore::default(), Key::from(&[0; 64]))
                    .cookie_secure(false)
                    .cookie_http_only(false)
                    .build(),
            )
            .wrap(Logger::default())
            .service(page_services())
    })
    .bind(("0.0.0.0", 5000))?
    .run()
    .await
}
