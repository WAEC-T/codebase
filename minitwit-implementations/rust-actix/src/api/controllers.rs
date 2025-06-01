use crate::{
    api::model::*,
    database::{
        repository::{
            create_msg, create_user, follow, get_followers, get_latest, get_public_messages,
            get_timeline, get_user_by_name, set_latest, unfollow,
        },
        PostgresConnection,
    },
    utils::datetime::convert_naive_to_utc,
};
use actix_web::{web, HttpResponse};
use chrono::Utc;

async fn get_user_id(conn: &mut PostgresConnection, username: &str) -> Option<i32> {
    get_user_by_name(conn, username)
        .await
        .map(|user| user.user_id)
}

async fn update_latest(conn: &mut PostgresConnection, query: web::Query<Latest>) {
    set_latest(conn, &query.latest).await;
}

pub async fn retrieve_latest(conn: &mut PostgresConnection) -> HttpResponse {
    let latest = get_latest(conn).await;
    HttpResponse::Ok().json(Latest { latest })
}

pub async fn register_new_user(
    conn: &mut PostgresConnection,
    info: RegisterInfo,
    query: web::Query<Latest>,
) -> HttpResponse {
    update_latest(conn, query).await;

    let user_exists = get_user_id(conn, &info.username).await;

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
        let _ = create_user(conn, &info.username, &info.email, &info.pwd).await;
        HttpResponse::NoContent().json(String::from(""))
    }
}

pub async fn list_feed_messages(
    conn: &mut PostgresConnection,
    amount: web::Query<MessageAmount>,
    query: web::Query<Latest>,
) -> HttpResponse {
    update_latest(conn, query).await;
    let messages: Vec<Message> = get_public_messages(conn, amount.no)
        .await
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
    conn: &mut PostgresConnection,
    username: String,
    amount: web::Query<MessageAmount>,
    query: web::Query<Latest>,
) -> HttpResponse {
    update_latest(conn, query).await;
    if let Some(user_id) = get_user_id(conn, &username).await {
        let messages: Vec<Message> = get_timeline(conn, user_id, amount.no)
            .await
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
    conn: &mut PostgresConnection,
    username: String,
    msg: MessageContent,
    query: web::Query<Latest>,
) -> HttpResponse {
    update_latest(conn, query).await;

    if let Some(user_id) = get_user_id(conn, &username).await {
        let _ = create_msg(conn, &user_id, &msg.content, Utc::now(), &0).await;
        HttpResponse::NoContent().json("")
    } else {
        HttpResponse::NotFound().json("")
    }
}

pub async fn list_user_followers(
    conn: &mut PostgresConnection,
    username: String,
    amount: web::Query<MessageAmount>,
    query: web::Query<Latest>,
) -> HttpResponse {
    update_latest(conn, query).await;

    if let Some(user_id) = get_user_id(conn, &username).await {
        let followers = get_followers(conn, user_id, amount.no).await;
        let followers = followers.into_iter().map(|user| user.username).collect();
        HttpResponse::Ok().json(Follows { follows: followers })
    } else {
        HttpResponse::NotFound().json("")
    }
}

pub async fn update_user_followers(
    conn: &mut PostgresConnection,
    username: String,
    follow_param: FollowParam,
    query: web::Query<Latest>,
) -> HttpResponse {
    update_latest(conn, query).await;

    if let Some(user_id) = get_user_id(conn, &username).await {
        if let Some(follow_username) = follow_param.follow {
            if let Some(follow_user_id) = get_user_id(conn, &follow_username).await {
                follow(conn, user_id, follow_user_id).await;
                return HttpResponse::NoContent().json("");
            }
        } else if let Some(unfollow_username) = follow_param.unfollow {
            if let Some(unfollow_user_id) = get_user_id(conn, &unfollow_username).await {
                unfollow(conn, user_id, unfollow_user_id).await;
                return HttpResponse::NoContent().json("");
            }
        }

        HttpResponse::BadRequest().json("")
    } else {
        HttpResponse::NotFound().json("")
    }
}
