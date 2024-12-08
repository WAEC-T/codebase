# frozen_string_literal: true

# Represents a message model in minitwit.
class Message < ActiveRecord::Base
  belongs_to :author, class_name: 'User'

  scope :authored_by, ->(users) { where(author: users) }

  validates_presence_of :text
  validates_presence_of :pub_date

  def sim_format
    {
      content: text,
      user: author.username,
      pub_date: pub_date
    }
  end
end
