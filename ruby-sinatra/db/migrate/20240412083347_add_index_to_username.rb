class AddIndexToUsername < ActiveRecord::Migration[7.1]
  disable_ddl_transaction!

  def change
    add_index :users, :username, algorithm: :concurrently
  end
end
