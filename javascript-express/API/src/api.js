const createError = require('http-errors');
const express = require('express');
const fs = require('fs');
const bcrypt = require('bcrypt');

// Configuration
const DEBUG = true;

// Import the sequlize functionality
const { Users, Messages, Followers, get_user_id } = require('../../utils/sequilize');
const { Sequelize } = require('sequelize');


fs.unlink("./latest_processed_sim_action_id.txt", (err) => {
    if (err && err.code !== 'ENOENT') {
        console.error('Error deleting latest processed file:', err);
    }
});

// Create our little application :)
const app = express();

// Setup
app.use(express.json()); // allows json in http request
app.use(express.urlencoded({ extended: false })); // same


function not_req_from_simulator(req) {
    const fromSimulator = req.headers.authorization;
    if (fromSimulator !== "Basic c2ltdWxhdG9yOnN1cGVyX3NhZmUh") {
        const error = "You are not authorized to use this resource!";
        return error;
    }
    return null;
}

async function update_latest(request) {
    const parsed_command_id = request.query.latest;
    if (parsed_command_id !== -1) {
        fs.writeFile("./latest_processed_sim_action_id.txt", parsed_command_id.toString(), (err) => {
            if (err) {
                console.error('Error writing file:', err);
            }
        });
    }
}

// --------------- Routes ------------------

// Get the latest value
app.get('/latest', async (req, res) => {
    fs.readFile('./latest_processed_sim_action_id.txt', 'utf8', (err, content) => {
        if (err) {
            res.json({ latest: -1 });
        } else {
            res.json({ latest: parseInt(content) });
        }
    });
});

// ------- Route to register a user ----------
app.post('/register', async (req, res) => {
    await update_latest(req);
    const request_data = req.body;

    let error = null;
    if (!request_data.username) {
        error = "You have to enter a username";
    } else if (!request_data.email || !request_data.email.includes('@')) {
        error = "You have to enter a valid email address";
    } else if (!request_data.pwd) {
        error = "You have to enter a password";
    } else {
        const hash_password = bcrypt.hashSync(request_data.pwd, 10)
        await Users.create({
            username: request_data.username,
            email: request_data.email,
            pw_hash: hash_password
        });
    }
    if (error) {
        res.status(400).json({ status: 400, error_msg: error });
    } else {
        res.sendStatus(204);
    }

});

// ---------- route to get messages by user -------------
app.get('/msgs/:username', async (req, res) => {
    await update_latest(req);

    const not_from_sim_response = not_req_from_simulator(req);
    if (not_from_sim_response) {
        return res.send(not_from_sim_response);
    }

    const username = req.params.username;
    const no_msgs = parseInt(req.query.no) || 100; // Default to 100 if 'no' parameter is not provided

    const user = await Users.findOne({
        attributes: ['user_id'],
        where: {
            username: username
        }
    });

    // Refactored using Sequelize
    const messages = await Messages.findAll({
        attributes: ['message_id', 'author_id', 'text', 'pub_date'],
        include: [{
            model: Users,
            as: 'User',
            attributes: ['user_id', 'username', 'email']
        }],
        where: {
            author_id: user.user_id
        },
        order: [['message_id', 'DESC']],
        limit: no_msgs,
        raw: true
    });

/*     const filtered_msgs = messages.map(msg => ({
        content: msg.text,
        pub_date: msg.pub_date,
        user: msg['User.username']
    })); */

    res.json(messages);
});

// ------------ Route to post a Messages by a given user --------------
app.post('/msgs/:username', async (req, res) => {
    await update_latest(req);

    const not_from_sim_response = not_req_from_simulator(req);
    if (not_from_sim_response) {
        return res.send(not_from_sim_response);
    }

    const username = req.params.username;
    const { content } = req.body;

    const user = await Users.findOne({
        attributes: ['user_id'],
        where: {
            username: username
        }
    });

    console.log("user: ", user);

    await Messages.create({
        author_id: user.user_id,
        text: content,
        pub_date: Math.floor(Date.now() / 1000),
        flagged: 0
    });
    res.sendStatus(204);
})

