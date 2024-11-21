use crate::api::controllers::*;
use crate::api::model::*;
use actix_web::{get, post, web, Responder};

#[get("/latest")]
pub async fn get_latest() -> impl Responder {
    retrieve_latest().await
}

#[post("/register")]
pub async fn post_register(
    info: web::Json<RegisterInfo>,
    query: web::Query<Latest>,
) -> impl Responder {
    register_new_user(info.into_inner(), query).await
}

#[get("/msgs")]
pub async fn get_messages(
    amount: web::Query<MessageAmount>,
    query: web::Query<Latest>,
) -> impl Responder {
    list_feed_messages(amount, query).await
}

#[get("/msgs/{username}")]
pub async fn get_messages_per_user(
    path: web::Path<String>,
    amount: web::Query<MessageAmount>,
    query: web::Query<Latest>,
) -> impl Responder {
    list_user_messages(path.into_inner(), amount, query).await
}

#[post("/msgs/{username}")]
pub async fn post_messages_per_user(
    path: web::Path<String>,
    msg: web::Json<MessageContent>,
    query: web::Query<Latest>,
) -> impl Responder {
    create_user_message(path.into_inner(), msg.into_inner(), query).await
}

#[get("/fllws/{username}")]
pub async fn get_followers(
    path: web::Path<String>,
    amount: web::Query<MessageAmount>,
    query: web::Query<Latest>,
) -> impl Responder {
    list_user_followers(path.into_inner(), amount, query).await
}

#[post("/fllws/{username}")]
pub async fn post_followers(
    path: web::Path<String>,
    follow_param: web::Json<FollowParam>,
    query: web::Query<Latest>,
) -> impl Responder {
    update_user_followers(path.into_inner(), follow_param.into_inner(), query).await
}
