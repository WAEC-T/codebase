use actix_web::middleware::Logger;
use actix_web::{App, HttpServer, Scope, web};
use std::collections::HashMap;
use crate::api::routes::*;

#[actix_web::main]
pub async fn start() -> std::io::Result<()> {
    let mut labels = HashMap::new();
    labels.insert("label1".to_string(), "value1".to_string());

    HttpServer::new(move || {
        App::new()
            .wrap(Logger::default())
            .service(get_latest)
            .service(post_register)
            .service(messages_per_user_get)
            .service(messages_per_user_post)
            .service(messages_api)
            .service(follows_get)
            .service(follows_post)
    })
    .bind(("0.0.0.0", 5001))?
    .run()
    .await
}

pub fn api_services() -> Scope {
    web::scope("/api")
        .service(get_latest)               
        .service(post_register)            
        .service(messages_per_user_get) 
        .service(messages_per_user_post)
        .service(messages_api)              
        .service(follows_get)            
        .service(follows_post)           
}
