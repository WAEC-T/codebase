use std::collections::HashMap;

use actix_files as fs;
use actix_identity::config::LogoutBehaviour;
use actix_identity::Identity;
use actix_identity::IdentityMiddleware;
use actix_session::Session;
use actix_session::{storage::CookieSessionStore, SessionMiddleware};
use actix_web::http::{header, StatusCode};
use actix_web::web::{self, Redirect};

use crate::create_msg;
use crate::create_user;
use crate::establish_connection;
use crate::follow;
use crate::frontend::flash_messages::*;
use crate::frontend::template_structs::*;
use crate::get_passwd_hash;
use crate::get_public_messages;
use crate::get_timeline;
use crate::get_user_by_id;
use crate::get_user_by_name;
use crate::get_user_timeline;
use crate::is_following;
use crate::unfollow;
use crate::Messages;
use crate::Users;
use actix_web::middleware::Logger;
use actix_web::HttpMessage;
use actix_web::HttpRequest;
use actix_web::{cookie::Key, get, post, App, HttpResponse, HttpServer, Responder};
use askama_actix::Template;
use chrono::Utc;
use md5::{Digest, Md5};
use pwhash::bcrypt;

fn regular_services() -> web::Scope {
    web::scope("")
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