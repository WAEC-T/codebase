# frozen_string_literal: true

require "sqlite3"
require "pg"
require "./myapp"

namespace :one_offs do
  task :pg_migrate do
    # rubocop:disable all
    # Rubocop is disabled for this task for the following reason:
    # This task has already been run and is really only here for archival purposes.
    # I don't think we should mess with this code as this is the code that was run in production.
    # - atro

    if User.any? || Follower.any? || Message.any?
      raise "Current database contains data"
    end

    sqlite_db = SQLite3::Database.new "./db/minitwit.db"
    sq_users = sqlite_db.execute "SELECT user_id, username, email, password_digest FROM user"
    sq_users.map! do |id, name, email, pwd|
      {
        id: id,
        username: name,
        email: email,
        password_digest: pwd || ""
      }
    end
    User.insert_all(sq_users)

    sq_messages = sqlite_db.execute "SELECT * FROM message"
    sq_messages.map! do |id, author, text, pubdate, flagged|
      time = Time.at(pubdate)
      {
        id: id,
        author_id: author,
        text: text,
        flagged: flagged == 1,
        created_at: time,
        updated_at: time
      }
    end
    Message.insert_all(sq_messages)

    sq_followers = sqlite_db.execute "SELECT * FROM follower"
    sq_followers.map! do |who_id, whom_id|
      {
        who_id: who_id,
        whom_id: whom_id
      }
    end
    Follower.insert_all(sq_followers)
  end
end
