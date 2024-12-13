const fs = require("fs");
const { Sequelize, Op, DataTypes } = require("sequelize");

const databaseUrl = process.env.DATABASE_URL;
if (!databaseUrl) {
  throw new Error("DATABASE_URL environment variable is not defined");
}

let sequelize = new Sequelize(databaseUrl);

//  -- Models --
const Users = sequelize.define(
  "Users",
  {
    user_id: {
      type: DataTypes.INTEGER,
      primaryKey: true,
      autoIncrement: true,
    },
    username: {
      type: DataTypes.STRING,
      allowNull: false,
    },
    email: {
      type: DataTypes.STRING,
      allowNull: false,
    },
    pw_hash: {
      type: DataTypes.STRING,
      allowNull: false,
    },
  },
  {
    tableName: "users",
    timestamps: false,
  }
);

const Followers = sequelize.define(
  "Followers",
  {
    who_id: {
      type: DataTypes.INTEGER,
      primaryKey: true,
    },
    whom_id: {
      type: DataTypes.INTEGER,
      primaryKey: true,
    },
  },
  {
    tableName: "followers",
    timestamps: false,
  }
);

const Messages = sequelize.define(
  "Messages",
  {
    message_id: {
      type: DataTypes.INTEGER,
      primaryKey: true,
      autoIncrement: true,
    },
    author_id: {
      type: DataTypes.INTEGER,
      allowNull: false,
    },
    text: {
      type: DataTypes.STRING,
      allowNull: false,
    },
    pub_date: {
      type: DataTypes.DATE,
    },
    flagged: {
      type: DataTypes.INTEGER,
    },
  },
  {
    tableName: "messages",
    timestamps: false,
  }
);

// Define associations between models
Messages.belongsTo(Users, { foreignKey: "author_id" });

// Define the many-to-many relationship between Users and Followers
// Define associations
// A user can follow many users
Users.belongsToMany(Users, {
  through: Followers,
  as: "UserFollowing", // Users that this user is following
  foreignKey: "who_id",
  otherKey: "whom_id",
});

// A user can be followed by many users
Users.belongsToMany(Users, {
  through: Followers,
  as: "UserFollowers", // Users that follow this user
  foreignKey: "whom_id",
  otherKey: "who_id",
});

const get_user_id = async (username) => {
  const user = await Users.findOne({
    attributes: ["user_id"],
    where: {
      username: username,
    },
  });
  return user ? user.user_id : null;
};

const validateUser = async (userid) => {
  const user = await Users.findOne({
    where: {
      user_id: userid,
    },
  });
  user !== null
};

const getUserTimelineMessages = async (userid, paginationLimit) => {
  const query = `
  ((SELECT users.user_id, users.username, users.email, 
          messages.message_id, messages.author_id, messages.text, messages.pub_date, messages.flagged 
    FROM followers
    INNER JOIN messages ON followers.whom_id = messages.author_id
    INNER JOIN users ON messages.author_id = users.user_id
    WHERE followers.who_id = :userId)
  UNION
  (SELECT users.user_id, users.username, users.email, 
          messages.message_id, messages.author_id, messages.text, messages.pub_date, messages.flagged 
    FROM messages
    INNER JOIN users ON messages.author_id = users.user_id
    WHERE users.user_id = :userId))
  ORDER BY pub_date DESC
  LIMIT :limit;
  `;

  const messages = await sequelize.query(query, {
    replacements: {
      userId: userid,
      limit: paginationLimit,
    },
    type: sequelize.QueryTypes.SELECT,
  });

  return messages;
}

module.exports = {
  Users,
  Followers,
  Messages,
  sequelize,
  Op,
  get_user_id,
  validateUser,
  getUserTimelineMessages
};
