# frozen_string_literal: true

require "sinatra/activerecord"
require "sinatra"
require "json"
require "newrelic_rpm"
Dir["./models/*.rb"].each { |file| require file }

DATABASE_URL = ENV.fetch("DATABASE_URL", nil)

NO_IS_REQUIRED = "no is required"

USER_NOT_FOUND = "User not found"

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
end

helpers do
  def request_is_not_from_simulator
    request.env["HTTP_AUTHORIZATION"] != "Basic c2ltdWxhdG9yOnN1cGVyX3NhZmUh"
  end
end

def update_latest(request)
  parsed_command_id = request.params["latest"]
  Latest.set(parsed_command_id.to_i) if parsed_command_id.present?
end

before do
  if request.path_info == "/latest"
    pass
  elsif request_is_not_from_simulator
    halt [403, { status: 403, error_msg: "You are not authorized to use this resource!" }.to_json]
  end

  update_latest(request)
end

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
      .order(created_at: :desc)
      .first(count)
      .map(&:sim_format)
      .to_json
end

get "/msgs/:username" do |username|
  user = User.find_by_username(username)
  return [400, USER_NOT_FOUND] if user.nil?
  return [400, NO_IS_REQUIRED] if params[:no].nil?

  count = params[:no].to_i

  return user.messages
      .unflagged
      .order(created_at: :desc)
      .first(count)
      .map(&:sim_format)
      .to_json
end

post "/msgs/:username" do |username|
  user = User.find_by_username(username)
  return [400, USER_NOT_FOUND] if user.nil?

  request_data = JSON.parse(request.body.read, symbolize_names: true)
  message = user.messages.create(
    text: request_data[:content],
    flagged: false,
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
