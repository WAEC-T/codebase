use std::sync::Arc;

use crate::api::controllers::*;
use crate::api::model::*;
use crate::database::DatabasePool;
use actix_web::{get, post, web, Responder};

#[get("/latest")]
pub async fn get_latest(pool: web::Data<Arc<DatabasePool>>) -> impl Responder {
    retrieve_latest(&mut pool.get().await.unwrap()).await
}

#[post("/register")]
pub async fn post_register(
    pool: web::Data<Arc<DatabasePool>>,
    info: web::Json<RegisterInfo>,
    query: web::Query<Latest>,
) -> impl Responder {
    register_new_user(&mut pool.get().await.unwrap(), info.into_inner(), query).await
}

#[get("/msgs")]
pub async fn get_messages(
    pool: web::Data<Arc<DatabasePool>>,
    amount: web::Query<MessageAmount>,
    query: web::Query<Latest>,
) -> impl Responder {
    list_feed_messages(&mut pool.get().await.unwrap(), amount, query).await
}

#[get("/msgs/{username}")]
pub async fn get_messages_per_user(
    pool: web::Data<Arc<DatabasePool>>,
    path: web::Path<String>,
    amount: web::Query<MessageAmount>,
    query: web::Query<Latest>,
) -> impl Responder {
    list_user_messages(
        &mut pool.get().await.unwrap(),
        path.into_inner(),
        amount,
        query,
    )
    .await
}

#[post("/msgs/{username}")]
pub async fn post_messages_per_user(
    pool: web::Data<Arc<DatabasePool>>,
    path: web::Path<String>,
    msg: web::Json<MessageContent>,
    query: web::Query<Latest>,
) -> impl Responder {
    create_user_message(
        &mut pool.get().await.unwrap(),
        path.into_inner(),
        msg.into_inner(),
        query,
    )
    .await
}

#[get("/fllws/{username}")]
pub async fn get_followers(
    pool: web::Data<Arc<DatabasePool>>,
    path: web::Path<String>,
    amount: web::Query<MessageAmount>,
    query: web::Query<Latest>,
) -> impl Responder {
    list_user_followers(
        &mut pool.get().await.unwrap(),
        path.into_inner(),
        amount,
        query,
    )
    .await
}

#[post("/fllws/{username}")]
pub async fn post_followers(
    pool: web::Data<Arc<DatabasePool>>,
    path: web::Path<String>,
    follow_param: web::Json<FollowParam>,
    query: web::Query<Latest>,
) -> impl Responder {
    update_user_followers(
        &mut pool.get().await.unwrap(),
        path.into_inner(),
        follow_param.into_inner(),
        query,
    )
    .await
}
