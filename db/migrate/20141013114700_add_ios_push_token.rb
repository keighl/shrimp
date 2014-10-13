class AddIosPushToken < ActiveRecord::Migration
  def change
    add_column :users, :ios_push_token, :string
  end

end
