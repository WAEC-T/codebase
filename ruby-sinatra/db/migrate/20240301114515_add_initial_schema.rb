class AddInitialSchema < ActiveRecord::Migration[7.1]
  def change
    create_table :users do |t|
      t.string :username, null: false
      t.string :email, null: false
      t.string :password_digest, null: false
    end

    create_table :followers do |t|
      t.integer :who_id
      t.integer :whom_id
    end

    create_table :messages do |t|
      t.references :author, foreign_key: { to_table: :users }
      t.string :text, null: false
      t.boolean :flagged

      t.timestamps
    end
  end
end
