# frozen_string_literal: true

class Message < ActiveRecord::Base
  belongs_to :author, class_name: "User"

  scope :unflagged, -> { where(flagged: false) }

  scope :authored_by, ->(users) { where(author: users) }

  validates_presence_of :text

  def sim_format
    {
      content: text,
      user: author.username,
      pub_date: created_at,
    }
  end
end
