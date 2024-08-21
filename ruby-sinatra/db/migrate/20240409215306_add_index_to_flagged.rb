class AddIndexToFlagged < ActiveRecord::Migration[7.1]
  disable_ddl_transaction!

  def change
    add_index :messages, :flagged, algorithm: :concurrently
  end
end
