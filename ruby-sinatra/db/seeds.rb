if ENV["APP_ENV"] == "test"
  user1 = User.create({username: "testuser", email: "test@user.dk", password: "1234", password_confirmation: "1234"})
  user2 = User.create({username: "otheruser", email: "other@user.dk", password: "1234", password_confirmation: "1234"})
  user3 = User.create({username: "thirduser", email: "third@user.dk", password: "1234", password_confirmation: "1234"})
  user1.following.append(user2)
  Message.create({author_id: user2.id, text: "This is a message written by user 2", flagged: false})
  Message.create({author_id: user3.id, text: "This is a message written by user 3", flagged: false})
end
