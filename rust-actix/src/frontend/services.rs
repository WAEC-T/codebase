use actix_files as fs;
use actix_identity::{config::LogoutBehaviour, IdentityMiddleware};
use actix_session::{storage::CookieSessionStore, SessionMiddleware};
use crate::frontend::routes::*;
use actix_web::{web::{self}, cookie::Key, middleware::{Logger, self}, Scope, dev::{ServiceFactory, ServiceRequest, ServiceResponse}, body::BoxBody};
pub fn page_services() -> Scope {
      return web::scope("")
        .service(fs::Files::new("/static", "./src/frontend/static/").index_file("index.html"))
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
        .service(add_message);
}

//TODO Fix this function!
pub fn page_middleware(scope: Scope) {
    scope
        .wrap(
            IdentityMiddleware::builder()
                .logout_behaviour(LogoutBehaviour::DeleteIdentityKeys)
                .build(),
        )
        .wrap(
            SessionMiddleware::builder(CookieSessionStore::default(), Key::from(&[0; 64]))
                .cookie_secure(false)
                .cookie_http_only(false)
                .build(),
        )
        .wrap(Logger::default());
}