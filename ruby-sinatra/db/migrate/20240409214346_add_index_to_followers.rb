class AddIndexToFollowers < ActiveRecord::Migration[7.1]
  # https://api.rubyonrails.org/v7.1.3.2/classes/ActiveRecord/ConnectionAdapters/SchemaStatements.html#method-i-add_index
  disable_ddl_transaction!

  def change
    add_index :followers, [:who_id, :whom_id], algorithm: :concurrently
  end
end
