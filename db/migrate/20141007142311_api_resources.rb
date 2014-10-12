class ApiResources < ActiveRecord::Migration

  def change

    create_table :api_clients do |t|
      t.string   :client_id
      t.string   :client_secret
      t.string   :name
      t.timestamps
    end

    create_table :api_sessions do |t|
      t.integer  :user_id
      t.integer  :api_client_id
      t.string   :session_token
      t.timestamps
    end

    add_index :api_sessions, :session_token
    add_index :api_sessions, :user_id
    add_index :api_clients, :client_id
  end
end
