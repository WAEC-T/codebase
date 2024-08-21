# This file is auto-generated from the current state of the database. Instead
# of editing this file, please use the migrations feature of Active Record to
# incrementally modify your database, and then regenerate this schema definition.
#
# This file is the source Rails uses to define your schema when running `bin/rails
# db:schema:load`. When creating a new database, `bin/rails db:schema:load` tends to
# be faster and is potentially less error prone than running all of your
# migrations from scratch. Old migrations may fail to apply correctly if those
# migrations use external dependencies or application code.
#
# It's strongly recommended that you check this file into your version control system.

ActiveRecord::Schema[7.1].define(version: 2024_04_12_083347) do
  # These are extensions that must be enabled in order to support this database
  enable_extension "plpgsql"

  create_table "followers", force: :cascade do |t|
    t.integer "who_id"
    t.integer "whom_id"
    t.index ["who_id", "whom_id"], name: "index_followers_on_who_id_and_whom_id"
  end

  create_table "latest", force: :cascade do |t|
    t.integer "value"
  end

  create_table "messages", force: :cascade do |t|
    t.bigint "author_id"
    t.string "text", null: false
    t.boolean "flagged"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["author_id"], name: "index_messages_on_author_id"
    t.index ["created_at"], name: "index_messages_on_created_at", order: :desc
    t.index ["flagged"], name: "index_messages_on_flagged"
  end

  create_table "users", force: :cascade do |t|
    t.string "username", null: false
    t.string "email", null: false
    t.string "password_digest", null: false
    t.index ["username"], name: "index_users_on_username"
  end

  add_foreign_key "messages", "users", column: "author_id"
end
