const { Sequelize, DataTypes } = require('sequelize');

const databaseUrl = String(process.env.DATABASE_URL);

if (!databaseUrl) {
    throw new Error('DATABASE_URL environment variable is not defined');
}

const connectionOptions = { logging: false };

const sequelize = new Sequelize(databaseUrl, {
    dialect: 'postgres',
    ...connectionOptions,
});

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

const Latest = sequelize.define(
    'Latest',
    {
        id: {
            type: DataTypes.INTEGER,
            primaryKey: true,
            autoIncrement: true,
        },
        value: {
            type: DataTypes.INTEGER,
            allowNull: false,
        },
    },
    {
        tableName: 'latest',
        timestamps: false,
    }
);

Users.hasMany(Followers, { foreignKey: 'whom_id' });
Messages.belongsTo(Users, { foreignKey: 'author_id' });
Followers.belongsTo(Users, { foreignKey: 'who_id' });

module.exports = {
    Users,
    Followers,
    Messages,
    Latest,
    sequelize,
};