// ------------ Route to get all messages in Database ----------------
app.get('/msgs', async (req, res) => {
    await update_latest(req);

    const not_from_sim_response = not_req_from_simulator(req);
    if (not_from_sim_response) {
        return res.status(403).json({ error: not_from_sim_response });
    }

    const no_msgs = parseInt(req.query.no, 10) || 100;

    // Query the database to get messages
    const messages = await Messages.findAll({
        attributes: [
            'message_id',
            'author_id',
            'text',
            'pub_date',
            [Sequelize.literal('"User"."user_id"'), 'user_id'],
            [Sequelize.literal('"User"."username"'), 'username'],
            [Sequelize.literal('"User"."email"'), 'email']
        ],
        include: [{
            model: Users,
            as: 'User', // Ensure this matches the alias used in your join
            attributes: [], // Don't fetch any additional attributes from the Users model
        }],
        where: { flagged: 0 },
        order: [['message_id', 'DESC']],
        limit: no_msgs,
        raw: true
    });

    const filtered_msgs = messages.map(msg => ({
        content: msg.text,
        pub_date: msg.pub_date,
        user: msg.username
    }));

    return res.json(filtered_msgs);
});

// -------------- Route to get the followers of a given user ------------------
app.get('/fllws/:username', async (req, res) => {
    await update_latest(req);

    const not_from_sim_response = not_req_from_simulator(req);
    if (not_from_sim_response) {
        return res.send(not_from_sim_response);
    }

    const username = req.params.username;
    const user = await Users.findOne({
        attributes: ['user_id'],
        where: { username: username }
    });

    if (!user) {
        return res.status(404).send("User not found");
    }

    try {
        const followers = await Followers.findAll({
            attributes: ['who_id'],
            where: {
                whom_id: user.user_id,
            },
            raw: true
        });
        const followerList = followers.map(x => x.who_id).concat();
        console.log("followerList: ", followerList);
        const users = await Users.findAll({
            attributes: ['user_id', 'username', 'email'],
            where: {
              user_id: followerList
            },
            raw: true
        });
        return res.json({ follows: users });
    } catch (error) {
        console.error('Error fetching followers:', error);
        return res.status(500).send('Internal Server Error');
    }
});

// ------------ Route to add/delete a follower --------------
app.post('/fllws/:username', async (req, res) => {
    await update_latest(req);

    const not_from_sim_response = not_req_from_simulator(req);
    if (not_from_sim_response) {
        return res.send(not_from_sim_response);
    }

    const username = req.params.username;
    const user = await Users.findOne({
        attributes: ['user_id'],
        where: {
            username: username
        }
    });
    if (!user) {
        return res.status(404).send('User not found');
    }

    const { follow, unfollow } = req.body;

    try {
        // ------------ CASE FOLLOW USER ------------------
        if (follow) {
            const follows_user = await Users.findOne({
                attributes: ['user_id'],
                where: {
                    username: follow
                }
            });
            if (!follows_user) {
                return res.status(404).send('User to follow not found');
            }
            console.log("follow: ", follows_user.user_id);
            console.log("follow: ", user.user_id);
            await Followers.create({
                who_id: user.user_id,
                whom_id: follows_user.user_id
            });
            return res.sendStatus(204);

            // ------------ CASE UNFOLLOW USER ------------------
        } else if (unfollow) {
            const unfollows_user = await Users.findOne({
                attributes: ['user_id'],
                where: {
                    username: unfollow
                }
            });
            if (!unfollows_user) {
                return res.status(404).send('User to unfollow not found');
            }

            await Followers.destroy({
                where: {
                    who_id: user.user_id,
                    whom_id: unfollows_user.user_id
                }
            });
            return res.sendStatus(204);
        } else {
            return res.status(400).send('Invalid request');
        }
    } catch (error) {
        console.error('Error processing follow/unfollow request:', error);
        return res.status(500).send('Internal Server Error');
    }
});

app.listen(5001, () => {
    console.log('Listening on port 5001')
})