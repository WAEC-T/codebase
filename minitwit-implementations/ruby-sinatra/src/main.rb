# frozen_string_literal: true

require 'sinatra/activerecord'
require 'sinatra'
require 'sinatra/namespace'
require 'sinatra/flash'
require 'json'
require 'newrelic_rpm'
require 'dotenv'

Dir['./models/*.rb'].sort.each { |file| require file }

set :logging, false

PR_PAGE = 30

DATABASE_URL = ENV.fetch('DATABASE_URL', nil)
puts "Loaded DATABASE_URL: #{DATABASE_URL}"

USER_NOT_FOUND = 'User not found'
NO_IS_REQUIRED = 'no is required'
API_MESSAGE_RESPONSE = 100

configure :production, :staging, :development, :test do
  db = URI.parse(DATABASE_URL)
  set :database, {
    adapter: db.scheme,
    host: db.host,
    port: db.port,
    database: db.path[1..],
    user: db.user,
    password: db.password,
    encoding: 'utf8'
  }
  if settings.environment == :development || settings.environment == :test
    enable :logging
    ActiveRecord.verbose_query_logs = true
    ActiveRecord::Base.logger = Logger.new($stdout)
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
    request.env['HTTP_AUTHORIZATION'] != 'Basic c2ltdWxhdG9yOnN1cGVyX3NhZmUh'
  end
end

def update_latest(request)
  parsed_command_id = request.params['latest']
  Latest.set(parsed_command_id.to_i) if parsed_command_id.present?
end

before '/api/*' do
  update_latest(request) if request.path != '/api/latest'
  if request_is_not_from_simulator
    halt [403, { status: 403, error_msg: 'You are not authorized to use this resource!' }.to_json]
  end
end

get '/' do
  redirect('/public') unless logged_in?

  @title = 'My Timeline'

  @messages = Message
              .joins('LEFT JOIN followers f ON f.whom_id = messages.author_id')
              .where('f.who_id = ? OR messages.author_id = ?', session[:user_id], session[:user_id])
              .includes(:author)
              .order(pub_date: :desc)
              .limit(PR_PAGE)

  erb :timeline, layout: :layout
end

# TODO: Figure out why it is only one query (probably the signout) - It needs to be 2 queries
get '/public' do
  @title = 'Public Timeline'
  @messages = Message
              .includes(:author)
              .order(pub_date: :desc)
              .first(PR_PAGE)

  erb :timeline, layout: :layout
end

get '/register' do
  erb :register, layout: :layout
end

post '/register' do
  existing_user = User.find_by(username: params[:username])
  if existing_user
    flash[:error] = 'The username is already taken'
    redirect('/register')
    return
  end

  user = User.new(
    username: params[:username],
    email: params[:email],
    password: params[:password],
    password_confirmation: params[:password2]
  )
  if user.save
    flash[:success] = 'You were successfully registered and can login now'
    redirect('/login')
  else
    errors = user.errors.map(&:full_message).join(', ')
    flash[:error] = errors
  end
  redirect('/register')
end

get '/login' do
  erb :login, layout: :layout
end

post '/login' do
  user = User.find_by_username(params[:username])
  if user.nil?
    error = 'Invalid username'
  elsif !user.authenticate(params[:password])
    error = 'Invalid password'
  else
    session[:user_id] = user.id
    flash[:success] = 'You were logged in'
    redirect('/')
  end
  flash[:error] = error
  redirect('/login')
end

get '/logout' do
  session[:user_id] = nil
  flash[:success] = 'You were logged out'
  redirect('/')
end

post '/add_message' do
  if !session[:user_id]
    return status 401
  elsif params[:text].nil? || params[:text].strip.empty?
    flash[:error] = 'Message cannot be empty!'
  elsif Message.create(
    author_id: session[:user_id],
    text: params[:text],
    flagged: false,
    pub_date: Time.now
  )
    flash[:success] = 'Your message was recorded'
  else
    flash[:error] = 'Something went wrong :('
  end

  redirect(request.referer)
