const { Sequelize, Op, DataTypes } = require('sequelize');

const databaseUrl = process.env.DATABASE_URL;

if (!databaseUrl) {
    throw new Error('DATABASE_URL environment variable is not defined');
}

let sequelize = new Sequelize(databaseUrl);

//  -- Models --
const Users = sequelize.define(
    'Users',
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
        tableName: 'users',
        timestamps: false,
    }
);

const Followers = sequelize.define(
    'Followers',
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
        tableName: 'followers',
        timestamps: false,
    }
);

const Messages = sequelize.define(
    'Messages',
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
        tableName: 'messages',
        timestamps: false,
    }
);

// Define associations between models
Messages.belongsTo(Users, { foreignKey: 'author_id' });

// Define the many-to-many relationship between Users and Followers
// Define associations
// A user can follow many users
Users.belongsToMany(Users, {
    through: Followers,
    as: 'UserFollowing', // Users that this user is following
    foreignKey: 'who_id',
    otherKey: 'whom_id',
});

// A user can be followed by many users
Users.belongsToMany(Users, {
    through: Followers,
    as: 'UserFollowers', // Users that follow this user
    foreignKey: 'whom_id',
    otherKey: 'who_id',
});

const getUserIdByName = async (username) => {
    const user = await Users.findOne({
        attributes: ['user_id'],
        where: {
            username: username,
        },
    });
    return user ? user.user_id : null;
};

const getUser = async (userId) => {
    const user = await Users.findOne({
        where: {
            user_id: userId,
        },
        raw: true,
    });
    return user;
};

const getUserTimelineMessages = async (userId, paginationLimit) => {
    // maximizing result minimizing effort ~~ ¯\_(ツ)_/¯
    const query = `
    ((SELECT users.user_id, users.username, users.email, messages.message_id, messages.author_id, messages.text, messages.pub_date, messages.flagged 
      FROM followers
      INNER JOIN messages ON followers.whom_id = messages.author_id
      INNER JOIN users ON messages.author_id = users.user_id
      WHERE followers.who_id = :userId)
    UNION
    (SELECT users.user_id, users.username, users.email, messages.message_id, messages.author_id, messages.text, messages.pub_date, messages.flagged 
      FROM messages
      INNER JOIN users ON messages.author_id = users.user_id
      WHERE users.user_id = :userId))
    ORDER BY pub_date DESC
    LIMIT :limit;
  `;

    const messages = await sequelize.query(query, {
        replacements: {
            userId: userId,
            limit: paginationLimit,
        },
        type: sequelize.QueryTypes.SELECT,
    });

    return messages;
};

const getPublicMessages = async (paginationLimit) => {
    const messages = await Messages.findAll({
        include: [
            {
                model: Users,
                required: true,
                attributes: ['user_id', 'username', 'email'],
            },
        ],
        order: [['pub_date', 'DESC']],
        limit: paginationLimit,
        raw: true,
    });

    return messages;
};

const getUserMessages = async (userId, paginationLimit) => {
    const messages = await Messages.findAll({
        where: {
            author_id: userId,
        },
        include: [
            {
                model: Users,
                required: true,
                attributes: ['user_id', 'username', 'email'],
            },
        ],
        order: [['pub_date', 'DESC']],
        limit: paginationLimit,
        raw: true,
    });

    return messages;
};

const isFollowing = async (sessionUser, profileUser) => {
    var follows = await Followers.findOne({
        where: { who_id: sessionUser.user_id, whom_id: profileUser.user_id },
    });

    return !!follows;
};

const getUserByName = async (username) => {
    const user = await Users.findOne({
        where: {
            username: username,
        },
        raw: true,
    });
    return user;
};

const followUser = async (whoId, whomId) => {
    return !!(await Followers.create({
        who_id: whoId,
        whom_id: whomId,
    }));
};

module.exports = {
    Users,
    Followers,
    Messages,
    sequelize,
    Op,
    getUserIdByName,
    getUser,
    getUserTimelineMessages,
    getPublicMessages,
    getUserMessages,
    isFollowing,
    getUserByName,
    followUser,
};
