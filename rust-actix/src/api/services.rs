use crate::api::routes::*;
use actix_web::{web, Scope};

pub fn api_services() -> Scope {
    web::scope("/api")
        .service(get_latest)
        .service(post_register)
        .service(get_messages)
        .service(get_messages_per_user)
        .service(post_messages_per_user)
        .service(get_followers)
        .service(post_followers)
}
