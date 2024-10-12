use actix_web::middleware::Logger;
use actix_web::{App, HttpServer};
use waect_rust::api::services::api_services;
use waect_rust::frontend::services::{page_services, page_middleware};

#[actix_web::main]
pub async fn start() -> std::io::Result<()> {
    HttpServer::new(move || {
        App::new()
            .wrap(Logger::default())
            .service(api_services())
            //TODO implement the middlewares!!
            .service(page_services())
    })
    .bind(("0.0.0.0", 5000))?
    .run()
    .await
}