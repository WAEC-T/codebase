# frozen_string_literal: true

class Latest < ActiveRecord::Base
  self.table_name = "latest"

  class << self
    def set(latest)
      Latest.first.update!(id: latest)
    end

    def get
      Latest.first.id
    end
  end
end
