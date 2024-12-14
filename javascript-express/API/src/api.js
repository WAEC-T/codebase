const express = require('express');

const { Users, Messages, Followers } = require('../../src/database/model');
const {
    updateLatest,
    getLatest,
    createNewUser,
    getUserIdByName,
    getPublicMessages,
    getUserMessages,
    getUserByName,
} = require('../../src/database/repository');
const { validateRegisterFields, formatMessages } = require('../../src/utils');

const API_DEFAULT_MESSAGES_AMOUNT = 100;

const router = express.Router();

const checkSimulatorToken = (authToken) => {
    return authToken !== 'Basic c2ltdWxhdG9yOnN1cGVyX3NhZmUh'
        ? 'You are not authorized to use this resource!'
        : null;
};

const validateRequest = async (req) => {
    const latestValue = req.query.latest;

    if (latestValue) await updateLatest(latestValue);

    return checkSimulatorToken(req.headers.authorization);
};

router.get('/latest', async (req, res) => {
    const latest = getLatest() ?? -1;
    res.json({ latest: parseInt(latest) });
});

router.post('/register', async (req, res) => {
    const unauthorized = validateRequest(req);

    if (unauthorized) {
        return res.status(401).send(unauthorized);
    } else if (await getUserIdByName(username)) {
        return res.status(400).send('The username is already taken');
    }

    const { username, email, pwd } = req.body;

    let error = validateRegisterFields(username, email, pwd, true);

    if (!error && (await createNewUser(username, email, pwd))) {
        res.sendStatus(204);
    } else {
        res.status(400).json({ status: 400, error_msg: error });
    }
});

router.get('/msgs', async (req, res) => {
    const unauthorized = validateRequest(req);

    if (unauthorized) {
        return res.status(401).send(unauthorized);
    }

    const messagesAmount =
        parseInt(req.query.no, 10) || API_DEFAULT_MESSAGES_AMOUNT;

    const messages = formatMessages(
        await getPublicMessages(messagesAmount),
        true
    );

    return res.json(messages);
});

router.get('/msgs/:username', async (req, res) => {
    const unauthorized = validateRequest(req);

    if (unauthorized) {
        return res.status(401).send(unauthorized);
    }

    const username = req.params.username;

    if (!username) {
        req.status(400).send('username not provided!');

        return;
    }

    if (getUserByName(username)) {
        res.status(404).send('User not found');

        return;
    }

    const messagesAmount =
        parseInt(req.query.no) || API_DEFAULT_MESSAGES_AMOUNT;

    const messages = formatMessages(
        getUserMessages(username, messagesAmount),
        true
    );

    res.json(messages);
});

// ------------ Route to post a Messages by a given user --------------
router.post('/msgs/:username', async (req, res) => {
    const unauthorized = validateRequest(req);

    if (unauthorized) {
        return res.status(401).send(unauthorized);
    }

    const username = req.params.username;

    const { content } = req.body;

    if (!username) {
        req.status(400).send('username or content not provided!');
    }

    const userId = getUserIdByName(username);

    const user = await Users.findOne({
        attributes: ['user_id'],
        where: {
            username: username,
        },
    });

    console.log('user: ', user);

    await Messages.create({
        author_id: user.user_id,
        text: content,
        pub_date: Math.floor(Date.now() / 1000),
        flagged: 0,
    });
    res.sendStatus(204);
});

// -------------- Route to get the followers of a given user ------------------
router.get('/fllws/:username', async (req, res) => {
    await update_latest(req);

    const not_from_sim_response = not_req_from_simulator(req);
    if (not_from_sim_response) {
        return res.send(not_from_sim_response);
    }

    const username = req.params.username;
    const user = await Users.findOne({
        attributes: ['user_id'],
        where: { username: username },
    });

    if (!user) {
        return res.status(404).send('User not found');
    }

    try {
        const followers = await Followers.findAll({
            attributes: ['who_id'],
            where: {
                whom_id: user.user_id,
            },
            raw: true,
        });
        const followerList = followers.map((x) => x.who_id).concat();
        console.log('followerList: ', followerList);
        const users = await Users.findAll({
            attributes: ['user_id', 'username', 'email'],
            where: {
                user_id: followerList,
            },
            raw: true,
        });
        return res.json({ follows: users });
    } catch (error) {
        console.error('Error fetching followers:', error);
        return res.status(500).send('Internal Server Error');
    }
});

// ------------ Route to add/delete a follower --------------
router.post('/fllws/:username', async (req, res) => {
    await update_latest(req);

    const not_from_sim_response = not_req_from_simulator(req);
    if (not_from_sim_response) {
        return res.send(not_from_sim_response);
    }

    const username = req.params.username;
    const user = await Users.findOne({
        attributes: ['user_id'],
        where: {
            username: username,
        },
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
                    username: follow,
                },
            });
            if (!follows_user) {
                return res.status(404).send('User to follow not found');
            }
            console.log('follow: ', follows_user.user_id);
            console.log('follow: ', user.user_id);
            await Followers.create({
                who_id: user.user_id,
                whom_id: follows_user.user_id,
            });
            return res.sendStatus(204);

            // ------------ CASE UNFOLLOW USER ------------------
        } else if (unfollow) {
            const unfollows_user = await Users.findOne({
                attributes: ['user_id'],
                where: {
                    username: unfollow,
                },
            });
            if (!unfollows_user) {
                return res.status(404).send('User to unfollow not found');
            }

            await Followers.destroy({
                where: {
                    who_id: user.user_id,
                    whom_id: unfollows_user.user_id,
                },
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

module.exports = router;
