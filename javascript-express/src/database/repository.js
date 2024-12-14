const { Users, Followers, Messages, Latest, sequelize } = require('./model');

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

const unfollowUser = async (whoId, whomId) =>
    !!(await Followers.destroy({
        where: {
            who_id: whoId,
            whom_id: whomId,
        },
    }));

const createNewUser = async (username, email, password) => {
    return !!(await Users.create({
        username: username,
        email: email,
        pw_hash: password,
    }));
};

const createMessage = async (userId, text) => {
    return !!(await Messages.create({
        author_id: userId,
        text: text,
        pub_date: new Date(),
        flagged: 0,
    }));
};

const getLatest = async () => {
    return await Latest.findOne({ where: { id: 1 } });
};

const updateLatest = async (latest) => {
    return !!(await Latest.update({ value: latest }, { where: { id: 1 } }));
};

const getUserFollows = async (userId) => {
    const followUsers = await Users.findAll({
        include: {
            model: Followers,
            where: {
                who_id: userId,
            },
        },
        attributes: ['user_id', 'username', 'email'],
        raw: true,
    });

    return followUsers.map((follow) => follow['username']);
};

module.exports = {
    getUserIdByName,
    getUser,
    getUserTimelineMessages,
    getPublicMessages,
    getUserMessages,
    isFollowing,
    getUserByName,
    followUser,
    unfollowUser,
    createNewUser,
    createMessage,
    updateLatest,
    getLatest,
    getUserFollows,
};
