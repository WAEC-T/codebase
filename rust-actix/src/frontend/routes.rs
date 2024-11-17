use crate::database::models::{Messages, Users};
use crate::database::repository::{
    create_msg, create_user, establish_connection, follow, get_passwd_hash, get_public_messages,
    get_timeline, get_user_by_id, get_user_by_name, get_user_timeline, is_following, unfollow,
};
use crate::frontend::flash_messages::*;
use crate::frontend::template_structs::*;
use crate::utils::datetime::format_datetime_to_message_string;
use actix_identity::Identity;
use actix_session::Session;
use actix_web::{
    get,
    http::{header, StatusCode},
    post,
    web::{self, Redirect},
    HttpMessage, HttpRequest, HttpResponse, Responder,
};
use askama_actix::Template;
use chrono::Utc;
use md5::{Digest, Md5};
use pwhash::bcrypt;

const PAGE_MESSAGES_LIMIT: i32 = 30;

fn get_user_id(username: &str) -> i32 {
    let diesel_conn = &mut establish_connection();
    let user = get_user_by_name(diesel_conn, username);
    if let Some(user) = user {
        user.user_id
    } else {
        -1
    }
}

fn get_user_template_by_name(username: &str) -> Option<UserTemplate> {
    let diesel_conn = &mut establish_connection();
    let user = get_user_by_name(diesel_conn, username);
    if let Some(user) = user {
        Some(UserTemplate {
            user_id: user.user_id,
            username: user.username,
            email: user.email,
        })
    } else {
        None
    }
}

fn get_user_template(user_id: i32) -> Option<UserTemplate> {
    let diesel_conn = &mut establish_connection();
    let user = get_user_by_id(diesel_conn, user_id);
    if let Some(user) = user {
        Some(UserTemplate {
            user_id: user.user_id,
            username: user.username,
            email: user.email,
        })
    } else {
        None
    }
}

fn get_user(user_option: Option<Identity>) -> Option<UserTemplate> {
    if let Some(user) = user_option {
        let user_id = user.id().unwrap().parse::<i32>().unwrap();
        get_user_template(user_id)
    } else {
        None
    }
}

fn gravatar_url(email: &str) -> String {
    let hash = Md5::digest(email.trim().to_lowercase().as_bytes());

    let hash_str = format!("{:x}", hash);

    format!(
        "https://www.gravatar.com/avatar/{}?d=identicon&s={}",
        hash_str, 48
    )
}

fn format_messages(messages: Vec<(Messages, Users)>) -> Vec<MessageTemplate> {
    let mut messages_for_template: Vec<MessageTemplate> = Vec::new();
    for (msg, user) in messages {
        let message = MessageTemplate {
            text: msg.text,
            username: user.username,
            gravatar_url: gravatar_url(&user.email),
            pub_date: format_datetime_to_message_string(Some(msg.pub_date)),
        };
        messages_for_template.push(message)
    }
    messages_for_template
}

#[get("/")]
async fn timeline(flash: Option<FlashMessages>, user: Option<Identity>) -> impl Responder {
    if let Some(user) = get_user(user) {
        let diesel_conn = &mut establish_connection();
        let messages =
            format_messages(get_timeline(diesel_conn, user.user_id, PAGE_MESSAGES_LIMIT));

        let rendered = TimelineTemplate {
            messages,
            request_endpoint: "timeline",
            profile_user: None,
            user: Some(user),
            followed: Some(false),
            flashes: flash.unwrap_or_default().messages,
            title: String::from("Timeline"),
        }
        .render()
        .unwrap();
        HttpResponse::Ok().body(rendered)
    } else {
        HttpResponse::TemporaryRedirect()
            .append_header((header::LOCATION, "/public"))
            .finish()
    }
}

#[get("/public")]
async fn public_timeline(
    flash_messages: Option<FlashMessages>,
    user: Option<Identity>,
) -> impl Responder {
    let user = get_user(user);
    let diesel_conn = &mut establish_connection();
    let messages = get_public_messages(diesel_conn, PAGE_MESSAGES_LIMIT);
    let messages_for_template = format_messages(messages);

    TimelineTemplate {
        messages: messages_for_template,
        request_endpoint: "/",
        profile_user: None,
        user,
        followed: Some(false),
        flashes: flash_messages.unwrap_or_default().messages,
        title: String::from(""),
    }
}

#[get("/{username}")]
async fn user_timeline(
    path: web::Path<String>,
    user: Option<Identity>,
    flash_messages: Option<FlashMessages>,
) -> impl Responder {
    let username = path.into_inner();
    let profile_user = get_user_template_by_name(&username);
    if let Some(profile_user) = profile_user {
        let mut followed = false;
        let user = get_user(user);
        let conn = &mut establish_connection();
        if let Some(user) = user.clone() {
            followed = is_following(conn, profile_user.user_id, user.user_id)
        }
        let messages = format_messages(get_user_timeline(
            conn,
            profile_user.user_id,
            PAGE_MESSAGES_LIMIT,
        ));
        let rendered = TimelineTemplate {
            messages,
            request_endpoint: "user_timeline",
            profile_user: Some(profile_user),
            user,
            followed: Some(followed),
            flashes: flash_messages.unwrap_or_default().messages,
            title: String::from("Timeline"),
        }
        .render()
        .unwrap();
        HttpResponse::Ok().body(rendered)
    } else {
        HttpResponse::NotFound().finish()
    }
}

