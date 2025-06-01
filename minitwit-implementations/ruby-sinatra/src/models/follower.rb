# frozen_string_literal: true

# Follower model in minitwit following the database pattern.
class Follower < ActiveRecord::Base
  validates :who_id, presence: true
  validates :whom_id, presence: true

  def self.follows?(who_id, whom_id)
    exists?(who_id: who_id, whom_id: whom_id)
  end
end
