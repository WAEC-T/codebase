use askama_actix::Template;
use serde::Deserialize;

#[derive(Clone)]
pub struct UserTemplate {
    pub user_id: i32,
    pub username: String,
    pub email: String,
}

#[derive(Debug)]
pub struct MessageTemplate {
    pub text: String,
    pub username: String,
    pub pub_date: String,
    pub gravatar_url: String,
}

#[derive(Template)]
#[template(path = "../templates/timeline.html")]
pub struct TimelineTemplate<'a> {
    pub messages: Vec<MessageTemplate>,
    pub user: Option<UserTemplate>,
    pub request_endpoint: &'a str,
    pub profile_user: Option<&'a UserTemplate>,
    pub followed: Option<bool>,
    pub flashes: Vec<String>,
    pub title: &'a str,
    pub error: &'a str,
}

#[derive(Template)]
#[template(path = "../templates/login.html")]
pub struct LoginTemplate<'a> {
    pub user: Option<UserTemplate>,
    pub error: &'a str,
    pub flashes: Vec<String>,
    pub username: &'a str,
}

#[derive(Template)]
#[template(path = "../templates/register.html")]
pub struct RegisterTemplate<'a> {
    pub user: Option<UserTemplate>,
    pub email: &'a str,
    pub username: &'a str,
    pub password: &'a str,
    pub flashes: Vec<String>,
    pub error: &'a str,
}

#[derive(Deserialize)]
pub struct MessageInfo {
    pub text: String,
}

#[derive(Deserialize)]
pub struct LoginInfo {
    pub username: String,
    pub password: String,
}

#[derive(Deserialize)]
pub struct RegisterInfo {
    pub username: String,
    pub email: String,
    pub password: String,
    pub password2: String,
}