#[get("/{username}/follow")]
async fn follow_user(
    user: Option<Identity>,
    path: web::Path<String>,
    _request: HttpRequest,
    session: Session,
) -> impl Responder {
    if let Some(_current_user) = user {
        let _target_username = path.clone();
        let _target_id = get_user_id(&_target_username);
        let conn = &mut establish_connection();
        follow(
            conn,
            _current_user.id().unwrap().parse::<i32>().unwrap(),
            _target_id,
        );

        let message = format!("You are now following \"{}\"", _target_username);
        add_flash(session, message.as_str());
    } else {
        return HttpResponse::Found()
            .append_header((header::LOCATION, "User not found"))
            .finish();
    }
    return HttpResponse::Found()
        .append_header((header::LOCATION, format!("/{}", path)))
        .finish();
}

#[get("/{username}/unfollow")]
async fn unfollow_user(
    user: Option<Identity>,
    path: web::Path<String>,
    _request: HttpRequest,
    session: Session,
) -> impl Responder {
    if let Some(_current_user) = user {
        let _target_username = path.clone();
        let _target_id = get_user_id(&_target_username);
        let conn = &mut establish_connection();
        unfollow(
            conn,
            _current_user.id().unwrap().parse::<i32>().unwrap(),
            _target_id,
        );
        let message = format!("You are no longer following \"{}\"", _target_username);
        add_flash(session, message.as_str());
    } else {
        return HttpResponse::Found()
            .append_header((header::LOCATION, "User not found"))
            .finish();
    }
    return HttpResponse::Found()
        .append_header((header::LOCATION, format!("/{}", path)))
        .finish();
}

#[post("/add_message")]
async fn add_message(
    user: Option<Identity>,
    msg: web::Form<MessageInfo>,
    session: Session,
) -> impl Responder {
    if let Some(user) = user {
        let conn = &mut establish_connection();
        let timestamp = Utc::now();
        let user_id = user.id().unwrap().parse::<i32>().unwrap();
        let _ = create_msg(conn, &user_id, &msg.text, timestamp, &0);
        add_flash(session, "Your message was recorded");
        return HttpResponse::Found()
            .append_header((header::LOCATION, "/"))
            .finish();
    }
    HttpResponse::Unauthorized()
        .status(StatusCode::UNAUTHORIZED)
        .finish()
}

#[get("/login")]
async fn login(
    flash_messages: Option<FlashMessages>,
    user: Option<Identity>,
    session: Session,
) -> impl Responder {
    if user.is_some() {
        add_flash(session, "You are already logged in");
        HttpResponse::TemporaryRedirect()
            .append_header((header::LOCATION, "/"))
            .finish()
    } else {
        let rendered = LoginTemplate {
            user: None,
            flashes: flash_messages.unwrap_or_default().messages,
            error: String::from(""),
            username: String::from(""),
        }
        .render()
        .unwrap();
        HttpResponse::Ok().body(rendered)
    }
}

#[post("/login")]
async fn post_login(
    info: web::Form<LoginInfo>,
    request: HttpRequest,
    session: Session,
) -> impl Responder {
    let conn = &mut establish_connection();
    let result = get_passwd_hash(conn, &info.username);
    if result.is_none() {
        add_flash(session, "Invalid username");
        return HttpResponse::Found()
            .append_header((header::LOCATION, "/login"))
            .finish();
    }
    //println!("{:?}", result);
    if let Some(stored_hash) = result {
        if bcrypt::verify(info.password.clone(), &stored_hash) {
            // Successful login
            let user_id = get_user_id(&info.username);
            let _ = Identity::login(&request.extensions(), user_id.to_string());
            add_flash(session, "You were logged in");

            return HttpResponse::Found()
                .append_header((header::LOCATION, "/"))
                .finish();
        }
    }

    // Password incorrect
    add_flash(session, "Invalid password");
    return HttpResponse::Found()
        .append_header((header::LOCATION, "/login"))
        .finish();
}

#[get("/register")]
async fn register(flash_messages: Option<FlashMessages>) -> impl Responder {
    RegisterTemplate {
        flashes: flash_messages.unwrap_or_default().messages,
        error: String::from(""),
        email: String::from(""),
        username: String::from(""),
        password: String::from(""),
        user: None,
    }
}

#[post("/register")]
async fn post_register(info: web::Form<RegisterInfo>, session: Session) -> impl Responder {
    if info.username.is_empty() {
        add_flash(session, "You have to enter a username");
        return Redirect::to("/register").see_other();
    } else if info.email.is_empty() || !info.email.contains('@') {
        add_flash(session, "You have to enter a valid email address");
        return Redirect::to("/register").see_other();
    } else if info.password.is_empty() {
        add_flash(session, "You have to enter a password");
        return Redirect::to("/register").see_other();
    } else if info.password != info.password2 {
        add_flash(session, "The two passwords do not match");
        return Redirect::to("/register").see_other();
    } else if get_user_id(&info.username) != -1 {
        add_flash(session, "The username is already taken");
        return Redirect::to("/register").see_other();
    }

    let hash = bcrypt::hash(info.password.clone()).unwrap();

    let conn = &mut establish_connection();
    let _ = create_user(conn, &info.username, &info.email, &hash);

    add_flash(
        session,
        "You were successfully registered and can login now",
    );
    Redirect::to("/login").see_other()
}
#[get("/logout")]
async fn logout(user: Identity, session: Session) -> impl Responder {
    add_flash(session, "You were logged out");
    user.logout();
    Redirect::to("/public").see_other()
}
