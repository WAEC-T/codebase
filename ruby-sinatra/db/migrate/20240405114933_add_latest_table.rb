class AddLatestTable < ActiveRecord::Migration[7.1]
  def change
    create_table :latest do |t|
      t.integer :value
    end

    begin
      content = File.read("latest_processed_sim_action_id.txt")
      latest_processed_command_id = content.to_i
      ActiveRecord::Base.connection.execute("INSERT INTO latest (value) VALUES (#{latest_processed_command_id})")
    rescue StandardError
      latest_processed_command_id = -1
      ActiveRecord::Base.connection.execute("INSERT INTO latest (value) VALUES (#{latest_processed_command_id})")
    end
  end
end