end

get '/user/:username' do
  @profile_user = User.find_by_username(params[:username])
  @title = "#{@profile_user.username}'s Timeline"
  @follows = Follower.follows?(current_user.id, @profile_user.id)
  @messages = Message
              .authored_by(@profile_user)
              .includes(:author)
              .order(pub_date: :desc)
              .first(PR_PAGE)

  erb :timeline, layout: :layout
end

get '/:username/follow' do
  return status 401 unless logged_in?

  whom = User.find_by_username(params[:username])
  return status 404 if whom.nil?

  Follower.find_or_create_by(who_id: session[:user_id], whom_id: whom.id)
  flash[:success] = "You are now following #{params[:username]}"
  redirect("/user/#{params[:username]}")
end

get '/:username/unfollow' do
  return status 401 unless logged_in?

  whom = User.find_by_username(params[:username])
  return status 404 if whom.nil?

  Follower.where(who_id: session[:user_id], whom_id: whom.id).delete_all
  flash[:success] = "You are no longer following #{params[:username]}"
  redirect("/user/#{params[:username]}")
end

namespace '/api' do
  get '/latest' do
    body({ latest: Latest.get }.to_json)
  end

  post '/register' do
    request_data = JSON.parse(request.body.read, symbolize_names: true)

    existing_user = User.find_by(username: request_data[:username])
    if existing_user
      status 400
      body({ status: 400, error_msg: ' a user with this email or username already exists.' }.to_json)
      return
    end

    user = User.new(
      username: request_data[:username],
      email: request_data[:email],
      password: request_data[:pwd],
      password_confirmation: request_data[:pwd]
    )

    if user.save
      status 204
    else
      status 400
      body({ status: 400, error_msg: user.errors.map(&:full_message) }.to_json)
    end
  end

  get '/msgs' do
    number_of_messages = params[:no]&.to_i&.positive? ? params[:no].to_i : API_MESSAGE_RESPONSE

    return Message
           .order(pub_date: :desc)
           .first(number_of_messages)
           .map(&:sim_format)
           .to_json
  end

  get '/msgs/:username' do |username|
    user = User.find_by_username(username)
    return [404, USER_NOT_FOUND] if user.nil?

    number_of_messages = params[:no]&.to_i&.positive? ? params[:no].to_i : API_MESSAGE_RESPONSE

    return user.messages
               .order(pub_date: :desc)
               .first(number_of_messages)
               .map(&:sim_format)
               .to_json
  end

  post '/msgs/:username' do |username|
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
      body 'Something went wrong'
    end
  end

  get '/fllws/:username' do |username|
    user = User.find_by_username(username)
    return [404, USER_NOT_FOUND] if user.nil?

    number_of_messages = params[:no]&.to_i&.positive? ? params[:no].to_i : API_MESSAGE_RESPONSE

    following_usernames = User.joins('JOIN followers ON followers.whom_id = users.user_id')
                              .where('followers.who_id = ?', user.id)
                              .limit(number_of_messages)
                              .pluck('users.username')

    status 200
    body({ follows: following_usernames }.to_json)
  end

  post '/fllws/:username' do |username|
    user = User.find_by_username(username)
    return [404, USER_NOT_FOUND] if user.nil?

    request_data = JSON.parse(request.body.read, symbolize_names: true)
    if request_data.key?(:follow)
      to_follow = User.find_by_username(request_data[:follow])
      return [400, 'User to follow not found'] if to_follow.nil?

      Follower.find_or_create_by(who_id: user.id, whom_id: to_follow.id)
      status 204
    elsif request_data.key?(:unfollow)
      to_unfollow = User.find_by_username(request_data[:unfollow])
      return [400, 'User to unfollow not found'] if to_unfollow.nil?

      Follower.where(who_id: user.id, whom_id: to_unfollow.id).delete_all
      status 204
    end
  end
end
