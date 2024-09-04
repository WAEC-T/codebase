const fs = require('fs');
const { Sequelize, Op, Model, DataTypes } = require('sequelize');

const dbConnectionType = process.env.DB_CONNECTION_TYPE;

let sequelize;
if (dbConnectionType === 'dev_db') {
  sequelize = new Sequelize('waect', 'user', 'pass', {
    host: 'local_database',
    post: 5432,
    dialect: 'postgres'
  });
} else if (dbConnectionType === 'prod') {

  const ca_file = fs.readFileSync('/express-docker/secrets/ca-certificate.crt');

  sequelize = new Sequelize(process.env.POSTGRES_NAME, process.env.POSTGRES_USER, process.env.POSTGRES_PASSWORD, {
    host: process.env.POSTGRES_HOST,
    port: process.env.POSTGRES_PORT,
    dialect: 'postgres',
    dialectOptions: {
      ssl: {
        ca: ca_file
      }
    }
  });
  //throw new Error('Not implemented yet');
}

//  -- Models --
const Users = sequelize.define('Users', {
  user_id: {
    type: DataTypes.INTEGER,
    primaryKey: true,
    autoIncrement: true // Indicates that it's a SERIAL column in PostgreSQL
  },
  username: {
    type: DataTypes.STRING,
    allowNull: false
  },
  email: {
    type: DataTypes.STRING,
    allowNull: false
  },
  pw_hash: {
    type: DataTypes.STRING,
    allowNull: false
  }
}, {
  tableName: 'users', // Ensure Sequelize uses the correct table name
  timestamps: false // Disable Sequelize's default timestamp fields
});

const Followers = sequelize.define('Followers', {
  who_id: {
    type: DataTypes.INTEGER,
    primaryKey: true
  },
  whom_id: {
    type: DataTypes.INTEGER,
    primaryKey: true
  }
}, {
  tableName: 'followers', // Ensure Sequelize uses the correct table name
  timestamps: false // Disable Sequelize's default timestamp fields
});

const Messages = sequelize.define('Messages', {
  message_id: {
    type: DataTypes.INTEGER,
    primaryKey: true,
    autoIncrement: true // Indicates that it's a SERIAL column in PostgreSQL
  },
  author_id: {
    type: DataTypes.INTEGER,
    allowNull: false
  },
  text: {
    type: DataTypes.STRING,
    allowNull: false
  },
  pub_date: {
    type: DataTypes.INTEGER
  },
  flagged: {
    type: DataTypes.INTEGER
  }
}, {
  tableName: 'messages', // Ensure Sequelize uses the correct table name
  timestamps: false // Disable Sequelize's default timestamp fields
});

// Define associations between models
Messages.belongsTo(Users, { foreignKey: 'author_id' });

// Define the many-to-many relationship between Users and Followers
// Define associations
// A user can follow many users
Users.belongsToMany(Users, {
  through: Followers,
  as: 'UserFollowing',     // Users that this user is following
  foreignKey: 'who_id',
  otherKey: 'whom_id'
});

// A user can be followed by many users
Users.belongsToMany(Users, {
  through: Followers,
  as: 'UserFollowers',     // Users that follow this user
  foreignKey: 'whom_id',
  otherKey: 'who_id'
});
// Test database operations
async function testDatabase() {
  // Inside the testDatabase function, after authenticating and before syncing models
  console.log('Defined Sequelize models:', Object.keys(sequelize.models));

  try {
    // Test connection by authenticatingf
    await sequelize.authenticate();
    console.log('Connection has been established successfully.');

    // Sync the model with the database
    await Users.sync();
    await Messages.sync();

    if (!(Users === sequelize.models.Users)) {
      throw new Error("Something is wrong with the Users model typing");
    }

    // Create a new user
    //const newUser = await Users.create({user_id: -1, username: 'testuser', email: 'test@example.com', pw_hash: "basj" });
    //console.log('New user created:', newUser.toJSON());
    const newMessage = await Messages.create({ author_id: -1, text: 'oggeli boggely', pub_date: 1, flagged: 0 });
    console.log('New Message created:', newMessage.toJSON());

    // Retrieve all messages from the database
    const messages = await Messages.findAll({ raw: true });

    // Display the retrieved messages
    console.log('All messages:', messages);

  } catch (error) {
    console.error('Unable to connect to the database:', error);
  } finally {
    // Close the Sequelize connection
    await sequelize.close();
  }
}

const get_user_id = async (username) => {
  const user = await Users.findOne({
    attributes: ['user_id'],
    where: {
      username: username
    }
  });
  return user ? user.user_id : null;
};


// Export the models
module.exports = { testDatabase, Users, Followers, Messages, sequelize, Op, get_user_id };