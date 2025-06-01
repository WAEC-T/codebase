const express = require('express');

const {
    updateLatest,
    getLatest,
    createNewUser,
    getUserIdByName,
    getPublicMessages,
    getUserMessages,
    getUserByName,
    createMessage,
    getUserFollows,
    followUser,
    unfollowUser,
} = require('../database/repository');
const { validateRegisterFields, formatMessages } = require('../utils/utils');

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
    const latest = (await getLatest()) ?? -1;
    res.json({ latest: latest.value });
});

router.post('/register', async (req, res) => {
    const unauthorized = await validateRequest(req);

    if (unauthorized) {
        return res.status(401).send(unauthorized);
    }

    const { username, email, pwd } = req.body;

    if (await getUserIdByName(username)) {
        return res.status(400).send('The username is already taken');
    }

    let error = await validateRegisterFields(
        username,
        email,
        pwd,
        'whatever',
        true
    );

    if (!error && (await createNewUser(username, email, pwd))) {
        res.sendStatus(204);
    } else {
        res.status(400).json({ status: 400, error_msg: error });
    }
});

router.get('/msgs', async (req, res) => {
    const unauthorized = await validateRequest(req);

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
    const unauthorized = await validateRequest(req);

    if (unauthorized) {
        return res.status(401).send(unauthorized);
    }

    const username = req.params.username;

    if (!username) {
        return res.status(400).send('Username not provided!');
    }

    const user = await getUserByName(username);

    if (!user) {
        return res.status(404).send('User not found');
    }

    const messagesAmount =
        parseInt(req.query.no) || API_DEFAULT_MESSAGES_AMOUNT;

    const messages = formatMessages(
        await getUserMessages(user.user_id, messagesAmount),
        true
    );

    res.json(messages);
});

router.post('/msgs/:username', async (req, res) => {
    const unauthorized = await validateRequest(req);

    if (unauthorized) {
        return res.status(401).send(unauthorized);
    }

    const username = req.params.username;

    const { content } = req.body;

    if (!username || !content) {
        return res.status(400).send('Username or content not provided');
    }

    const userId = await getUserIdByName(username);

    if (!userId) {
        return res.status(404).send('User not found');
    }

    const success = await createMessage(userId, content);

    if (success) {
        return res.sendStatus(204);
    } else {
        return res.status(500).send('Internal Server Error');
    }
});

router.get('/fllws/:username', async (req, res) => {
    const unauthorized = await validateRequest(req);

    if (unauthorized) {
        return res.status(401).send(unauthorized);
    }

    const username = req.params.username;

    if (!username) {
        return res.status(400).send('Username not provided');
    }

    const userId = await getUserIdByName(username);

    if (!userId) {
        return res.status(404).send('User not found');
    }

    const followers = await getUserFollows(userId);

    res.json({ follows: followers });
});

router.post('/fllws/:username', async (req, res) => {
    const unauthorized = await validateRequest(req);

    if (unauthorized) {
        return res.status(401).send(unauthorized);
    }

    const username = req.params.username;

    if (!username) {
        return res.status(400).send('Username not provided');
    }

    const followerId = await getUserIdByName(username);

    if (!followerId) {
        return res.status(404).send('User not found');
    }

    const { follow, unfollow } = req.body;

    if (follow && unfollow) {
        return res.status(400).send('You should pass only one operation');
    }

    if (follow) {
        const followedId = await getUserIdByName(follow);
        if (!followedId) {
            return res.status(404).send('User to follow not found');
        }

        const success = await followUser(followerId, followedId);

        if (success) {
            return res.sendStatus(204);
        } else {
            return res.status(500).send('Internal Server Error');
        }
    } else if (unfollow) {
        const unfollowedId = await getUserIdByName(unfollow);

        if (!unfollowedId) {
            return res.status(404).send('User to follow not found');
        }

        const success = await unfollowUser(followerId, unfollowedId);

        if (success) {
            res.sendStatus(204);
        } else {
            res.status(500).send('Internal Server Error');
        }
    } else {
        res.status(400).send(
            'Follow or unfollow should be provided in the request body'
        );
    }
});

module.exports = router;
