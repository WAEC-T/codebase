class AddIndexToCreatedAt < ActiveRecord::Migration[7.1]
  disable_ddl_transaction!

  def change
    add_index :messages, :created_at, order: :desc, algorithm: :concurrently
  end
end
