use crate::api::controllers::*;
use crate::api::model::*;
use crate::database::DatabasePool;
use actix_web::{get, post, web, Responder};

#[get("/latest")]
pub async fn get_latest(pool: web::Data<DatabasePool>) -> impl Responder {
    retrieve_latest(pool).await
}

#[post("/register")]
pub async fn post_register(
    pool: web::Data<DatabasePool>,
    info: web::Json<RegisterInfo>,
    query: web::Query<Latest>,
) -> impl Responder {
    register_new_user(pool, info.into_inner(), query).await
}

#[get("/msgs")]
pub async fn get_messages(
    pool: web::Data<DatabasePool>,
    amount: web::Query<MessageAmount>,
    query: web::Query<Latest>,
) -> impl Responder {
    list_feed_messages(pool, amount, query).await
}

#[get("/msgs/{username}")]
pub async fn get_messages_per_user(
    pool: web::Data<DatabasePool>,
    path: web::Path<String>,
    amount: web::Query<MessageAmount>,
    query: web::Query<Latest>,
) -> impl Responder {
    list_user_messages(pool, path.into_inner(), amount, query).await
}

#[post("/msgs/{username}")]
pub async fn post_messages_per_user(
    pool: web::Data<DatabasePool>,
    path: web::Path<String>,
    msg: web::Json<MessageContent>,
    query: web::Query<Latest>,
) -> impl Responder {
    create_user_message(pool, path.into_inner(), msg.into_inner(), query).await
}

#[get("/fllws/{username}")]
pub async fn get_followers(
    pool: web::Data<DatabasePool>,
    path: web::Path<String>,
    amount: web::Query<MessageAmount>,
    query: web::Query<Latest>,
) -> impl Responder {
    list_user_followers(pool, path.into_inner(), amount, query).await
}

#[post("/fllws/{username}")]
pub async fn post_followers(
    pool: web::Data<DatabasePool>,
    path: web::Path<String>,
    follow_param: web::Json<FollowParam>,
    query: web::Query<Latest>,
) -> impl Responder {
    update_user_followers(pool, path.into_inner(), follow_param.into_inner(), query).await
}
