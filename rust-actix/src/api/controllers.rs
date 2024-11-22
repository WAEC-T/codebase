use crate::{
    api::model::*,
    database::{repository::{
        create_msg, create_user, follow, get_followers, get_latest,
        get_public_messages, get_timeline, get_user_by_name, set_latest, unfollow,
    }, PostgresConnection},
    utils::datetime::convert_naive_to_utc,
};
use actix_web::{web, HttpResponse};
use chrono::Utc;
use pwhash::bcrypt;
use crate::database::pool::DatabasePool;

async fn get_user_id(pool: web::Data<DatabasePool>, username: &str) -> Option<i32> {
    let mut conn = pool.get().await.unwrap();
    get_user_by_name(&mut conn, username).await.map(|user| user.user_id)
}

async fn update_latest(conn: &mut PostgresConnection, query: web::Query<Latest>) {
    set_latest(conn, query.latest).await;
}

pub async fn retrieve_latest(pool: web::Data<DatabasePool>) -> HttpResponse {
    let mut conn = pool.get().await.unwrap();
    let latest = get_latest(&mut conn).await;
    HttpResponse::Ok().json(Latest { latest })
}

pub async fn register_new_user(pool: web::Data<DatabasePool>, info: RegisterInfo, query: web::Query<Latest>) -> HttpResponse {
    let mut conn = pool.get().await.unwrap();
    update_latest(&mut conn, query).await;

    let user_exists = get_user_id(pool.clone(), &info.username).await;

    let error = if info.username.is_empty() {
        Some(String::from("You have to enter a username"))
    } else if info.email.is_empty() {
        Some(String::from("You have to enter a valid email address"))
    } else if info.pwd.is_empty() {
        Some(String::from("You have to enter a password"))
    } else if user_exists.is_some() {
        Some(String::from("The username is already taken"))
    } else {
        None
    };

    if let Some(err_msg) = error {
        let reg_err = RegisterError {
            status: 400,
            error_msg: err_msg.to_string(),
        };
        HttpResponse::BadRequest().json(reg_err)
    } else {
        let hash = bcrypt::hash(info.pwd.clone()).unwrap();
        let _ = create_user(&mut conn, &info.username, &info.email, &hash);
        HttpResponse::NoContent().json(String::from(""))
    }
}

pub async fn list_feed_messages(
    pool: web::Data<DatabasePool>,
    amount: web::Query<MessageAmount>,
    query: web::Query<Latest>,
) -> HttpResponse {
    let mut conn = pool.get().await.unwrap();
    update_latest(&mut conn, query).await;
    let messages: Vec<Message> = get_public_messages(&mut conn, amount.no).await
        .into_iter()
        .map(|(msg, user)| Message {
            content: msg.text,
            user: user.username,
            pub_date: convert_naive_to_utc(msg.pub_date),
        })
        .collect();

    HttpResponse::Ok().json(messages)
}

pub async fn list_user_messages(
    pool: web::Data<DatabasePool>,
    username: String,
    amount: web::Query<MessageAmount>,
    query: web::Query<Latest>,
) -> HttpResponse {
    let mut conn = pool.get().await.unwrap();
    update_latest(&mut conn, query).await;

    if let Some(user_id) = get_user_id(pool.clone(), &username).await {
        let messages: Vec<Message> = get_timeline(&mut conn, user_id, amount.no).await
            .into_iter()
            .map(|(msg, user)| Message {
                content: msg.text,
                user: user.username,
                pub_date: convert_naive_to_utc(msg.pub_date),
            })
            .collect();

        HttpResponse::Ok().json(messages)
    } else {
        HttpResponse::NotFound().json("")
    }
}

pub async fn create_user_message(
    pool: web::Data<DatabasePool>,
    username: String,
    msg: MessageContent,
    query: web::Query<Latest>,
) -> HttpResponse {
    let mut conn = pool.get().await.unwrap();
    update_latest(&mut conn, query).await;

    if let Some(user_id) = get_user_id(pool.clone(), &username).await {
        let _ = create_msg(&mut conn, &user_id, &msg.content, Utc::now(), &0);
        HttpResponse::NoContent().json("")
    } else {
        HttpResponse::NotFound().json("")
    }
}

pub async fn list_user_followers(
    pool: web::Data<DatabasePool>,
    username: String,
    amount: web::Query<MessageAmount>,
    query: web::Query<Latest>,
) -> HttpResponse {
    let mut conn = pool.get().await.unwrap();
    update_latest(&mut conn, query).await;

    if let Some(user_id) = get_user_id(pool.clone(), &username).await {
        let followers = get_followers(&mut conn, user_id, amount.no).await;
        let followers = followers.into_iter().map(|user| user.username).collect();
        HttpResponse::Ok().json(Follows { follows: followers })
    } else {
        HttpResponse::NotFound().json("")
    }
}

pub async fn update_user_followers(
    pool: web::Data<DatabasePool>,
    username: String,
    follow_param: FollowParam,
    query: web::Query<Latest>,
) -> HttpResponse {
    let mut conn = pool.get().await.unwrap();
    update_latest(&mut conn, query).await;

    if let Some(user_id) = get_user_id(pool.clone(), &username).await {
        if let Some(follow_username) = follow_param.follow {
            if let Some(follow_user_id) = get_user_id(pool.clone(), &follow_username).await {
                follow(&mut conn, user_id, follow_user_id).await;
                return HttpResponse::NoContent().json("");
            }
        } else if let Some(unfollow_username) = follow_param.unfollow {
            if let Some(unfollow_user_id) = get_user_id(pool.clone(), &unfollow_username).await {
                unfollow(&mut conn, user_id, unfollow_user_id).await;
                return HttpResponse::NoContent().json("");
            }
        }

        HttpResponse::BadRequest().json("")
    } else {
        HttpResponse::NotFound().json("")
    }
}
