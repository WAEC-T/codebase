use actix_files as fs;
use actix_identity::config::LogoutBehaviour;
use actix_identity::IdentityMiddleware;
use actix_session::{storage::CookieSessionStore, SessionMiddleware};
use actix_web::cookie::Key;
// use actix_web::middleware::Logger;
use actix_web::{web, App, HttpServer};
use waect_rust::api::middleware::AuthMiddleware;
use waect_rust::api::services::api_services;
use waect_rust::database;
use waect_rust::frontend::services::page_services;

#[actix_web::main]
pub async fn start() -> std::io::Result<()> {
    let pool = database::establish_pool().await.unwrap();
    HttpServer::new(move || {
        App::new()
            // .wrap(Logger::default())
            .app_data(web::Data::new(pool.clone()))
            .wrap(AuthMiddleware)
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
            .service(page_services())
    })
    .bind(("0.0.0.0", 5000))?
    .run()
    .await
}
