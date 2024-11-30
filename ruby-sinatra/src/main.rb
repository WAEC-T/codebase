# frozen_string_literal: true

require "sinatra/activerecord"
require "sinatra"
require "sinatra/namespace"
require "sinatra/flash"
require "json"
require "newrelic_rpm"
require 'dotenv'

Dir["./models/*.rb"].each { |file| require file }

USER_NOT_FOUND = "User not found"

PR_PAGE = 30

DATABASE_URL = ENV.fetch("DATABASE_URL", nil)

puts "Loaded DATABASE_URL: #{DATABASE_URL}"

NO_IS_REQUIRED = "no is required"


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


configure :development, :test do
  if DATABASE_URL
    set :database, DATABASE_URL
  else
    # Fall back to specific configurations for development and test environments
    case settings.environment
    when :development
      set :database, { adapter: "postgresql", database: "waect" }
      ActiveRecord.verbose_query_logs = true
    when :test
      set :database, { adapter: "postgresql", database: "minitwit_test" }
      enable :logging
      ActiveRecord::Base.logger = Logger.new($stdout)
    end
  end
  set :public_folder, "#{__dir__}/static"
  enable :sessions
end

helpers do
  def logged_in?
    !!session[:user_id]
  end

  def current_user
    return unless logged_in?

    @current_user ||= User.find(session[:user_id])
  end

  def request_is_not_from_simulator
    request.env["HTTP_AUTHORIZATION"] != "Basic c2ltdWxhdG9yOnN1cGVyX3NhZmUh"
  end
  
end

def update_latest(request)
  parsed_command_id = request.params["latest"]
  Latest.set(parsed_command_id.to_i) if parsed_command_id.present?
end

before "/api/*" do
  update_latest(request)
  if request_is_not_from_simulator
    halt [403, { status: 403, error_msg: "You are not authorized to use this resource!" }.to_json]
  end
end

get "/" do
  redirect("/public") unless logged_in?

  @title = "My Timeline"

  @messages = Message
    .unflagged
    .authored_by(current_user.following + [current_user])
    .includes(:author)
    .order(pub_date: :desc)
    .first(PR_PAGE)

  erb :timeline, layout: :layout
end

get "/public" do
  @title = "Public Timeline"
  @messages = Message
    .unflagged
    .includes(:author)
    .order(pub_date: :desc)
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
  flash[:success] = "You were logged out"
  redirect("/")
end

post "/add_message" do
  if !session[:user_id]
    return status 401
  elsif params[:text].nil? || params[:text].strip.empty?
    flash[:error] = "Message cannot be empty!"
  else
    if Message.create(
      author_id: session[:user_id],
      text: params[:text],
      flagged: false,
      pub_date: Time.now,
    )
      flash[:success] = "Your message was recorded"
    else
      flash[:error] = "Something went wrong :("
    end
  end

  redirect("/")
end

get "/:username" do
  @profile_user = User.find_by_username(params[:username])
  @title = "#{@profile_user.username}'s Timeline"
  @messages = Message
    .unflagged
    .authored_by(@profile_user)
    .includes(:author)
    .order(pub_date: :desc)
    .first(PR_PAGE)

  erb :timeline, layout: :layout
end

get "/:username/follow" do
  return status 401 unless logged_in?

  whom = User.find_by_username(params[:username])

  return status 404 if whom.nil?

  current_user.following << whom
  flash[:success] = "You are now following #{params[:username]}"
  redirect("/#{params[:username]}")
end

get "/:username/unfollow" do
  return status 401 unless logged_in?

  whom = User.find_by_username(params[:username])

  return status 404 if whom.nil?

  current_user.following.delete(whom)
  flash[:success] = "You are no longer following #{params[:username]}"
  redirect("/#{params[:username]}")
end

namespace "/api" do

  get "/latest" do
    body({ latest: Latest.get }.to_json)
  end

  post "/register" do
    request_data = JSON.parse(request.body.read, symbolize_names: true)

    user = User.new(
      username: request_data[:username],
      email: request_data[:email],
      password: request_data[:pwd],
      password_confirmation: request_data[:pwd],
    )

    if user.save
      status 204
    else
      status 400
      body({ status: 400, error_msg: user.errors.map(&:full_message) }.to_json)
    end
  end

  get "/msgs" do
    return [400, NO_IS_REQUIRED] if params[:no].nil?

    count = params[:no].to_i

    return Message
        .unflagged
        .order(pub_date: :desc)
        .first(count)
        .map(&:sim_format)
        .to_json
  end

  get "/msgs/:username" do |username|
    user = User.find_by_username(username)
    return [404, USER_NOT_FOUND] if user.nil?
    return [400, NO_IS_REQUIRED] if params[:no].nil?

    count = params[:no].to_i

    return user.messages
        .unflagged
        .order(pub_date: :desc)
        .first(count)
        .map(&:sim_format)
        .to_json
  end

  post "/msgs/:username" do |username|
    user = User.find_by_username(username)
    return [404, USER_NOT_FOUND] if user.nil?

    request_data = JSON.parse(request.body.read, symbolize_names: true)
    message = user.messages.create(
      text: request_data[:content],
      flagged: false,
      pub_date: Time.now
    )

    if message
      status 204
    else
      status 400
      body "Something went wrong"
    end
  end

  get "/fllws/:username" do |username|
    user = User.find_by_username(username)
    return [404, USER_NOT_FOUND] if user.nil?
    return [400, NO_IS_REQUIRED] if params[:no].nil?

    count = params[:no].to_i

    following = user.following
      .first(count)
      .pluck(:username)

    status 200
    body({ follows: following }.to_json)
  end

  post "/fllws/:username" do |username|
    user = User.find_by_username(username)
    return [404, USER_NOT_FOUND] if user.nil?

    request_data = JSON.parse(request.body.read, symbolize_names: true)
    if request_data.key?(:follow)
      to_follow = User.find_by_username(request_data[:follow])
      return [400, "User to follow not found"] if to_follow.nil?

      user.following.append(to_follow)
      status 204
    elsif request_data.key?(:unfollow)
      to_unfollow = User.find_by_username(request_data[:unfollow])
      return [400, "User to unfollow not found"] if to_unfollow.nil?

      user.following.delete(to_unfollow)
      status 204
    end
  end
end