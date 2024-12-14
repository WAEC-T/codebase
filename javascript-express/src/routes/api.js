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
    getUserFollowers,
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
        return req.status(400).send('Username not provided!');
    }

    if (getUserByName(username)) {
        return res.status(404).send('User not found');
    }

    const messagesAmount =
        parseInt(req.query.no) || API_DEFAULT_MESSAGES_AMOUNT;

    const messages = formatMessages(
        getUserMessages(username, messagesAmount),
        true
    );

    res.json(messages);
});

router.post('/msgs/:username', async (req, res) => {
    const unauthorized = validateRequest(req);

    if (unauthorized) {
        return res.status(401).send(unauthorized);
    }

    const username = req.params.username;

    const { content } = req.body;

    if (!username || !content) {
        req.status(400).send('Username or content not provided');
    }

    const userId = getUserIdByName(username);

    if (!userId) {
        req.status(404).send('User not found');
    }

    const success = await createMessage(userId, content);

    if (success) {
        res.sendStatus(204);
    } else {
        res.status(500).send('Internal Server Error');
    }
});

router.get('/fllws/:username', async (req, res) => {
    const unauthorized = validateRequest(req);

    if (unauthorized) {
        return res.status(401).send(unauthorized);
    }

    const username = req.params.username;

    if (!username) {
        req.status(400).send('Username not provided');
    }

    const userId = getUserIdByName(username);

    if (!userId) {
        req.status(404).send('User not found');
    }

    const followers = getUserFollowers(userId);

    res.json({ follows: followers });
});

router.post('/fllws/:username', async (req, res) => {
    const unauthorized = validateRequest(req);

    if (unauthorized) {
        return res.status(401).send(unauthorized);
    }

    const username = req.params.username;

    if (!username) {
        req.status(400).send('Username not provided');
    }

    const followerId = getUserIdByName(username);

    if (!followerId) {
        req.status(404).send('User not found');
    }

    const { follow, unfollow } = req.body;

    if (follow && unfollow) {
        req.status(400).send('You should pass only one operation');
    }

    if (follow) {
        const followedId = getUserIdByName(follow);
        if (!followedId) {
            return res.status(404).send('User to follow not found');
        }

        console.log('follow: ', followerId);
        console.log('follow: ', followedId);

        const success = await followUser(followerId, followedId);

        if (success) {
            return res.sendStatus(204);
        } else {
            return res.status(500).send('Internal Server Error');
        }
    } else if (unfollow) {
        const unfollowedId = getUserIdByName(unfollow);

        if (!unfollowedId) {
            return res.status(404).send('User to follow not found');
        }

        const success = await unfollowUser(followerId, unfollowedId);

        if (success) {
            return res.sendStatus(204);
        } else {
            return res.status(500).send('Internal Server Error');
        }
    } else {
        req.status(400).send(
            'Follow or unfollow should be provided in the request body'
        );
    }
});

module.exports = router;
