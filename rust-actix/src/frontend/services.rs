use crate::frontend::routes::*;
use actix_web::{
    web::{self},
    Scope,
};

pub fn page_services() -> Scope {
    web::scope("")
        .service(register)
        .service(post_register)
        .service(timeline)
        .service(public_timeline)
        .service(login)
        .service(post_login)
        .service(logout)
        .service(user_timeline)
        .service(follow_user)
        .service(unfollow_user)
        .service(add_message)
}
