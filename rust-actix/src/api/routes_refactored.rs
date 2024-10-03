use crate::controllers::*;
use actix_web::{web, get, post};

#[get("/latest")]
pub async fn get_latest() -> impl Responder {
    get_latest().await
}

#[post("/register")]
pub async fn post_register(info: web::Json<RegisterInfo>, query: web::Query<Latest>) -> impl Responder {
    post_register(info.into_inner(), query).await
}

#[get("/msgs")]
pub async fn messages_api(amount: web::Query<MessageAmount>, query: web::Query<Latest>) -> impl Responder {
    messages_api(amount, query).await
}

#[get("/msgs/{username}")]
pub async fn messages_per_user_get(path: web::Path<String>, amount: web::Query<MessageAmount>, query: web::Query<Latest>) -> impl Responder {
    messages_per_user_get(path.into_inner(), amount, query).await
}

#[post("/msgs/{username}")]
pub async fn messages_per_user_post(path: web::Path<String>, msg: web::Json<MessageContent>, query: web::Query<Latest>) -> impl Responder {
    messages_per_user_post(path.into_inner(), msg.into_inner(), query).await
}

#[get("/fllws/{username}")]
pub async fn follows_get(path: web::Path<String>, amount: web::Query<MessageAmount>, query: web::Query<Latest>) -> impl Responder {
    follows_get(path.into_inner(), amount, query).await
}

#[post("/fllws/{username}")]
pub async fn follows_post(path: web::Path<String>, follow_param: web::Json<FollowParam>, query: web::Query<Latest>) -> impl Responder {
    follows_post(path.into_inner(), follow_param.into_inner(), query).await
}
