# frozen_string_literal: true

require "sinatra/activerecord"
require "sinatra"
require "sinatra/flash"
require "newrelic_rpm"
Dir["./models/*.rb"].each { |file| require file }

PR_PAGE = 30

DATABASE_URL = ENV.fetch("DATABASE_URL", nil)

configure :production, :staging do
  db = URI.parse(DATABASE_URL)
  set :database, {
    adapter: db.scheme,
    host: db.host,
    port: db.port,
    database: db.path[1..],
    user: db.user,
    password: db.password,
    encoding: "utf8",
  }
  set :public_folder, "#{__dir__}/static"
  enable :sessions
end

configure :development do
  set :database, { adapter: "postgresql", database: "minitwit_development" }
  set :public_folder, "#{__dir__}/static"
  ActiveRecord.verbose_query_logs = true
  enable :sessions
end

configure :test do
  if DATABASE_URL
    db = URI.parse(DATABASE_URL)
    set :database, {
      adapter: db.scheme,
      host: db.host,
      port: db.port,
      database: db.path[1..],
      user: db.user,
      password: db.password,
      encoding: "utf8",
    }
  else
    set :database, { adapter: "postgresql", database: "minitwit_test" }
  end
  enable :sessions
  enable :logging
  ActiveRecord::Base.logger = Logger.new($stdout)
  set :public_folder, "#{__dir__}/static"
end

helpers do
  def logged_in?
    !!session[:user_id]
  end

  def current_user
    return unless logged_in?

    @current_user ||= User.find(session[:user_id])
  end
end

get "/" do
  redirect("/public") unless logged_in?

  @messages = Message
    .unflagged
    .authored_by(current_user.following + [current_user])
    .includes(:author)
    .order(created_at: :desc)
    .first(PR_PAGE)

  erb :timeline, layout: :layout
end

get "/public" do
  @messages = Message
    .unflagged
    .includes(:author)
    .order(created_at: :desc)
    .first(PR_PAGE)

  erb :timeline, layout: :layout
end

get "/register" do
  erb :register, layout: :layout
end

post "/register" do
  user = User.new(
    username: params[:username],
    email: params[:email],
    password: params[:password],
    password_confirmation: params[:password2],
  )
  if user.save
    flash[:success] = "You were successfully registered and can login now"
    redirect("/login")
  else
    errors = user.errors.map(&:full_message).join(", ")
    flash[:error] = errors
  end
  redirect("/register")
end

get "/login" do
  erb :login, layout: :layout
end

post "/login" do
  user = User.find_by_username(params[:username])
  if user.nil?
    error = "Invalid username"
  elsif !user.authenticate(params[:password])
    error = "Invalid password"
  else
    session[:user_id] = user.id
    flash[:success] = "You were logged in"
    redirect("/")
  end
  flash[:error] = error
  redirect("/login")
end

get "/logout" do
  session[:user_id] = nil

  "You were logged out"
end

post "/add_message" do
  if !session[:user_id]
    return status 401
  elsif params[:text]
    if Message.create(
      author_id: session[:user_id],
      text: params[:text],
      flagged: false,
    )
      flash[:success] = "Your message was recorded"
    end
  else
    flash[:error] = "Something went wrong :("
  end

  redirect("/")
end

get "/:username" do
  @profile_user = User.find_by_username(params[:username])
  @messages = Message
    .unflagged
    .authored_by(@profile_user)
    .includes(:author)
    .order(created_at: :desc)
    .first(PR_PAGE)

  erb :timeline, layout: :layout
end

get "/:username/follow" do
  return status 401 unless logged_in?

  whom = User.find_by_username(params[:username])

  return status 404 if whom.nil?

  current_user.following << whom
  flash[:success] = "You are now following &#34;#{params[:username]}&#34;"
  redirect("/#{params[:username]}")
end

get "/:username/unfollow" do
  return status 401 unless logged_in?

  whom = User.find_by_username(params[:username])

  return status 404 if whom.nil?

  current_user.following.delete(whom)
  flash[:success] = "You are no longer following &#34;#{params[:username]}&#34;"
  redirect("/#{params[:username]}")
end
