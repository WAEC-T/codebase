use crate::{
    api::model::*,
    database::repository::{
        create_msg, create_user, establish_connection, follow, get_followers, get_latest,
        get_public_messages, get_timeline, get_user_by_name, set_latest, unfollow,
    },
    utils::datetime::convert_naive_to_utc,
};
use actix_web::{web, HttpResponse};
use chrono::Utc;
use diesel::PgConnection;
use pwhash::bcrypt;

fn get_user_id(username: &str) -> Option<i32> {
    let conn = &mut establish_connection();
    get_user_by_name(conn, username).map(|user| user.user_id)
}

fn update_latest(conn: &mut PgConnection, query: web::Query<Latest>) {
    set_latest(conn, query.latest);
}

pub async fn retrieve_latest() -> HttpResponse {
    let conn: &mut PgConnection = &mut establish_connection();
    let latest = get_latest(conn);
    HttpResponse::Ok().json(Latest { latest })
}

pub async fn register_new_user(info: RegisterInfo, query: web::Query<Latest>) -> HttpResponse {
    let conn = &mut establish_connection();
    update_latest(conn, query);

    let user_exists = get_user_id(&info.username);

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
        let _ = create_user(conn, &info.username, &info.email, &hash);
        HttpResponse::NoContent().json(String::from(""))
    }
}

pub async fn list_feed_messages(
    amount: web::Query<MessageAmount>,
    query: web::Query<Latest>,
) -> HttpResponse {
    let conn = &mut establish_connection();
    update_latest(conn, query);
    let messages: Vec<Message> = get_public_messages(conn, amount.no)
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
    username: String,
    amount: web::Query<MessageAmount>,
    query: web::Query<Latest>,
) -> HttpResponse {
    let conn = &mut establish_connection();
    update_latest(conn, query);

    if let Some(user_id) = get_user_id(&username) {
        let messages: Vec<Message> = get_timeline(conn, user_id, amount.no)
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
    username: String,
    msg: MessageContent,
    query: web::Query<Latest>,
) -> HttpResponse {
    let conn = &mut establish_connection();
    update_latest(conn, query);

    if let Some(user_id) = get_user_id(&username) {
        let _ = create_msg(conn, &user_id, &msg.content, Utc::now(), &0);
        HttpResponse::NoContent().json("")
    } else {
        HttpResponse::NotFound().json("")
    }
}

pub async fn list_user_followers(
    username: String,
    amount: web::Query<MessageAmount>,
    query: web::Query<Latest>,
) -> HttpResponse {
    let conn = &mut establish_connection();
    update_latest(conn, query);

    if let Some(user_id) = get_user_id(&username) {
        let followers = get_followers(conn, user_id, amount.no);
        let followers = followers.into_iter().map(|user| user.username).collect();
        HttpResponse::Ok().json(Follows { follows: followers })
    } else {
        HttpResponse::NotFound().json("")
    }
}

pub async fn update_user_followers(
    username: String,
    follow_param: FollowParam,
    query: web::Query<Latest>,
) -> HttpResponse {
    let conn = &mut establish_connection();
    update_latest(conn, query);

    if let Some(user_id) = get_user_id(&username) {
        if let Some(follow_username) = follow_param.follow {
            if let Some(follow_user_id) = get_user_id(&follow_username) {
                follow(conn, user_id, follow_user_id);
                return HttpResponse::NoContent().json("");
            }
        } else if let Some(unfollow_username) = follow_param.unfollow {
            if let Some(unfollow_user_id) = get_user_id(&unfollow_username) {
                unfollow(conn, user_id, unfollow_user_id);
                return HttpResponse::NoContent().json("");
            }
        }

        HttpResponse::BadRequest().json("")
    } else {
        HttpResponse::NotFound().json("")
    }
}
