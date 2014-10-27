# encoding: UTF-8
# This file is auto-generated from the current state of the database. Instead
# of editing this file, please use the migrations feature of Active Record to
# incrementally modify your database, and then regenerate this schema definition.
#
# Note that this schema.rb definition is the authoritative source for your
# database schema. If you need to create the application database on another
# system, you should be using db:schema:load, not running all the migrations
# from scratch. The latter is a flawed and unsustainable approach (the more migrations
# you'll amass, the slower it'll run and the greater likelihood for issues).
#
# It's strongly recommended to check this file into your version control system.

ActiveRecord::Schema.define(:version => 20141027193434) do

  create_table "api_clients", :force => true do |t|
    t.string   "client_id"
    t.string   "client_secret"
    t.string   "name"
    t.datetime "created_at",    :null => false
    t.datetime "updated_at",    :null => false
  end

  add_index "api_clients", ["client_id"], :name => "index_api_clients_on_client_id"

  create_table "api_sessions", :force => true do |t|
    t.integer  "user_id"
    t.integer  "api_client_id"
    t.string   "session_token"
    t.datetime "created_at",    :null => false
    t.datetime "updated_at",    :null => false
  end

  add_index "api_sessions", ["session_token"], :name => "index_api_sessions_on_session_token"
  add_index "api_sessions", ["user_id"], :name => "index_api_sessions_on_user_id"

  create_table "password_resets", :force => true do |t|
    t.string   "token"
    t.integer  "user_id"
    t.boolean  "active",     :default => true
    t.datetime "created_at"
    t.datetime "updated_at"
    t.datetime "expires_at"
  end

  add_index "password_resets", ["active", "token"], :name => "index_password_resets_on_active_and_token"

  create_table "todos", :force => true do |t|
    t.string   "title"
    t.integer  "user_id"
    t.boolean  "complete",   :default => false
    t.datetime "created_at"
    t.datetime "updated_at"
  end

  create_table "users", :force => true do |t|
    t.string   "email"
    t.string   "name_first"
    t.string   "name_last"
    t.string   "crypted_password"
    t.string   "salt"
    t.datetime "created_at"
    t.datetime "updated_at"
    t.string   "ios_push_token"
  end

end
