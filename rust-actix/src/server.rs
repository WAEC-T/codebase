use actix_web::middleware::Logger;
use actix_web::{App, HttpServer, Scope, web};
use std::collections::HashMap;
use crate::api::routes::*;

// TODO: right now it has only the api ! Add all the frontend as a service here too! [figure out how to not use the middleware and stuff fr the api]
#[actix_web::main]
pub async fn start() -> std::io::Result<()> {
    HttpServer::new(move || {
        App::new()
            .wrap(Logger::default())
            .configure(api_services)
    })
    .bind(("0.0.0.0", 5000))?
    .run()
    .await
}