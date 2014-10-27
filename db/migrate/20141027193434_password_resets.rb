class PasswordResets < ActiveRecord::Migration
  def change
    create_table :password_resets do |t|
      t.string :token
      t.integer :user_id
      t.boolean :active, default: true
      t.datetime :created_at
      t.datetime :updated_at
      t.datetime :expires_at
    end

    add_index :password_resets, [:active, :token]
  end
end

