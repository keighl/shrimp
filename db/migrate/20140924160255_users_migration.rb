class UsersMigration < ActiveRecord::Migration
  def change
    create_table :users do |t|
      t.string :email
      t.string :name_first
      t.string :name_last
      t.string :crypted_password
      t.string :salt
      t.datetime :created_at
      t.datetime :updated_at
      t.string :api_token
    end

    add_index :users, :api_token
  end
end
