class CreateTodos < ActiveRecord::Migration
  def change
    create_table :todos do |t|
      t.string :title
      t.integer :user_id
      t.boolean :complete, default: false
      t.datetime :created_at
      t.datetime :updated_at
    end

    add_index :todos, [:id, :user_id]
  end
end
